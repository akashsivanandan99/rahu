// Package parser implements a recursive descent parser for Python source code.
//
// The parser transforms a stream of tokens from the lexer into an Abstract Syntax Tree (AST)
// that represents the syntactic structure of Python programs. It handles Python's indentation-based
// syntax through INDENT and DEDENT tokens provided by the lexer.
//
// Features:
//   - Recursive descent parsing with operator precedence (Pratt parsing)
//   - Support for Python statements: assignments, function definitions, control flow (if/elif/else,
//     for, while), return, break, and continue
//   - Support for Python expressions: arithmetic operations, comparisons (with chaining),
//     boolean operations (and/or), unary operations, function calls, lists, and tuples
//   - Proper operator precedence handling for arithmetic, comparison, and boolean operators
//   - Position tracking for all AST nodes (line and column information)
//   - Advanced features: function default arguments, tuple unpacking in for loops, comparison chaining
//
// The parser produces an AST consisting of Statement and Expression nodes. Each node type
// implements the appropriate interface and includes position information for error reporting
// and language server features.
//
// Example usage:
//
//	input := `
//	def fibonacci(n):
//	    if n <= 1:
//	        return n
//	    return fibonacci(n-1) + fibonacci(n-2)
//	`
//	p := parser.New(input)
//	module := p.Parse()
//
// The resulting AST can be traversed for semantic analysis, type checking, or code generation.
package parser

import (
	"fmt"

	"rahu/lexer"
)

const (
	LOWEST = iota
	OR
	AND
	NOT
	COMPARE
	SUM
	PRODUCT
	PREFIX
)

type Parser struct {
	lexer   *lexer.Lexer
	current lexer.Token
	peek    lexer.Token
}

func New(input string) *Parser {
	l := lexer.New(input)
	p := &Parser{
		lexer: l,
	}

	p.current = l.NextToken()
	p.peek = l.NextToken()
	return p
}

func (p *Parser) advance() {
	p.current = p.peek
	p.peek = p.lexer.NextToken()
}

func (p *Parser) advanceBy(count int) {
	for range count {
		p.advance()
	}
}

func (p *Parser) Parse() *Module {
	statements := []Statement{}
	for p.current.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}
	return &Module{Body: statements}
}

func (p *Parser) parseStatement() Statement {
	if p.current.Type == lexer.NAME && p.peek.Type == lexer.EQUAL {
		return p.parseAssignment()
	}

	if p.current.Type == lexer.IF {
		return p.parseIf()
	}

	if p.current.Type == lexer.NAME || p.current.Type == lexer.NUMBER ||
		p.current.Type == lexer.STRING || p.current.Type == lexer.LPAR {
		expr := p.parseExpression(LOWEST)
		if p.current.Type == lexer.NEWLINE {
			p.advance()
		}
		return &ExprStmt{Value: expr, Pos: expr.Position()}
	}

	if p.current.Type == lexer.DEF {
		funcdef := p.parseFunc()
		return funcdef
	}

	if p.current.Type == lexer.BREAK {
		pos := Range{
			Start: Position{Line: p.current.Line, Col: p.current.Col},
			End:   Position{Line: p.current.Line, Col: p.current.EndCol},
		}
		p.advance()
		if p.current.Type == lexer.NEWLINE {
			p.advance()
		}
		return &Break{Pos: pos}
	}

	if p.current.Type == lexer.CONTINUE {
		pos := Range{
			Start: Position{Line: p.current.Line, Col: p.current.Col},
			End:   Position{Line: p.current.Line, Col: p.current.EndCol},
		}
		p.advance()
		if p.current.Type == lexer.NEWLINE {
			p.advance()
		}
		return &Continue{Pos: pos}
	}

	if p.current.Type == lexer.RETURN {
		return p.parseReturn()
	}

	if p.current.Type == lexer.FOR {
		return p.parseFor()
	}

	if p.current.Type == lexer.WHILE {
		return p.parseWhile()
	}

	p.advance()
	return nil
}

func (p *Parser) parseForTarget() Expression {
	if p.current.Type != lexer.NAME {
		panic(`expected variable name`)
	}

	first := &Name{
		Id: p.current.Literal,
		Pos: Range{
			Start: Position{Line: p.current.Line, Col: p.current.Col},
			End:   Position{Line: p.current.Line, Col: p.current.EndCol},
		},
	}
	p.advance()

	if p.current.Type == lexer.COMMA {
		targets := []Expression{first}

		for p.current.Type == lexer.COMMA {
			p.advance()

			if p.current.Type != lexer.NAME {
				panic(`expected variable name`)
			}
			targets = append(
				targets, &Name{
					Id: p.current.Literal,
					Pos: Range{
						Start: Position{Line: p.current.Line, Col: p.current.Col},
						End:   Position{Line: p.current.Line, Col: p.current.EndCol},
					},
				})
			p.advance()
		}
		lastName := targets[len(targets)-1]
		return &Tuple{
			Elts: targets,
			Pos: Range{
				Start: first.Pos.Start,
				End:   lastName.Position().End,
			},
		}
	}
	return first
}

func (p *Parser) parseWhile() Statement {
	startPos := Position{Line: p.current.Line, Col: p.current.Col}
	p.advance()
	whileStmt := &WhileLoop{}

	whileStmt.Test = p.parseExpression(LOWEST)

	if p.current.Type != lexer.COLON {
		panic(`expected ':' after while condition`)
	}
	p.advance()

	if p.current.Type != lexer.NEWLINE {
		panic(`expected newline after ':'`)
	}
	p.advance()

	if p.current.Type != lexer.INDENT {
		panic(`expected indent after while:`)
	}
	p.advance()

	body := []Statement{}
	for p.current.Type != lexer.DEDENT && p.current.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			body = append(body, stmt)
		}
	}
	whileStmt.Body = body

	endPos := Position{Line: p.current.Line, Col: p.current.Col}
	if p.current.Type == lexer.DEDENT {
		p.advance()
	}
	whileStmt.Pos = Range{Start: startPos, End: endPos}

	return whileStmt
}

func (p *Parser) parseFor() Statement {
	startPos := Position{Line: p.current.Line, Col: p.current.Col}
	forStmt := &For{Pos: Range{Start: startPos}}
	p.advance()
	forStmt.Target = p.parseForTarget()

	if p.current.Type != lexer.IN {
		panic(`expected 'in' after loop variable`)
	}

	p.advance()

	forStmt.Iter = p.parseExpression(LOWEST)

	if p.current.Type != lexer.COLON {
		fmt.Printf("%vfailed for token v\n", p.current)
		panic(`expected ':' after for clause`)
	}
	p.advance()

	if p.current.Type != lexer.NEWLINE {
		panic(`expected newline after ':'`)
	}

	p.advance()

	if p.current.Type != lexer.INDENT {
		panic(`expected indent after for statement`)
	}

	p.advance()

	body := []Statement{}

	for p.current.Type != lexer.DEDENT && p.current.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			body = append(body, stmt)
		}
	}

	forStmt.Body = body

	if p.current.Type == lexer.DEDENT {
		p.advance()
	}
	forStmt.Pos.End = Position{Line: p.current.Line, Col: p.current.Col}
	return forStmt
}

func (p *Parser) parseReturn() Statement {
	startPos := Position{Line: p.current.Line, Col: p.current.Col}
	p.advance()
	if p.current.Type == lexer.NEWLINE || p.current.Type == lexer.EOF {
		endPos := Position{Line: p.current.Line, Col: p.current.Col}
		p.advance()
		return &Return{Value: nil, Pos: Range{Start: startPos, End: endPos}}
	}

	value := p.parseExpression(LOWEST)
	endPos := Position{Line: p.current.Line, Col: p.current.Col}

	if p.current.Type == lexer.NEWLINE {
		p.advance()
	}

	return &Return{Value: value, Pos: Range{Start: startPos, End: endPos}}
}

func (p *Parser) parseFunc() Statement {
	startPos := Position{Line: p.current.Line, Col: p.current.Col}
	p.advance()
	funcDef := &FunctionDef{}

	if p.current.Type != lexer.NAME {
		fmt.Printf("%v\n", p.current)
		panic(`expected function name after "def"`)
	}

	funcDef.Name = p.current.Literal
	p.advance() // advance past func name

	if p.current.Type != lexer.LPAR {
		panic(`expected '(' after function name`)
	}

	p.advance() // advanced past left parantheses

	args := []FuncArg{}
	seenDefault := false

	if p.current.Type != lexer.RPAR {
		for {
			if p.current.Type != lexer.NAME {
				panic("expected param name")
			}

			arg := FuncArg{
				Name: p.current.Literal,
			}

			p.advance()

			if p.current.Type == lexer.EQUAL {
				seenDefault = true
				p.advance()
				arg.Default = p.parseExpression(LOWEST)
			} else if seenDefault {
				panic("non-default argumentn follows default argument")
			}

			args = append(args, arg)

			if p.current.Type == lexer.COMMA {
				p.advance()
				continue
			}
			break
		}
	}

	funcDef.Args = args

	if p.current.Type != lexer.RPAR {
		panic(`expected ')' after params`)
	}

	p.advance()

	if p.current.Type != lexer.COLON {
		panic(`expected ':' after function signature`)
	}
	p.advance()

	if p.current.Type != lexer.NEWLINE {
		panic(`expected newline after ':'`)
	}

	p.advance()

	if p.current.Type != lexer.INDENT {
		panic(`expected indentation after function signature`)
	}
	p.advance()

	funcBody := []Statement{}

	for p.current.Type != lexer.DEDENT && p.current.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			funcBody = append(funcBody, stmt)
		}
	}

	if p.current.Type == lexer.DEDENT {
		p.advance() // skip past dedent
	}

	funcDef.Body = funcBody
	endPos := Position{Line: p.current.Line, Col: p.current.Col}
	funcDef.Pos = Range{Start: startPos, End: endPos}

	return funcDef
}

func (p *Parser) parseIf() Statement {
	startPos := Position{
		Line: p.current.Line,
		Col:  p.current.Col,
	}

	p.advance()
	ifExpr := &If{}
	testCond := p.parseExpression(LOWEST)

	ifExpr.Test = testCond

	if p.current.Type != lexer.COLON {
		fmt.Printf("panicked with curr type : %v and testCond: %v\n", p.current, testCond)
		panic(`expected ':' after if condition`)
	}

	p.advance() // skip colon

	if p.current.Type != lexer.NEWLINE {
		panic(`expected newline after ":"`)
	}
	p.advance() // moving ahead of newline

	if p.current.Type != lexer.INDENT {
		panic("expected indentation block after if condtion")
	}
	p.advance() // moving ahead of indent

	body := []Statement{}

	for p.current.Type != lexer.DEDENT && p.current.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			body = append(body, stmt)
		}
	}

	ifExpr.Body = body

	endPos := Position{
		Line: p.current.Line,
		Col:  p.current.Col,
	}

	if p.current.Type == lexer.DEDENT {
		p.advance() // skip past dedent
	}

	orelse := []Statement{}
	if p.current.Type == lexer.ELIF {
		elifStmt := p.parseIf()
		orelse = append(orelse, elifStmt)

		if ifNode, ok := elifStmt.(*If); ok {
			endPos = ifNode.Pos.End
		}
	} else if p.current.Type == lexer.ELSE {
		p.advance()

		endPos.Col = p.current.Col
		endPos.Line = p.current.Line
		if p.current.Type != lexer.COLON {
			panic(`expected ':' after else`)
		}

		p.advance() // advancing past 'else' block

		if p.current.Type != lexer.NEWLINE {
			panic(`expected NEWLINE after 'else:'`)
		}

		p.advance() // move past newline

		if p.current.Type != lexer.INDENT {
			panic(`expected indent block after "else:"`)
		}

		p.advance() // move past indent

		for p.current.Type != lexer.DEDENT && p.current.Type != lexer.EOF {
			stmt := p.parseStatement()
			if stmt != nil {
				orelse = append(orelse, stmt)
			}
		}

		endPos = Position{
			Line: p.current.Line,
			Col:  p.current.Col,
		}

		if p.current.Type == lexer.DEDENT {
			p.advance() // advance past dedent
		}
	}

	ifExpr.Orelse = orelse
	ifExpr.Pos = Range{
		Start: startPos,
		End:   endPos,
	}

	return ifExpr
}

func (p *Parser) parseAssignment() Statement {
	assgnStart := Position{Line: p.current.Line, Col: p.current.Col}
	targetStart := Position{Line: p.current.Line, Col: p.current.Col}

	target := &Name{Id: p.current.Literal}

	p.advance()

	target.Pos = Range{
		Start: targetStart,
		End:   Position{Line: p.current.Line, Col: p.current.Col},
	}

	p.advance()

	// moved past '='

	value := p.parseExpression(LOWEST)

	assgnEnd := Position{
		Line: p.current.Line,
		Col:  p.current.Col,
	}

	if p.current.Type == lexer.NEWLINE {
		p.advance()
	}

	return &Assign{
		Targets: []Expression{target},
		Value:   value,
		Pos:     Range{Start: assgnStart, End: assgnEnd},
	}
}

func (p *Parser) parseExpression(minBP int) Expression {
	left := p.parsePrimary()

	for {
		bp := infixBindingPower(p.current.Type)
		if bp <= minBP {
			break
		}

		opTok := p.current

		if isCompareOp(opTok.Type) {
			startPos := left.Position().Start
			ops := []CompareOp{}
			rights := []Expression{}

			for isCompareOp(p.current.Type) {
				op := tokenTypeToCompareOp(p.current.Type)
				p.advance()

				right := p.parseExpression(COMPARE + 1)
				ops = append(ops, op)
				rights = append(rights, right)
			}
			lastRight := rights[len(rights)-1]
			endPos := lastRight.Position().End

			left = &Compare{
				Left:  left,
				Ops:   ops,
				Right: rights,
				Pos:   Range{Start: startPos, End: endPos},
			}
			continue
		}

		if opTok.Type == lexer.AND || opTok.Type == lexer.OR {
			p.advance()
			right := p.parseExpression(bp)

			if boolOp, ok := left.(*BooleanOp); ok {
				if (opTok.Type == lexer.AND && boolOp.Operator == And) || (opTok.Type == lexer.OR && boolOp.Operator == Or) {
					boolOp.Values = append(boolOp.Values, right)
					continue
				}
			}

			op := And
			if opTok.Type == lexer.OR {
				op = Or
			}
			left = &BooleanOp{
				Operator: op,
				Values:   []Expression{left, right},
				Pos: Range{
					Start: left.Position().Start,
					End:   right.Position().End,
				},
			}
			continue
		}
		p.advance()

		right := p.parseExpression(bp)

		left = &BinOp{
			Left:  left,
			Op:    p.tokenTypeToOperator(opTok.Type),
			Right: right,
			Pos: Range{
				Start: left.Position().Start,
				End:   right.Position().End,
			},
		}
	}

	return left
}

func isCompareOp(t lexer.TokenType) bool {
	switch t {
	case lexer.EQEQUAL, lexer.NOTEQUAL, lexer.LESS, lexer.LESSEQUAL, lexer.GREATER, lexer.GREATEREQUAL:
		return true
	default:
		return false
	}
}

func infixBindingPower(t lexer.TokenType) int {
	switch t {
	case lexer.PLUS, lexer.MINUS:
		return SUM
	case lexer.STAR, lexer.SLASH, lexer.DOUBLESLASH, lexer.PERCENT:
		return PRODUCT
	case lexer.EQEQUAL, lexer.NOTEQUAL, lexer.LESS, lexer.LESSEQUAL, lexer.GREATER, lexer.GREATEREQUAL:
		return COMPARE
	case lexer.OR:
		return OR
	case lexer.AND:
		return AND
	default:
		return LOWEST
	}
}

func (p *Parser) parsePrimary() Expression {
	switch p.current.Type {
	case lexer.NUMBER:
		n := &Number{
			Value: p.current.Literal,
			Pos: Range{
				Start: Position{Line: p.current.Line, Col: p.current.Col},
				End:   Position{Line: p.current.Line, Col: p.current.EndCol},
			},
		}
		p.advance()
		return n
	case lexer.TRUE:
		ret := &Boolean{
			Value: true,
			Pos: Range{
				Start: Position{Line: p.current.Line, Col: p.current.Col},
				End:   Position{Line: p.current.Line, Col: p.current.EndCol},
			},
		}
		p.advance()
		return ret

	case lexer.FALSE:
		ret := &Boolean{
			Value: false,
			Pos: Range{
				Start: Position{Line: p.current.Line, Col: p.current.Col},
				End:   Position{Line: p.current.Line, Col: p.current.EndCol},
			},
		}
		p.advance()
		return ret
	case lexer.NAME:
		n := &Name{
			Id: p.current.Literal,
			Pos: Range{
				Start: Position{Line: p.current.Line, Col: p.current.Col},
				End:   Position{Line: p.current.Line, Col: p.current.EndCol},
			},
		}
		p.advance()

		if p.current.Type == lexer.LPAR {
			return p.parseCall(n)
		}
		return n
	case lexer.LPAR:
		p.advance()
		expr := p.parseExpression(LOWEST)
		if p.current.Type != lexer.RPAR {
			panic(`expected )`)
		}
		p.advance()
		return expr
	case lexer.LSQB:
		return p.parseList()
	case lexer.STRING:
		s := &String{
			Value: p.current.Literal,
			Pos: Range{
				Start: Position{Line: p.current.Line, Col: p.current.Col},
				End:   Position{Line: p.current.Line, Col: p.current.EndCol},
			},
		}
		p.advance()
		return s

	case lexer.MINUS:
		startPos := Position{Line: p.current.Line, Col: p.current.Col}
		p.advance()
		operand := p.parseExpression(PREFIX)
		endPos := operand.Position().End
		return &UnaryOp{
			Op:      USub,
			Operand: operand,
			Pos:     Range{Start: startPos, End: endPos},
		}

	case lexer.PLUS:
		startPos := Position{Line: p.current.Line, Col: p.current.Col}
		p.advance()
		operand := p.parseExpression(PREFIX)
		endPos := operand.Position().End
		return &UnaryOp{
			Op:      UAdd,
			Operand: operand,
			Pos:     Range{Start: startPos, End: endPos},
		}

	case lexer.NOT:
		startPos := Position{Line: p.current.Line, Col: p.current.Col}
		p.advance()
		expr := p.parseExpression(PREFIX)
		endPos := expr.Position().End
		return &UnaryOp{
			Op:      Not,
			Operand: expr,
			Pos:     Range{Start: startPos, End: endPos},
		}

	}
	panic(fmt.Sprintf("unexpected token %v\n", p.current))
}

func (p *Parser) parseList() Expression {
	startPos := Position{Line: p.current.Line, Col: p.current.Col}
	p.advance()
	elts := []Expression{}

	if p.current.Type != lexer.RSQB {
		elts = append(elts, p.parseExpression(LOWEST))

		for p.current.Type == lexer.COMMA {
			p.advance()

			if p.current.Type == lexer.RSQB {
				break
			}

			elts = append(elts, p.parseExpression(LOWEST))
		}
	}

	if p.current.Type != lexer.RSQB {
		panic(`expected ']' after list elements`)
	}

	endPos := Position{
		Line: p.current.Line,
		Col:  p.current.Col,
	}
	p.advance()
	return &List{Elts: elts, Pos: Range{Start: startPos, End: endPos}}
}

func (p *Parser) parseCall(funcExpr Expression) Expression {
	startPos := funcExpr.(*Name).Pos.Start
	p.advance()
	args := []Expression{}

	if p.current.Type != lexer.RPAR {
		args = append(args, p.parseExpression(LOWEST))

		for p.current.Type == lexer.COMMA {
			p.advance()
			args = append(args, p.parseExpression(LOWEST))
		}
	}

	if p.current.Type != lexer.RPAR {
		panic("expected ) after function arguments")
	}
	endPos := Position{Line: p.current.Line, Col: p.current.Col}
	p.advance()

	return &Call{
		Func: funcExpr,
		Args: args,
		Pos:  Range{Start: startPos, End: endPos},
	}
}

func (p *Parser) isOperator(t lexer.TokenType) bool {
	return t == lexer.PLUS || t == lexer.MINUS || t == lexer.STAR || t == lexer.SLASH || t == lexer.DOUBLESLASH || t == lexer.PERCENT
}

func tokenTypeToCompareOp(t lexer.TokenType) CompareOp {
	switch t {
	case lexer.EQEQUAL:
		return Eq
	case lexer.NOTEQUAL:
		return NotEq
	case lexer.LESS:
		return Lt
	case lexer.LESSEQUAL:
		return LtE
	case lexer.GREATER:
		return Gt
	case lexer.GREATEREQUAL:
		return GtE
	default:
		panic("invalid compare operator")
	}
}

func (p *Parser) tokenTypeToOperator(t lexer.TokenType) Operator {
	switch t {
	case lexer.PLUS:
		return Add
	case lexer.MINUS:
		return Sub
	case lexer.STAR:
		return Mult
	case lexer.SLASH:
		return Div
	case lexer.DOUBLESLASH:
		return FloorDiv
	case lexer.PERCENT:
		return Mod
	default:
		return Add
	}
}

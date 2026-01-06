package parser

import (
	"rahu/lexer"
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

	p.advance()
	return nil
}

func (p *Parser) parseAssignment() Statement {
	if p.current.Type == lexer.INDENT || p.current.Type == lexer.DEDENT {
		p.advance()
		return nil
	}
	target := &Name{Id: p.current.Literal}
	// moving past name and equal
	p.advanceBy(2)

	value := p.parseExpression()

	if p.current.Type == lexer.NEWLINE {
		p.advance()
	}

	return &Assign{
		Targets: []Expression{target},
		Value:   value,
	}
}

func (p *Parser) parseExpression() Expression {
	left := p.parsePrimary()

	if p.isOperator(p.current.Type) {
		op := p.tokenTypeToOperator(p.current.Type)
		p.advance()

		right := p.parsePrimary()
		return &BinOp{
			Left:  left,
			Op:    op,
			Right: right,
		}
	}

	return left
}

func (p *Parser) parsePrimary() Expression {
	if p.current.Type == lexer.NUMBER {
		num := &Number{Value: p.current.Literal}
		p.advance()
		return num
	}

	if p.current.Type == lexer.NAME {
		name := &Name{Id: p.current.Literal}
		p.advance()
		return name
	}
	return nil
}

func (p *Parser) isOperator(t lexer.TokenType) bool {
	return t == lexer.PLUS || t == lexer.MINUS || t == lexer.STAR || t == lexer.SLASH || t == lexer.DOUBLESLASH
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
		return Mod
	default:
		return Add
	}
}

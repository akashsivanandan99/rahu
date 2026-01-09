package parser

type Node any

type Operator int

type Position struct {
	Line int
	Col int
}

type Range struct {
	Start Position
	End Position
}

const (
	Add Operator = iota
	Sub
	Mult
	Div
	FloorDiv
	Mod
	Pow
)

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Number struct {
	Value string
	Pos Range
}

func (n *Number) expressionNode() {}

type Name struct {
	Id string
	Pos Range
}

func (n *Name) expressionNode() {}

type BinOp struct {
	Left  Expression
	Op    Operator
	Right Expression
	Pos Range
}

func (b *BinOp) expressionNode() {}

type Assign struct {
	Targets []Expression
	Value   Expression
	Pos Range
}

func (a *Assign) statementNode() {}

type Module struct {
	Body []Statement
}

type FuncArg struct {
	Name    string
	Default Expression
	Pos Range
}

type FunctionDef struct {
	Name string
	Args []FuncArg
	Body []Statement
	Pos Range
}

func (f *FunctionDef) statementNode() {}

type Return struct {
	Value Expression
	Pos Range
}

func (r *Return) statementNode() {}

type If struct {
	Test   Expression
	Body   []Statement
	Orelse []Statement
	Pos Range
}

func (i *If) statementNode() {}

type Tuple struct {
	Elts []Expression
	Pos Range
}

func (t *Tuple) expressionNode() {}

type For struct {
	Target Expression
	Iter   Expression
	Body   []Statement
	Pos Range
}

func (f *For) statementNode() {}

type WhileLoop struct {
	Test Expression
	Body []Statement
	Pos Range
}

func (w *WhileLoop) statementNode() {}

type Call struct {
	Func Expression
	Args []Expression
	Pos Range
}

func (c *Call) expressionNode() {}

type Compare struct {
	Left  Expression
	Ops   []CompareOp
	Right []Expression
	Pos Range
}

func (c *Compare) expressionNode() {}

type CompareOp int

const (
	Eq    CompareOp = iota // ==
	NotEq                  // !=
	Lt                     // <
	LtE                    // <=
	Gt                     // >
	GtE                    // >=
)

type UnaryOperator int

const (
	UAdd UnaryOperator = iota // +x
	USub                      // -x
	Not                       // not x
)


type UnaryOp struct {
	Op      UnaryOperator
	Operand Expression
	Pos Range
}

func (u *UnaryOp) expressionNode() {}

type String struct {
	Value string
	Pos Range
}

func (s *String) expressionNode() {}

type ExprStmt struct {
	Value Expression
	Pos Range
}

func (e *ExprStmt) statementNode() {}

type Boolean struct {
	Value bool
	Pos Range
}

func (b *Boolean) expressionNode() {}

type BooleanOperator int

const (
	And BooleanOperator = iota
	Or
)

type BooleanOp struct {
	Operator BooleanOperator
	Values   []Expression
	Pos Range
}

func (bo *BooleanOp) expressionNode() {}

type List struct {
	Elts []Expression
	Pos Range
}

func (l *List) expressionNode() {}

type Break struct{
	Pos Range
}

func (b *Break) statementNode() {}

type Continue struct{
	Pos Range
}

func (c *Continue) statementNode() {}

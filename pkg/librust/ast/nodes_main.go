package ast

type Node interface {
	Accept(v RusterBaseVisitor)
}
type Statement Node
type Expression Statement
type ExpressionWithBlock Expression
type Item Statement
type Type Item
type Pattern Item
type NonRangePattern Pattern

type Literal string
type PathSegments []string

const (
	String  Literal = "str"
	Char    Literal = "char"
	Integer Literal = "int"
	Boolean Literal = "bool"
)

type Crate struct {
	Node
	Items []Item
}

func (c *Crate) Accept(v RusterBaseVisitor) {
	v.VisitCrate(c)
}

type BlockExpression struct {
	Node
	Statements []Statement
	Expr       Expression
}

func (be *BlockExpression) Accept(v RusterBaseVisitor) {
	v.VisitBlockExpression(be)
}

type UseDecl struct {
	Item
	All  bool
	Path SimplePath
}

func (ud *UseDecl) Accept(v RusterBaseVisitor) {
	v.VisitUseDecl(ud)
}

type SimplePath struct {
	Node
	Segments PathSegments
}

func (sp *SimplePath) Accept(v RusterBaseVisitor) {
	v.VisitSimplePath(sp)
}

type Function struct {
	Item
	ID         string
	ReturnType Type
	Params     []Parameter
	Body       BlockExpression
}

func (f *Function) Accept(v RusterBaseVisitor) {
	v.VisitFunction(f)
}

type Parameter struct {
	Node
	ID      interface{}
	VarType Type
}

func (p *Parameter) Accept(v RusterBaseVisitor) {
	v.VisitParameter(p)
}

type LetStatement struct {
	Statement
	Ptrn    Pattern
	VarType Type
	Expr    Expression
}

func (ls *LetStatement) Accept(v RusterBaseVisitor) {
	v.VisitLetStatement(ls)
}

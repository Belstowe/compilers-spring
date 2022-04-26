package ast_nodes

import "github.com/Compiler2022/compilers-1-Belstowe/parser"

type Crate struct {
	Items []Item
}

type Node interface {
	Accept(parser.RustParserVisitor)
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

type BlockExpression struct {
	Node
	Statements []Statement
	Expr       Expression
}

type UseDecl struct {
	Item
	All  bool
	Path SimplePath
}

type SimplePath struct {
	Node
	Segments PathSegments
}

type Function struct {
	Item
	ID         string
	ReturnType Type
	Params     []Parameter
	Body       BlockExpression
}

type Parameter struct {
	Node
	ID      interface{}
	VarType Type
}

type LetStatement struct {
	Statement
	Ptrn    Pattern
	VarType Type
	Expr    Expression
}

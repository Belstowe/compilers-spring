package ast

import (
	"reflect"

	"github.com/iancoleman/strcase"
)

type Node interface{}
type Statement Node
type Expression Statement
type ExpressionWithBlock Expression
type Item Statement
type Type Item
type Pattern Item
type NonRangePattern Pattern

type Literal string
type PathSegments []string
type SimplePath []string

type Terminal struct {
	n interface{}
}

func (t Terminal) MarshalYAML() (interface{}, error) {
	m := make(map[string]interface{})
	m[strcase.ToCamel(reflect.TypeOf(t.n).String())] = t.n
	return m, nil
}

func Wrap(n interface{}) Terminal {
	return Terminal{n: n}
}

const (
	String  Literal = "str"
	Char    Literal = "char"
	Integer Literal = "int"
	Boolean Literal = "bool"
)

type Crate struct {
	Items []Item
}

func (c *Crate) Accept(v RusterBaseVisitor) {
	v.VisitCrate(c)
}

type BlockExpression struct {
	Statements []Terminal
	Expr       Expression `yaml:"return,omitempty"`
}

func (be *BlockExpression) Accept(v RusterBaseVisitor) {
	v.VisitBlockExpression(be)
}

type UseDecl struct {
	All  bool       `yaml:"includeAllItems,omitempty"`
	Path SimplePath `yaml:"path,flow"`
}

func (ud *UseDecl) Accept(v RusterBaseVisitor) {
	v.VisitUseDecl(ud)
}

func (sp *SimplePath) Accept(v RusterBaseVisitor) {
	v.VisitSimplePath(sp)
}

type Function struct {
	ID         string          `yaml:"ID"`
	ReturnType Type            `yaml:"ReturnType,flow,omitempty"`
	Params     []Parameter     `yaml:"Params,omitempty"`
	Body       BlockExpression `yaml:"Body"`
}

func (f *Function) Accept(v RusterBaseVisitor) {
	v.VisitFunction(f)
}

type Parameter struct {
	ID      string `yaml:"id"`
	VarType Type   `yaml:"type,flow"`
}

func (p *Parameter) Accept(v RusterBaseVisitor) {
	v.VisitParameter(p)
}

type LetStatement struct {
	Ptrn    Pattern    `yaml:"assignee"`
	VarType Type       `yaml:"type,attr,omitempty"`
	Expr    Expression `yaml:"expression"`
}

func (ls *LetStatement) Accept(v RusterBaseVisitor) {
	v.VisitLetStatement(ls)
}

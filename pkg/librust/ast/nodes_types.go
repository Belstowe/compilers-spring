package ast

type TypePath []TypePathSegment

func (tp *TypePath) Accept(v RusterBaseVisitor) {
	v.VisitTypePath(tp)
}

type TypePathSegment struct {
	ID string `yaml:"id"`
	Fn Node   `yaml:"function-body,omitempty"`
}

func (tps *TypePathSegment) Accept(v RusterBaseVisitor) {
	v.VisitTypePathSegment(tps)
}

type TypePathFunction struct {
	Inputs     []Terminal
	ReturnType Type
}

func (tpf *TypePathFunction) Accept(v RusterBaseVisitor) {
	v.VisitTypePathFunction(tpf)
}

type ParenthesizedType struct {
	VarType Type `yaml:"type"`
}

func (pt *ParenthesizedType) Accept(v RusterBaseVisitor) {
	v.VisitParenthesizedType(pt)
}

type PointerType struct {
	IsMut   bool `yaml:"mutable,omitempty"`
	VarType Type `yaml:"type"`
}

func (pt *PointerType) Accept(v RusterBaseVisitor) {
	v.VisitPointerType(pt)
}

type ReferenceType struct {
	IsMut   bool `yaml:"mutable,omitempty"`
	VarType Type `yaml:"type"`
}

func (rt *ReferenceType) Accept(v RusterBaseVisitor) {
	v.VisitReferenceType(rt)
}

type TupleType struct {
	Types []Terminal `yaml:"type"`
}

func (tt *TupleType) Accept(v RusterBaseVisitor) {
	v.VisitTupleType(tt)
}

type ArrayType struct {
	VarType Type       `yaml:"type"`
	Expr    Expression `yaml:"patternExpression"`
}

func (at *ArrayType) Accept(v RusterBaseVisitor) {
	v.VisitArrayType(at)
}

type SliceType struct {
	VarType Type `yaml:"type"`
}

func (st *SliceType) Accept(v RusterBaseVisitor) {
	v.VisitSliceType(st)
}

type NeverType struct{ Type }

func (nt *NeverType) Accept(v RusterBaseVisitor) {
	v.VisitNeverType(nt)
}

type InferredType struct{ Type }

func (it *InferredType) Accept(v RusterBaseVisitor) {
	v.VisitInferredType(it)
}
package ast

type TypePath struct {
	Type
	Segments []TypePathSegment
}

func (tp *TypePath) Accept(v RusterBaseVisitor) {
	v.VisitTypePath(tp)
}

type TypePathSegment struct {
	Node
	ID string
	Fn Node
}

func (tps *TypePathSegment) Accept(v RusterBaseVisitor) {
	v.VisitTypePathSegment(tps)
}

type TypePathFunction struct {
	Node
	Inputs     []Type
	ReturnType Type
}

func (tpf *TypePathFunction) Accept(v RusterBaseVisitor) {
	v.VisitTypePathFunction(tpf)
}

type ParenthesizedType struct {
	Type
	VarType Type
}

func (pt *ParenthesizedType) Accept(v RusterBaseVisitor) {
	v.VisitParenthesizedType(pt)
}

type PointerType struct {
	Type
	IsMut   bool
	VarType Type
}

func (pt *PointerType) Accept(v RusterBaseVisitor) {
	v.VisitPointerType(pt)
}

type ReferenceType struct {
	Type
	IsMutable bool
	VarType   Type
}

func (rt *ReferenceType) Accept(v RusterBaseVisitor) {
	v.VisitReferenceType(rt)
}

type TupleType struct {
	Type
	Types []Type
}

func (tt *TupleType) Accept(v RusterBaseVisitor) {
	v.VisitTupleType(tt)
}

type ArrayType struct {
	Type
	VarType Type
	Expr    Expression
}

func (at *ArrayType) Accept(v RusterBaseVisitor) {
	v.VisitArrayType(at)
}

type SliceType struct {
	Type
	VarType Type
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

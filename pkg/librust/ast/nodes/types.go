package ast_nodes

type TypePath struct {
	Type
	Segments []TypePathSegment
}

type TypePathSegment struct {
	Node
	ID string
	Fn Node
}

type TypePathFunction struct {
	Node
	Inputs     []Type
	ReturnType Type
}

type ParenthesizedType struct {
	Type
	VarType Type
}

type PointerType struct {
	Type
	IsMut   bool
	VarType Type
}

type ReferenceType struct {
	Type
	IsMutable bool
	VarType   Type
}

type TupleType struct {
	Type
	Types []Type
}

type ArrayType struct {
	Type
	VarType Type
	Expr    Expression
}

type SliceType struct {
	Type
	VarType Type
}

type NeverType struct{ Type }
type InferredType struct{ Type }

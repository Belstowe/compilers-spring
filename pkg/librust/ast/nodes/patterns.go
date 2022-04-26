package ast_nodes

type LiteralPattern struct {
	NonRangePattern
	Tp  Literal
	Val string
}

type ReferencePattern struct {
	NonRangePattern
	IsDoubleRef bool
	Ptrn        NonRangePattern
}

type IdentifierPattern struct {
	NonRangePattern
	IsRef bool
	IsMut bool
	ID    string
}

type PathPattern struct {
	NonRangePattern
	Segments []string
}

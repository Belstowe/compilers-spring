package ast

type LiteralPattern struct {
	Tp  Literal
	Val string
}

func (lp LiteralPattern) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitLiteralPattern(&lp)
}

type ReferencePattern struct {
	IsDoubleRef bool
	Ptrn        NonRangePattern
}

func (rp ReferencePattern) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitReferencePattern(&rp)
}

type IdentifierPattern struct {
	IsRef bool `yaml:"ref,omitempty"`
	IsMut bool `yaml:"mut,omitempty"`
	ID    string
}

func (ip IdentifierPattern) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitIdentifierPattern(&ip)
}

type PathPattern struct {
	Segments []string `yaml:"segments"`
}

func (pp PathPattern) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitPathPattern(&pp)
}

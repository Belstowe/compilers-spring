package ast

type LiteralPattern struct {
	NonRangePattern
	Tp  Literal
	Val string
}

func (lp *LiteralPattern) Accept(v RusterBaseVisitor) {
	v.VisitLiteralPattern(lp)
}

type ReferencePattern struct {
	NonRangePattern
	IsDoubleRef bool
	Ptrn        NonRangePattern
}

func (rp *ReferencePattern) Accept(v RusterBaseVisitor) {
	v.VisitReferencePattern(rp)
}

type IdentifierPattern struct {
	NonRangePattern
	IsRef bool
	IsMut bool
	ID    string
}

func (ip *IdentifierPattern) Accept(v RusterBaseVisitor) {
	v.VisitIdentifierPattern(ip)
}

type PathPattern struct {
	NonRangePattern
	Segments []string
}

func (pp *PathPattern) Accept(v RusterBaseVisitor) {
	v.VisitPathPattern(pp)
}

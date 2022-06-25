package ast

type CallParams []Expression
type ArrayElements []Expression
type TupleElements []Expression

type LiteralExpression struct {
	Tp  Literal `yaml:"Type"`
	Val string  `yaml:"Value"`
}

func (le LiteralExpression) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitLiteralExpression(&le)
}

type PathExpression struct {
	Segments PathSegments
}

func (pe PathExpression) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitPathExpression(&pe)
}

type IfExpression struct {
	Expr    Expression      `yaml:"If"`
	IfTrue  BlockExpression `yaml:"Do"`
	IfFalse interface{}     `yaml:"Else,omitempty"`
}

func (ie IfExpression) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitIfExpression(&ie)
}

type MatchExpression struct {
	Expr  Expression
	Cases []MatchArm
}

func (me MatchExpression) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitMatchExpression(&me)
}

type MatchArm struct {
	Patterns []Pattern
	Body     BlockExpression
}

func (ma MatchArm) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitMatchArm(&ma)
}

type InfiniteLoopExpression struct {
	Body BlockExpression `yaml:"Do"`
}

func (ile InfiniteLoopExpression) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitInfiniteLoopExpression(&ile)
}

type PredicateLoopExpression struct {
	Expr Expression      `yaml:"While"`
	Body BlockExpression `yaml:"Do"`
}

func (ple PredicateLoopExpression) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitPredicateLoopExpression(&ple)
}

type IteratorLoopExpression struct {
	Ptrn Pattern         `yaml:"For"`
	Expr Expression      `yaml:"In"`
	Body BlockExpression `yaml:"Do"`
}

func (ile IteratorLoopExpression) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitIteratorLoopExpression(&ile)
}

type UnaryOperator struct {
	Op  string
	Val Expression
}

func (uo UnaryOperator) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitUnaryOperator(&uo)
}

type BinaryOperator struct {
	Op  string     `yaml:"Operand"`
	LHS Expression `yaml:"lhs"`
	RHS Expression `yaml:"rhs"`
}

func (bo BinaryOperator) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitBinaryOperator(&bo)
}

type RHSRangeOperator struct {
	Op  string
	Val Expression
}

func (rro RHSRangeOperator) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitRHSRangeOperator(&rro)
}

type RangeOperator struct {
	Op  string     `yaml:"Operand"`
	LHS Expression `yaml:"lhs,flow"`
	RHS Expression `yaml:"rhs,flow"`
}

func (ro RangeOperator) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitRangeOperator(&ro)
}

type ReturnExpression struct {
	Expr Expression `yaml:"return,flow,omitempty"`
}

func (re ReturnExpression) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitReturnExpression(&re)
}

type ContinueExpression struct {
	Expr Expression `yaml:"continue,flow,omitempty"`
}

func (ce ContinueExpression) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitContinueExpression(&ce)
}

type BreakExpression struct {
	Expr Expression `yaml:"break,flow,omitempty"`
}

func (be BreakExpression) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitBreakExpression(&be)
}

type TypeCastExpression struct {
	Expr Expression
	Tp   Type
}

func (tce TypeCastExpression) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitTypeCastExpression(&tce)
}

type CallExpression struct {
	FnHeader Expression `yaml:"Function,flow"`
	Params   CallParams
}

func (ce CallExpression) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitCallExpression(&ce)
}

type MethodCallExpression struct {
	FnHeader Expression `yaml:"Function,flow"`
	Params   CallParams
	Method   string
}

func (mce MethodCallExpression) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitMethodCallExpression(&mce)
}

type BorrowExpression struct {
	IsMut       bool       `yaml:"IsMut,omitempty"`
	IsDoubleRef bool       `yaml:"IsDoubleRef,omitempty"`
	Expr        Expression `yaml:"Expression,flow"`
}

func (be BorrowExpression) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitBorrowExpression(&be)
}

type ArrayIndexExpression struct {
	Object Expression `yaml:"Array"`
	Index  Expression `yaml:"Index,flow"`
}

func (aie ArrayIndexExpression) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitArrayIndexExpression(&aie)
}

type TupleIndexExpression struct {
	Object Expression `yaml:"Tuple"`
	Index  string     `yaml:"Index,flow"`
}

func (tie TupleIndexExpression) Accept(v RusterBaseVisitor) interface{} {
	return v.VisitTupleIndexExpression(&tie)
}

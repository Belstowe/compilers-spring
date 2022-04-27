package ast

type CallParams []Terminal

type LiteralExpression struct {
	Tp  Literal `yaml:"Type"`
	Val string  `yaml:"Value"`
}

func (le *LiteralExpression) Accept(v RusterBaseVisitor) {
	v.VisitLiteralExpression(le)
}

type PathExpression struct {
	Segments PathSegments
}

func (pe *PathExpression) Accept(v RusterBaseVisitor) {
	v.VisitPathExpression(pe)
}

type IfExpression struct {
	Expr    Expression      `yaml:"If"`
	IfTrue  BlockExpression `yaml:"Do"`
	IfFalse Node            `yaml:"Else,omitempty"`
}

func (ie *IfExpression) Accept(v RusterBaseVisitor) {
	v.VisitIfExpression(ie)
}

type MatchExpression struct {
	Expr  Expression
	Cases []MatchArm
}

func (me *MatchExpression) Accept(v RusterBaseVisitor) {
	v.VisitMatchExpression(me)
}

type MatchArm struct {
	Patterns []Pattern
	Body     BlockExpression
}

func (ma *MatchArm) Accept(v RusterBaseVisitor) {
	v.VisitMatchArm(ma)
}

type InfiniteLoopExpression struct {
	Body BlockExpression `yaml:"Do"`
}

func (ile *InfiniteLoopExpression) Accept(v RusterBaseVisitor) {
	v.VisitInfiniteLoopExpression(ile)
}

type PredicateLoopExpression struct {
	Expr Expression      `yaml:"While"`
	Body BlockExpression `yaml:"Do"`
}

func (ple *PredicateLoopExpression) Accept(v RusterBaseVisitor) {
	v.VisitPredicateLoopExpression(ple)
}

type IteratorLoopExpression struct {
	Ptrn Pattern         `yaml:"For"`
	Expr Expression      `yaml:"In"`
	Body BlockExpression `yaml:"Do"`
}

func (ile *IteratorLoopExpression) Accept(v RusterBaseVisitor) {
	v.VisitIteratorLoopExpression(ile)
}

type UnaryOperator struct {
	Op  string
	Val Expression
}

func (uo *UnaryOperator) Accept(v RusterBaseVisitor) {
	v.VisitUnaryOperator(uo)
}

type BinaryOperator struct {
	Op  string     `yaml:"Operand"`
	LHS Expression `yaml:",flow"`
	RHS Expression `yaml:",flow"`
}

func (bo *BinaryOperator) Accept(v RusterBaseVisitor) {
	v.VisitBinaryOperator(bo)
}

type ReturnExpression struct {
	Expr Expression `yaml:"return,flow,omitempty"`
}

func (re *ReturnExpression) Accept(v RusterBaseVisitor) {
	v.VisitReturnExpression(re)
}

type ContinueExpression struct {
	Expr Expression `yaml:"continue,flow,omitempty"`
}

func (ce *ContinueExpression) Accept(v RusterBaseVisitor) {
	v.VisitContinueExpression(ce)
}

type BreakExpression struct {
	Expr Expression `yaml:"break,flow,omitempty"`
}

func (be *BreakExpression) Accept(v RusterBaseVisitor) {
	v.VisitBreakExpression(be)
}

type TypeCastExpression struct {
	Expr Expression
	Tp   Type
}

func (tce *TypeCastExpression) Accept(v RusterBaseVisitor) {
	v.VisitTypeCastExpression(tce)
}

type CallExpression struct {
	FnHeader Expression `yaml:"Function,flow"`
	Params   CallParams
}

func (ce *CallExpression) Accept(v RusterBaseVisitor) {
	v.VisitCallExpression(ce)
}

type MethodCallExpression struct {
	FnHeader Expression `yaml:"Function,flow"`
	Params   CallParams
	Method   string
}

func (mce *MethodCallExpression) Accept(v RusterBaseVisitor) {
	v.VisitMethodCallExpression(mce)
}

type BorrowExpression struct {
	IsMut       bool       `yaml:"IsMut,omitempty"`
	IsDoubleRef bool       `yaml:"IsDoubleRef,omitempty"`
	Expr        Expression `yaml:"Expression,flow"`
}

func (be *BorrowExpression) Accept(v RusterBaseVisitor) {
	v.VisitBorrowExpression(be)
}

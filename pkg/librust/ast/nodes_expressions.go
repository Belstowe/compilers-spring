package ast

type CallParams []Expression

type LiteralExpression struct {
	Expression
	Tp  Literal
	Val string
}

func (le *LiteralExpression) Accept(v RusterBaseVisitor) {
	v.VisitLiteralExpression(le)
}

type PathExpression struct {
	Expression
	Segments PathSegments
}

func (pe *PathExpression) Accept(v RusterBaseVisitor) {
	v.VisitPathExpression(pe)
}

type IfExpression struct {
	ExpressionWithBlock
	Expr    Expression
	IfTrue  BlockExpression
	IfFalse Node
}

func (ie *IfExpression) Accept(v RusterBaseVisitor) {
	v.VisitIfExpression(ie)
}

type MatchExpression struct {
	ExpressionWithBlock
	Expr  Expression
	Cases []MatchArm
}

func (me *MatchExpression) Accept(v RusterBaseVisitor) {
	v.VisitMatchExpression(me)
}

type MatchArm struct {
	Node
	Patterns []Pattern
	Body     BlockExpression
}

func (ma *MatchArm) Accept(v RusterBaseVisitor) {
	v.VisitMatchArm(ma)
}

type InfiniteLoopExpression struct {
	ExpressionWithBlock
	Body BlockExpression
}

func (ile *InfiniteLoopExpression) Accept(v RusterBaseVisitor) {
	v.VisitInfiniteLoopExpression(ile)
}

type PredicateLoopExpression struct {
	ExpressionWithBlock
	Expr Expression
	Body BlockExpression
}

func (ple *PredicateLoopExpression) Accept(v RusterBaseVisitor) {
	v.VisitPredicateLoopExpression(ple)
}

type IteratorLoopExpression struct {
	ExpressionWithBlock
	Ptrn Pattern
	Expr Expression
	Body BlockExpression
}

func (ile *IteratorLoopExpression) Accept(v RusterBaseVisitor) {
	v.VisitIteratorLoopExpression(ile)
}

type UnaryOperator struct {
	Expression
	Op  string
	Val Expression
}

func (uo *UnaryOperator) Accept(v RusterBaseVisitor) {
	v.VisitUnaryOperator(uo)
}

type BinaryOperator struct {
	Expression
	Op  string
	LHS Expression
	RHS Expression
}

func (bo *BinaryOperator) Accept(v RusterBaseVisitor) {
	v.VisitBinaryOperator(bo)
}

type ReturnExpression struct {
	Expression
	Expr Expression
}

func (re *ReturnExpression) Accept(v RusterBaseVisitor) {
	v.VisitReturnExpression(re)
}

type ContinueExpression struct {
	Expression
	Expr Expression
}

func (ce *ContinueExpression) Accept(v RusterBaseVisitor) {
	v.VisitContinueExpression(ce)
}

type BreakExpression struct {
	Expression
	Expr Expression
}

func (be *BreakExpression) Accept(v RusterBaseVisitor) {
	v.VisitBreakExpression(be)
}

type TypeCastExpression struct {
	Expression
	Expr Expression
	Tp   Type
}

func (tce *TypeCastExpression) Accept(v RusterBaseVisitor) {
	v.VisitTypeCastExpression(tce)
}

type CallExpression struct {
	Expression
	FnHeader Expression
	Params   CallParams
}

func (ce *CallExpression) Accept(v RusterBaseVisitor) {
	v.VisitCallExpression(ce)
}

type MethodCallExpression struct {
	Expression
	*CallExpression
	Method string
}

func (mce *MethodCallExpression) Accept(v RusterBaseVisitor) {
	v.VisitMethodCallExpression(mce)
}

type BorrowExpression struct {
	Expression
	IsMut       bool
	IsDoubleRef bool
	Expr        Expression
}

func (be *BorrowExpression) Accept(v RusterBaseVisitor) {
	v.VisitBorrowExpression(be)
}

package ast_nodes

type CallParams []Expression

type LiteralExpression struct {
	Expression
	Tp  Literal
	Val string
}

type PathExpression struct {
	Expression
	Segments PathSegments
}

type IfExpression struct {
	ExpressionWithBlock
	Expr    Expression
	IfTrue  BlockExpression
	IfFalse Node
}

type MatchExpression struct {
	ExpressionWithBlock
	Expr  Expression
	Cases []MatchArm
}

type MatchArm struct {
	Node
	Patterns []Pattern
	Body     BlockExpression
}

type InfiniteLoopExpression struct {
	ExpressionWithBlock
	Body BlockExpression
}

type PredicateLoopExpression struct {
	ExpressionWithBlock
	Expr Expression
	Body BlockExpression
}

type IteratorLoopExpression struct {
	ExpressionWithBlock
	Ptrn Pattern
	Expr Expression
	Body BlockExpression
}

type UnaryOperator struct {
	Expression
	Op  string
	Val Expression
}

type BinaryOperator struct {
	Expression
	Op  string
	LHS Expression
	RHS Expression
}

type ReturnExpression struct {
	Expression
	Expr Expression
}

type ContinueExpression struct {
	Expression
	Expr Expression
}

type BreakExpression struct {
	Expression
	Expr Expression
}

type TypeCastExpression struct {
	Expression
	Expr Expression
	Tp   Type
}

type CallExpression struct {
	Expression
	FnHeader Expression
	Params   CallParams
}

type MethodCallExpression struct {
	Expression
	*CallExpression
	Method string
}

type BorrowExpression struct {
	Expression
	IsMut       bool
	IsDoubleRef bool
	Expr        Expression
}

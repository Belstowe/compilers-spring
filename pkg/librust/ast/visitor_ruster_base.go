package ast

type RusterBaseVisitor interface {
	DumpLogs() interface{}

	Visit(tree Node) interface{}

	VisitCrate(c *Crate) interface{}
	VisitBlockExpression(be *BlockExpression) interface{}
	VisitUseDecl(ud *UseDecl) interface{}
	VisitSimplePath(sp *SimplePath) interface{}
	VisitParameter(p *Parameter) interface{}
	VisitFunction(f *Function) interface{}
	VisitLetStatement(ls *LetStatement) interface{}

	VisitLiteralExpression(le *LiteralExpression) interface{}
	VisitPathExpression(pe *PathExpression) interface{}
	VisitIfExpression(ie *IfExpression) interface{}
	VisitMatchExpression(me *MatchExpression) interface{}
	VisitMatchArm(ma *MatchArm) interface{}
	VisitInfiniteLoopExpression(ile *InfiniteLoopExpression) interface{}
	VisitPredicateLoopExpression(ple *PredicateLoopExpression) interface{}
	VisitIteratorLoopExpression(ile *IteratorLoopExpression) interface{}
	VisitUnaryOperator(uo *UnaryOperator) interface{}
	VisitBinaryOperator(bo *BinaryOperator) interface{}
	VisitRHSRangeOperator(rro *RHSRangeOperator) interface{}
	VisitRangeOperator(ro *RangeOperator) interface{}
	VisitReturnExpression(e *ReturnExpression) interface{}
	VisitContinueExpression(ce *ContinueExpression) interface{}
	VisitBreakExpression(be *BreakExpression) interface{}
	VisitTypeCastExpression(tce *TypeCastExpression) interface{}
	VisitCallExpression(ce *CallExpression) interface{}
	VisitMethodCallExpression(mce *MethodCallExpression) interface{}
	VisitBorrowExpression(be *BorrowExpression) interface{}
	VisitArrayIndexExpression(aie *ArrayIndexExpression) interface{}
	VisitTupleIndexExpression(tie *TupleIndexExpression) interface{}

	VisitLiteralPattern(lp *LiteralPattern) interface{}
	VisitReferencePattern(rp *ReferencePattern) interface{}
	VisitIdentifierPattern(ip *IdentifierPattern) interface{}
	VisitPathPattern(pp *PathPattern) interface{}

	VisitTypePath(tp *TypePath) interface{}
	VisitTypePathSegment(tps *TypePathSegment) interface{}
	VisitTypePathFunction(tpf *TypePathFunction) interface{}
	VisitParenthesizedType(pt *ParenthesizedType) interface{}
	VisitPointerType(pt *PointerType) interface{}
	VisitReferenceType(rt *ReferenceType) interface{}
	VisitTupleType(tt *TupleType) interface{}
	VisitArrayType(at *ArrayType) interface{}
	VisitSliceType(st *SliceType) interface{}
	VisitNeverType(nt *NeverType) interface{}
	VisitInferredType(it *InferredType) interface{}
}

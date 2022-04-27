package ast

type RusterXMLVisitor struct {
	RusterBaseVisitor
}

func (v RusterXMLVisitor) Visit(tree Node) interface{} {
	switch val := tree.(type) {
	case *Crate:
		return v.VisitCrate(val)
	case *BlockExpression:
		return v.VisitBlockExpression(val)
	case *UseDecl:
		return v.VisitUseDecl(val)
	case *SimplePath:
		return v.VisitSimplePath(val)
	case *Parameter:
		return v.VisitParameter(val)
	case *Function:
		return v.VisitFunction(val)
	case *LetStatement:
		return v.VisitLetStatement(val)
	case *LiteralExpression:
		return v.VisitLiteralExpression(val)
	case *PathExpression:
		return v.VisitPathExpression(val)
	case *IfExpression:
		return v.VisitIfExpression(val)
	case *MatchExpression:
		return v.VisitMatchExpression(val)
	case *MatchArm:
		return v.VisitMatchArm(val)
	case *InfiniteLoopExpression:
		return v.VisitInfiniteLoopExpression(val)
	case *PredicateLoopExpression:
		return v.VisitPredicateLoopExpression(val)
	case *IteratorLoopExpression:
		return v.VisitIteratorLoopExpression(val)
	case *UnaryOperator:
		return v.VisitUnaryOperator(val)
	case *BinaryOperator:
		return v.VisitBinaryOperator(val)
	case *ReturnExpression:
		return v.VisitReturnExpression(val)
	case *ContinueExpression:
		return v.VisitContinueExpression(val)
	case *BreakExpression:
		return v.VisitBreakExpression(val)
	case *TypeCastExpression:
		return v.VisitTypeCastExpression(val)
	case *CallExpression:
		return v.VisitCallExpression(val)
	case *MethodCallExpression:
		return v.VisitMethodCallExpression(val)
	case *BorrowExpression:
		return v.VisitBorrowExpression(val)
	case *LiteralPattern:
		return v.VisitLiteralPattern(val)
	case *ReferencePattern:
		return v.VisitReferencePattern(val)
	case *IdentifierPattern:
		return v.VisitIdentifierPattern(val)
	case *PathPattern:
		return v.VisitPathPattern(val)
	case *TypePath:
		return v.VisitTypePath(val)
	case *TypePathSegment:
		return v.VisitTypePathSegment(val)
	case *TypePathFunction:
		return v.VisitTypePathFunction(val)
	case *ParenthesizedType:
		return v.VisitParenthesizedType(val)
	case *PointerType:
		return v.VisitPointerType(val)
	case *ReferenceType:
		return v.VisitReferenceType(val)
	case *TupleType:
		return v.VisitTupleType(val)
	case *ArrayType:
		return v.VisitArrayType(val)
	case *SliceType:
		return v.VisitSliceType(val)
	case *NeverType:
		return v.VisitNeverType(val)
	case *InferredType:
		return v.VisitInferredType(val)
	default:
		panic("unknown context")
	}
}

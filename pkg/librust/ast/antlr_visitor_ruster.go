package ast

import (
	"github.com/Compiler2022/compilers-1-Belstowe/parser"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type TokenVocabulary []string

var vocabulary TokenVocabulary = parser.NewRustLexer(antlr.NewInputStream("")).GetSymbolicNames()

type ANTLRRusterVisitor struct {
	*parser.BaseRustParserVisitor
	crate Crate
}

func NewANTLRRusterVisitor() ANTLRRusterVisitor {
	var v ANTLRRusterVisitor
	return v
}

func (v *ANTLRRusterVisitor) VisitCrate(ctx *parser.CrateContext) Crate {
	for _, element := range ctx.AllItem() {
		item := v.Visit(element)
		v.crate.Items = append(v.crate.Items, item.(Item))
	}
	return v.crate
}

func (v *ANTLRRusterVisitor) VisitUseTree(ctx *parser.UseTreeContext) interface{} {
	var useDecl UseDecl
	if ctx.RuleIndex == 0 {
		useDecl.All = true
	} else {
		useDecl.All = false
	}

	useDecl.Path = v.Visit(ctx.SimplePath()).(SimplePath)

	return useDecl
}

func (v *ANTLRRusterVisitor) VisitFunction(ctx *parser.FunctionContext) interface{} {
	var fn Function
	fn.ID = v.Visit(ctx.Identifier()).(string)

	if ctx.FunctionReturnType().IsEmpty() {
		fn.ReturnType = nil
	} else {
		fn.ReturnType = v.Visit(ctx.FunctionReturnType()).(Type)
	}

	if ctx.FunctionParameters().IsEmpty() {
		fn.Params = make([]Parameter, 0)
	} else {
		fn.Params = v.Visit(ctx.FunctionParameters()).([]Parameter)
	}

	fn.Body = v.Visit(ctx.BlockExpression()).(BlockExpression)

	return fn
}

func (v *ANTLRRusterVisitor) VisitFunctionParameters(ctx *parser.FunctionParametersContext) interface{} {
	params := make([]Parameter, 0)
	for _, e := range ctx.AllFunctionParam() {
		params = append(params, v.Visit(e).(Parameter))
	}
	return params
}

func (v *ANTLRRusterVisitor) VisitFunctionParam(ctx *parser.FunctionParamContext) interface{} {
	var param Parameter
	if ctx.Identifier().IsEmpty() {
		param.ID = nil
	} else {
		param.ID = v.Visit(ctx.Identifier()).(string)
	}

	param.VarType = v.Visit(ctx.Type()).(Type)

	return param
}

func (v *ANTLRRusterVisitor) VisitFunctionReturnType(ctx *parser.FunctionReturnTypeContext) interface{} {
	return v.Visit(ctx.Type()).(Type)
}

func (v *ANTLRRusterVisitor) VisitStruct(ctx *parser.StructContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitStructFields(ctx *parser.StructFieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitStructField(ctx *parser.StructFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitTypeAlias(ctx *parser.TypeAliasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitConstantItem(ctx *parser.ConstantItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitLetStatement(ctx *parser.LetStatementContext) interface{} {
	var statement LetStatement

	statement.Ptrn = v.Visit(ctx.Pattern()).(Pattern)

	if ctx.Type().IsEmpty() {
		statement.VarType = nil
	} else {
		statement.VarType = v.Visit(ctx.Type()).(Type)
	}

	if ctx.Expression().IsEmpty() {
		statement.Expr = nil
	} else {
		statement.Expr = v.Visit(ctx.Expression()).(Expression)
	}

	return statement
}

func (v *ANTLRRusterVisitor) VisitTypeCastExpression(ctx *parser.TypeCastExpressionContext) interface{} {
	var expr TypeCastExpression

	expr.Tp = v.Visit(ctx.Type()).(Type)
	expr.Expr = v.Visit(ctx.Expression()).(Expression)

	return expr
}

func (v *ANTLRRusterVisitor) VisitTupleExpression(ctx *parser.TupleExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitIndexExpression(ctx *parser.IndexExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitRangeExpression(ctx *parser.RangeExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitReturnExpression(ctx *parser.ReturnExpressionContext) interface{} {
	var expr ReturnExpression
	if ctx.Expression().IsEmpty() {
		expr.Expression = nil
	} else {
		expr.Expression = v.Visit(ctx.Expression()).(Expression)
	}
	return expr
}

func (v *ANTLRRusterVisitor) VisitErrorPropagationExpression(ctx *parser.ErrorPropagationExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitContinueExpression(ctx *parser.ContinueExpressionContext) interface{} {
	var expr ContinueExpression
	if ctx.Expression().IsEmpty() {
		expr.Expression = nil
	} else {
		expr.Expression = v.Visit(ctx.Expression()).(Expression)
	}
	return expr
}

func (v *ANTLRRusterVisitor) VisitAssignmentExpression(ctx *parser.AssignmentExpressionContext) interface{} {
	var expr BinaryOperator
	expr.LHS = v.Visit(ctx.LHS).(Expression)
	expr.Op = ctx.EQ().GetText()
	expr.RHS = v.Visit(ctx.RHS).(Expression)
	return expr
}

func (v *ANTLRRusterVisitor) VisitMethodCallExpression(ctx *parser.MethodCallExpressionContext) interface{} {
	var expr MethodCallExpression
	expr.Params = v.Visit(ctx.CallParams()).([]Expression)
	expr.FnHeader = v.Visit(ctx.Expression()).(Expression)
	expr.Method = v.Visit(ctx.SimplePathSegment()).(string)
	return expr
}

func (v *ANTLRRusterVisitor) VisitLiteralExpression_(ctx *parser.LiteralExpression_Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitStructExpression_(ctx *parser.StructExpression_Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitTupleIndexingExpression(ctx *parser.TupleIndexingExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitCallExpression(ctx *parser.CallExpressionContext) interface{} {
	var expr CallExpression
	expr.Params = v.Visit(ctx.CallParams()).([]Expression)
	expr.FnHeader = v.Visit(ctx.Expression()).(Expression)
	return expr
}

func (v *ANTLRRusterVisitor) VisitDereferenceExpression(ctx *parser.DereferenceExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitNestedExpression(ctx *parser.NestedExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitUnaryOpExpression(ctx *parser.UnaryOpExpressionContext) interface{} {
	var expr UnaryOperator
	expr.Op = ctx.Op.GetText()
	expr.Val = v.Visit(ctx.Val).(Expression)
	return expr
}

func (v *ANTLRRusterVisitor) VisitBinaryOpExpression(ctx *parser.BinaryOpExpressionContext) interface{} {
	var expr BinaryOperator
	expr.LHS = v.Visit(ctx.LHS).(Expression)
	expr.Op = ctx.Op.GetText()
	expr.RHS = v.Visit(ctx.RHS).(Expression)
	return expr
}

func (v *ANTLRRusterVisitor) VisitBreakExpression(ctx *parser.BreakExpressionContext) interface{} {
	var expr BreakExpression
	if ctx.Expression().IsEmpty() {
		expr.Expression = nil
	} else {
		expr.Expression = v.Visit(ctx.Expression()).(Expression)
	}
	return expr
}

func (v *ANTLRRusterVisitor) VisitFieldExpression(ctx *parser.FieldExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitBorrowExpression(ctx *parser.BorrowExpressionContext) interface{} {
	var expr BorrowExpression
	expr.IsDoubleRef = ctx.RefToken.GetText() == "&&"
	expr.IsMut = ctx.MutToken.GetText() == "mut"
	expr.Expression = v.Visit(ctx.Expression()).(Expression)
	return expr
}

func (v *ANTLRRusterVisitor) VisitCompoundAssignmentExpression(ctx *parser.CompoundAssignmentExpressionContext) interface{} {
	var expr BinaryOperator
	expr.LHS = v.Visit(ctx.LHS).(Expression)
	expr.Op = v.Visit(ctx.CompoundAssignOperator()).(string)
	expr.RHS = v.Visit(ctx.RHS).(Expression)
	return expr
}

func (v *ANTLRRusterVisitor) VisitArrayExpression(ctx *parser.ArrayExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitCompoundAssignOperator(ctx *parser.CompoundAssignOperatorContext) interface{} {
	return ctx.GetText()
}

func (v *ANTLRRusterVisitor) VisitLiteralExpression(ctx *parser.LiteralExpressionContext) interface{} {
	var literal LiteralExpression
	literal.Val = ctx.GetText()
	switch vocabulary[ctx.Literal.GetTokenType()] {
	case "STRING_LITERAL":
		literal.Tp = String
	case "CHAR_LITERAL":
		literal.Tp = Char
	case "INTEGER_LITERAL":
		literal.Tp = Integer
	case "KW_TRUE":
		literal.Tp = Boolean
	case "KW_FALSE":
		literal.Tp = Boolean
	}
	return literal
}

func (v *ANTLRRusterVisitor) VisitStatements(ctx *parser.StatementsContext) interface{} {
	var block BlockExpression

	block.Statements = make([]Statement, 0)
	for _, e := range ctx.AllStatement() {
		block.Statements = append(block.Statements, v.Visit(e).(Statement))
	}

	if ctx.Expression().IsEmpty() {
		block.Expr = nil
	} else {
		block.Expr = v.Visit(ctx.Expression()).(Expression)
	}

	return block
}

func (v *ANTLRRusterVisitor) VisitArrayElements(ctx *parser.ArrayElementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitTupleElements(ctx *parser.TupleElementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitTupleIndex(ctx *parser.TupleIndexContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitStructExpression(ctx *parser.StructExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitStructExprFields(ctx *parser.StructExprFieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitStructExprField(ctx *parser.StructExprFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitCallParams(ctx *parser.CallParamsContext) interface{} {
	segments := make([]Expression, 0)
	for _, e := range ctx.AllExpression() {
		segments = append(segments, v.Visit(e).(Expression))
	}
	return segments
}

func (v *ANTLRRusterVisitor) VisitPathExpression(ctx *parser.PathExpressionContext) interface{} {
	segments := make([]string, 0)
	for _, e := range ctx.AllSimplePathSegment() {
		segments = append(segments, v.Visit(e).(string))
	}
	return segments
}

func (v *ANTLRRusterVisitor) VisitIfExpression(ctx *parser.IfExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitMatchExpression(ctx *parser.MatchExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitMatchArms(ctx *parser.MatchArmsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitMatchArmExpression(ctx *parser.MatchArmExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitMatchArm(ctx *parser.MatchArmContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitInfiniteLoopExpression(ctx *parser.InfiniteLoopExpressionContext) interface{} {
	var expr InfiniteLoopExpression
	expr.Body = v.Visit(ctx.BlockExpression()).(BlockExpression)
	return expr
}

func (v *ANTLRRusterVisitor) VisitPredicateLoopExpression(ctx *parser.PredicateLoopExpressionContext) interface{} {
	var expr PredicateLoopExpression
	expr.Expr = v.Visit(ctx.Expression()).(Expression)
	expr.Body = v.Visit(ctx.BlockExpression()).(BlockExpression)
	return expr
}

func (v *ANTLRRusterVisitor) VisitIteratorLoopExpression(ctx *parser.IteratorLoopExpressionContext) interface{} {
	var expr IteratorLoopExpression
	expr.Ptrn = v.Visit(ctx.Pattern()).(Pattern)
	expr.Expr = v.Visit(ctx.Expression()).(Expression)
	expr.Body = v.Visit(ctx.BlockExpression()).(BlockExpression)
	return expr
}

func (v *ANTLRRusterVisitor) VisitLiteralPattern(ctx *parser.LiteralPatternContext) interface{} {
	var literal LiteralPattern
	literal.Val = ctx.GetText()
	switch vocabulary[ctx.Literal.GetTokenType()] {
	case "STRING_LITERAL":
		literal.Tp = String
	case "CHAR_LITERAL":
		literal.Tp = Char
	case "INTEGER_LITERAL":
		literal.Tp = Integer
	case "KW_TRUE":
		literal.Tp = Boolean
	case "KW_FALSE":
		literal.Tp = Boolean
	}
	return literal
}

func (v *ANTLRRusterVisitor) VisitIdentifierPattern(ctx *parser.IdentifierPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitWildcardPattern(ctx *parser.WildcardPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitRestPattern(ctx *parser.RestPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitReferencePattern(ctx *parser.ReferencePatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitStructPattern(ctx *parser.StructPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitStructPatternElements(ctx *parser.StructPatternElementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitStructPatternFields(ctx *parser.StructPatternFieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitStructPatternField(ctx *parser.StructPatternFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitTuplePattern(ctx *parser.TuplePatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitTuplePatternItems(ctx *parser.TuplePatternItemsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitGroupedPattern(ctx *parser.GroupedPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitSlicePattern(ctx *parser.SlicePatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitSlicePatternItems(ctx *parser.SlicePatternItemsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitPathPattern(ctx *parser.PathPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitRangePattern(ctx *parser.RangePatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitRangePatternBound(ctx *parser.RangePatternBoundContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitNestedType(ctx *parser.NestedTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitNeverType(ctx *parser.NeverTypeContext) interface{} {
	var tp NeverType
	return tp
}

func (v *ANTLRRusterVisitor) VisitInferredType(ctx *parser.InferredTypeContext) interface{} {
	var tp InferredType
	return tp
}

func (v *ANTLRRusterVisitor) VisitTupleType(ctx *parser.TupleTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitArrayType(ctx *parser.ArrayTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitSliceType(ctx *parser.SliceTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitReferenceType(ctx *parser.ReferenceTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitPointerType(ctx *parser.PointerTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitFunctionType(ctx *parser.FunctionTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRRusterVisitor) VisitTypePath(ctx *parser.TypePathContext) interface{} {
	var path TypePath
	path.Segments = make([]TypePathSegment, 0)
	for _, e := range ctx.AllTypePathSegment() {
		path.Segments = append(path.Segments, v.Visit(e).(TypePathSegment))
	}
	return path
}

func (v *ANTLRRusterVisitor) VisitTypePathSegment(ctx *parser.TypePathSegmentContext) interface{} {
	var segment TypePathSegment
	segment.ID = v.Visit(ctx.SimplePathSegment()).(string)
	if ctx.TypePathFn().IsEmpty() {
		segment.Fn = nil
	} else {
		segment.Fn = v.Visit(ctx.TypePathFn()).(*TypePathFunction)
	}
	return segment
}

func (v *ANTLRRusterVisitor) VisitTypePathFn(ctx *parser.TypePathFnContext) interface{} {
	var fn TypePathFunction
	fn.Inputs = v.Visit(ctx.TypePathInputs()).([]Type)
	if ctx.ReturnType.IsEmpty() {
		fn.ReturnType = nil
	} else {
		fn.ReturnType = v.Visit(ctx.ReturnType).(Type)
	}
	return fn
}

func (v *ANTLRRusterVisitor) VisitTypePathInputs(ctx *parser.TypePathInputsContext) interface{} {
	inputs := make([]Type, 0)
	for _, e := range ctx.AllType() {
		inputs = append(inputs, v.Visit(e).(Type))
	}
	return inputs
}

func (v *ANTLRRusterVisitor) VisitSimplePath(ctx *parser.SimplePathContext) interface{} {
	path := make([]string, 0)
	for _, e := range ctx.AllSimplePathSegment() {
		path = append(path, v.Visit(e).(string))
	}
	return path
}

func (v *ANTLRRusterVisitor) VisitSimplePathSegment(ctx *parser.SimplePathSegmentContext) interface{} {
	return ctx.GetText()
}

func (v *ANTLRRusterVisitor) VisitIdentifier(ctx *parser.IdentifierContext) interface{} {
	return ctx.GetText()
}

func (v *ANTLRRusterVisitor) VisitKeyword(ctx *parser.KeywordContext) interface{} {
	return ctx.GetText()
}

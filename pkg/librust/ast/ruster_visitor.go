package ast

import (
	"github.com/Compiler2022/compilers-1-Belstowe/parser"
	nodes "github.com/Compiler2022/compilers-1-Belstowe/pkg/librust/ast/nodes"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type TokenVocabulary []string

var vocabulary TokenVocabulary = parser.NewRustLexer(antlr.NewInputStream("")).GetSymbolicNames()

type RusterVisitor struct {
	*parser.BaseRustParserVisitor
	crate nodes.Crate
}

func (v *RusterVisitor) VisitCrate(ctx *parser.CrateContext) []nodes.Item {
	for _, element := range ctx.AllItem() {
		item := v.Visit(element)
		v.crate.Items = append(v.crate.Items, item.(nodes.Item))
	}
	return v.crate.Items
}

func (v *RusterVisitor) VisitUseTree(ctx *parser.UseTreeContext) interface{} {
	var useDecl nodes.UseDecl
	if ctx.RuleIndex == 0 {
		useDecl.All = true
	} else {
		useDecl.All = false
	}

	useDecl.Path = v.Visit(ctx.SimplePath()).(nodes.SimplePath)

	return useDecl
}

func (v *RusterVisitor) VisitFunction(ctx *parser.FunctionContext) interface{} {
	var fn nodes.Function
	fn.ID = v.Visit(ctx.Identifier()).(string)

	if ctx.FunctionReturnType().IsEmpty() {
		fn.ReturnType = nil
	} else {
		fn.ReturnType = v.Visit(ctx.FunctionReturnType()).(nodes.Type)
	}

	if ctx.FunctionParameters().IsEmpty() {
		fn.Params = make([]nodes.Parameter, 0)
	} else {
		fn.Params = v.Visit(ctx.FunctionParameters()).([]nodes.Parameter)
	}

	fn.Body = v.Visit(ctx.BlockExpression()).(nodes.BlockExpression)

	return fn
}

func (v *RusterVisitor) VisitFunctionParameters(ctx *parser.FunctionParametersContext) interface{} {
	params := make([]nodes.Parameter, 0)
	for _, e := range ctx.AllFunctionParam() {
		params = append(params, v.Visit(e).(nodes.Parameter))
	}
	return params
}

func (v *RusterVisitor) VisitFunctionParam(ctx *parser.FunctionParamContext) interface{} {
	var param nodes.Parameter
	if ctx.Identifier().IsEmpty() {
		param.ID = nil
	} else {
		param.ID = v.Visit(ctx.Identifier()).(string)
	}

	param.VarType = v.Visit(ctx.Type()).(nodes.Type)

	return param
}

func (v *RusterVisitor) VisitFunctionReturnType(ctx *parser.FunctionReturnTypeContext) interface{} {
	return v.Visit(ctx.Type()).(nodes.Type)
}

func (v *RusterVisitor) VisitStruct(ctx *parser.StructContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitStructFields(ctx *parser.StructFieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitStructField(ctx *parser.StructFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitTypeAlias(ctx *parser.TypeAliasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitConstantItem(ctx *parser.ConstantItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitLetStatement(ctx *parser.LetStatementContext) interface{} {
	var statement nodes.LetStatement

	statement.Ptrn = v.Visit(ctx.Pattern()).(nodes.Pattern)

	if ctx.Type().IsEmpty() {
		statement.VarType = nil
	} else {
		statement.VarType = v.Visit(ctx.Type()).(nodes.Type)
	}

	if ctx.Expression().IsEmpty() {
		statement.Expr = nil
	} else {
		statement.Expr = v.Visit(ctx.Expression()).(nodes.Expression)
	}

	return statement
}

func (v *RusterVisitor) VisitTypeCastExpression(ctx *parser.TypeCastExpressionContext) interface{} {
	var expr nodes.TypeCastExpression

	expr.Tp = v.Visit(ctx.Type()).(nodes.Type)
	expr.Expr = v.Visit(ctx.Expression()).(nodes.Expression)

	return expr
}

func (v *RusterVisitor) VisitTupleExpression(ctx *parser.TupleExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitIndexExpression(ctx *parser.IndexExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitRangeExpression(ctx *parser.RangeExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitReturnExpression(ctx *parser.ReturnExpressionContext) interface{} {
	var expr nodes.ReturnExpression
	if ctx.Expression().IsEmpty() {
		expr.Expression = nil
	} else {
		expr.Expression = v.Visit(ctx.Expression()).(nodes.Expression)
	}
	return expr
}

func (v *RusterVisitor) VisitErrorPropagationExpression(ctx *parser.ErrorPropagationExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitContinueExpression(ctx *parser.ContinueExpressionContext) interface{} {
	var expr nodes.ContinueExpression
	if ctx.Expression().IsEmpty() {
		expr.Expression = nil
	} else {
		expr.Expression = v.Visit(ctx.Expression()).(nodes.Expression)
	}
	return expr
}

func (v *RusterVisitor) VisitAssignmentExpression(ctx *parser.AssignmentExpressionContext) interface{} {
	var expr nodes.BinaryOperator
	expr.LHS = v.Visit(ctx.LHS).(nodes.Expression)
	expr.Op = ctx.EQ().GetText()
	expr.RHS = v.Visit(ctx.RHS).(nodes.Expression)
	return expr
}

func (v *RusterVisitor) VisitMethodCallExpression(ctx *parser.MethodCallExpressionContext) interface{} {
	var expr nodes.MethodCallExpression
	expr.Params = v.Visit(ctx.CallParams()).([]nodes.Expression)
	expr.FnHeader = v.Visit(ctx.Expression()).(nodes.Expression)
	expr.Method = v.Visit(ctx.SimplePathSegment()).(string)
	return expr
}

func (v *RusterVisitor) VisitLiteralExpression_(ctx *parser.LiteralExpression_Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitStructExpression_(ctx *parser.StructExpression_Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitTupleIndexingExpression(ctx *parser.TupleIndexingExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitCallExpression(ctx *parser.CallExpressionContext) interface{} {
	var expr nodes.CallExpression
	expr.Params = v.Visit(ctx.CallParams()).([]nodes.Expression)
	expr.FnHeader = v.Visit(ctx.Expression()).(nodes.Expression)
	return expr
}

func (v *RusterVisitor) VisitDereferenceExpression(ctx *parser.DereferenceExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitNestedExpression(ctx *parser.NestedExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitUnaryOpExpression(ctx *parser.UnaryOpExpressionContext) interface{} {
	var expr nodes.UnaryOperator
	expr.Op = ctx.Op.GetText()
	expr.Val = v.Visit(ctx.Val).(nodes.Expression)
	return expr
}

func (v *RusterVisitor) VisitBinaryOpExpression(ctx *parser.BinaryOpExpressionContext) interface{} {
	var expr nodes.BinaryOperator
	expr.LHS = v.Visit(ctx.LHS).(nodes.Expression)
	expr.Op = ctx.Op.GetText()
	expr.RHS = v.Visit(ctx.RHS).(nodes.Expression)
	return expr
}

func (v *RusterVisitor) VisitBreakExpression(ctx *parser.BreakExpressionContext) interface{} {
	var expr nodes.BreakExpression
	if ctx.Expression().IsEmpty() {
		expr.Expression = nil
	} else {
		expr.Expression = v.Visit(ctx.Expression()).(nodes.Expression)
	}
	return expr
}

func (v *RusterVisitor) VisitFieldExpression(ctx *parser.FieldExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitBorrowExpression(ctx *parser.BorrowExpressionContext) interface{} {
	var expr nodes.BorrowExpression
	expr.IsDoubleRef = ctx.RefToken.GetText() == "&&"
	expr.IsMut = ctx.MutToken.GetText() == "mut"
	expr.Expression = v.Visit(ctx.Expression()).(nodes.Expression)
	return expr
}

func (v *RusterVisitor) VisitCompoundAssignmentExpression(ctx *parser.CompoundAssignmentExpressionContext) interface{} {
	var expr nodes.BinaryOperator
	expr.LHS = v.Visit(ctx.LHS).(nodes.Expression)
	expr.Op = v.Visit(ctx.CompoundAssignOperator()).(string)
	expr.RHS = v.Visit(ctx.RHS).(nodes.Expression)
	return expr
}

func (v *RusterVisitor) VisitArrayExpression(ctx *parser.ArrayExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitCompoundAssignOperator(ctx *parser.CompoundAssignOperatorContext) interface{} {
	return ctx.GetText()
}

func (v *RusterVisitor) VisitLiteralExpression(ctx *parser.LiteralExpressionContext) interface{} {
	var literal nodes.LiteralExpression
	literal.Val = ctx.GetText()
	switch vocabulary[ctx.Literal.GetTokenType()] {
	case "STRING_LITERAL":
		literal.Tp = nodes.String
	case "CHAR_LITERAL":
		literal.Tp = nodes.Char
	case "INTEGER_LITERAL":
		literal.Tp = nodes.Integer
	case "KW_TRUE":
		literal.Tp = nodes.Boolean
	case "KW_FALSE":
		literal.Tp = nodes.Boolean
	}
	return literal
}

func (v *RusterVisitor) VisitStatements(ctx *parser.StatementsContext) interface{} {
	var block nodes.BlockExpression

	block.Statements = make([]nodes.Statement, 0)
	for _, e := range ctx.AllStatement() {
		block.Statements = append(block.Statements, v.Visit(e).(nodes.Statement))
	}

	if ctx.Expression().IsEmpty() {
		block.Expr = nil
	} else {
		block.Expr = v.Visit(ctx.Expression()).(nodes.Expression)
	}

	return block
}

func (v *RusterVisitor) VisitArrayElements(ctx *parser.ArrayElementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitTupleElements(ctx *parser.TupleElementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitTupleIndex(ctx *parser.TupleIndexContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitStructExpression(ctx *parser.StructExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitStructExprFields(ctx *parser.StructExprFieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitStructExprField(ctx *parser.StructExprFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitCallParams(ctx *parser.CallParamsContext) interface{} {
	segments := make([]nodes.Expression, 0)
	for _, e := range ctx.AllExpression() {
		segments = append(segments, v.Visit(e).(nodes.Expression))
	}
	return segments
}

func (v *RusterVisitor) VisitPathExpression(ctx *parser.PathExpressionContext) interface{} {
	segments := make([]string, 0)
	for _, e := range ctx.AllSimplePathSegment() {
		segments = append(segments, v.Visit(e).(string))
	}
	return segments
}

func (v *RusterVisitor) VisitIfExpression(ctx *parser.IfExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitMatchExpression(ctx *parser.MatchExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitMatchArms(ctx *parser.MatchArmsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitMatchArmExpression(ctx *parser.MatchArmExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitMatchArm(ctx *parser.MatchArmContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitInfiniteLoopExpression(ctx *parser.InfiniteLoopExpressionContext) interface{} {
	var expr nodes.InfiniteLoopExpression
	expr.Body = v.Visit(ctx.BlockExpression()).(nodes.BlockExpression)
	return expr
}

func (v *RusterVisitor) VisitPredicateLoopExpression(ctx *parser.PredicateLoopExpressionContext) interface{} {
	var expr nodes.PredicateLoopExpression
	expr.Expr = v.Visit(ctx.Expression()).(nodes.Expression)
	expr.Body = v.Visit(ctx.BlockExpression()).(nodes.BlockExpression)
	return expr
}

func (v *RusterVisitor) VisitIteratorLoopExpression(ctx *parser.IteratorLoopExpressionContext) interface{} {
	var expr nodes.IteratorLoopExpression
	expr.Ptrn = v.Visit(ctx.Pattern()).(nodes.Pattern)
	expr.Expr = v.Visit(ctx.Expression()).(nodes.Expression)
	expr.Body = v.Visit(ctx.BlockExpression()).(nodes.BlockExpression)
	return expr
}

func (v *RusterVisitor) VisitLiteralPattern(ctx *parser.LiteralPatternContext) interface{} {
	var literal nodes.LiteralPattern
	literal.Val = ctx.GetText()
	switch vocabulary[ctx.Literal.GetTokenType()] {
	case "STRING_LITERAL":
		literal.Tp = nodes.String
	case "CHAR_LITERAL":
		literal.Tp = nodes.Char
	case "INTEGER_LITERAL":
		literal.Tp = nodes.Integer
	case "KW_TRUE":
		literal.Tp = nodes.Boolean
	case "KW_FALSE":
		literal.Tp = nodes.Boolean
	}
	return literal
}

func (v *RusterVisitor) VisitIdentifierPattern(ctx *parser.IdentifierPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitWildcardPattern(ctx *parser.WildcardPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitRestPattern(ctx *parser.RestPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitReferencePattern(ctx *parser.ReferencePatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitStructPattern(ctx *parser.StructPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitStructPatternElements(ctx *parser.StructPatternElementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitStructPatternFields(ctx *parser.StructPatternFieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitStructPatternField(ctx *parser.StructPatternFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitTuplePattern(ctx *parser.TuplePatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitTuplePatternItems(ctx *parser.TuplePatternItemsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitGroupedPattern(ctx *parser.GroupedPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitSlicePattern(ctx *parser.SlicePatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitSlicePatternItems(ctx *parser.SlicePatternItemsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitPathPattern(ctx *parser.PathPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitRangePattern(ctx *parser.RangePatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitRangePatternBound(ctx *parser.RangePatternBoundContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitNestedType(ctx *parser.NestedTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitNeverType(ctx *parser.NeverTypeContext) interface{} {
	var tp nodes.NeverType
	return tp
}

func (v *RusterVisitor) VisitInferredType(ctx *parser.InferredTypeContext) interface{} {
	var tp nodes.InferredType
	return tp
}

func (v *RusterVisitor) VisitTupleType(ctx *parser.TupleTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitArrayType(ctx *parser.ArrayTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitSliceType(ctx *parser.SliceTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitReferenceType(ctx *parser.ReferenceTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitPointerType(ctx *parser.PointerTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitFunctionType(ctx *parser.FunctionTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *RusterVisitor) VisitTypePath(ctx *parser.TypePathContext) interface{} {
	var path nodes.TypePath
	path.Segments = make([]nodes.TypePathSegment, 0)
	for _, e := range ctx.AllTypePathSegment() {
		path.Segments = append(path.Segments, v.Visit(e).(nodes.TypePathSegment))
	}
	return path
}

func (v *RusterVisitor) VisitTypePathSegment(ctx *parser.TypePathSegmentContext) interface{} {
	var segment nodes.TypePathSegment
	segment.ID = v.Visit(ctx.SimplePathSegment()).(string)
	if ctx.TypePathFn().IsEmpty() {
		segment.Fn = nil
	} else {
		segment.Fn = v.Visit(ctx.TypePathFn()).(nodes.TypePathFunction)
	}
	return segment
}

func (v *RusterVisitor) VisitTypePathFn(ctx *parser.TypePathFnContext) interface{} {
	var fn nodes.TypePathFunction
	fn.Inputs = v.Visit(ctx.TypePathInputs()).([]nodes.Type)
	if ctx.ReturnType.IsEmpty() {
		fn.ReturnType = nil
	} else {
		fn.ReturnType = v.Visit(ctx.ReturnType).(nodes.Type)
	}
	return fn
}

func (v *RusterVisitor) VisitTypePathInputs(ctx *parser.TypePathInputsContext) interface{} {
	inputs := make([]nodes.Type, 0)
	for _, e := range ctx.AllType() {
		inputs = append(inputs, v.Visit(e).(nodes.Type))
	}
	return inputs
}

func (v *RusterVisitor) VisitSimplePath(ctx *parser.SimplePathContext) interface{} {
	path := make([]string, 0)
	for _, e := range ctx.AllSimplePathSegment() {
		path = append(path, v.Visit(e).(string))
	}
	return path
}

func (v *RusterVisitor) VisitSimplePathSegment(ctx *parser.SimplePathSegmentContext) interface{} {
	return ctx.GetText()
}

func (v *RusterVisitor) VisitIdentifier(ctx *parser.IdentifierContext) interface{} {
	return ctx.GetText()
}

func (v *RusterVisitor) VisitKeyword(ctx *parser.KeywordContext) interface{} {
	return ctx.GetText()
}

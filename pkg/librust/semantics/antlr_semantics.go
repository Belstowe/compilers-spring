package semantics

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Compiler2022/compilers-1-Belstowe/parser"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

var vocabulary []string = parser.NewRustLexer(antlr.NewInputStream("")).GetSymbolicNames()

type MessageType int

const (
	INFO MessageType = iota
	WARN
	ERROR
)

func (mt MessageType) String() string {
	switch mt {
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	}
	return "???"
}

type Message struct {
	Type MessageType
	Desc string
}

func (m Message) String() string {
	return fmt.Sprintf("[%s] %s", m.Type.String(), m.Desc)
}

type PatternReturn struct {
	Name     string
	Type     string
	Volatile bool
	Size     int
}

type IDAttr struct {
	Volatile  bool
	Size      int
	Type      string
	BaseType  string
	NumOfElem int
}

func (id IDAttr) String() string {
	return fmt.Sprintf("{Volatile: %t, Size: %d, Type: %s, BaseType: %s, NumOfElem: %d}", id.Volatile, id.Size, id.Type, id.BaseType, id.NumOfElem)
}

type Scope map[string]IDAttr

type ANTLRSymtabVisitor struct {
	antlr.ParseTreeVisitor
	scopes []Scope
	logs   []Message
}

func NewANTLRSymtabVisitor() parser.RustParserVisitor {
	return &ANTLRSymtabVisitor{
		ParseTreeVisitor: &parser.BaseRustParserVisitor{},
		scopes:           make([]Scope, 0),
		logs:             make([]Message, 0),
	}
}

func (v *ANTLRSymtabVisitor) log(tp MessageType, msg string) {
	v.logs = append(v.logs, Message{
		Type: tp,
		Desc: msg,
	})
}

func (v ANTLRSymtabVisitor) GetLogs() []Message {
	return v.logs
}

func (v *ANTLRSymtabVisitor) enterScope() {
	v.scopes = append(v.scopes, Scope{})
	v.log(INFO, fmt.Sprintf("Entering scope %d...", len(v.scopes)))
}

func (v *ANTLRSymtabVisitor) exitScope() {
	v.log(INFO, fmt.Sprintf("Leaving scope %d...", len(v.scopes)))
	v.scopes = v.scopes[:len(v.scopes)-1]
}

func (v ANTLRSymtabVisitor) hasInScope(varName string, index int) bool {
	if _, ok := v.scopes[index][varName]; ok {
		return true
	}
	return false
}

func (v *ANTLRSymtabVisitor) declare(varName string, attr IDAttr) {
	v.log(INFO, fmt.Sprintf("{Scope %d} Declaring var %s %s...", len(v.scopes), varName, attr.String()))
	if v.hasInScope(varName, len(v.scopes)-1) {
		v.log(ERROR, fmt.Sprintf("'%s': already defined in the same scope %d!", varName, len(v.scopes)))
	} else {
		for i := len(v.scopes) - 2; i >= 0; i-- {
			if v.hasInScope(varName, i) {
				v.log(WARN, fmt.Sprintf("'%s': redefined in scope %d (earlier definition in scope %d: %s)!", varName, len(v.scopes), i+1, v.scopes[i][varName].String()))
			}
		}
	}
	v.scopes[len(v.scopes)-1][varName] = attr
}

func (v *ANTLRSymtabVisitor) seeDeclaredType(varName string) string {
	for i := len(v.scopes) - 1; i >= 0; i-- {
		if v.hasInScope(varName, i) {
			return v.scopes[i][varName].Type
		}
	}
	return "nil"
}

func (v *ANTLRSymtabVisitor) seeDeclaredBaseType(varName string) string {
	for i := len(v.scopes) - 1; i >= 0; i-- {
		if v.hasInScope(varName, i) {
			return v.scopes[i][varName].BaseType
		}
	}
	return "nil"
}

func (v *ANTLRSymtabVisitor) assertDeclared(varName string) {
	if v.seeDeclaredType(varName) == "nil" {
		v.log(ERROR, fmt.Sprintf("%s: undeclared", varName))
	}
}

func (v *ANTLRSymtabVisitor) assertTypeDeclared(varName string, tpName string) {
	if v.seeDeclaredType(varName) != tpName {
		v.log(ERROR, fmt.Sprintf("%s: undeclared (expected %s)", varName, tpName))
	}
}

func (v *ANTLRSymtabVisitor) assertTypesDeclared(varName string, tpNames []string) {
	for _, tpName := range tpNames {
		if v.seeDeclaredType(varName) == tpName {
			return
		}
	}
	tpNamesJoined := strings.Join(tpNames, "; ")
	v.log(ERROR, fmt.Sprintf("%s: undeclared (expected %s)", varName, tpNamesJoined))
}

func (v *ANTLRSymtabVisitor) assertTupleDeclared(varName string, index int) {
	for i := len(v.scopes) - 1; i >= 0; i-- {
		if v.hasInScope(varName, len(v.scopes)-1) {
			baseTp := v.seeDeclaredBaseType(varName)
			if baseTp != "tuple" && baseTp != "struct" {
				v.log(ERROR, fmt.Sprintf("%s: not a tuple/struct", varName))
			}
			if index >= v.scopes[i][varName].NumOfElem {
				v.log(ERROR, fmt.Sprintf("%s.%d: out of boundaries (declared size: %d)", varName, index, v.scopes[i][varName].NumOfElem))
			}
			return
		}
	}
	v.log(ERROR, fmt.Sprintf("%s: undeclared", varName))
}

func (v *ANTLRSymtabVisitor) combineIDAttr(first IDAttr, second IDAttr) IDAttr {
	var result IDAttr = second
	if first.Type != "" {
		result.Type = first.Type
	}
	if first.BaseType != "" {
		result.BaseType = first.BaseType
	}
	if first.Volatile {
		result.Volatile = first.Volatile
	}
	result.Size *= first.Size
	result.NumOfElem += first.NumOfElem
	return result
}

func (v *ANTLRSymtabVisitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

func (v *ANTLRSymtabVisitor) VisitTerminal(node antlr.TerminalNode) interface{} {
	return node.GetText()
}

func (v *ANTLRSymtabVisitor) VisitChildren(node antlr.RuleNode) interface{} {
	result := make([]interface{}, 0)
	for _, child := range node.GetChildren() {
		switch child := child.(type) {
		case antlr.ErrorNode:
			v.VisitErrorNode(child)
		case antlr.RuleNode:
			result = append(result, v.Visit(child))
		case antlr.TerminalNode:
			result = append(result, v.VisitTerminal(child))
		}
	}
	if len(result) == 1 {
		return result[0]
	}
	return result
}

func (v *ANTLRSymtabVisitor) VisitCrate(ctx *parser.CrateContext) interface{} {
	v.enterScope()
	// built-in types
	v.declare("i8", IDAttr{
		Type:     "type",
		BaseType: "byte",
		Size:     1,
	})
	v.declare("i16", IDAttr{
		Type:     "type",
		BaseType: "short",
		Size:     2,
	})
	v.declare("i32", IDAttr{
		Type:     "type",
		BaseType: "int",
		Size:     4,
	})
	v.declare("i64", IDAttr{
		Type:     "type",
		BaseType: "long",
		Size:     8,
	})
	v.declare("isize", IDAttr{
		Type:     "typedef",
		BaseType: "i64",
	})
	v.declare("u8", IDAttr{
		Type:     "type",
		BaseType: "ubyte",
		Size:     1,
	})
	v.declare("u16", IDAttr{
		Type:     "type",
		BaseType: "ushort",
		Size:     2,
	})
	v.declare("u32", IDAttr{
		Type:     "type",
		BaseType: "uint",
		Size:     4,
	})
	v.declare("u64", IDAttr{
		Type:     "type",
		BaseType: "ulong",
		Size:     8,
	})
	v.declare("usize", IDAttr{
		Type:     "typedef",
		BaseType: "u64",
	})
	v.declare("f32", IDAttr{
		Type:     "type",
		BaseType: "float",
		Size:     4,
	})
	v.declare("f64", IDAttr{
		Type:     "type",
		BaseType: "double",
		Size:     8,
	})
	v.declare("char", IDAttr{
		Type:     "typedef",
		BaseType: "u32",
	})
	v.declare("str", IDAttr{
		Type:     "type",
		BaseType: "array",
	})
	v.declare("bool", IDAttr{
		Type:     "typedef",
		BaseType: "u8",
	})
	v.declare("len", IDAttr{
		Type:     "function",
		BaseType: "isize",
	})
	v.declare("iter", IDAttr{
		Type:     "function",
		BaseType: "pointer",
	})
	v.declare("String", IDAttr{
		Type:     "type",
		BaseType: "struct",
	})
	v.declare("Some", IDAttr{
		Type:     "type",
		BaseType: "struct",
	})
	v.declare("None", IDAttr{
		Type:     "type",
		BaseType: "struct",
	})
	v.declare("std", IDAttr{
		Type:     "namespace",
		BaseType: "namespace",
	})
	for _, element := range ctx.AllItem() {
		v.Visit(element)
	}
	v.exitScope()
	return v.logs
}

func (v *ANTLRSymtabVisitor) VisitUseDeclaration(ctx *parser.UseDeclarationContext) interface{} {
	return v.Visit(ctx.UseTree())
}

func (v *ANTLRSymtabVisitor) VisitUseTree(ctx *parser.UseTreeContext) interface{} {
	path := v.Visit(ctx.SimplePath()).([]string)
	for i := range path {
		pathSegment := path[i:]
		v.declare(strings.Join(pathSegment, "::"), IDAttr{
			Type:     "namespace",
			BaseType: "namespace",
		})
	}
	return nil
}

func (v *ANTLRSymtabVisitor) VisitFunction(ctx *parser.FunctionContext) interface{} {
	id := v.Visit(ctx.Identifier()).(string)
	var returnTp string = "void"
	if ctx.FunctionReturnType() != nil {
		returnTp = v.Visit(ctx.FunctionReturnType()).(IDAttr).BaseType
	}
	v.declare(id, IDAttr{
		Type:     "function",
		BaseType: returnTp,
	})

	v.enterScope()
	v.Visit(ctx.FunctionParameters())
	v.Visit(ctx.BlockExpression())
	v.exitScope()

	return nil
}

func (v *ANTLRSymtabVisitor) VisitFunctionParameters(ctx *parser.FunctionParametersContext) interface{} {
	for _, e := range ctx.AllFunctionParam() {
		v.Visit(e)
	}
	return nil
}

func (v *ANTLRSymtabVisitor) VisitFunctionParam(ctx *parser.FunctionParamContext) interface{} {
	vartp := v.Visit(ctx.Type()).(IDAttr)
	id := v.Visit(ctx.Identifier()).(string)
	if vartp.Type == "value" {
		vartp.Volatile = true
	}
	v.declare(id, vartp)
	return nil
}

func (v *ANTLRSymtabVisitor) VisitFunctionReturnType(ctx *parser.FunctionReturnTypeContext) interface{} {
	return v.Visit(ctx.Type())
}

func (v *ANTLRSymtabVisitor) VisitLetStatement(ctx *parser.LetStatementContext) interface{} {
	patternReturn := v.Visit(ctx.Pattern()).(PatternReturn)

	if ctx.Expression() != nil {
		v.Visit(ctx.Expression())
	}

	var vartp IDAttr
	if ctx.Type() != nil {
		vartp = v.Visit(ctx.Type()).(IDAttr)
	}
	v.declare(patternReturn.Name, v.combineIDAttr(vartp, IDAttr{Type: "value", BaseType: patternReturn.Type, Volatile: patternReturn.Volatile}))

	return nil
}

func (v *ANTLRSymtabVisitor) VisitTypeCastExpression(ctx *parser.TypeCastExpressionContext) interface{} {
	v.Visit(ctx.Type())
	v.Visit(ctx.Expression())
	return nil
}

func (v *ANTLRSymtabVisitor) VisitTupleExpression(ctx *parser.TupleExpressionContext) interface{} {
	v.Visit(ctx.TupleElements())
	return nil
}

func (v *ANTLRSymtabVisitor) VisitIndexExpression(ctx *parser.IndexExpressionContext) interface{} {
	v.Visit(ctx.Index)
	obj := v.Visit(ctx.Object).(string)
	return obj
}

func (v *ANTLRSymtabVisitor) VisitRangeExpression(ctx *parser.RangeExpressionContext) interface{} {
	v.Visit(ctx.LHS)
	v.Visit(ctx.RHS)
	return nil
}

func (v *ANTLRSymtabVisitor) VisitReturnExpression(ctx *parser.ReturnExpressionContext) interface{} {
	if ctx.Expression() != nil {
		return v.Visit(ctx.Expression())
	}
	return nil
}

func (v *ANTLRSymtabVisitor) VisitContinueExpression(ctx *parser.ContinueExpressionContext) interface{} {
	if ctx.Expression() != nil {
		return v.Visit(ctx.Expression())
	}
	return nil
}

func (v *ANTLRSymtabVisitor) VisitAssignmentExpression(ctx *parser.AssignmentExpressionContext) interface{} {
	v.Visit(ctx.RHS)
	return v.Visit(ctx.LHS)
}

func (v *ANTLRSymtabVisitor) VisitMethodCallExpression(ctx *parser.MethodCallExpressionContext) interface{} {
	if ctx.CallParams() != nil {
		v.Visit(ctx.CallParams())
	}
	v.Visit(ctx.SimplePathSegment())
	//method := v.Visit(ctx.SimplePathSegment()).(string)
	//v.assertTypeDeclared(method, "function")
	return v.Visit(ctx.Expression())
}

func (v *ANTLRSymtabVisitor) VisitLiteralExpression_(ctx *parser.LiteralExpression_Context) interface{} {
	return v.Visit(ctx.LiteralExpression())
}

func (v *ANTLRSymtabVisitor) VisitStructExpression_(ctx *parser.StructExpression_Context) interface{} {
	return v.Visit(ctx.StructExpression())
}

func (v *ANTLRSymtabVisitor) VisitTupleIndexingExpression(ctx *parser.TupleIndexingExpressionContext) interface{} {
	index := v.Visit(ctx.TupleIndex()).(int)
	id := v.Visit(ctx.Expression()).(string)
	v.assertTupleDeclared(id, index)
	return id
}

func (v *ANTLRSymtabVisitor) VisitCallExpression(ctx *parser.CallExpressionContext) interface{} {
	if ctx.CallParams() != nil {
		v.Visit(ctx.CallParams())
	}
	header := v.Visit(ctx.Expression()).(string)
	pathSegments := strings.Split(header, "::")
	if len(pathSegments) > 1 {
		namespaceToFind := strings.Join(pathSegments[:len(pathSegments)-1], "::")
		v.assertTypeDeclared(namespaceToFind, "namespace")
	} else {
		v.assertTypesDeclared(header, []string{"function", "type"})
	}
	return header
}

func (v *ANTLRSymtabVisitor) VisitDereferenceExpression(ctx *parser.DereferenceExpressionContext) interface{} {
	return v.Visit(ctx.Expression())
}

func (v *ANTLRSymtabVisitor) VisitNestedExpression(ctx *parser.NestedExpressionContext) interface{} {
	return v.Visit(ctx.Expression())
}

func (v *ANTLRSymtabVisitor) VisitUnaryOpExpression(ctx *parser.UnaryOpExpressionContext) interface{} {
	return v.Visit(ctx.Val)
}

func (v *ANTLRSymtabVisitor) VisitBinaryOpExpression(ctx *parser.BinaryOpExpressionContext) interface{} {
	v.Visit(ctx.RHS)
	return v.Visit(ctx.LHS)
}

func (v *ANTLRSymtabVisitor) VisitBreakExpression(ctx *parser.BreakExpressionContext) interface{} {
	v.Visit(ctx.Expression())
	return nil
}

func (v *ANTLRSymtabVisitor) VisitBorrowExpression(ctx *parser.BorrowExpressionContext) interface{} {
	return v.Visit(ctx.Expression())
}

func (v *ANTLRSymtabVisitor) VisitCompoundAssignmentExpression(ctx *parser.CompoundAssignmentExpressionContext) interface{} {
	v.Visit(ctx.RHS)
	return v.Visit(ctx.LHS)
}

func (v *ANTLRSymtabVisitor) VisitArrayExpression(ctx *parser.ArrayExpressionContext) interface{} {
	v.Visit(ctx.ArrayElements())
	return nil
}

func (v *ANTLRSymtabVisitor) VisitCompoundAssignOperator(ctx *parser.CompoundAssignOperatorContext) interface{} {
	return ctx.GetText()
}

func (v *ANTLRSymtabVisitor) VisitLiteralExpression(ctx *parser.LiteralExpressionContext) interface{} {
	switch vocabulary[ctx.Literal.GetTokenType()] {
	case "STRING_LITERAL":
		return "str"
	case "CHAR_LITERAL":
		return "char"
	case "INTEGER_LITERAL":
		return "i64"
	case "KW_TRUE":
	case "KW_FALSE":
		return "bool"
	}
	v.log(ERROR, fmt.Sprintf("'%s': unknown literal type '%s'", ctx.Literal.GetText(), vocabulary[ctx.Literal.GetTokenType()]))
	return "unknown"
}

func (v *ANTLRSymtabVisitor) VisitStatements(ctx *parser.StatementsContext) interface{} {
	v.enterScope()

	for _, statement := range ctx.AllStatement() {
		v.Visit(statement)
	}

	if ctx.Expression() != nil {
		returnResult := v.Visit(ctx.Expression())
		v.exitScope()
		return returnResult
	}
	v.exitScope()
	return nil
}

func (v *ANTLRSymtabVisitor) VisitArrayElements(ctx *parser.ArrayElementsContext) interface{} {
	for _, e := range ctx.AllExpression() {
		v.Visit(e)
	}
	return nil
}

func (v *ANTLRSymtabVisitor) VisitTupleElements(ctx *parser.TupleElementsContext) interface{} {
	for _, e := range ctx.AllExpression() {
		v.Visit(e)
	}
	return nil
}

func (v *ANTLRSymtabVisitor) VisitTupleIndex(ctx *parser.TupleIndexContext) interface{} {
	literalString := ctx.INTEGER_LITERAL().GetText()
	literalInteger, err := strconv.Atoi(literalString)
	if err != nil {
		v.log(ERROR, fmt.Sprintf("Tuple indexing: got string '%s' instead of integer", literalString))
		return 0
	}
	return literalInteger
}

func (v *ANTLRSymtabVisitor) VisitCallParams(ctx *parser.CallParamsContext) interface{} {
	for _, e := range ctx.AllExpression() {
		v.Visit(e)
	}
	return nil
}

func (v *ANTLRSymtabVisitor) VisitPathExpression(ctx *parser.PathExpressionContext) interface{} {
	segments := make([]string, 0)
	for _, e := range ctx.AllSimplePathSegment() {
		segments = append(segments, v.Visit(e).(string))
	}
	return strings.Join(segments, "::")
}

func (v *ANTLRSymtabVisitor) VisitIfExpression(ctx *parser.IfExpressionContext) interface{} {
	v.Visit(ctx.Cond)
	v.Visit(ctx.IfBody)
	if ctx.ElseIf != nil {
		v.Visit(ctx.ElseIf)
	} else if ctx.ElseBody != nil {
		v.Visit(ctx.ElseBody)
	}
	return nil
}

func (v *ANTLRSymtabVisitor) VisitInfiniteLoopExpression(ctx *parser.InfiniteLoopExpressionContext) interface{} {
	v.Visit(ctx.BlockExpression())
	return nil
}

func (v *ANTLRSymtabVisitor) VisitPredicateLoopExpression(ctx *parser.PredicateLoopExpressionContext) interface{} {
	v.Visit(ctx.Expression())
	v.Visit(ctx.BlockExpression())
	return nil
}

func (v *ANTLRSymtabVisitor) VisitIteratorLoopExpression(ctx *parser.IteratorLoopExpressionContext) interface{} {
	v.Visit(ctx.Pattern())
	v.Visit(ctx.Expression())
	v.Visit(ctx.BlockExpression())
	return nil
}

func (v *ANTLRSymtabVisitor) VisitLiteralPattern(ctx *parser.LiteralPatternContext) interface{} {
	switch vocabulary[ctx.Literal.GetTokenType()] {
	case "STRING_LITERAL":
		return PatternReturn{Type: "str"}
	case "CHAR_LITERAL":
		return PatternReturn{Type: "char"}
	case "INTEGER_LITERAL":
		return PatternReturn{Type: "i64"}
	case "KW_TRUE":
	case "KW_FALSE":
		return PatternReturn{Type: "bool"}
	}
	v.log(ERROR, fmt.Sprintf("%s: unknown literal type '%s'", ctx.Literal.GetText(), vocabulary[ctx.Literal.GetTokenType()]))
	return PatternReturn{}
}

func (v *ANTLRSymtabVisitor) VisitIdentifierPattern(ctx *parser.IdentifierPatternContext) interface{} {
	var ptrn PatternReturn
	ptrn.Name = ctx.Identifier().GetText()
	ptrn.Volatile = ctx.KW_MUT() != nil
	return ptrn
}

func (v *ANTLRSymtabVisitor) VisitNeverType(ctx *parser.NeverTypeContext) interface{} {
	return IDAttr{Type: ctx.GetText()}
}

func (v *ANTLRSymtabVisitor) VisitInferredType(ctx *parser.InferredTypeContext) interface{} {
	return IDAttr{Type: ctx.GetText()}
}

func (v *ANTLRSymtabVisitor) VisitTupleType(ctx *parser.TupleTypeContext) interface{} {
	var tp IDAttr
	tp.Type = "tuple"
	for _, elem := range ctx.AllType() {
		elemAttr := v.Visit(elem).(IDAttr)
		tp.Size += elemAttr.Size
		tp.NumOfElem++
	}
	return tp
}

func (v *ANTLRSymtabVisitor) VisitArrayType(ctx *parser.ArrayTypeContext) interface{} {
	var tp IDAttr
	tp.Type = "array"
	elemAttr := v.Visit(ctx.Type()).(IDAttr)
	tp.BaseType = elemAttr.Type
	return tp
}

func (v *ANTLRSymtabVisitor) VisitSliceType(ctx *parser.SliceTypeContext) interface{} {
	var tp IDAttr
	tp.Type = "array"
	elemAttr := v.Visit(ctx.Type()).(IDAttr)
	tp.BaseType = elemAttr.Type
	return tp
}

func (v *ANTLRSymtabVisitor) VisitReferenceType(ctx *parser.ReferenceTypeContext) interface{} {
	tp := v.Visit(ctx.Type()).(IDAttr)
	tp.BaseType = tp.Type
	tp.Type = "reference"
	tp.Volatile = ctx.Mutable != nil
	return tp
}

func (v *ANTLRSymtabVisitor) VisitPointerType(ctx *parser.PointerTypeContext) interface{} {
	tp := v.Visit(ctx.Type()).(IDAttr)
	tp.BaseType = tp.Type
	tp.Type = "pointer"
	tp.Volatile = ctx.Mutable != nil
	return tp
}

func (v *ANTLRSymtabVisitor) VisitTypePath(ctx *parser.TypePathContext) interface{} {
	segments := make([]string, 0)
	for _, e := range ctx.AllTypePathSegment() {
		segments = append(segments, v.Visit(e).(string))
	}
	fullName := strings.Join(segments, "::")
	v.assertTypesDeclared(fullName, []string{"type", "typedef"})
	return IDAttr{Type: "value", BaseType: fullName}
}

func (v *ANTLRSymtabVisitor) VisitTypePathSegment(ctx *parser.TypePathSegmentContext) interface{} {
	return v.Visit(ctx.SimplePathSegment()).(string)
}

func (v *ANTLRSymtabVisitor) VisitTypePathFn(ctx *parser.TypePathFnContext) interface{} {
	return nil
}

func (v *ANTLRSymtabVisitor) VisitTypePathInputs(ctx *parser.TypePathInputsContext) interface{} {
	return nil
}

func (v *ANTLRSymtabVisitor) VisitSimplePath(ctx *parser.SimplePathContext) interface{} {
	path := make([]string, 0)
	for _, e := range ctx.AllSimplePathSegment() {
		path = append(path, v.Visit(e).(string))
	}
	return path
}

func (v *ANTLRSymtabVisitor) VisitSimplePathSegment(ctx *parser.SimplePathSegmentContext) interface{} {
	return ctx.GetText()
}

func (v *ANTLRSymtabVisitor) VisitIdentifier(ctx *parser.IdentifierContext) interface{} {
	return ctx.GetText()
}

func (v *ANTLRSymtabVisitor) VisitKeyword(ctx *parser.KeywordContext) interface{} {
	return ctx.GetText()
}

func (v *ANTLRSymtabVisitor) VisitBlockExpression(ctx *parser.BlockExpressionContext) interface{} {
	return v.Visit(ctx.Statements())
}

func (v *ANTLRSymtabVisitor) VisitExpressionStatement(ctx *parser.ExpressionStatementContext) interface{} {
	if ctx.Expression() != nil {
		return v.Visit(ctx.Expression())
	} else if ctx.ExpressionWithBlock() != nil {
		return v.Visit(ctx.ExpressionWithBlock())
	}
	return nil
}

func (v *ANTLRSymtabVisitor) VisitRHSRangeExpression(ctx *parser.RHSRangeExpressionContext) interface{} {
	return v.Visit(ctx.Val)
}

func (v *ANTLRSymtabVisitor) VisitStatement(ctx *parser.StatementContext) interface{} {
	if ctx.Item() != nil {
		return v.Visit(ctx.Item())
	} else if ctx.LetStatement() != nil {
		return v.Visit(ctx.LetStatement())
	} else if ctx.ExpressionStatement() != nil {
		return v.Visit(ctx.ExpressionStatement())
	}
	return ctx.GetText()
}

/*
 * TBD
 */

func (v *ANTLRSymtabVisitor) VisitErrorPropagationExpression(ctx *parser.ErrorPropagationExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitWildcardPattern(ctx *parser.WildcardPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitRestPattern(ctx *parser.RestPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitReferencePattern(ctx *parser.ReferencePatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitStructPattern(ctx *parser.StructPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitStructPatternElements(ctx *parser.StructPatternElementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitStructPatternFields(ctx *parser.StructPatternFieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitStructPatternField(ctx *parser.StructPatternFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitTuplePattern(ctx *parser.TuplePatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitTuplePatternItems(ctx *parser.TuplePatternItemsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitGroupedPattern(ctx *parser.GroupedPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitSlicePattern(ctx *parser.SlicePatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitSlicePatternItems(ctx *parser.SlicePatternItemsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitPathPattern(ctx *parser.PathPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitRangePattern(ctx *parser.RangePatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitRangePatternBound(ctx *parser.RangePatternBoundContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitFunctionType(ctx *parser.FunctionTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitNestedType(ctx *parser.NestedTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitType(ctx *parser.TypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitStruct(ctx *parser.StructContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitStructFields(ctx *parser.StructFieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitStructField(ctx *parser.StructFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitTypeAlias(ctx *parser.TypeAliasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitConstantItem(ctx *parser.ConstantItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitExpressionWithBlock(ctx *parser.ExpressionWithBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitExpressionWithBlock_(ctx *parser.ExpressionWithBlock_Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitItem(ctx *parser.ItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitLoopExpression(ctx *parser.LoopExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitNonRangePattern(ctx *parser.NonRangePatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitPathExpression_(ctx *parser.PathExpression_Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitFieldExpression(ctx *parser.FieldExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitStructExpression(ctx *parser.StructExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitStructExprFields(ctx *parser.StructExprFieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitStructExprField(ctx *parser.StructExprFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitPattern(ctx *parser.PatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitMatchExpression(ctx *parser.MatchExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitMatchArms(ctx *parser.MatchArmsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitMatchArmExpression(ctx *parser.MatchArmExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ANTLRSymtabVisitor) VisitMatchArm(ctx *parser.MatchArmContext) interface{} {
	return v.VisitChildren(ctx)
}

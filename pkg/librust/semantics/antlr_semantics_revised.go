package semantics

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Compiler2022/compilers-1-Belstowe/pkg/librust/ast"
	"github.com/llir/llvm/ir/types"
)

type ANTLRSemVisitor struct {
	ast.RusterBaseVisitor
	scopes []Scope
	logs   []Message
}

func NewANTLRSemVisitor() ast.RusterBaseVisitor {
	return &ANTLRSemVisitor{
		scopes: make([]Scope, 0),
		logs:   make([]Message, 0),
	}
}

func (v *ANTLRSemVisitor) Visit(tree ast.Node) interface{} {
	return tree.Accept(v)
}

func (v *ANTLRSemVisitor) VisitCrate(c *ast.Crate) interface{} {
	v.enterScope()
	v.declare(IDAttr{
		Name:      "",
		TypeParam: TypeAttr{BaseType: types.Void},
	})
	v.declare(IDAttr{
		Name: "u8",
		TypeParam: TypeAttr{
			BaseType: types.I8,
		},
	})
	v.declare(IDAttr{
		Name: "u16",
		TypeParam: TypeAttr{
			BaseType: types.I16,
		},
	})
	v.declare(IDAttr{
		Name: "u32",
		TypeParam: TypeAttr{
			BaseType: types.I32,
		},
	})
	v.declare(IDAttr{
		Name: "u64",
		TypeParam: TypeAttr{
			BaseType: types.I64,
		},
	})
	v.declare(IDAttr{
		Name: "i8",
		TypeParam: TypeAttr{
			BaseType: types.I8,
		},
	})
	v.declare(IDAttr{
		Name: "i16",
		TypeParam: TypeAttr{
			BaseType: types.I16,
		},
	})
	v.declare(IDAttr{
		Name: "i32",
		TypeParam: TypeAttr{
			BaseType: types.I32,
		},
	})
	v.declare(IDAttr{
		Name: "i64",
		TypeParam: TypeAttr{
			BaseType: types.I64,
		},
	})
	v.declare(IDAttr{
		Name: "isize",
		TypeParam: TypedefAttr{
			Type: v.seeDeclared("i64"),
		},
	})
	v.declare(IDAttr{
		Name: "char",
		TypeParam: TypedefAttr{
			Type: v.seeDeclared("i8"),
		},
	})
	v.declare(IDAttr{
		Name: "str",
		TypeParam: PointerAttr{
			Volatile: false,
			Type:     v.seeDeclared("char"),
		},
	})
	v.declare(IDAttr{
		Name: "ruster::writeln_i64",
		TypeParam: FuncAttr{
			ReturnType: TypeAttr{BaseType: types.Void},
			CallParam: []TypeDef{
				TypeAttr{BaseType: types.I64},
			},
		},
	})

	for _, item := range c.Items {
		item.Accept(v)
	}
	v.exitScope()
	return nil
}

func (v *ANTLRSemVisitor) VisitBlockExpression(be *ast.BlockExpression) interface{} {
	v.enterScope()

	returnTypes := make([]TypeDef, 0)
	for _, ex := range be.Statements {
		exReturn := ex.Accept(v)
		if exReturn != nil {
			exReturnAttr := exReturn.(IDAttr)
			switch toReturnAttr := exReturnAttr.TypeParam.(type) {
			case ReturnAttr:
				returnTypes = append(returnTypes, GetToType(toReturnAttr))
			}
		}
	}
	if be.Expr != nil {
		returnTypes = append(returnTypes, GetToType(be.Expr.Accept(v).(IDAttr).TypeParam))
	}
	if len(returnTypes) == 0 {
		returnTypes = append(returnTypes, v.seeDeclared(""))
	}

	v.exitScope()

	returnTypeSet := make(map[TypeDef]struct{})
	for _, rtp := range returnTypes {
		returnTypeSet[rtp] = struct{}{}
	}
	if len(returnTypeSet) > 1 {
		var setRepr string
		setRepr += "{ "
		for _, rtp := range returnTypes {
			setRepr += rtp.String() + "; "
		}
		setRepr += " }"
		v.log(ERROR, fmt.Sprintf("block expression returns several different types: %s", setRepr))
	}
	return IDAttr{
		TypeParam: returnTypes[len(returnTypes)-1],
	}
}

func (v *ANTLRSemVisitor) VisitUseDecl(ud *ast.UseDecl) interface{} {
	sp := ud.Path.Accept(v).([]string)
	v.declare(IDAttr{
		Name: sp[len(sp)-1],
		TypeParam: NamespaceAttr{
			External: strings.Join(sp, "::"),
		},
	})
	return nil
}

func (v *ANTLRSemVisitor) VisitSimplePath(sp *ast.SimplePath) interface{} {
	return sp
}

func (v *ANTLRSemVisitor) VisitParameter(p *ast.Parameter) interface{} {
	return IDAttr{
		Name: p.ID,
		TypeParam: ValueAttr{
			BaseType: p.VarType.Accept(v).(TypeDef),
		},
	}
}

func (v *ANTLRSemVisitor) VisitFunction(f *ast.Function) interface{} {
	callParams := make([]TypeDef, len(f.Params))
	for i, param := range f.Params {
		paramAttr := param.Accept(v).(IDAttr)
		callParams[i] = GetToType(paramAttr.TypeParam)
	}

	var returnType TypeDef = v.seeDeclared("")
	if f.ReturnType != nil {
		returnType = f.ReturnType.Accept(v).(TypeDef)
	}

	v.declare(IDAttr{
		Name: f.ID,
		TypeParam: FuncAttr{
			ReturnType: returnType,
			CallParam:  callParams,
		},
	})

	v.enterScope()

	for _, param := range f.Params {
		paramAttr := param.Accept(v).(IDAttr)
		v.declare(paramAttr)
	}

	bodyReturnType := f.Body.Accept(v).(IDAttr)
	if returnType != bodyReturnType.TypeParam {
		v.log(ERROR, fmt.Sprintf("function %s: claimed return type %s doesn't correlate with body return type %s", f.ID, returnType.String(), bodyReturnType.TypeParam.String()))
	}

	v.exitScope()
	return nil
}

func (v *ANTLRSemVisitor) VisitLetStatement(ls *ast.LetStatement) interface{} {
	idDeclared := ls.Ptrn.Accept(v).(IDAttr)
	lastType := &idDeclared.TypeParam
out:
	for {
		switch attr := (*lastType).(type) {
		case RefAttr:
			lastType = &attr.BaseType
		case ValueAttr:
			lastType = &attr.BaseType
		case ArrayAttr:
			lastType = &attr.Type
		case nil:
			break out
		default:
			v.log(ERROR, fmt.Sprintf("%s: unexpected attribute type '%T'", idDeclared.Name, attr))
			return nil
		}
	}

	if ls.Expr == nil && ls.VarType == nil {
		v.logf(ERROR, "No fields to determine var %s type!", idDeclared.Name)
		return nil
	}

	var exprReturnType TypeDef = nil
	if ls.Expr != nil {
		exprReturnType = GetToType(ls.Expr.Accept(v).(IDAttr))
	}
	var declaredType TypeDef = nil
	if ls.VarType != nil {
		declaredType = GetToType(ls.VarType.Accept(v).(TypeDef))
	}
	if declaredType != nil && exprReturnType != nil {
		if declaredType != exprReturnType {
			v.logf(ERROR, "%s type contradiction: declared type is %s; expression return type is %s", idDeclared.Name, declaredType.String(), exprReturnType.String())
			return nil
		}
	}
	if declaredType == nil {
		declaredType = exprReturnType
	}
	*lastType = declaredType

	v.declare(idDeclared)
	return nil
}

func (v *ANTLRSemVisitor) VisitLiteralExpression(le *ast.LiteralExpression) interface{} {
	var tp TypeDef
	switch le.Tp {
	case ast.String:
		tp = ArrayAttr{
			NumOfElem: len(le.Val),
			Type:      v.seeDeclared("i8"),
		}
	case ast.Boolean:
		tp = v.seeDeclared("i8")
	case ast.Char:
		tp = v.seeDeclared("i8")
	case ast.Integer:
		tp = v.seeDeclared("i64")
	}

	return IDAttr{
		TypeParam: ValueAttr{
			Volatile: false,
			BaseType: tp,
		},
	}
}

func (v *ANTLRSemVisitor) VisitPathExpression(pe *ast.PathExpression) interface{} {
	varId := strings.Join([]string(pe.Segments), "::")
	decl := v.seeDeclared(varId)
	if decl == nil {
		v.log(ERROR, fmt.Sprintf("%s undeclared!", varId))
	}

	return IDAttr{
		Name:      varId,
		TypeParam: decl,
	}
}

func (v *ANTLRSemVisitor) VisitIfExpression(ie *ast.IfExpression) interface{} {
	v.enterScope()

	ie.Expr.Accept(v)
	doReturnType := ie.IfTrue.Accept(v).(IDAttr)
	var elseReturnType IDAttr
	switch elseExpr := ie.IfFalse.(type) {
	case ast.BlockExpression:
	case ast.IfExpression:
		elseReturnType = elseExpr.Accept(v).(IDAttr)
	case nil:
		elseReturnType.TypeParam = nil
	default:
		v.log(ERROR, fmt.Sprintf("invalid expression %T for else block!", elseExpr))
	}
	if doReturnType.TypeParam != elseReturnType.TypeParam {
		if doReturnType.TypeParam != nil && elseReturnType.TypeParam != nil {
			v.log(ERROR, fmt.Sprintf("return type contradiction: do block %s; else block %s", doReturnType.TypeParam.String(), elseReturnType.TypeParam.String()))
		}
	}

	v.exitScope()

	return doReturnType
}

func (v *ANTLRSemVisitor) VisitMatchExpression(me *ast.MatchExpression) interface{} {
	return nil
}

func (v *ANTLRSemVisitor) VisitMatchArm(ma *ast.MatchArm) interface{} {
	return nil
}

func (v *ANTLRSemVisitor) VisitInfiniteLoopExpression(ile *ast.InfiniteLoopExpression) interface{} {
	return ile.Body.Accept(v)
}

func (v *ANTLRSemVisitor) VisitPredicateLoopExpression(ple *ast.PredicateLoopExpression) interface{} {
	v.enterScope()

	ple.Expr.Accept(v)
	rtp := ple.Body.Accept(v)

	v.exitScope()

	return rtp
}

func (v *ANTLRSemVisitor) VisitIteratorLoopExpression(ile *ast.IteratorLoopExpression) interface{} {
	v.enterScope()

	idDeclared := ile.Ptrn.Accept(v).(IDAttr)
	exprReturnType := ile.Expr.Accept(v).(IDAttr)

	v.declare(IDAttr{
		Name:      idDeclared.Name,
		TypeParam: exprReturnType.TypeParam,
	})

	bodyReturnType := ile.Body.Accept(v).(IDAttr)

	v.exitScope()

	return bodyReturnType
}

func (v *ANTLRSemVisitor) VisitUnaryOperator(uo *ast.UnaryOperator) interface{} {
	exprReturnType := uo.Val.Accept(v).(IDAttr)
	if _, ok := allowedUnaryOpTypes[uo.Op]; !ok {
		v.log(ERROR, fmt.Sprintf("unknown operand %s", uo.Op))
		return nil
	}
	switch tp := exprReturnType.TypeParam.(type) {
	case ValueAttr:
		for _, allowedType := range allowedUnaryOpTypes[uo.Op] {
			var allowedTypeAttr TypeDef = v.seeDeclared(allowedType)
			if GetToType(tp.BaseType) == allowedTypeAttr {
				return exprReturnType
			}
		}
	}
	v.log(ERROR, fmt.Sprintf("unsupported type %s", exprReturnType.TypeParam.String()))
	return nil
}

func (v *ANTLRSemVisitor) VisitBinaryOperator(bo *ast.BinaryOperator) interface{} {
	lhsReturnType := bo.LHS.Accept(v).(IDAttr)
	rhsReturnType := bo.RHS.Accept(v).(IDAttr)
	if _, ok := allowedBinaryOpTypes[bo.Op]; !ok {
		v.log(ERROR, fmt.Sprintf("unknown operand %s", bo.Op))
		return nil
	}
	switch lhsTp := lhsReturnType.TypeParam.(type) {
	case ValueAttr:
		switch rhsTp := rhsReturnType.TypeParam.(type) {
		case ValueAttr:
			for _, allowedTypes := range allowedBinaryOpTypes[bo.Op] {
				var lhsAllowedAttr = v.seeDeclared(allowedTypes[0])
				var rhsAllowedAttr = v.seeDeclared(allowedTypes[1])
				if GetToType(lhsTp.BaseType) == lhsAllowedAttr && GetToType(rhsTp.BaseType) == rhsAllowedAttr {
					return lhsReturnType
				}
			}
		}
	}
	v.log(ERROR, fmt.Sprintf("binary operation unsupported types: %s %s %s", lhsReturnType.String(), bo.Op, rhsReturnType.String()))
	return nil
}

func (v *ANTLRSemVisitor) VisitRHSRangeOperator(rro *ast.RHSRangeOperator) interface{} {
	return nil
}

func (v *ANTLRSemVisitor) VisitRangeOperator(ro *ast.RangeOperator) interface{} {
	return nil
}

func (v *ANTLRSemVisitor) VisitReturnExpression(e *ast.ReturnExpression) interface{} {
	return IDAttr{
		TypeParam: ReturnAttr{
			Type: e.Expr.Accept(v).(IDAttr).TypeParam,
		}}
}

func (v *ANTLRSemVisitor) VisitContinueExpression(ce *ast.ContinueExpression) interface{} {
	return ce.Expr.Accept(v)
}

func (v *ANTLRSemVisitor) VisitBreakExpression(be *ast.BreakExpression) interface{} {
	return be.Expr.Accept(v)
}

func (v *ANTLRSemVisitor) VisitTypeCastExpression(tce *ast.TypeCastExpression) interface{} {
	exprReturnType := tce.Expr.Accept(v).(IDAttr)
	convertToType := tce.Tp.Accept(v).(TypeDef)
	return IDAttr{
		Name:      exprReturnType.Name,
		TypeParam: convertToType,
	}
}

func (v *ANTLRSemVisitor) VisitCallExpression(ce *ast.CallExpression) interface{} {
	fnId := ce.FnHeader.Accept(v).(IDAttr).Name
	fnParams := make([]TypeDef, len(ce.Params))
	for i, param := range []ast.Expression(ce.Params) {
		fnParams[i] = GetToType(param.Accept(v).(IDAttr).TypeParam)
	}

	if fnId == "" {
		v.log(ERROR, "empty fn header")
		return nil
	}

	fnCalled := IDAttr{
		Name: fnId,
		TypeParam: FuncAttr{
			CallParam: fnParams,
		},
	}
	varFound := v.seeDeclared(fnId)

	if varFound == nil {
		v.logf(ERROR, "function called not found in symtable; id is %s", fnId)
		return nil
	}
	switch varFound.(type) {
	case FuncAttr:
		break
	default:
		v.logf(ERROR, "function %s got called for different type; its attributes are %s", fnId, varFound.String())
		return nil
	}
	if !reflect.DeepEqual(fnCalled.TypeParam.(FuncAttr).CallParam, varFound.(FuncAttr).CallParam) {
		v.logf(ERROR, "function %s called parameters %s differ from ones declared %s", fnId, fnCalled.TypeParam.String(), varFound.String())
		return nil
	}

	return IDAttr{
		TypeParam: ValueAttr{
			BaseType: varFound.(FuncAttr).ReturnType,
		},
	}
}

func (v *ANTLRSemVisitor) VisitMethodCallExpression(mce *ast.MethodCallExpression) interface{} {
	return nil
}

func (v *ANTLRSemVisitor) VisitBorrowExpression(be *ast.BorrowExpression) interface{} {
	return IDAttr{
		TypeParam: RefAttr{
			Volatile: be.IsMut,
			BaseType: be.Expr.Accept(v).(IDAttr).TypeParam,
		},
	}
}

func (v *ANTLRSemVisitor) VisitArrayIndexExpression(aie *ast.ArrayIndexExpression) interface{} {
	arrayCalled := aie.Object.Accept(v).(IDAttr)
	indexTaken := aie.Index.Accept(v).(IDAttr).TypeParam

	correctIndexType := false
	for _, allowedIndexType := range []string{"i8", "i16", "i32", "i64", "u8", "u16", "u32", "u64"} {
		allowedAttr := v.seeDeclared(allowedIndexType)
		if GetToType(indexTaken) == allowedAttr {
			correctIndexType = true
			break
		}
	}
	if !correctIndexType {
		v.logf(ERROR, "%s is not correct indexing type", indexTaken.String())
		return nil
	}

	switch GetToType(arrayCalled.TypeParam).(type) {
	case ArrayAttr:
		break
	default:
		v.logf(ERROR, "%s is not an array", arrayCalled.String())
		return nil
	}

	return IDAttr{
		TypeParam: arrayCalled.TypeParam.(ArrayAttr).Type,
	}
}

func (v *ANTLRSemVisitor) VisitTupleIndexExpression(tie *ast.TupleIndexExpression) interface{} {
	return nil
}

func (v *ANTLRSemVisitor) VisitLiteralPattern(lp *ast.LiteralPattern) interface{} {
	var tp TypeDef
	switch lp.Tp {
	case ast.String:
		tp = ArrayAttr{
			NumOfElem: len(lp.Val),
			Type:      v.seeDeclared("i8"),
		}
	case ast.Boolean:
		tp = v.seeDeclared("i8")
	case ast.Char:
		tp = v.seeDeclared("i8")
	case ast.Integer:
		tp = v.seeDeclared("i64")
	}

	return IDAttr{
		TypeParam: ValueAttr{
			Volatile: false,
			BaseType: tp,
		},
	}
}

func (v *ANTLRSemVisitor) VisitReferencePattern(rp *ast.ReferencePattern) interface{} {
	return nil
}

func (v *ANTLRSemVisitor) VisitIdentifierPattern(ip *ast.IdentifierPattern) interface{} {
	return IDAttr{
		Name: ip.ID,
		TypeParam: ValueAttr{
			Volatile: ip.IsMut,
			BaseType: nil,
		},
	}
}

func (v *ANTLRSemVisitor) VisitPathPattern(pp *ast.PathPattern) interface{} {
	return nil
}

func (v *ANTLRSemVisitor) VisitTypePath(tp *ast.TypePath) interface{} {
	segments := make([]string, len(*tp))
	for i, segment := range *tp {
		segments[i] = segment.Accept(v).(string)
	}

	typeId := strings.Join(segments, "::")
	if v.seeDeclared(typeId) == nil {
		v.logf(ERROR, "unknown type %s", typeId)
		return nil
	}
	return v.seeDeclared(typeId)
}

func (v *ANTLRSemVisitor) VisitTypePathSegment(tps *ast.TypePathSegment) interface{} {
	return tps.ID
}

func (v *ANTLRSemVisitor) VisitTypePathFunction(tpf *ast.TypePathFunction) interface{} {
	return nil
}

func (v *ANTLRSemVisitor) VisitParenthesizedType(pt *ast.ParenthesizedType) interface{} {
	return pt.VarType.Accept(v)
}

func (v *ANTLRSemVisitor) VisitPointerType(pt *ast.PointerType) interface{} {
	return PointerAttr{
		Volatile: pt.IsMut,
		Type:     pt.VarType.Accept(v).(TypeDef),
	}
}

func (v *ANTLRSemVisitor) VisitReferenceType(rt *ast.ReferenceType) interface{} {
	return nil
}

func (v *ANTLRSemVisitor) VisitTupleType(tt *ast.TupleType) interface{} {
	return nil
}

func (v *ANTLRSemVisitor) VisitArrayType(at *ast.ArrayType) interface{} {
	return ArrayAttr{
		Type: at.VarType.Accept(v).(TypeDef),
	}
}

func (v *ANTLRSemVisitor) VisitSliceType(st *ast.SliceType) interface{} {
	return nil
}

func (v *ANTLRSemVisitor) VisitNeverType(nt *ast.NeverType) interface{} {
	return nil
}

func (v *ANTLRSemVisitor) VisitInferredType(it *ast.InferredType) interface{} {
	return nil
}

func (v ANTLRSemVisitor) DumpLogs() interface{} {
	return v.logs
}

func (v *ANTLRSemVisitor) log(tp MessageType, msg string) {
	v.logs = append(v.logs, Message{
		Type: tp,
		Desc: msg,
	})
}

func (v *ANTLRSemVisitor) logf(tp MessageType, format string, a ...interface{}) {
	v.log(tp, fmt.Sprintf(format, a...))
}

func (v *ANTLRSemVisitor) enterScope() {
	v.scopes = append(v.scopes, Scope{})
	v.log(INFO, fmt.Sprintf("Entering scope %d...", len(v.scopes)))
}

func (v *ANTLRSemVisitor) exitScope() {
	v.log(INFO, fmt.Sprintf("Leaving scope %d...", len(v.scopes)))
	v.scopes = v.scopes[:len(v.scopes)-1]
}

func (v ANTLRSemVisitor) hasInScope(varName string, index int) bool {
	if _, ok := v.scopes[index][varName]; ok {
		return true
	}
	return false
}

func (v *ANTLRSemVisitor) declare(attr IDAttr) {
	v.log(INFO, fmt.Sprintf("{Scope %d} Declaring var %s %s...", len(v.scopes), attr.Name, attr.TypeParam.String()))
	if v.hasInScope(attr.Name, len(v.scopes)-1) {
		v.log(ERROR, fmt.Sprintf("'%s': already defined in the same scope %d!", attr.Name, len(v.scopes)))
	} else {
		for i := len(v.scopes) - 2; i >= 0; i-- {
			if v.hasInScope(attr.Name, i) {
				v.log(WARN, fmt.Sprintf("'%s': redefined in scope %d (earlier definition in scope %d: %s)!", attr.Name, len(v.scopes), i+1, v.scopes[i][attr.Name].String()))
			}
		}
	}
	v.scopes[len(v.scopes)-1][attr.Name] = attr.TypeParam
}

func (v *ANTLRSemVisitor) seeDeclared(varName string) TypeDef {
	for i := len(v.scopes) - 1; i >= 0; i-- {
		if v.hasInScope(varName, i) {
			return v.scopes[i][varName]
		}
	}
	return nil
}

var allowedUnaryOpTypes = map[string][]string{
	"!": {"u8", "u16", "u32", "u64", "i8", "i16", "i32", "i64"},
	"-": {"i8", "i16", "i32", "i64"},
}

var integerOperands = [][2]string{{"i8", "i8"}, {"i16", "i16"}, {"i32", "i32"}, {"i64", "i64"}, {"u8", "u8"}, {"u16", "u16"}, {"u32", "u32"}, {"u64", "u64"}}

var allowedBinaryOpTypes = map[string][][2]string{
	"-":  integerOperands,
	"+":  integerOperands,
	"*":  integerOperands,
	"/":  integerOperands,
	"<":  integerOperands,
	">":  integerOperands,
	"<=": integerOperands,
	">=": integerOperands,
	"==": integerOperands,
}

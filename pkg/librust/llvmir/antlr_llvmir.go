package llvmir

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Compiler2022/compilers-1-Belstowe/pkg/librust/ast"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type LLVMTopContext struct {
	*ir.Module
	fn map[string]*ir.Func
}

type LLVMContext struct {
	*ir.Block
	*LLVMTopContext
	parent     *LLVMContext
	vars       map[string]value.Value
	leaveBlock *ir.Block
	backBlock  *ir.Block
}

func NewLLVMTopContext() *LLVMTopContext {
	tc := &LLVMTopContext{
		Module: ir.NewModule(),
		fn:     make(map[string]*ir.Func),
	}
	printf := tc.NewFunc("printf", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	printf.Sig.Variadic = true
	writeln_i64 := tc.NewFunction("ruster::writeln_i64", types.Void, ir.NewParam("x", types.I64))
	writeln_i64.NewCall(printf, constant.NewCharArrayFromString("%ld"), writeln_i64.NewLoad(types.I64, writeln_i64.Parent.Params[0]))
	return tc
}

func (tc *LLVMTopContext) NewFunction(name string, retType types.Type, params ...*ir.Param) *LLVMContext {
	tc.fn[name] = ir.NewFunc(name, retType, params...)
	b := tc.fn[name].NewBlock("")
	return &LLVMContext{
		Block:          b,
		LLVMTopContext: tc,
		parent:         nil,
		vars:           make(map[string]value.Value),
	}
}

func NewLLVMContext() *LLVMContext {
	return &LLVMContext{
		Block:          nil,
		LLVMTopContext: NewLLVMTopContext(),
		parent:         nil,
		vars:           make(map[string]value.Value),
	}
}

func (c *LLVMContext) NewLLVMContext(b *ir.Block) *LLVMContext {
	return &LLVMContext{
		Block:          b,
		LLVMTopContext: c.LLVMTopContext,
		parent:         c,
		vars:           make(map[string]value.Value),
	}
}

func (c LLVMContext) lookupVariable(name string) value.Value {
	if v, ok := c.vars[name]; ok {
		return v
	} else if c.parent != nil {
		return c.parent.lookupVariable(name)
	} else {
		panic(fmt.Sprintf("no such variable: %s", name))
	}
}

func (c LLVMContext) callFunction(name string, args ...value.Value) value.Value {
	if _, ok := c.fn[name]; !ok {
		panic(fmt.Sprintf("no such function: %s", name))
	}
	return c.NewCall(c.fn[name], args...)
}

func (c *LLVMContext) Visit(tree ast.Node) interface{} {
	return tree.Accept(c)
}

func (c *LLVMContext) DumpLogs() interface{} {
	return nil
}

func (ctx *LLVMContext) VisitCrate(c *ast.Crate) interface{} {
	for _, item := range c.Items {
		ctx.Visit(item)
	}
	return ctx.String()
}

func (ctx *LLVMContext) VisitBlockExpression(be *ast.BlockExpression) interface{} {
	for _, stmt := range be.Statements {
		ctx.Visit(stmt)
	}
	if be.Expr != nil {
		return ctx.NewRet(ctx.Visit(be.Expr).(value.Value))
	}
	return nil
}

func (ctx *LLVMContext) VisitSimplePath(sp *ast.SimplePath) interface{} {
	return sp
}

func (ctx *LLVMContext) VisitParameter(p *ast.Parameter) interface{} {
	return ir.NewParam(p.ID, ctx.Visit(p.VarType).(types.Type))
}

func (c *LLVMContext) VisitFunction(f *ast.Function) interface{} {
	params := make([]*ir.Param, 0)
	for _, param := range f.Params {
		params = append(params, c.Visit(param).(*ir.Param))
	}
	newc := c.NewFunction(f.ID, c.Visit(f.ReturnType).(types.Type), params...)
	newc.Visit(f.Body)
	return c.fn[f.ID]
}

func (c *LLVMContext) VisitLetStatement(ls *ast.LetStatement) interface{} {
	v := c.NewAlloca(c.Visit(ls.VarType).(types.Type))
	if ls.Expr != nil {
		c.NewStore(c.Visit(ls.Expr).(value.Value), v)
	}
	c.vars[c.Visit(ls.Ptrn).(LLVMPatternStruct).ID] = v
	return nil
}

func (c *LLVMContext) VisitLiteralExpression(le *ast.LiteralExpression) interface{} {
	switch le.Tp {
	case ast.String:
		return constant.NewCharArrayFromString(le.Val)
	case ast.Boolean:
		return constant.NewBool(func() bool {
			if strings.ToLower(le.Val) == "true" || le.Val == "1" {
				return true
			}
			return false
		}())
	case ast.Char:
		return constant.NewInt(types.I8, int64([]byte(le.Val)[0]))
	case ast.Integer:
		if intVal, ok := constant.NewIntFromString(types.I64, le.Val); ok == nil {
			return intVal
		}
		panic(fmt.Sprintf("couldn't convert literal '%s' into integer!", le.Val))
	}
	return nil
}

func (c *LLVMContext) VisitPathExpression(pe *ast.PathExpression) interface{} {
	return c.lookupVariable(strings.Join(pe.Segments, "::"))
}

func (c *LLVMContext) VisitIfExpression(ie *ast.IfExpression) interface{} {
	thenCtx := c.NewLLVMContext(c.Parent.NewBlock("if.then"))
	thenCtx.Visit(ie.IfTrue)
	switch elseExpr := ie.IfFalse.(type) {
	case ast.BlockExpression:
	case ast.IfExpression:
		elseCtx := c.NewLLVMContext(c.Parent.NewBlock("if.else"))
		elseCtx.Visit(elseExpr)
		c.NewCondBr(c.Visit(ie.Expr).(value.Value), thenCtx.Block, elseCtx.Block)
		thenCtx.NewBr(c.Parent.NewBlock("leave.if"))
	default:
		elseContinueBlock := c.Parent.NewBlock("leave.if")
		c.NewCondBr(c.Visit(ie.Expr).(value.Value), thenCtx.Block, elseContinueBlock)
	}
	return nil
}

func (c *LLVMContext) VisitInfiniteLoopExpression(ile *ast.InfiniteLoopExpression) interface{} {
	backBlock := c.Parent.NewBlock("inf.loop.body")
	loopCtx := c.NewLLVMContext(backBlock)
	leaveBlock := c.Parent.NewBlock("leave.inf.loop")
	loopCtx.leaveBlock = leaveBlock
	loopCtx.backBlock = backBlock
	loopCtx.Visit(ile.Body)
	loopCtx.NewBr(backBlock)
	return nil
}

func (c *LLVMContext) VisitPredicateLoopExpression(ple *ast.PredicateLoopExpression) interface{} {
	backBlock := c.Parent.NewBlock("while.loop.cond")
	condCtx := c.NewLLVMContext(backBlock)
	c.NewBr(condCtx.Block)
	loopCtx := c.NewLLVMContext(c.Parent.NewBlock("while.loop.body"))
	leaveBlock := c.Parent.NewBlock("leave.while.loop")
	condCtx.NewCondBr(condCtx.Visit(ple.Expr).(value.Value), loopCtx.Block, leaveBlock)
	condCtx.leaveBlock = leaveBlock
	condCtx.backBlock = backBlock
	loopCtx.leaveBlock = leaveBlock
	loopCtx.backBlock = backBlock
	loopCtx.Visit(ple.Body)
	loopCtx.NewBr(condCtx.Block)
	return nil
}

func (c *LLVMContext) VisitIteratorLoopExpression(ile *ast.IteratorLoopExpression) interface{} {
	return nil
}

func (c *LLVMContext) VisitUnaryOperator(uo *ast.UnaryOperator) interface{} {
	val := c.Visit(uo.Val).(value.Value)
	switch uo.Op {
	case "!":
		val = c.NewXor(val, constant.NewInt(types.I1, 1))
	case "-":
		val = c.NewSub(constant.NewInt(types.I64, 0), val)
	}
	return val
}

func (c *LLVMContext) VisitBinaryOperator(bo *ast.BinaryOperator) interface{} {
	lhs := c.Visit(bo.LHS).(value.Value)
	rhs := c.Visit(bo.RHS).(value.Value)
	if types.IsFloat(rhs.Type()) {
		switch bo.Op {
		case "-":
			return c.NewFSub(lhs, rhs)
		case "+":
			return c.NewFAdd(lhs, rhs)
		case "/":
			return c.NewFDiv(lhs, rhs)
		case "*":
			return c.NewFMul(lhs, rhs)
		case "<":
			return c.NewFCmp(enum.FPredOLT, lhs, rhs)
		case ">":
			return c.NewFCmp(enum.FPredOGT, lhs, rhs)
		case "<=":
			return c.NewFCmp(enum.FPredOLE, lhs, rhs)
		case ">=":
			return c.NewFCmp(enum.FPredOGE, lhs, rhs)
		case "==":
			return c.NewFCmp(enum.FPredOEQ, lhs, rhs)
		case "!=":
			return c.NewFCmp(enum.FPredONE, lhs, rhs)
		}
	} else if types.IsInt(rhs.Type()) {
		switch bo.Op {
		case "-":
			return c.NewSub(lhs, rhs)
		case "+":
			return c.NewAdd(lhs, rhs)
		case "/":
			return c.NewSDiv(lhs, rhs)
		case "*":
			return c.NewMul(lhs, rhs)
		case "<":
			return c.NewICmp(enum.IPredSLT, lhs, rhs)
		case ">":
			return c.NewICmp(enum.IPredSGT, lhs, rhs)
		case "<=":
			return c.NewICmp(enum.IPredSLE, lhs, rhs)
		case ">=":
			return c.NewICmp(enum.IPredSGE, lhs, rhs)
		case "==":
			return c.NewICmp(enum.IPredEQ, lhs, rhs)
		case "!=":
			return c.NewICmp(enum.IPredNE, lhs, rhs)
		}
	}
	return lhs
}

func (c *LLVMContext) VisitReturnExpression(e *ast.ReturnExpression) interface{} {
	return c.NewRet(c.Visit(e.Expr).(value.Value))
}

func (c *LLVMContext) VisitBreakExpression(be *ast.BreakExpression) interface{} {
	return c.NewBr(c.leaveBlock)
}

func (c *LLVMContext) VisitContinueExpression(ce *ast.ContinueExpression) interface{} {
	return c.NewBr(c.backBlock)
}

func (c *LLVMContext) VisitTypeCastExpression(tce *ast.TypeCastExpression) interface{} {
	return nil
}

func (c *LLVMContext) VisitCallExpression(ce *ast.CallExpression) interface{} {
	fnName := c.Visit(ce.FnHeader).(value.Value)
	paramValues := make([]value.Value, 0)
	for _, param := range ce.Params {
		paramValues = append(paramValues, c.Visit(param).(value.Value))
	}
	return c.callFunction(fnName.Ident(), paramValues...)
}

func (c *LLVMContext) VisitBorrowExpression(be *ast.BorrowExpression) interface{} {
	return nil
}

func (c *LLVMContext) VisitArrayIndexExpression(aie *ast.ArrayIndexExpression) interface{} {
	indexStr := c.Visit(aie.Index).(value.Value).Ident()
	if index, err := strconv.Atoi(indexStr); err == nil {
		obj := c.Visit(aie.Index).(value.Value)
		return c.NewExtractValue(obj, uint64(index))
	}
	panic(fmt.Sprintf("couldn't convert index '%s' into integer!", indexStr))
}

func (c *LLVMContext) VisitLiteralPattern(lp *ast.LiteralPattern) interface{} {
	return LLVMPatternStruct{
		ID:    lp.Val,
		IsMut: false,
	}
}

func (c *LLVMContext) VisitIdentifierPattern(ip *ast.IdentifierPattern) interface{} {
	return LLVMPatternStruct{
		ID:    ip.ID,
		IsMut: ip.IsMut,
	}
}

func (c *LLVMContext) VisitTypePath(tp *ast.TypePath) interface{} {
	nodes := make([]string, 0)
	for _, segment := range *tp {
		nodes = append(nodes, c.Visit(segment).(string))
	}
	finalType := strings.Join(nodes, "::")
	switch finalType {
	case "i8":
	case "u8":
	case "char":
		return types.I8
	case "i16":
	case "u16":
		return types.I16
	case "i32":
	case "u32":
		return types.I32
	case "i64":
	case "u64":
	case "isize":
	case "usize":
		return types.I64
	case "i128":
	case "u128":
		return types.I128
	case "str":
		return types.NewArray(64, types.I8)
	}
	return types.Void
}

func (c *LLVMContext) VisitTypePathSegment(tps *ast.TypePathSegment) interface{} {
	return tps.ID
}

func (c *LLVMContext) VisitParenthesizedType(pt *ast.ParenthesizedType) interface{} {
	return c.Visit(pt.VarType)
}

func (c *LLVMContext) VisitPointerType(pt *ast.PointerType) interface{} {
	return nil
}

func (c *LLVMContext) VisitArrayType(at *ast.ArrayType) interface{} {
	lengthString := c.Visit(at.Expr).(value.Value).Ident()
	if length, ok := strconv.Atoi(lengthString); ok != nil {
		return types.NewArray(uint64(length), c.Visit(at.Expr).(types.Type))
	}
	panic(fmt.Sprintf("couldn't convert '%s' into integer!", lengthString))
}

func (c *LLVMContext) VisitMatchExpression(me *ast.MatchExpression) interface{} {
	return nil
}

func (c *LLVMContext) VisitMatchArm(ma *ast.MatchArm) interface{} {
	return nil
}

func (c *LLVMContext) VisitRHSRangeOperator(rro *ast.RHSRangeOperator) interface{} {
	return nil
}

func (c *LLVMContext) VisitRangeOperator(ro *ast.RangeOperator) interface{} {
	return nil
}

func (c *LLVMContext) VisitMethodCallExpression(mce *ast.MethodCallExpression) interface{} {
	return nil
}

func (c *LLVMContext) VisitTupleIndexExpression(tie *ast.TupleIndexExpression) interface{} {
	return nil
}

func (c *LLVMContext) VisitReferencePattern(rp *ast.ReferencePattern) interface{} {
	return nil
}

func (c *LLVMContext) VisitPathPattern(pp *ast.PathPattern) interface{} {
	return nil
}

func (c *LLVMContext) VisitReferenceType(rt *ast.ReferenceType) interface{} {
	return nil
}

func (c *LLVMContext) VisitTupleType(tt *ast.TupleType) interface{} {
	return nil
}

func (c *LLVMContext) VisitSliceType(st *ast.SliceType) interface{} {
	return nil
}

func (c *LLVMContext) VisitNeverType(nt *ast.NeverType) interface{} {
	return nil
}

func (c *LLVMContext) VisitInferredType(it *ast.InferredType) interface{} {
	return nil
}

func (c *LLVMContext) VisitTypePathFunction(tpf *ast.TypePathFunction) interface{} {
	return nil
}

func (c *LLVMContext) VisitUseDecl(ud *ast.UseDecl) interface{} {
	return nil
}

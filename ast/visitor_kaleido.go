package ast

import (
	"fmt"
	"log"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type VisitorKaleido struct {
	Module       *ir.Module
	namedValues  map[string]interface{}
	currentBlock *ir.Block
}

func NewVisitorKaleido() VisitorKaleido {
	return VisitorKaleido{Module: ir.NewModule()}
}

func (v *VisitorKaleido) VisitNumberExprAST(node *NumberExprAST) interface{} {
	log.Println("VisitNumberExprAST")
	res, err := constant.NewFloatFromString(&types.FloatType{TypeName: "", Kind: types.FloatKindDouble}, string(*node))
	if err != nil {
		panic(err)
	}
	return res
}

func (v *VisitorKaleido) VisitBinaryExprAST(node *BinaryExprAST) interface{} {
	log.Println("VisitBinaryExprAST")
	lhsValue := node.LHS.Accept(v).(value.Value)
	rhsValue := node.RHS.Accept(v).(value.Value)
	switch node.Op {
	case '+':
		return v.currentBlock.NewFAdd(lhsValue, rhsValue)
	case '-':
		return v.currentBlock.NewSub(lhsValue, rhsValue)
	case '*':
		return v.currentBlock.NewFMul(lhsValue, rhsValue)
	case '<':
		res := v.currentBlock.NewFCmp(enum.FPredULT, lhsValue, rhsValue)
		return v.currentBlock.NewUIToFP(res, types.Double)
	}
	panic(fmt.Sprintf("Unknown operator: %v", node.Op))
}

func (v *VisitorKaleido) VisitVariableExprAST(node *VariableExprAST) interface{} {
	log.Println("VisitVariableExprAST")
	if res, found := v.namedValues[string(*node)]; found {
		return res
	}
	panic(fmt.Sprintf("Variable %v not found", string(*node)))
}

func (v *VisitorKaleido) VisitCallExprAST(node *CallExprAST) interface{} {
	log.Println("VisitCallExprAST")
	funcRef, found := searchLLVMFunc(node.FunctionName, v.Module.Funcs)
	if !found {
		panic("Function " + node.FunctionName + " does not exist")
	}
	if len(funcRef.Params) != len(node.Args) {
		panic("Function " + node.FunctionName + ": incorrect number of arguments")
	}
	llvmArgs := make([]value.Value, 0, len(node.Args))
	for _, arg := range node.Args {
		evaluatedArg := arg.Accept(v).(value.Value)
		llvmArgs = append(llvmArgs, evaluatedArg)
	}
	return v.currentBlock.NewCall(funcRef, llvmArgs...)
}

func searchLLVMFunc(funcName string, funcs []*ir.Func) (result *ir.Func, found bool) {
	for _, candidate := range funcs {
		if funcName == candidate.GlobalName {
			return candidate, true
		}
	}
	return nil, false
}

func (v *VisitorKaleido) VisitPrototypeAST(node *PrototypeAST) interface{} {
	log.Println("VisitPrototypeAST")
	params := make([]*ir.Param, 0, len(node.Args))
	for _, name := range node.Args {
		params = append(params, ir.NewParam(name, types.Double))
	}
	llvmFunc := v.Module.NewFunc(node.FunctionName, types.Double, params...)
	llvmFunc.Linkage = enum.LinkageExternal
	return llvmFunc
}

func (v *VisitorKaleido) VisitFunctionAST(node *FunctionAST) interface{} {
	log.Println("VisitFunctionAST")
	llvmFunc, found := searchLLVMFunc(node.Prototype.FunctionName, v.Module.Funcs)
	if !found {
		llvmFunc = node.Prototype.Accept(v).(*ir.Func)
	}
	if llvmFunc == nil {
		return nil
	}
	if len(llvmFunc.Blocks) != 0 {
		panic("Function " + node.Prototype.FunctionName + " cannot be redefined")
	}

	v.namedValues = make(map[string]interface{})
	for _, param := range llvmFunc.Params {
		v.namedValues[param.LocalName] = param
	}
	v.currentBlock = llvmFunc.NewBlock("entry")
	bodyValue := node.Body.Accept(v).(value.Value)
	v.currentBlock.NewRet(bodyValue)
	return llvmFunc
}

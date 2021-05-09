package parser

import (
	"fmt"
	"log"

	"github.com/llvm/llvm-project/llvm/bindings/go/llvm"
)

type VisitorKaleido struct {
	context     *llvm.Context
	Module      *llvm.Module
	builder     *llvm.Builder
	namedValues map[string]interface{}
}

func NewVisitorKaleido() VisitorKaleido {
	context := llvm.NewContext()
	module := llvm.NewModule("")
	builder := context.NewBuilder()
	return VisitorKaleido{context: &context, Module: &module, builder: &builder}
}

func (v *VisitorKaleido) VisitNumberExprAST(node *NumberExprAST) interface{} {
	log.Println("VisitNumberExprAST")
	value := llvm.ConstFloatFromString(llvm.DoubleType(), string(*node))
	return value
}

func (v *VisitorKaleido) VisitBinaryExprAST(node *BinaryExprAST) interface{} {
	log.Println("VisitBinaryExprAST")
	lhsValue := node.LHS.Accept(v).(llvm.Value)
	rhsValue := node.RHS.Accept(v).(llvm.Value)
	switch node.Op {
	case '+':
		return v.builder.CreateFAdd(lhsValue, rhsValue, "addtmp")
	case '-':
		return v.builder.CreateFSub(lhsValue, rhsValue, "subtmp")
	case '*':
		return v.builder.CreateFMul(lhsValue, rhsValue, "multmp")
	case '<':
		res := v.builder.CreateFCmp(llvm.FloatULT, lhsValue, rhsValue, "cmptmp")
		return v.builder.CreateUIToFP(res, llvm.DoubleType(), "booltmp")
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
	funcRef := v.Module.NamedFunction(node.FunctionName)
	if funcRef.IsNil() {
		panic("Function " + node.FunctionName + " does not exist")
	}
	if funcRef.ParamsCount() != len(node.Args) {
		panic("Function " + node.FunctionName + ": incorrect number of arguments")
	}
	llvmArgs := make([]llvm.Value, 0, len(node.Args))
	for _, arg := range node.Args {
		evaluatedArg := arg.Accept(v).(llvm.Value)
		llvmArgs = append(llvmArgs, evaluatedArg)
	}
	return v.builder.CreateCall(funcRef, llvmArgs, "calltmp")
}

func (v *VisitorKaleido) VisitPrototypeAST(node *PrototypeAST) interface{} {
	log.Println("VisitPrototypeAST")
	paramTypes := make([]llvm.Type, 0, len(node.Args))
	for range node.Args {
		paramTypes = append(paramTypes, llvm.DoubleType())
	}
	functionType := llvm.FunctionType(llvm.DoubleType(), paramTypes, false)
	llvmFunc := llvm.AddFunction(*v.Module, node.FunctionName, functionType)
	llvmFunc.SetLinkage(llvm.ExternalLinkage)
	for i, argName := range node.Args {
		llvmFunc.Params()[i].SetName(argName)
	}
	return llvmFunc
}

func (v *VisitorKaleido) VisitFunctionAST(node *FunctionAST) interface{} {
	log.Println("VisitFunctionAST")
	llvmFunc := v.Module.NamedFunction(node.Prototype.FunctionName)
	if llvmFunc.IsNil() {
		llvmFunc = node.Prototype.Accept(v).(llvm.Value)
	}
	if llvmFunc.IsNil() {
		panic("Function " + node.Prototype.FunctionName + " does not exist")
	}
	if llvmFunc.BasicBlocksCount() != 0 {
		panic("Function " + node.Prototype.FunctionName + " cannot be redefined")
	}

	v.namedValues = make(map[string]interface{})
	for _, param := range llvmFunc.Params() {
		v.namedValues[param.Name()] = param
	}
	basicBlock := v.context.AddBasicBlock(llvmFunc, "entry")
	v.builder.SetInsertPointAtEnd(basicBlock)
	bodyValue := node.Body.Accept(v).(llvm.Value)
	v.builder.CreateRet(bodyValue)
	//llvm.VerifyFunction()
	// missing from tuto :
	// Error reading body, remove function.
	// TheFunction->eraseFromParent();
	// return nullptr;
	return llvmFunc
}

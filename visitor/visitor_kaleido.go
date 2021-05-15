/*
MIT License

Copyright (c) 2021 Vincent Hiribarren

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package visitor

import (
	"errors"
	"fmt"
	"log"

	"github.com/llvm/llvm-project/llvm/bindings/go/llvm"
	"github.com/vhiribarren/tuto-llvm-kaleidoscope-golang/parser"
)

func initModuleAndPassManager() (*llvm.Module, *llvm.PassManager) {
	module := llvm.NewModule("main")
	passManager := llvm.NewFunctionPassManagerForModule(module)
	passManager.AddInstructionCombiningPass()
	passManager.AddReassociatePass()
	passManager.AddGVNPass()
	passManager.AddCFGSimplificationPass()
	passManager.InitializeFunc()
	return &module, &passManager
}

type VisitorKaleido struct {
	context     *llvm.Context
	module      *llvm.Module
	builder     *llvm.Builder
	passManager *llvm.PassManager
	namedValues map[string]interface{}
}

func NewVisitorKaleido() VisitorKaleido {
	context := llvm.NewContext()
	module, passManager := initModuleAndPassManager()
	builder := context.NewBuilder()
	return VisitorKaleido{context: &context, module: module, passManager: passManager, builder: &builder}
}

func (v *VisitorKaleido) FeedAST(node *parser.ProgramAST) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch recoveredError := r.(type) {
			case error:
				err = recoveredError
			case string:
				err = errors.New(recoveredError)
			default:
				err = errors.New("Panic occured, recovered")
			}
		}
	}()
	node.Accept(v)
	return nil
}

func (v *VisitorKaleido) GenerateModuleIR() (irCode string) {
	return v.module.String()
}

func (v *VisitorKaleido) EvalutateMain() (float64, error) {
	jit := NewKaleidoJIT(v.module)
	return jit.Run("__main__")
}

func (v *VisitorKaleido) VisitNumberExprAST(node *parser.NumberExprAST) interface{} {
	log.Println("VisitNumberExprAST")
	value := llvm.ConstFloatFromString(llvm.DoubleType(), string(*node))
	return value
}

func (v *VisitorKaleido) VisitBinaryExprAST(node *parser.BinaryExprAST) interface{} {
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

func (v *VisitorKaleido) VisitVariableExprAST(node *parser.VariableExprAST) interface{} {
	log.Println("VisitVariableExprAST")
	if res, found := v.namedValues[string(*node)]; found {
		return res
	}
	panic(fmt.Sprintf("Variable %v not found", string(*node)))
}

func (v *VisitorKaleido) VisitCallExprAST(node *parser.CallExprAST) interface{} {
	log.Println("VisitCallExprAST")
	funcRef := v.module.NamedFunction(node.FunctionName)
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

func (v *VisitorKaleido) VisitPrototypeAST(node *parser.PrototypeAST) interface{} {
	log.Println("VisitPrototypeAST")
	paramTypes := make([]llvm.Type, 0, len(node.Args))
	for range node.Args {
		paramTypes = append(paramTypes, llvm.DoubleType())
	}
	functionType := llvm.FunctionType(llvm.DoubleType(), paramTypes, false)
	llvmFunc := llvm.AddFunction(*v.module, node.FunctionName, functionType)
	llvmFunc.SetLinkage(llvm.ExternalLinkage)
	for i, argName := range node.Args {
		llvmFunc.Params()[i].SetName(argName)
	}
	return llvmFunc
}

func (v *VisitorKaleido) VisitFunctionAST(node *parser.FunctionAST) interface{} {
	log.Println("VisitFunctionAST")
	llvmFunc := v.module.NamedFunction(node.Prototype.FunctionName)
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
	if bodyValue.IsNil() {
		// Error reading body, remove function.
		llvmFunc.EraseFromParentAsFunction()
		panic("Error reading body")
	}
	v.builder.CreateRet(bodyValue)
	if err := llvm.VerifyFunction(llvmFunc, llvm.PrintMessageAction); err != nil {
		llvmFunc.EraseFromParentAsFunction()
		panic(err)
	}
	v.passManager.RunFunc(llvmFunc)
	return llvmFunc
}

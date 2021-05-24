package visitor

import (
	"errors"

	"github.com/llvm/llvm-project/llvm/bindings/go/llvm"
)

type KaleidoscopeJIT struct {
	executionEngine llvm.ExecutionEngine
}

func init() {
	llvm.InitializeNativeTarget()
	llvm.InitializeNativeAsmPrinter()
}

func NewKaleidoJIT(module *llvm.Module) KaleidoscopeJIT {
	compilerOptions := llvm.NewMCJITCompilerOptions()
	executionEngine, err := llvm.NewMCJITCompiler(*module, compilerOptions)
	if err != nil {
		panic(err)
	}
	return KaleidoscopeJIT{executionEngine}
}

func (j *KaleidoscopeJIT) AddModule(module llvm.Module) {
	j.executionEngine.AddModule(module)
}

func (j *KaleidoscopeJIT) Run(name string, args ...float64) (float64, error) {
	f := j.executionEngine.FindFunction(name)
	if f.IsNil() {
		return 0, errors.New("Function " + name + " does not exist")
	}
	if f.ParamsCount() != len(args) {
		return 0, errors.New("Bad number of arguments")
	}
	genericValues := make([]llvm.GenericValue, f.ParamsCount())
	for _, arg := range args {
		genericValues = append(genericValues, llvm.NewGenericValueFromFloat(llvm.DoubleType(), arg))
	}
	result := j.executionEngine.RunFunction(f, genericValues)
	return result.Float(llvm.DoubleType()), nil
}

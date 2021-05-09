package parser

type Visitor interface {
	VisitNumberExprAST(*NumberExprAST) interface{}
	VisitBinaryExprAST(*BinaryExprAST) interface{}
	VisitVariableExprAST(*VariableExprAST) interface{}
	VisitCallExprAST(*CallExprAST) interface{}
	VisitPrototypeAST(*PrototypeAST) interface{}
	VisitFunctionAST(*FunctionAST) interface{}
}

type Visitable interface {
	Accept(Visitor) interface{}
}

type ProgramAST struct {
	Funcs  []FunctionAST
	Protos []PrototypeAST
}

func (p *ProgramAST) Accept(visitor Visitor) interface{} {
	for _, e := range p.Protos {
		e.Accept(visitor)
	}
	for _, e := range p.Funcs {
		e.Accept(visitor)
	}
	return nil
}

type ExprAST interface {
	Visitable
}

type NumberExprAST string

func (n NumberExprAST) Accept(visitor Visitor) interface{} {
	return visitor.VisitNumberExprAST(&n)
}

type BinaryExprAST struct {
	LHS ExprAST
	RHS ExprAST
	Op  rune
}

func (b *BinaryExprAST) Accept(visitor Visitor) interface{} {
	return visitor.VisitBinaryExprAST(b)
}

type VariableExprAST string

func (v VariableExprAST) Accept(visitor Visitor) interface{} {
	return visitor.VisitVariableExprAST(&v)
}

type CallExprAST struct {
	FunctionName string
	Args         []ExprAST
}

func (c *CallExprAST) Accept(visitor Visitor) interface{} {
	return visitor.VisitCallExprAST(c)
}

type PrototypeAST struct {
	FunctionName string
	Args         []string
}

func (p *PrototypeAST) Accept(visitor Visitor) interface{} {
	return visitor.VisitPrototypeAST(p)
}

type FunctionAST struct {
	Prototype PrototypeAST
	Body      ExprAST
}

func (f *FunctionAST) Accept(visitor Visitor) interface{} {
	return visitor.VisitFunctionAST(f)
}

type ArgList []string

type ExprList []ExprAST

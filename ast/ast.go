package ast

type ExprAST interface {
}

type NumberExprAST int

type BinaryExprAST struct {
	LHS ExprAST
	RHS ExprAST
	Op  rune
}

type VariableExprAST string

type CallExprAST struct {
	FunctionName string
	Args         []ExprAST
}

type PrototypeAST struct {
	FunctionName string
	Args         []string
}

type FunctionAST struct {
	Prototype PrototypeAST
	Body      ExprAST
}

type ArgList []string

type ExprList []ExprAST

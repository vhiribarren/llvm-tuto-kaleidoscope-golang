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

%{
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

package yacc

import(
    "log"
    "errors"
    "unicode/utf8"
    "github.com/vhiribarren/tuto-llvm-kaleidoscope-golang/parser"
    "github.com/vhiribarren/tuto-llvm-kaleidoscope-golang/lexer"
)

%}

%union{
    token lexer.KaleidoTokenContext
	proto parser.PrototypeAST
    function parser.FunctionAST
	expr  parser.ExprAST
    argList parser.ArgList
    exprList parser.ExprList
    program parser.ProgramAST
    number parser.NumberExprAST
    variable parser.VariableExprAST
}

%token DEF
%token EXTERN
%token<token> NUMBER

%left '<'
%left '+' '-'
%left '*'

%left<token> IDENTIFIER
%left '('

%type<expr> Expr FuncExpr
%type<argList> ProtoArgList
%type<exprList> ExprList ExprListContinuation
%type<proto> Prototype Ext
%type<function> Def TopLevelExpr
%type<program> TopLevel Program


%%

Program : TopLevel
    {
        $$ = $1
        yylex.(*parserContext).result = &$$
    }

TopLevel: TopLevel Def Delimiter
    {
        $1.Funcs = append($1.Funcs, $2)
        $$ = $1
    };
TopLevel: TopLevel Ext Delimiter
    {
        $1.Protos = append($1.Protos, $2)
        $$ = $1
    };
TopLevel: TopLevel TopLevelExpr Delimiter
    {
        $1.Funcs = append($1.Funcs, $2)
        $$ = $1
    };
TopLevel: /* Empty */ 
    {
        $$ = parser.ProgramAST{}
    };
Delimiter: ';' ;
Delimiter: /* Empty */ ;

Def: DEF Prototype Expr
    {
        $$ = parser.FunctionAST{Prototype: $2, Body: $3}
    };
Ext: EXTERN Prototype ';'
    {
        $$ = $2
    };
TopLevelExpr: Expr
    {
        $$ = parser.FunctionAST{Prototype: parser.PrototypeAST{FunctionName: "__main__", Args: []string{}},Body: $1}
    };

Expr: IDENTIFIER
    { $$ = parser.VariableExprAST($1.Value) };
Expr: NUMBER
    { $$ = parser.NumberExprAST($1.Value) };
Expr: FuncExpr ;
Expr: '(' Expr ')'
    { $$ = $2 };
Expr:  Expr '+' Expr
    { $$ = &parser.BinaryExprAST{LHS: $1, RHS: $3, Op: '+'} };
Expr:  Expr '-' Expr
    { $$ = &parser.BinaryExprAST{LHS: $1, RHS: $3, Op: '-'} };
Expr:  Expr '<' Expr
    { $$ = &parser.BinaryExprAST{LHS: $1, RHS: $3, Op: '<'} };
Expr:  Expr '*' Expr
    { $$ = &parser.BinaryExprAST{LHS: $1, RHS: $3, Op: '*'} };

FuncExpr: IDENTIFIER '(' ExprList ')'
    {
        log.Println("Parsed rule: FuncExpr")
        $$ = &parser.CallExprAST{FunctionName: $1.Value, Args: $3}
    };
ExprList: ExprListContinuation ;
ExprList:  /* Empty */
        { $$ = []parser.ExprAST {} };
ExprListContinuation: ExprListContinuation ',' Expr
    { $$ = append($1, $3) };
ExprListContinuation: Expr
    { $$ = []parser.ExprAST { $1 } };

Prototype: IDENTIFIER '(' ProtoArgList ')'
    {
        $$ = parser.PrototypeAST{FunctionName: $1.Value, Args: $3}
    };
ProtoArgList: ProtoArgList IDENTIFIER
    { $$ = append($1, $2.Value) };
ProtoArgList:  /* Empty */
    { $$ = []string {  } } ; 

%%

const EOF = 0

type parserContext struct {
    lexer.KaleidoLexer
    result * parser.ProgramAST
    err error
}

func (s *parserContext) Lex(lval *yySymType) int {
    tokenContext := s.NextToken()
    lval.token = *tokenContext
    switch tokenContext.Token {
    case lexer.KTokenEOF:
        return EOF
    case lexer.KTokenDef:
        return DEF
    case lexer.KTokenExtern:
        return EXTERN
    case lexer.KTokenIdentifier:
        return IDENTIFIER
    case lexer.KTokenNumber:
        return NUMBER
	default:
		val, _ := utf8.DecodeRuneInString(tokenContext.Value)
		return int(val)
    }
}

func (s *parserContext) Error(e string) {
    s.result = nil
    s.err = errors.New(e)
}

func BuildKaleidoAST(buffer string) (*parser.ProgramAST, error) {
    context := &parserContext{KaleidoLexer: lexer.NewKaleidoLexer(buffer)}
    yyParse(context)
    if context.result == nil {
        return nil, context.err
    }
    return context.result, nil
}

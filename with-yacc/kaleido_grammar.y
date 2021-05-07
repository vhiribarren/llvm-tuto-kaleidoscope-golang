%{
package main

import "log"
import "unicode/utf8"
import "alea.net/xp/llvm/kaleidoscope/ast"
import "alea.net/xp/llvm/kaleidoscope/lexer"

%}

%union{
    token lexer.KaleidoTokenContext
	proto ast.PrototypeAST
    function ast.FunctionAST
	expr  ast.ExprAST
    argList ast.ArgList
    exprList ast.ExprList
    program ast.ProgramAST
    number ast.NumberExprAST
    variable ast.VariableExprAST
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
%type<argList> ProtoArgList ProtoArgListContinuation
%type<exprList> ExprList ExprListContinuation
%type<proto> Prototype Ext
%type<function> Def TopLevelExpr
%type<program> TopLevel Program


%%

Program : TopLevel
    {
        $$ = $1
        yylex.(*ParserContext).result = &$$
    }

TopLevel: TopLevel Def
    {
        $1.Funcs = append($1.Funcs, $2)
        $$ = $1
    };
TopLevel: TopLevel Ext
    {
        $1.Protos = append($1.Protos, $2)
        $$ = $1
    };
TopLevel: TopLevel TopLevelExpr
    {
        $1.Funcs = append($1.Funcs, $2)
        $$ = $1
    };
TopLevel: ';' 
{
    $$ = ast.ProgramAST{}
};
TopLevel: /* Empty */ 
    {
        $$ = ast.ProgramAST{}
    };

Def: DEF Prototype Expr
    {
        $$ = ast.FunctionAST{Prototype: $2, Body: $3}
    };
Ext: EXTERN Prototype ';'
    {
        $$ = $2
    };
TopLevelExpr: Expr
    {
        $$ = ast.FunctionAST{Prototype: ast.PrototypeAST{FunctionName: "__main__", Args: []string{}},Body: $1}
    };

Expr: IDENTIFIER
    { $$ = ast.VariableExprAST($1.Value) };
Expr: NUMBER
    { $$ = ast.NumberExprAST($1.Value) };
Expr: FuncExpr ;
Expr: '(' Expr ')'
    { $$ = $2 };
Expr:  Expr '+' Expr
    { $$ = &ast.BinaryExprAST{LHS: $1, RHS: $3, Op: '+'} };
Expr:  Expr '-' Expr
    { $$ = &ast.BinaryExprAST{LHS: $1, RHS: $3, Op: '-'} };
Expr:  Expr '<' Expr
    { $$ = &ast.BinaryExprAST{LHS: $1, RHS: $3, Op: '<'} };
Expr:  Expr '*' Expr
    { $$ = &ast.BinaryExprAST{LHS: $1, RHS: $3, Op: '*'} };

FuncExpr: IDENTIFIER '(' ExprList ')'
    {
        log.Println("Parsed rule: FuncExpr")
        $$ = &ast.CallExprAST{FunctionName: $1.Value, Args: $3}
    };
ExprList: ExprListContinuation ;
ExprList:  /* Empty */
        { $$ = []ast.ExprAST {} };
ExprListContinuation: ExprListContinuation ',' Expr
    { $$ = append($1, $3) };
ExprListContinuation: Expr
    { $$ = []ast.ExprAST { $1 } };

Prototype: IDENTIFIER '(' ProtoArgList ')'
    {
        $$ = ast.PrototypeAST{FunctionName: $1.Value, Args: $3}
    };
ProtoArgList: ProtoArgListContinuation;
ProtoArgList:  /* Empty */
    { $$ = []string {  } } ; 
ProtoArgListContinuation: ProtoArgListContinuation ',' IDENTIFIER
    { $$ = append($1, $3.Value) };
ProtoArgListContinuation: IDENTIFIER
    { $$ = []string { $1.Value } };

%%

const EOF = 0

type ParserContext struct {
    lexer.KaleidoLexer
    result * ast.ProgramAST
}

func (s *ParserContext) Lex(lval *yySymType) int {
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

func (s *ParserContext) Error(e string) {
    panic(e)
}

func (s *ParserContext) Result() (*ast.ProgramAST) {
    return s.result
}

func Parse(buffer string) (* ParserContext){
    parserContext := &ParserContext{KaleidoLexer: lexer.NewKaleidoLexer(buffer)}
    yyParse(parserContext)
    return parserContext
}

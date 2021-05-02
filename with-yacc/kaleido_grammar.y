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



%%

TopLevel: TopLevel Def;
TopLevel: TopLevel Ext ;
TopLevel: TopLevel TopLevelExpr;
TopLevel: ';' | /* Empty */ ;

Def: DEF Prototype Expr
    {
        log.Println("Parsed rule: Def")
        $$ = ast.FunctionAST{Prototype: $2, Body: $3}
        log.Printf("%#v", $$)
    };
Ext: EXTERN Prototype ';'
    {
        log.Println("Parsed rule: Ext")
        $$ = $2
        log.Printf("%#v", $$)
    };
TopLevelExpr: Expr
    {
        log.Println("Parsed rule: TopLevelExpr")
        $$ = ast.FunctionAST{Prototype: ast.PrototypeAST{FunctionName: "__main__", Args: []string{}},Body: $1}
        log.Printf("%#v", $$)
    };

Expr: IDENTIFIER
    { $$ = $1.Value };
Expr: NUMBER
    { $$ = $1.Value };
Expr: FuncExpr ;
Expr: '(' Expr ')'
    { $$ = $2 };
Expr:  Expr '+' Expr
    { $$ = ast.BinaryExprAST{LHS: $1, RHS: $3, Op: '+'} };
Expr:  Expr '-' Expr
    { $$ = ast.BinaryExprAST{LHS: $1, RHS: $3, Op: '-'} };
Expr:  Expr '<' Expr
    { $$ = ast.BinaryExprAST{LHS: $1, RHS: $3, Op: '<'} };
Expr:  Expr '*' Expr
    { $$ = ast.BinaryExprAST{LHS: $1, RHS: $3, Op: '*'} };

FuncExpr: IDENTIFIER '(' ExprList ')'
    {
        log.Println("Parsed rule: FuncExpr")
        $$ = ast.CallExprAST{FunctionName: $1.Value, Args: $3}
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
        log.Println("Parsed rule: Prototype")
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

type Scanner struct {
    lexer.KaleidoLexer
}

func (s *Scanner) Lex(lval *yySymType) int {
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

func (s *Scanner) Error(e string) {
    panic(e)
}

func Parse(buffer string) {
    scanner := Scanner{KaleidoLexer: lexer.NewKaleidoLexer(buffer)}
    yyParse(&scanner)
}

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
/*
%type<expr> Expr FuncExpr
%type<argList> ProtoArgList ProtoArgListContinuation
%type<exprList> ExprList ExprListContinuation
%type<proto> Prototype Ext
%type<function> Def TopLevelExpr
*/


%%

TopLevel: TopLevel Def;
TopLevel: TopLevel Ext ;
TopLevel: TopLevel TopLevelExpr;
TopLevel: ';' | /* Empty */ ;

Def: DEF Prototype Expr { log.Println("Parsed rule: Def") };
Ext: EXTERN Prototype ';' { log.Println("Parsed rule: Ext") };

TopLevelExpr: Expr { log.Println("Parsed rule: TopLevelExpr") };

Expr: IDENTIFIER | NUMBER;
Expr: FuncExpr ;
Expr: '(' Expr ')'  ;
Expr:  Expr '+' Expr;
Expr:  Expr '-' Expr;
Expr:  Expr '<' Expr;
Expr:  Expr '*' Expr;

FuncExpr: IDENTIFIER '(' ExprList ')' { log.Println("Parsed rule: FuncExpr") };
ExprList: ExprListContinuation | /* Empty */ ;
ExprListContinuation: ExprListContinuation ',' Expr | Expr ;

Prototype: IDENTIFIER '(' ProtoArgList ')'  { log.Println("Parsed rule: Prototype") };
ProtoArgList: ProtoArgListContinuation | /* Empty */ ; 
ProtoArgListContinuation: ProtoArgListContinuation ',' IDENTIFIER | IDENTIFIER ;

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

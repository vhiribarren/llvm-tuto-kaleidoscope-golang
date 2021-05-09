%{
package yacc

import "log"
import "unicode/utf8"
import "alea.net/xp/llvm/kaleidoscope/parser"
import "alea.net/xp/llvm/kaleidoscope/lexer"

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
    $$ = parser.ProgramAST{}
};
TopLevel: /* Empty */ 
    {
        $$ = parser.ProgramAST{}
    };

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
    result * parser.ProgramAST
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

func (s *ParserContext) Result() (*parser.ProgramAST) {
    return s.result
}

func Parse(buffer string) (* ParserContext){
    parserContext := &ParserContext{KaleidoLexer: lexer.NewKaleidoLexer(buffer)}
    yyParse(parserContext)
    return parserContext
}

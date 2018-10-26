%{
package parser

import "github.com/BenLubar/bit/ast"
%}

%union {
        line *ast.Line
        stmt ast.Stmt
	expr ast.Expr
        bit2 [2]ast.Bits
        bits ast.Bits
        bit  ast.Bit
}

%start program

%type <line> line
%type <stmt> statement
%type <expr> expression expression0
%type <bit2> goto
%type <bits> number
%type <bit> bit

%left tNand
%token tZero tOne
%token tLineNumber tCode
%token tOpen tClose tParenthesis
%token tThe tValue tAt tBeyond tAddressOf
%token tJumpRegister tVariable
%token tGoto tIf tIs
%token tEquals tPrint tRead

%%

program: line { yylex.(*parser).line($1) }
       | program line { yylex.(*parser).line($2) };

line: tLineNumber number tCode statement goto
    { $$ = &ast.Line{Num: $2, Stmt: $4, Goto0: $5[0], Goto1: $5[1]} };

statement: tRead { $$ = &ast.Read{} }
         | tPrint bit { $$ = &ast.Print{Val: $2} }
         | expression tEquals expression
         {
             if !ast.CanAssign($1) {
                 yylex.Error("cannot assign to " + $1.String())
             }
             $$ = &ast.Equals{Left: $1, Right: $3}
         };

expression: expression0 tNand expression
	  {
	      if !ast.CanVal($1) {
	          yylex.Error("cannot take value of " + $1.String())
              }
              if !ast.CanVal($3) {
                  yylex.Error("cannot take value of " + $3.String())
              }
              $$ = &ast.Nand{Left: $1, Right: $3}
          }
          | expression0 { $$ = $1 };

expression0: tOpen tParenthesis expression tClose tParenthesis { $$ = $3 }
           | tThe tValue tAt expression0
           {
               if !ast.CanDeref($4) {
                   yylex.Error("cannot dereference " + $4.String())
               }
               $$ = &ast.ValueAt{Ptr: $4}
           }
           | tThe tValue tBeyond expression0
           {
               if !ast.CanDeref($4) {
                   yylex.Error("cannot dereference " + $4.String())
               }
               $$ = &ast.ValueBeyond{Ptr: $4}
           }
           | tThe tAddressOf expression0
           {
               if !ast.CanAddr($3) {
                   yylex.Error("cannot take address of " + $3.String())
               }
               $$ = &ast.AddressOf{Val: $3}
           }
	   | tThe tJumpRegister
	   { $$ = &ast.JumpRegister{} }
           | tVariable number { $$ = &ast.Variable{Num: $2} }
           | bit { $$ = &ast.Constant{Val: $1} };

goto: /* empty */
    { $$ = [2]ast.Bits{nil, nil} }
    | tGoto number
    { $$ = [2]ast.Bits{$2, $2} }
    | tGoto number tIf tThe tJumpRegister tIs tZero
    { $$ = [2]ast.Bits{$2, nil} }
    | tGoto number tIf tThe tJumpRegister tIs tOne
    { $$ = [2]ast.Bits{nil, $2} }
    | tGoto number tIf tThe tJumpRegister tIs tZero tGoto number tIf tThe tJumpRegister tIs tOne
    { $$ = [2]ast.Bits{$2, $9} }
    | tGoto number tIf tThe tJumpRegister tIs tOne tGoto number tIf tThe tJumpRegister tIs tZero
    { $$ = [2]ast.Bits{$9, $2} };

number: bit { $$ = ast.Bits{$1} }
      | number bit { $$ = append($1, $2) };

bit: tZero { $$ = ast.Zero }
   | tOne { $$ = ast.One };

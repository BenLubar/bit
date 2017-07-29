%{
//go:generate go get golang.org/x/tools/cmd/goyacc
//go:generate goyacc syntax.y

package main

import "github.com/BenLubar/bit/bitnum"
%}

%union {
	program *Program
	lines   []*Line
	line    *Line
	bit     bool
	num     *bitnum.Number
	stmt    Stmt
	expr    Expr
}

%type<program> program
%type<lines> lines
%type<line> line goto
%type<bit> bit
%type<num> number
%type<stmt> statement
%type<expr> lvalue rvalue rvalue1

%%

program: lines
{
	$$ = &Program{
		Lines: $1,
	}
	yylex.(*lex).program = $$
}

lines: line
{
	$$ = []*Line{$1}
}

lines: lines line
{
	$$ = append($1, $2)
}

line: LINE NUMBER number CODE statement goto
{
	$$ = $6
	$$.Num = $3
	$$.Stmt = $5
}

bit: ZERO
{
	$$ = false
}

bit: ONE
{
	$$ = true
}

statement: READ
{
	$$ = &ReadStmt{}
}

statement: READ number
{
	if !yylex.(*lex).ext {
		yylex.Error("READ with a line number is an extension (run bitc with the -ext option)")
	}
	$$ = &ReadStmt{
		EOFLine: $2,
	}
}

statement: PRINT bit
{
	$$ = &PrintStmt{
		Bit: $2,
	}
}

goto: /* nothing */
{
	$$ = &Line{}
}

goto: GOTO number
{
	$$ = &Line{
		Zero: $2,
		One:  $2,
	}
}

goto: GOTO number IF THE JUMP REGISTER IS ZERO
{
	$$ = &Line{
		Zero: $2,
	}
}

goto: GOTO number IF THE JUMP REGISTER IS ONE
{
	$$ = &Line{
		One: $2,
	}
}

goto: GOTO number IF THE JUMP REGISTER IS ZERO GOTO number IF THE JUMP REGISTER IS ONE
{
	$$ = &Line{
		Zero: $2,
		One:  $10,
	}
}

goto: GOTO number IF THE JUMP REGISTER IS ONE GOTO number IF THE JUMP REGISTER IS ZERO
{
	$$ = &Line{
		Zero: $10,
		One:  $2,
	}
}

lvalue: VARIABLE number
{
	$$ = &UnknownVariable{
		Num: $2,
	}
}

lvalue: THE JUMP REGISTER
{
	$$ = &JumpRegister{}
}

rvalue1: lvalue
{
	$$ = $1
}

lvalue: THE VALUE AT rvalue1
{
	if !$4.Pointer() {
		yylex.Error("not a pointer: " + $4.String())
	}
	$$ = &ValueAt{
		Target: $4,
		Offset: 0,
	}
}

lvalue: THE VALUE BEYOND rvalue1
{
	if !$4.Pointer() {
		yylex.Error("not a pointer: " + $4.String())
	}
	$$ = &ValueAt{
		Target: $4,
		Offset: 1,
	}
}

rvalue1: THE ADDRESS OF rvalue1
{
	if !$4.Addressable() {
		yylex.Error("not addressable: " + $4.String())
	}
	$$ = &AddressOf{
		Variable: $4,
	}
}

rvalue: rvalue1 NAND rvalue1
{
	if !$1.Value() {
		yylex.Error("not a value: " + $1.String())
	}
	if !$3.Value() {
		yylex.Error("not a value: " + $3.String())
	}
	$$ = &Nand{
		Left:  $1,
		Right: $3,
	}
}

rvalue1: OPEN PARENTHESIS rvalue CLOSE PARENTHESIS
{
	$$ = $3
}

rvalue: rvalue1
{
	$$ = $1
}

rvalue1: bit
{
	$$ = &BitValue{
		Bit: $1,
	}
}

statement: lvalue EQUALS rvalue
{
	if (!$1.Pointer() || !$3.Pointer()) && (!$1.Value() || !$3.Value()) {
		yylex.Error("invalid assignment: " + $1.String() + " EQUALS " + $3.String())
	}
	$$ = &EqualsStmt{
		Left:  $1,
		Right: $3,
	}
}

number: bit
{
	$$ = &bitnum.Number{}
	$$.Append($1)
}

number: number bit
{
	$$ = $1
	$$.Append($2)
}

LINE : 'L' 'I' 'N' 'E';
NUMBER: 'N' 'U' 'M' 'B' 'E' 'R';
ZERO: 'Z' 'E' 'R' 'O';
ONE: 'O' 'N' 'E';
CODE: 'C' 'O' 'D' 'E';
READ: 'R' 'E' 'A' 'D';
PRINT: 'P' 'R' 'I' 'N' 'T';
GOTO: 'G' 'O' 'T' 'O';
IF: 'I' 'F';
THE: 'T' 'H' 'E';
JUMP: 'J' 'U' 'M' 'P';
REGISTER: 'R' 'E' 'G' 'I' 'S' 'T' 'E' 'R';
IS: 'I' 'S';
VARIABLE: 'V' 'A' 'R' 'I' 'A' 'B' 'L' 'E';
VALUE: 'V' 'A' 'L' 'U' 'E';
AT: 'A' 'T';
BEYOND: 'B' 'E' 'Y' 'O' 'N' 'D';
ADDRESS: 'A' 'D' 'D' 'R' 'E' 'S' 'S';
OF: 'O' 'F';
NAND: 'N' 'A' 'N' 'D';
OPEN: 'O' 'P' 'E' 'N';
PARENTHESIS: 'P' 'A' 'R' 'E' 'N' 'T' 'H' 'E' 'S' 'I' 'S';
CLOSE: 'C' 'L' 'O' 'S' 'E';
EQUALS: 'E' 'Q' 'U' 'A' 'L' 'S';

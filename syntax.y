%{
//go:generate go get golang.org/x/tools/cmd/goyacc
//go:generate goyacc syntax.y

package bit

import "fmt"
%}

%union {
	program *Program
	expr    Expr
	stmt    Stmt
	goto0   *uint64
	goto1   *uint64
	number  uint64
	line    struct {
		stmt    Stmt
		goto0   *uint64
		goto1   *uint64
		number  uint64
	}
	numberbits struct {
		number uint64
		bits    uint8
	}
}

%type <program>    program
%type <line>       line goto
%type <stmt>       stmt
%type <expr>       expr expr1 expr2
%type <number>     bit_constant
%type <numberbits> number

%%

program
: line
	{
		$$ = new(Program)
		$$.AddLine($1.number, $1.stmt, $1.goto0, $1.goto1)
		yylex.(*lex).prog = $$
	}
| program line
	{
		$$ = $1
		$$.AddLine($2.number, $2.stmt, $2.goto0, $2.goto1)
		yylex.(*lex).prog = $$
	}
;

bit_constant
: ZERO
	{
		$$ = 0
	}
| ONE
	{
		$$ = 1
	}
;

number
: bit_constant
	{
		$$.number = $1
		$$.bits = 1
	}
| number bit_constant
	{
		$$.number = $1.number<<1 | $2
		$$.bits = $1.bits + 1
		if $$.bits == 64 {
			panic(fmt.Errorf("bit: integer overflow"))
		}
	}
;

goto
: GOTO number
	{
		$$.goto0 = new(uint64)
		$$.goto1 = new(uint64)
		*$$.goto0, *$$.goto1 = $2.number, $2.number
	}
| GOTO number IF THE JUMP REGISTER IS bit_constant
	{
		if $8 == 0 {
			$$.goto0 = new(uint64)
			*$$.goto0 = $2.number
			$$.goto1 = nil
		} else {
			$$.goto1 = new(uint64)
			*$$.goto1 = $2.number
			$$.goto0 = nil
		}
	}
;

expr
: expr NAND expr1
	{
		$$ = NandExpr{$1, $3}
	}
| expr1
	{
		$$ = $1
	}
;

expr1
: THE ADDRESS OF expr1
	{
		$$ = AddrExpr{$4}
	}
| THE VALUE BEYOND expr1
	{
		$$ = NextExpr{$4, 0}
	}
| THE VALUE AT expr1
	{
		$$ = StarExpr{$4}
	}
| expr2
	{
		$$ = $1
	}
;

expr2
: OPEN PARENTHESIS expr CLOSE PARENTHESIS
	{
		$$ = $3
	}
| VARIABLE number
	{
		$$ = VarExpr($2.number)
	}
| bit_constant
	{
		if $1 == 0 {
			$$ = BitExpr(false)
		} else {
			$$ = BitExpr(true)
		}
	}
;

stmt
: expr EQUALS expr
	{
		$$ = AssignStmt{$1, $3}
	}
| THE JUMP REGISTER EQUALS expr
	{
		$$ = JumpRegisterStmt{$5}
	}
| PRINT bit_constant
	{
		if $2 == 0 {
			$$ = PrintStmt(false)
		} else {
			$$ = PrintStmt(true)
		}
	}
| READ
	{
		$$ = ReadStmt{}
	}
| READ number /* non-standard */
	{
		pc := new(uint64)
		*pc = $2.number
		$$ = ReadStmt{
			pc: pc,
		}
	}
;

line
: LINE NUMBER number CODE stmt
	{
		$$.number, $$.stmt = $3.number, $5
		$$.goto0, $$.goto1 = nil, nil
	}
| line goto
	{
		$$.number, $$.stmt = $1.number, $1.stmt
		$$.goto0, $$.goto1 = $1.goto0, $1.goto1

		if $2.goto0 != nil {
			if $$.goto0 != nil {
				panic(fmt.Errorf("bit: duplicate goto on line %v", $$.number))
			}
			$$.goto0 = $2.goto0
		}

		if $2.goto1 != nil {
			if $$.goto1 != nil {
				panic(fmt.Errorf("bit: duplicate goto on line %v", $$.number))
			}
			$$.goto1 = $2.goto1
		}
	}
;

ZERO: 'Z' 'E' 'R' 'O';
ONE: 'O' 'N' 'E';
GOTO: 'G' 'O' 'T' 'O';
LINE: 'L' 'I' 'N' 'E';
NUMBER: 'N' 'U' 'M' 'B' 'E' 'R';
CODE: 'C' 'O' 'D' 'E';
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
EQUALS: 'E' 'Q' 'U' 'A' 'L' 'S';
OPEN: 'O' 'P' 'E' 'N';
CLOSE: 'C' 'L' 'O' 'S' 'E';
PARENTHESIS: 'P' 'A' 'R' 'E' 'N' 'T' 'H' 'E' 'S' 'I' 'S';
PRINT: 'P' 'R' 'I' 'N' 'T';
READ: 'R' 'E' 'A' 'D';

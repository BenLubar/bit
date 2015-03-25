%{
//go:generate go tool yacc syntax.y

package bit

import "fmt"
%}

%union {
	program Program
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

%token ZERO ONE GOTO LINE NUMBER CODE IF THE JUMP REGISTER IS VARIABLE THE VALUE AT BEYOND ADDRESS OF NAND EQUALS OPEN CLOSE PARENTHESIS PRINT READ

%type <program>    program
%type <line>       line goto
%type <stmt>       stmt
%type <expr>       variable expr expr1 expr2 expr3 expr4
%type <number>     bit_constant
%type <numberbits> number

%%

program
: line
	{
		$$ = make(Program)
		if err := $$.AddLine($1.number, $1.stmt, $1.goto0, $1.goto1); err != nil {
			panic(err)
		}
		yylex.(*lex).prog = $$
	}
| program line
	{
		$$ = $1
		if err := $$.AddLine($2.number, $2.stmt, $2.goto0, $2.goto1); err != nil {
			panic(err)
		}
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

variable
: VARIABLE number
	{
		$$ = VarExpr($2.number)
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
| expr2
	{
		$$ = $1
	}
;

expr2
: THE VALUE BEYOND expr2
	{
		$$ = NextExpr{$4}
	}
| expr3
	{
		$$ = $1
	}
;

expr3
: THE VALUE AT expr3
	{
		$$ = StarExpr{$4}
	}
| expr4
	{
		$$ = $1
	}
;

expr4
: OPEN PARENTHESIS expr CLOSE PARENTHESIS
	{
		$$ = $3
	}
| variable
	{
		$$ = $1
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

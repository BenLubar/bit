%{
//go:generate go tool yacc syntax.y

package main

import "github.com/BenLubar/bit/cmd/brainfuckc/bf"
%}

%union {
	many []bf.Command
	one  bf.Command
	tok  bf.BF
}

%type<many> program stmts
%type<one> stmt
%type<tok> Right Left
%type<tok> Increment Decrement
%type<tok> Output Input
%type<tok> Begin End

%%

program
: stmts
	{
		$$ = $1
		yylex.(*lex).prog = $$
	}
;

stmts
:
	{
		$$ = nil
	}
| stmts stmt
	{
		$$ = append($1, $2)
	}
;

stmt
: Begin stmts End
	{
		$$ = bf.Command{
			Token: $1,
			Loop:  $2,
		}
	}
| Right
	{
		$$ = bf.Command{
			Token: $1,
		}
	}
| Left
	{
		$$ = bf.Command{
			Token: $1,
		}
	}
| Increment
	{
		$$ = bf.Command{
			Token: $1,
		}
	}
| Decrement
	{
		$$ = bf.Command{
			Token: $1,
		}
	}
| Output
	{
		$$ = bf.Command{
			Token: $1,
		}
	}
| Input
	{
		$$ = bf.Command{
			Token: $1,
		}
	}
;

Ook
: 'O' 'o' 'k'
;

Right
: Ook '.' Ook '?'
	{
		$$ = bf.Right
	}
;

Left
: Ook '?' Ook '.'
	{
		$$ = bf.Left
	}
;

Increment
: Ook '.' Ook '.'
	{
		$$ = bf.Increment
	}
;

Decrement
: Ook '!' Ook '!'
	{
		$$ = bf.Decrement
	}
;

Input
: Ook '.' Ook '!'
	{
		$$ = bf.Input
	}
;

Output
: Ook '!' Ook '.'
	{
		$$ = bf.Output
	}
;

Begin
: Ook '!' Ook '?'
	{
		$$ = bf.Begin
	}
;

End
: Ook '?' Ook '!'
	{
		$$ = bf.End
	}
;

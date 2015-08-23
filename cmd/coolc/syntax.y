%{
package main
%}

%union {
	id  ID
	typ TYPE
	n   int32
	s   string

	cls *ClassDecl
	vd  []*VarDecl
	ftr []Feature
	exp Expr
	act []Expr
	cas *Cases
}

%token<id>  tokID
%token<typ> tokTYPE
%token<n>   tokINTEGER
%token<s>   tokSTRING

%token tokCLASS tokEXTENDS tokLPAREN tokRPAREN tokVAR tokCOLON tokCOMMA
%token tokLBRACE tokRBRACE tokOVERRIDE tokDEF tokSEMICOLON tokELSE tokSUPER
%token tokNEW tokNULL tokTRUE tokFALSE tokTHIS tokCASE tokARROW

%token tokINVALID /* returned by the lexer when an error occurs */

%left tokDOT
%left tokNEGATE
%left tokMULTIPLY tokDIVIDE
%left tokPLUS tokMINUS
%left tokEQUALEQUAL
%left tokLESSEQUAL tokLESSTHAN
%left tokMATCH
%left tokIF tokWHILE
%left tokASSIGN

%type<cls> classdecl
%type<vd>  varformals varformals0 formals formals0
%type<ftr> classbody features
%type<exp> block block0 expr primary
%type<act> actuals actuals0
%type<cas> cases case

%%

program
: /* empty */
| program classdecl
	{
		yylex.(*lexer).ast.Classes = append(yylex.(*lexer).ast.Classes, $2)
	}
;

classdecl
: tokCLASS tokTYPE varformals classbody
	{
		$$ = &ClassDecl{
			Name: $2,
			Args: $3,
			Body: $4,
		}
	}
| tokCLASS tokTYPE varformals tokEXTENDS tokTYPE actuals classbody
	{
		$$ = &ClassDecl{
			Name: $2,
			Args: $3,
			Extends: &ExtendsDecl{
				Type: $5,
				Args: $6,
			},
			Body: $7,
		}
	}
;

varformals
: tokLPAREN tokRPAREN
	{
		$$ = nil
	}
| tokLPAREN varformals0 tokRPAREN
	{
		$$ = $2
	}
;

varformals0
: tokVAR tokID tokCOLON tokTYPE
	{
		$$ = append([]*VarDecl(nil), &VarDecl{
			Name: $2,
			Type: $4,
		})
	}
| varformals0 tokCOMMA tokVAR tokID tokCOLON tokTYPE
	{
		$$ = append($1, &VarDecl{
			Name: $4,
			Type: $6,
		})
	}
;

classbody
: tokLBRACE features tokRBRACE
	{
		$$ = $2
	}
;

features
: /* empty */
	{
		$$ = nil
	}
| features tokOVERRIDE tokDEF tokID formals tokCOLON tokTYPE tokASSIGN expr tokSEMICOLON
	{
		$$ = append($1, &MethodFeature{
			Override: true,
			Name:     $4,
			Args:     $5,
			Return:   $7,
			Body:     $9,
		})
	}
| features tokDEF tokID formals tokCOLON tokTYPE tokASSIGN expr tokSEMICOLON
	{
		$$ = append($1, &MethodFeature{
			Name:   $3,
			Args:   $4,
			Return: $6,
			Body:   $8,
		})
	}
| features tokVAR tokID tokCOLON tokTYPE tokASSIGN expr tokSEMICOLON
	{
		$$ = append($1, &VarFeature{
			VarDecl: VarDecl{
				Name: $3,
				Type: $5,
			},
			Value: $7,
		})
	}
| features tokLBRACE block tokRBRACE tokSEMICOLON
	{
		$$ = append($1, &BlockFeature{
			Expr: $3,
		})
	}
;

formals
: tokLPAREN tokRPAREN
	{
		$$ = nil
	}
| tokLPAREN formals0 tokRPAREN
	{
		$$ = $2
	}
;

formals0
: tokID tokCOLON tokTYPE
	{
		$$ = append([]*VarDecl(nil), &VarDecl{
			Name: $1,
			Type: $3,
		})
	}
| formals0 tokCOMMA tokID tokCOLON tokTYPE
	{
		$$ = append($1, &VarDecl{
			Name: $3,
			Type: $5,
		})
	}
;

actuals
: tokLPAREN tokRPAREN
	{
		$$ = nil
	}
| tokLPAREN actuals0 tokRPAREN
	{
		$$ = $2
	}
;

actuals0
: expr
	{
		$$ = append([]Expr(nil), $1)
	}
| actuals0 tokCOMMA expr
	{
		$$ = append($1, $3)
	}
;

block
: /* empty */
	{
		$$ = &UnitExpr{}
	}
| block0
	{
		$$ = $1
	}
;

block0
: expr
	{
		$$ = $1
	}
| expr tokSEMICOLON block0
	{
		$$ = &ChainExpr{
			Pre:  $1,
			Expr: $3,
		}
	}
| tokVAR tokID tokCOLON tokTYPE tokASSIGN expr tokSEMICOLON block0
	{
		$$ = &VarExpr{
			Var: &VarFeature{
				
			},
			Expr: $8,
		}
	}
;

expr
: primary
	{
		$$ = $1
	}
| tokID tokASSIGN expr
	{
		$$ = &AssignExpr{
			Left:  $1,
			Right: $3,
		}
	}
| tokNEGATE expr
	{
		$$ = &NotExpr{
			Right: $2,
		}
	}
| tokMINUS expr %prec tokNEGATE
	{
		$$ = &NegativeExpr{
			Right: $2,
		}
	}
| tokIF tokLPAREN expr tokRPAREN expr tokELSE expr %prec tokIF
	{
		$$ = &IfExpr{
			Condition: $3,
			Then:      $5,
			Else:      $7,
		}
	}
| tokWHILE tokLPAREN expr tokRPAREN expr %prec tokWHILE
	{
		$$ = &WhileExpr{
			Condition: $3,
			Do:        $5,
		}
	}
| expr tokLESSEQUAL expr
	{
		$$ = &LessThanOrEqualExpr{
			Left:  $1,
			Right: $3,
		}
	}
| expr tokLESSTHAN expr
	{
		$$ = &LessThanExpr{
			Left:  $1,
			Right: $3,
		}
	}
| expr tokEQUALEQUAL expr
	{
		$$ = &EqualEqualExpr{
			Left:  $1,
			Right: $3,
		}
	}
| expr tokMULTIPLY expr
	{
		$$ = &MultiplyExpr{
			Left:  $1,
			Right: $3,
		}
	}
| expr tokDIVIDE expr
	{
		$$ = &DivideExpr{
			Left:  $1,
			Right: $3,
		}
	}
| expr tokPLUS expr
	{
		$$ = &AddExpr{
			Left:  $1,
			Right: $3,
		}
	}
| expr tokMINUS expr
	{
		$$ = &SubtractExpr{
			Left:  $1,
			Right: $3,
		}
	}
| expr tokMATCH tokLBRACE cases tokRBRACE /* note: this slightly differs from the diagrams in cool-manual.pdf */
	{
		$$ = &MatchExpr{
			Left:  $1,
			Cases: $4,
		}
	}
| expr tokDOT tokID actuals
	{
		$$ = &CallExpr{
			Left: $1,
			Name: $3,
			Args: $4,
		}
	}
;

primary
: tokSUPER tokDOT tokID actuals
	{
		$$ = &SelfCallExpr{
			Super: true,
			Name:  $3,
			Args:  $4,
		}
	}
| tokID actuals
	{
		$$ = &SelfCallExpr{
			Name: $1,
			Args: $2,
		}
	}
| tokNEW tokTYPE actuals
	{
		$$ = &NewExpr{
			Type: $2,
			Args: $3,
		}
	}
| tokLBRACE block tokRBRACE
	{
		$$ = $2
	}
| tokLPAREN expr tokRPAREN
	{
		$$ = $2
	}
| tokNULL
	{
		$$ = &NullExpr{}
	}
| tokLPAREN tokRPAREN
	{
		$$ = &UnitExpr{}
	}
| tokID
	{
		$$ = &NameExpr{
			Name: $1,
		}
	}
| tokINTEGER
	{
		$$ = &IntegerExpr{
			N: $1,
		}
	}
| tokSTRING
	{
		$$ = &StringExpr{
			S: $1,
		}
	}
| tokTRUE
	{
		$$ = &BooleanExpr{
			B: true,
		}
	}
| tokFALSE
	{
		$$ = &BooleanExpr{
			B: false,
		}
	}
| tokTHIS
	{
		$$ = &ThisExpr{}
	}
;

cases
: case
	{
		$$ = $1
	}
| cases case
	{
		$$ = $1
		if $2.Null != nil {
			if $$.Null == nil {
				$$.Null = $2.Null
			} else {
				yylex.Error("duplicate null case")
			}
		}
		$$.Cases = append($$.Cases, $2.Cases...)
	}
;

case
: tokCASE tokID tokCOLON tokTYPE tokARROW block
	{
		$$ = &Cases{
			Cases: []*Case{
				{
					Name: $2,
					Type: $4,
					Body: $6,
				},
			},
		}
	}
| tokCASE tokNULL tokARROW block
	{
		$$ = &Cases{
			Null: $4,
		}
	}
;

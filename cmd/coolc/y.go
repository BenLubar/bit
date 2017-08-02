//line syntax.y:2
package main

import __yyfmt__ "fmt"

//line syntax.y:2
//line syntax.y:5
type yySymType struct {
	yys int
	id  ID
	typ TYPE
	n   int32
	s   string

	cls *ClassDecl
	vd  []*VarDecl
	ftr []Feature
	exp Expr
	act []Expr
	cas *Case
	css []*Case
}

const tokID = 57346
const tokTYPE = 57347
const tokINTEGER = 57348
const tokSTRING = 57349
const tokCLASS = 57350
const tokEXTENDS = 57351
const tokLPAREN = 57352
const tokRPAREN = 57353
const tokVAR = 57354
const tokCOLON = 57355
const tokCOMMA = 57356
const tokLBRACE = 57357
const tokRBRACE = 57358
const tokOVERRIDE = 57359
const tokDEF = 57360
const tokSEMICOLON = 57361
const tokELSE = 57362
const tokSUPER = 57363
const tokNEW = 57364
const tokTRUE = 57365
const tokFALSE = 57366
const tokTHIS = 57367
const tokCASE = 57368
const tokARROW = 57369
const tokNULL = 57370
const tokINVALID = 57371
const tokASSIGN = 57372
const tokIF = 57373
const tokWHILE = 57374
const tokMATCH = 57375
const tokLESSEQUAL = 57376
const tokLESSTHAN = 57377
const tokEQUALEQUAL = 57378
const tokPLUS = 57379
const tokMINUS = 57380
const tokMULTIPLY = 57381
const tokDIVIDE = 57382
const tokNEGATE = 57383
const tokDOT = 57384

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"tokID",
	"tokTYPE",
	"tokINTEGER",
	"tokSTRING",
	"tokCLASS",
	"tokEXTENDS",
	"tokLPAREN",
	"tokRPAREN",
	"tokVAR",
	"tokCOLON",
	"tokCOMMA",
	"tokLBRACE",
	"tokRBRACE",
	"tokOVERRIDE",
	"tokDEF",
	"tokSEMICOLON",
	"tokELSE",
	"tokSUPER",
	"tokNEW",
	"tokTRUE",
	"tokFALSE",
	"tokTHIS",
	"tokCASE",
	"tokARROW",
	"tokNULL",
	"tokINVALID",
	"tokASSIGN",
	"tokIF",
	"tokWHILE",
	"tokMATCH",
	"tokLESSEQUAL",
	"tokLESSTHAN",
	"tokEQUALEQUAL",
	"tokPLUS",
	"tokMINUS",
	"tokMULTIPLY",
	"tokDIVIDE",
	"tokNEGATE",
	"tokDOT",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 316

var yyAct = [...]int{

	52, 51, 114, 50, 18, 79, 66, 60, 63, 64,
	61, 62, 32, 66, 42, 43, 73, 142, 40, 77,
	30, 61, 62, 39, 66, 19, 138, 134, 125, 37,
	38, 44, 45, 46, 69, 70, 41, 68, 154, 35,
	36, 76, 129, 75, 140, 67, 34, 127, 115, 33,
	63, 64, 61, 62, 109, 66, 101, 115, 86, 87,
	88, 89, 90, 91, 92, 93, 130, 47, 96, 7,
	82, 94, 97, 98, 8, 32, 9, 42, 43, 100,
	9, 40, 144, 53, 103, 110, 39, 122, 27, 25,
	123, 139, 37, 38, 44, 45, 46, 124, 56, 41,
	116, 57, 35, 36, 119, 120, 111, 15, 104, 34,
	16, 85, 33, 155, 81, 26, 128, 23, 131, 132,
	24, 20, 21, 22, 10, 12, 137, 65, 58, 59,
	60, 63, 64, 61, 62, 143, 66, 19, 80, 146,
	72, 71, 149, 150, 148, 32, 107, 42, 43, 6,
	152, 40, 28, 105, 3, 156, 39, 147, 157, 136,
	133, 126, 37, 38, 44, 45, 46, 121, 112, 41,
	108, 74, 35, 36, 32, 135, 42, 43, 55, 34,
	40, 13, 33, 4, 1, 39, 99, 95, 84, 78,
	54, 37, 38, 44, 45, 46, 49, 48, 41, 153,
	17, 35, 36, 113, 29, 31, 14, 106, 34, 151,
	11, 33, 5, 65, 58, 59, 60, 63, 64, 61,
	62, 145, 66, 65, 58, 59, 60, 63, 64, 61,
	62, 2, 66, 141, 118, 65, 58, 59, 60, 63,
	64, 61, 62, 0, 66, 117, 65, 58, 59, 60,
	63, 64, 61, 62, 0, 66, 65, 58, 59, 60,
	63, 64, 61, 62, 102, 66, 0, 65, 58, 59,
	60, 63, 64, 61, 62, 0, 66, 0, 0, 0,
	0, 0, 83, 0, 0, 0, 65, 58, 59, 60,
	63, 64, 61, 62, 0, 66, 65, 58, 59, 60,
	63, 64, 61, 62, 0, 66, 65, 58, 59, 60,
	63, 64, 61, 62, 0, 66,
}
var yyPact = [...]int{

	-1000, 146, -1000, 178, 139, 65, 113, -1000, 176, -1000,
	-1000, 96, 196, 127, 105, -1000, 77, 102, 61, 141,
	-1000, 49, 193, 192, 71, 186, 173, -1000, -1000, 87,
	273, -1000, 15, 170, 170, 131, 130, -26, 166, 71,
	8, -1000, -1000, -1000, -1000, -1000, -1000, 185, 128, 101,
	54, -1000, 263, 184, 98, -1000, -1000, 170, 170, 170,
	170, 170, 170, 170, 170, 56, 183, 170, -1000, -36,
	-36, 170, 170, 182, 127, 40, 253, -1000, 128, 95,
	142, 165, 35, 71, 93, 163, 273, -29, -29, 13,
	-36, -36, -18, -18, 22, 127, 273, 234, 223, 127,
	-1000, -1000, -1000, 92, 162, -1000, 76, 84, -2, -1000,
	-1000, 156, -1000, 31, -1000, 38, -1000, 170, 170, -1000,
	155, -3, -1000, 171, 154, 170, -4, -1000, -1000, 78,
	17, 213, 273, -13, 170, 69, -1000, 202, 170, 152,
	71, 170, 170, 190, 145, -1000, 180, 11, -1000, 273,
	94, -1000, -1000, 71, 71, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 231, 212, 210, 5, 207, 69, 206, 3, 1,
	0, 205, 4, 204, 203, 2, 184,
}
var yyR1 = [...]int{

	0, 16, 16, 1, 1, 2, 2, 3, 3, 6,
	7, 7, 7, 7, 7, 4, 4, 5, 5, 12,
	12, 13, 13, 8, 8, 9, 9, 9, 10, 10,
	10, 10, 10, 10, 10, 10, 10, 10, 10, 10,
	10, 10, 10, 11, 11, 11, 11, 11, 11, 11,
	11, 11, 11, 11, 11, 11, 14, 14, 15, 15,
}
var yyR2 = [...]int{

	0, 0, 2, 4, 7, 2, 3, 4, 6, 3,
	0, 10, 9, 8, 5, 2, 3, 3, 5, 2,
	3, 1, 3, 0, 1, 1, 3, 8, 1, 3,
	2, 2, 7, 5, 3, 3, 3, 3, 3, 3,
	3, 5, 4, 4, 2, 3, 3, 3, 1, 2,
	1, 1, 1, 1, 1, 1, 1, 2, 6, 4,
}
var yyChk = [...]int{

	-1000, -16, -1, 8, 5, -2, 10, -6, 9, 15,
	11, -3, 12, 5, -7, 11, 14, 4, -12, 10,
	16, 17, 18, 12, 15, 12, 13, -6, 11, -13,
	-10, -11, 4, 41, 38, 31, 32, 21, 22, 15,
	10, 28, 6, 7, 23, 24, 25, 18, 4, 4,
	-8, -9, -10, 12, 4, 5, 11, 14, 34, 35,
	36, 39, 40, 37, 38, 33, 42, 30, -12, -10,
	-10, 10, 10, 42, 5, -8, -10, 11, 4, -4,
	10, 13, 16, 19, 4, 13, -10, -10, -10, -10,
	-10, -10, -10, -10, 15, 4, -10, -10, -10, 4,
	-12, 16, 11, -4, 13, 11, -5, 4, 5, 19,
	-9, 13, 5, -14, -15, 26, -12, 11, 11, -12,
	13, 5, 11, 14, 13, 30, 5, 16, -15, 4,
	28, -10, -10, 5, 30, 4, 5, -10, 30, 13,
	27, 20, 30, -10, 13, 19, -10, 5, -8, -10,
	-10, 19, 5, 19, 27, 19, -9, -8,
}
var yyDef = [...]int{

	1, -2, 2, 0, 0, 0, 0, 3, 0, 10,
	5, 0, 0, 0, 0, 6, 0, 0, 0, 0,
	9, 0, 0, 0, 23, 0, 0, 4, 19, 0,
	21, 28, 50, 0, 0, 0, 0, 0, 0, 23,
	0, 48, 51, 52, 53, 54, 55, 0, 0, 0,
	0, 24, 25, 0, 0, 7, 20, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 44, 30,
	31, 0, 0, 0, 0, 0, 0, 49, 0, 0,
	0, 0, 0, 0, 0, 0, 22, 34, 35, 36,
	37, 38, 39, 40, 0, 0, 29, 0, 0, 0,
	45, 46, 47, 0, 0, 15, 0, 0, 0, 14,
	26, 0, 8, 0, 56, 0, 42, 0, 0, 43,
	0, 0, 16, 0, 0, 0, 0, 41, 57, 0,
	0, 0, 33, 0, 0, 0, 17, 0, 0, 0,
	23, 0, 0, 0, 0, 13, 0, 0, 59, 32,
	0, 12, 18, 0, 23, 11, 27, 58,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:55
		{
			yylex.(*lexer).ast.Classes = append(yylex.(*lexer).ast.Classes, yyDollar[2].cls)
		}
	case 3:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line syntax.y:62
		{
			yyVAL.cls = &ClassDecl{
				Name: yyDollar[2].typ,
				Args: yyDollar[3].vd,
				Body: yyDollar[4].ftr,
			}
		}
	case 4:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line syntax.y:70
		{
			yyVAL.cls = &ClassDecl{
				Name: yyDollar[2].typ,
				Args: yyDollar[3].vd,
				Extends: &ExtendsDecl{
					Type: yyDollar[5].typ,
					Args: yyDollar[6].act,
				},
				Body: yyDollar[7].ftr,
			}
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:85
		{
			yyVAL.vd = nil
		}
	case 6:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:89
		{
			yyVAL.vd = yyDollar[2].vd
		}
	case 7:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line syntax.y:96
		{
			yyVAL.vd = append([]*VarDecl(nil), &VarDecl{
				Name: yyDollar[2].id,
				Type: yyDollar[4].typ,
			})
		}
	case 8:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line syntax.y:103
		{
			yyVAL.vd = append(yyDollar[1].vd, &VarDecl{
				Name: yyDollar[4].id,
				Type: yyDollar[6].typ,
			})
		}
	case 9:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:113
		{
			yyVAL.ftr = yyDollar[2].ftr
		}
	case 10:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line syntax.y:120
		{
			yyVAL.ftr = nil
		}
	case 11:
		yyDollar = yyS[yypt-10 : yypt+1]
		//line syntax.y:124
		{
			yyVAL.ftr = append(yyDollar[1].ftr, &MethodFeature{
				Override: true,
				Name:     yyDollar[4].id,
				Args:     yyDollar[5].vd,
				Return:   yyDollar[7].typ,
				Body:     yyDollar[9].exp,
			})
		}
	case 12:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line syntax.y:134
		{
			yyVAL.ftr = append(yyDollar[1].ftr, &MethodFeature{
				Name:   yyDollar[3].id,
				Args:   yyDollar[4].vd,
				Return: yyDollar[6].typ,
				Body:   yyDollar[8].exp,
			})
		}
	case 13:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line syntax.y:143
		{
			yyVAL.ftr = append(yyDollar[1].ftr, &VarFeature{
				VarDecl: VarDecl{
					Name: yyDollar[3].id,
					Type: yyDollar[5].typ,
				},
				Value: yyDollar[7].exp,
			})
		}
	case 14:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line syntax.y:153
		{
			yyVAL.ftr = append(yyDollar[1].ftr, &BlockFeature{
				Expr: yyDollar[3].exp,
			})
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:162
		{
			yyVAL.vd = nil
		}
	case 16:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:166
		{
			yyVAL.vd = yyDollar[2].vd
		}
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:173
		{
			yyVAL.vd = append([]*VarDecl(nil), &VarDecl{
				Name: yyDollar[1].id,
				Type: yyDollar[3].typ,
			})
		}
	case 18:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line syntax.y:180
		{
			yyVAL.vd = append(yyDollar[1].vd, &VarDecl{
				Name: yyDollar[3].id,
				Type: yyDollar[5].typ,
			})
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:190
		{
			yyVAL.act = nil
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:194
		{
			yyVAL.act = yyDollar[2].act
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:201
		{
			yyVAL.act = append([]Expr(nil), yyDollar[1].exp)
		}
	case 22:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:205
		{
			yyVAL.act = append(yyDollar[1].act, yyDollar[3].exp)
		}
	case 23:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line syntax.y:212
		{
			yyVAL.exp = &UnitExpr{}
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:216
		{
			yyVAL.exp = yyDollar[1].exp
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:223
		{
			yyVAL.exp = yyDollar[1].exp
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:227
		{
			yyVAL.exp = &ChainExpr{
				Pre:  yyDollar[1].exp,
				Expr: yyDollar[3].exp,
			}
		}
	case 27:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line syntax.y:234
		{
			yyVAL.exp = &VarExpr{
				VarFeature: VarFeature{
					VarDecl: VarDecl{
						Name: yyDollar[2].id,
						Type: yyDollar[4].typ,
					},
					Value: yyDollar[6].exp,
				},
				Expr: yyDollar[8].exp,
			}
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:250
		{
			yyVAL.exp = yyDollar[1].exp
		}
	case 29:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:254
		{
			yyVAL.exp = &AssignExpr{
				Left:  yyDollar[1].id,
				Right: yyDollar[3].exp,
			}
		}
	case 30:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:261
		{
			yyVAL.exp = &CallExpr{
				Left: yyDollar[2].exp,
				Name: yyDollar[1].id,
			}
		}
	case 31:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:268
		{
			yyDollar[1].id.Name = "_negative"
			yyVAL.exp = &CallExpr{
				Left: yyDollar[2].exp,
				Name: yyDollar[1].id,
			}
		}
	case 32:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line syntax.y:276
		{
			yyVAL.exp = &IfExpr{
				Pos:       yyDollar[1].typ.Pos,
				Condition: yyDollar[3].exp,
				Then:      yyDollar[5].exp,
				Else:      yyDollar[7].exp,
			}
		}
	case 33:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line syntax.y:285
		{
			yyVAL.exp = &WhileExpr{
				Pos:       yyDollar[1].typ.Pos,
				Condition: yyDollar[3].exp,
				Do:        yyDollar[5].exp,
			}
		}
	case 34:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:293
		{
			yyVAL.exp = &CallExpr{
				Left: yyDollar[1].exp,
				Name: yyDollar[2].id,
				Args: []Expr{
					yyDollar[3].exp,
				},
			}
		}
	case 35:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:303
		{
			yyVAL.exp = &CallExpr{
				Left: yyDollar[1].exp,
				Name: yyDollar[2].id,
				Args: []Expr{
					yyDollar[3].exp,
				},
			}
		}
	case 36:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:313
		{
			yyVAL.exp = &CallExpr{
				Left: yyDollar[1].exp,
				Name: yyDollar[2].id,
				Args: []Expr{
					yyDollar[3].exp,
				},
			}
		}
	case 37:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:323
		{
			yyVAL.exp = &CallExpr{
				Left: yyDollar[1].exp,
				Name: yyDollar[2].id,
				Args: []Expr{
					yyDollar[3].exp,
				},
			}
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:333
		{
			yyVAL.exp = &CallExpr{
				Left: yyDollar[1].exp,
				Name: yyDollar[2].id,
				Args: []Expr{
					yyDollar[3].exp,
				},
			}
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:343
		{
			yyVAL.exp = &CallExpr{
				Left: yyDollar[1].exp,
				Name: yyDollar[2].id,
				Args: []Expr{
					yyDollar[3].exp,
				},
			}
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:353
		{
			yyVAL.exp = &CallExpr{
				Left: yyDollar[1].exp,
				Name: yyDollar[2].id,
				Args: []Expr{
					yyDollar[3].exp,
				},
			}
		}
	case 41:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line syntax.y:363
		{
			yyVAL.exp = &MatchExpr{
				Left:  yyDollar[1].exp,
				Cases: yyDollar[4].css,
			}
		}
	case 42:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line syntax.y:370
		{
			yyVAL.exp = &CallExpr{
				Left: yyDollar[1].exp,
				Name: yyDollar[3].id,
				Args: yyDollar[4].act,
			}
		}
	case 43:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line syntax.y:381
		{
			yyVAL.exp = &StaticCallExpr{
				Name: yyDollar[3].id,
				Args: yyDollar[4].act,
			}
		}
	case 44:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:388
		{
			yyVAL.exp = &CallExpr{
				Left: &ThisExpr{},
				Name: yyDollar[1].id,
				Args: yyDollar[2].act,
			}
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:396
		{
			yyVAL.exp = &CallExpr{
				Left: &NewExpr{
					Type: yyDollar[2].typ,
				},
				Name: ID{
					Name: yyDollar[2].typ.Name,
					Pos:  yyDollar[2].typ.Pos,
				},
				Args: yyDollar[3].act,
			}
		}
	case 46:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:409
		{
			yyVAL.exp = yyDollar[2].exp
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:413
		{
			yyVAL.exp = yyDollar[2].exp
		}
	case 48:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:417
		{
			yyVAL.exp = &NullExpr{}
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:421
		{
			yyVAL.exp = &UnitExpr{}
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:425
		{
			yyVAL.exp = &NameExpr{
				Name: yyDollar[1].id,
			}
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:431
		{
			yyVAL.exp = &IntegerExpr{
				N: yyDollar[1].n,
			}
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:437
		{
			yyVAL.exp = &StringExpr{
				S: yyDollar[1].s,
			}
		}
	case 53:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:443
		{
			yyVAL.exp = &BooleanExpr{
				B: true,
			}
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:449
		{
			yyVAL.exp = &BooleanExpr{
				B: false,
			}
		}
	case 55:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:455
		{
			yyVAL.exp = &ThisExpr{}
		}
	case 56:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:462
		{
			yyVAL.css = []*Case{yyDollar[1].cas}
		}
	case 57:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:466
		{
			yyVAL.css = append(yyDollar[1].css, yyDollar[2].cas)
		}
	case 58:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line syntax.y:473
		{
			yyVAL.cas = &Case{
				Name: yyDollar[2].id,
				Type: yyDollar[4].typ,
				Body: yyDollar[6].exp,
			}
		}
	case 59:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line syntax.y:481
		{
			yyVAL.cas = &Case{
				Name: ID{
					Name: "null",
					Pos:  yyDollar[2].typ.Pos,
				},
				Type: yyDollar[2].typ,
				Body: yyDollar[4].exp,
			}
		}
	}
	goto yystack /* stack new state and value */
}

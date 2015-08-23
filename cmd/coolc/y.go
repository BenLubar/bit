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
	cas *Cases
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
const tokNULL = 57365
const tokTRUE = 57366
const tokFALSE = 57367
const tokTHIS = 57368
const tokCASE = 57369
const tokARROW = 57370
const tokINVALID = 57371
const tokDOT = 57372
const tokNEGATE = 57373
const tokMULTIPLY = 57374
const tokDIVIDE = 57375
const tokPLUS = 57376
const tokMINUS = 57377
const tokEQUALEQUAL = 57378
const tokLESSEQUAL = 57379
const tokLESSTHAN = 57380
const tokMATCH = 57381
const tokIF = 57382
const tokWHILE = 57383
const tokASSIGN = 57384

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
	"tokNULL",
	"tokTRUE",
	"tokFALSE",
	"tokTHIS",
	"tokCASE",
	"tokARROW",
	"tokINVALID",
	"tokDOT",
	"tokNEGATE",
	"tokMULTIPLY",
	"tokDIVIDE",
	"tokPLUS",
	"tokMINUS",
	"tokEQUALEQUAL",
	"tokLESSEQUAL",
	"tokLESSTHAN",
	"tokMATCH",
	"tokIF",
	"tokWHILE",
	"tokASSIGN",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 60
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 321

var yyAct = [...]int{

	52, 51, 114, 50, 18, 79, 142, 19, 138, 134,
	125, 32, 65, 42, 43, 73, 154, 40, 140, 53,
	30, 115, 39, 60, 58, 59, 65, 109, 37, 38,
	41, 44, 45, 46, 69, 70, 47, 68, 33, 67,
	101, 76, 34, 75, 82, 25, 127, 35, 36, 61,
	62, 63, 64, 60, 58, 59, 65, 115, 86, 87,
	88, 89, 90, 91, 92, 93, 94, 9, 96, 58,
	59, 65, 97, 98, 144, 32, 139, 42, 43, 100,
	129, 40, 77, 8, 103, 110, 39, 3, 7, 9,
	124, 120, 37, 38, 41, 44, 45, 46, 122, 130,
	116, 123, 33, 111, 119, 56, 34, 27, 57, 104,
	15, 35, 36, 16, 85, 81, 128, 26, 131, 132,
	63, 64, 60, 58, 59, 65, 137, 10, 12, 19,
	32, 107, 42, 43, 80, 143, 40, 28, 105, 146,
	72, 39, 149, 150, 148, 71, 6, 37, 38, 41,
	44, 45, 46, 152, 147, 156, 136, 33, 157, 133,
	126, 34, 32, 1, 42, 43, 35, 36, 40, 121,
	112, 23, 108, 39, 24, 20, 21, 22, 74, 37,
	38, 41, 44, 45, 46, 55, 113, 13, 4, 33,
	155, 135, 99, 34, 95, 84, 78, 29, 35, 36,
	153, 66, 54, 61, 62, 63, 64, 60, 58, 59,
	65, 66, 151, 61, 62, 63, 64, 60, 58, 59,
	65, 49, 145, 66, 48, 61, 62, 63, 64, 60,
	58, 59, 65, 66, 141, 61, 62, 63, 64, 60,
	58, 59, 65, 118, 66, 17, 61, 62, 63, 64,
	60, 58, 59, 65, 117, 31, 14, 106, 11, 5,
	2, 0, 66, 0, 61, 62, 63, 64, 60, 58,
	59, 65, 102, 66, 0, 61, 62, 63, 64, 60,
	58, 59, 65, 0, 0, 0, 0, 0, 0, 0,
	83, 66, 0, 61, 62, 63, 64, 60, 58, 59,
	65, 66, 0, 61, 62, 63, 64, 60, 58, 59,
	65, 66, 0, 61, 62, 63, 64, 60, 58, 59,
	65,
}
var yyPact = [...]int{

	-1000, 79, -1000, 183, 136, 74, 116, -1000, 182, -1000,
	-1000, 99, 241, 119, 159, -1000, 33, 104, 52, 126,
	-1000, 18, 220, 217, 7, 198, 180, -1000, -1000, 94,
	281, -1000, -3, 158, 158, 135, 130, -15, 173, 7,
	71, -1000, -1000, -1000, -1000, -1000, -1000, 192, 124, 102,
	28, -1000, 271, 191, 101, -1000, -1000, 158, 158, 158,
	158, 158, 158, 158, 158, 51, 190, 158, -1000, 17,
	17, 158, 158, 188, 119, 24, 261, -1000, 124, 96,
	127, 167, 8, 7, 90, 165, 281, -27, -27, 32,
	86, 86, -13, -13, -6, 119, -1000, 243, 232, 119,
	-1000, -1000, -1000, 78, 164, -1000, 87, 77, -32, -1000,
	-1000, 155, -1000, 30, -1000, 76, -1000, 158, 158, -1000,
	154, -33, -1000, 187, 151, 158, -34, -1000, -1000, 63,
	-10, 214, -1000, -36, 158, 61, -1000, 203, 158, 149,
	7, 158, 158, 193, 148, -1000, 181, -12, -1000, -1000,
	171, -1000, -1000, 7, 7, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 260, 259, 258, 5, 257, 88, 256, 3, 1,
	0, 255, 4, 197, 186, 2, 163,
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
	-10, -11, 4, 31, 35, 40, 41, 21, 22, 15,
	10, 23, 6, 7, 24, 25, 26, 18, 4, 4,
	-8, -9, -10, 12, 4, 5, 11, 14, 37, 38,
	36, 32, 33, 34, 35, 39, 30, 42, -12, -10,
	-10, 10, 10, 30, 5, -8, -10, 11, 4, -4,
	10, 13, 16, 19, 4, 13, -10, -10, -10, -10,
	-10, -10, -10, -10, 15, 4, -10, -10, -10, 4,
	-12, 16, 11, -4, 13, 11, -5, 4, 5, 19,
	-9, 13, 5, -14, -15, 27, -12, 11, 11, -12,
	13, 5, 11, 14, 13, 42, 5, 16, -15, 4,
	23, -10, -10, 5, 42, 4, 5, -10, 42, 13,
	28, 20, 42, -10, 13, 19, -10, 5, -8, -10,
	-10, 19, 5, 19, 28, 19, -9, -8,
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
	lookahead func() int
}

func (p *yyParserImpl) Lookahead() int {
	return p.lookahead()
}

func yyNewParser() yyParser {
	p := &yyParserImpl{
		lookahead: func() int { return -1 },
	}
	return p
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
	var yylval yySymType
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yytoken := -1 // yychar translated into internal numbering
	yyrcvr.lookahead = func() int { return yychar }
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yychar = -1
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
	if yychar < 0 {
		yychar, yytoken = yylex1(yylex, &yylval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yychar = -1
		yytoken = -1
		yyVAL = yylval
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
		if yychar < 0 {
			yychar, yytoken = yylex1(yylex, &yylval)
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
			yychar = -1
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
		//line syntax.y:52
		{
			yylex.(*lexer).ast.Classes = append(yylex.(*lexer).ast.Classes, yyDollar[2].cls)
		}
	case 3:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line syntax.y:59
		{
			yyVAL.cls = &ClassDecl{
				Name: yyDollar[2].typ,
				Args: yyDollar[3].vd,
				Body: yyDollar[4].ftr,
			}
		}
	case 4:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line syntax.y:67
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
		//line syntax.y:82
		{
			yyVAL.vd = nil
		}
	case 6:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:86
		{
			yyVAL.vd = yyDollar[2].vd
		}
	case 7:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line syntax.y:93
		{
			yyVAL.vd = append([]*VarDecl(nil), &VarDecl{
				Name: yyDollar[2].id,
				Type: yyDollar[4].typ,
			})
		}
	case 8:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line syntax.y:100
		{
			yyVAL.vd = append(yyDollar[1].vd, &VarDecl{
				Name: yyDollar[4].id,
				Type: yyDollar[6].typ,
			})
		}
	case 9:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:110
		{
			yyVAL.ftr = yyDollar[2].ftr
		}
	case 10:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line syntax.y:117
		{
			yyVAL.ftr = nil
		}
	case 11:
		yyDollar = yyS[yypt-10 : yypt+1]
		//line syntax.y:121
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
		//line syntax.y:131
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
		//line syntax.y:140
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
		//line syntax.y:150
		{
			yyVAL.ftr = append(yyDollar[1].ftr, &BlockFeature{
				Expr: yyDollar[3].exp,
			})
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:159
		{
			yyVAL.vd = nil
		}
	case 16:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:163
		{
			yyVAL.vd = yyDollar[2].vd
		}
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:170
		{
			yyVAL.vd = append([]*VarDecl(nil), &VarDecl{
				Name: yyDollar[1].id,
				Type: yyDollar[3].typ,
			})
		}
	case 18:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line syntax.y:177
		{
			yyVAL.vd = append(yyDollar[1].vd, &VarDecl{
				Name: yyDollar[3].id,
				Type: yyDollar[5].typ,
			})
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:187
		{
			yyVAL.act = nil
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:191
		{
			yyVAL.act = yyDollar[2].act
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:198
		{
			yyVAL.act = append([]Expr(nil), yyDollar[1].exp)
		}
	case 22:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:202
		{
			yyVAL.act = append(yyDollar[1].act, yyDollar[3].exp)
		}
	case 23:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line syntax.y:209
		{
			yyVAL.exp = &UnitExpr{}
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:213
		{
			yyVAL.exp = yyDollar[1].exp
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:220
		{
			yyVAL.exp = yyDollar[1].exp
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:224
		{
			yyVAL.exp = &ChainExpr{
				Pre:  yyDollar[1].exp,
				Expr: yyDollar[3].exp,
			}
		}
	case 27:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line syntax.y:231
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
		//line syntax.y:247
		{
			yyVAL.exp = yyDollar[1].exp
		}
	case 29:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:251
		{
			yyVAL.exp = &AssignExpr{
				Left:  yyDollar[1].id,
				Right: yyDollar[3].exp,
			}
		}
	case 30:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:258
		{
			yyVAL.exp = &NotExpr{
				Right: yyDollar[2].exp,
			}
		}
	case 31:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:264
		{
			yyVAL.exp = &NegativeExpr{
				Right: yyDollar[2].exp,
			}
		}
	case 32:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line syntax.y:270
		{
			yyVAL.exp = &IfExpr{
				Condition: yyDollar[3].exp,
				Then:      yyDollar[5].exp,
				Else:      yyDollar[7].exp,
			}
		}
	case 33:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line syntax.y:278
		{
			yyVAL.exp = &WhileExpr{
				Condition: yyDollar[3].exp,
				Do:        yyDollar[5].exp,
			}
		}
	case 34:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:285
		{
			yyVAL.exp = &LessThanOrEqualExpr{
				Left:  yyDollar[1].exp,
				Right: yyDollar[3].exp,
			}
		}
	case 35:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:292
		{
			yyVAL.exp = &LessThanExpr{
				Left:  yyDollar[1].exp,
				Right: yyDollar[3].exp,
			}
		}
	case 36:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:299
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
		//line syntax.y:309
		{
			yyVAL.exp = &MultiplyExpr{
				Left:  yyDollar[1].exp,
				Right: yyDollar[3].exp,
			}
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:316
		{
			yyVAL.exp = &DivideExpr{
				Left:  yyDollar[1].exp,
				Right: yyDollar[3].exp,
			}
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:323
		{
			yyVAL.exp = &AddExpr{
				Left:  yyDollar[1].exp,
				Right: yyDollar[3].exp,
			}
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:330
		{
			yyVAL.exp = &SubtractExpr{
				Left:  yyDollar[1].exp,
				Right: yyDollar[3].exp,
			}
		}
	case 41:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line syntax.y:337
		{
			yyVAL.exp = &MatchExpr{
				Left:  yyDollar[1].exp,
				Cases: yyDollar[4].cas,
			}
		}
	case 42:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line syntax.y:344
		{
			yyVAL.exp = &CallExpr{
				Left: yyDollar[1].exp,
				Name: yyDollar[3].id,
				Args: yyDollar[4].act,
			}
		}
	case 43:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line syntax.y:355
		{
			yyVAL.exp = &SelfCallExpr{
				Super: true,
				Name:  yyDollar[3].id,
				Args:  yyDollar[4].act,
			}
		}
	case 44:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:363
		{
			yyVAL.exp = &SelfCallExpr{
				Name: yyDollar[1].id,
				Args: yyDollar[2].act,
			}
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:370
		{
			yyVAL.exp = &NewExpr{
				Type: yyDollar[2].typ,
				Args: yyDollar[3].act,
			}
		}
	case 46:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:377
		{
			yyVAL.exp = yyDollar[2].exp
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:381
		{
			yyVAL.exp = yyDollar[2].exp
		}
	case 48:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:385
		{
			yyVAL.exp = &NullExpr{}
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:389
		{
			yyVAL.exp = &UnitExpr{}
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:393
		{
			yyVAL.exp = &NameExpr{
				Name: yyDollar[1].id,
			}
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:399
		{
			yyVAL.exp = &IntegerExpr{
				N: yyDollar[1].n,
			}
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:405
		{
			yyVAL.exp = &StringExpr{
				S: yyDollar[1].s,
			}
		}
	case 53:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:411
		{
			yyVAL.exp = &BooleanExpr{
				B: true,
			}
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:417
		{
			yyVAL.exp = &BooleanExpr{
				B: false,
			}
		}
	case 55:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:423
		{
			yyVAL.exp = &ThisExpr{}
		}
	case 56:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:430
		{
			yyVAL.cas = yyDollar[1].cas
		}
	case 57:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:434
		{
			yyVAL.cas = yyDollar[1].cas
			if yyDollar[2].cas.Null != nil {
				if yyVAL.cas.Null == nil {
					yyVAL.cas.Null = yyDollar[2].cas.Null
				} else {
					yylex.Error("duplicate null case")
				}
			}
			yyVAL.cas.Cases = append(yyVAL.cas.Cases, yyDollar[2].cas.Cases...)
		}
	case 58:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line syntax.y:449
		{
			yyVAL.cas = &Cases{
				Cases: []*Case{
					{
						Name: yyDollar[2].id,
						Type: yyDollar[4].typ,
						Body: yyDollar[6].exp,
					},
				},
			}
		}
	case 59:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line syntax.y:461
		{
			yyVAL.cas = &Cases{
				Null: yyDollar[4].exp,
			}
		}
	}
	goto yystack /* stack new state and value */
}

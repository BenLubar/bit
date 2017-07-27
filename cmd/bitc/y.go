//line syntax.y:2

//go:generate go get golang.org/x/tools/cmd/goyacc
//go:generate goyacc syntax.y

package main

import __yyfmt__ "fmt"

//line syntax.y:5
//line syntax.y:8
type yySymType struct {
	yys     int
	program *Program
	lines   []*Line
	line    *Line
	bit     bool
	num     *Number
	stmt    Stmt
	expr    Expr
}

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"'L'",
	"'I'",
	"'N'",
	"'E'",
	"'U'",
	"'M'",
	"'B'",
	"'R'",
	"'Z'",
	"'O'",
	"'C'",
	"'D'",
	"'A'",
	"'P'",
	"'T'",
	"'G'",
	"'F'",
	"'H'",
	"'J'",
	"'S'",
	"'V'",
	"'Y'",
	"'Q'",
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

const yyLast = 207

var yyAct = [...]int{

	12, 145, 13, 11, 49, 63, 82, 40, 70, 89,
	61, 60, 29, 67, 19, 14, 66, 117, 30, 34,
	177, 34, 88, 171, 32, 33, 10, 33, 51, 51,
	52, 43, 51, 151, 52, 162, 157, 156, 155, 142,
	141, 164, 54, 124, 104, 41, 19, 116, 65, 62,
	34, 158, 19, 42, 150, 115, 98, 84, 48, 22,
	90, 118, 19, 120, 114, 75, 112, 58, 28, 136,
	91, 74, 107, 77, 68, 53, 65, 65, 14, 15,
	20, 15, 125, 96, 97, 110, 93, 83, 103, 65,
	55, 14, 15, 65, 14, 15, 106, 127, 148, 132,
	109, 111, 105, 59, 56, 35, 14, 71, 121, 165,
	137, 128, 65, 81, 78, 47, 36, 134, 38, 123,
	100, 23, 119, 92, 76, 16, 168, 163, 159, 152,
	135, 147, 140, 133, 113, 99, 138, 95, 80, 79,
	57, 46, 45, 37, 24, 21, 153, 149, 154, 143,
	129, 122, 94, 86, 22, 17, 8, 146, 174, 131,
	102, 160, 161, 69, 9, 144, 139, 130, 101, 5,
	19, 19, 44, 169, 170, 172, 173, 175, 176, 126,
	180, 181, 183, 182, 178, 179, 3, 166, 167, 6,
	64, 85, 108, 87, 73, 72, 50, 31, 27, 26,
	18, 7, 4, 25, 39, 2, 1,
}
var yyPact = [...]int{

	165, -1000, 165, -1000, 150, 159, -1000, 79, 117, 149,
	66, -1000, -1000, -1000, 138, 148, 112, 137, 1, -1000,
	92, 105, 136, 108, -1000, 26, 79, 79, 135, 134,
	104, 79, 10, 59, 21, 75, 91, -1000, 133, -1000,
	79, 90, 79, -1000, 3, -13, 58, 158, 79, 96,
	55, 116, 57, 103, 132, 131, -1000, 102, 82, 39,
	-1000, 147, -1000, 6, 43, -1000, 53, 115, 71, 146,
	-1000, 130, 3, 3, 38, 128, 111, 164, 155, -1000,
	-1000, -1000, 32, 24, 89, 3, 56, 87, 70, 3,
	50, 127, 48, -1000, 37, 28, -1000, -1000, -1000, -8,
	44, 114, 47, 7, -1000, -1000, -1000, 145, 3, 23,
	67, 83, 100, 144, 163, -1000, 154, 86, -1000, 126,
	107, 96, 54, -1000, -1000, 99, 43, 162, 125, -1000,
	17, 16, 143, -1000, 161, 152, -1000, 124, -1000, 85,
	141, -1000, 36, 18, 122, 79, 15, 14, 13, 33,
	121, -1000, -1000, 26, 26, -1000, 12, 120, 20, 98,
	79, 79, -1000, -1000, 119, -1000, 82, 82, 0, 32,
	32, 153, 7, 7, -3, 96, 96, -1000, 152, 152,
	68, 94, -1000, -1000,
}
var yyPgo = [...]int{

	0, 206, 205, 186, 204, 3, 26, 203, 49, 11,
	10, 202, 201, 200, 0, 2, 199, 198, 7, 6,
	5, 4, 8, 1, 197, 196, 195, 194, 193, 192,
	191, 190, 9, 179, 172,
}
var yyR1 = [...]int{

	0, 1, 2, 2, 3, 5, 5, 7, 7, 7,
	4, 4, 4, 4, 4, 4, 8, 8, 10, 8,
	8, 10, 9, 10, 9, 10, 7, 6, 6, 11,
	12, 14, 15, 13, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34,
}
var yyR2 = [...]int{

	0, 1, 1, 2, 6, 1, 1, 1, 2, 2,
	0, 2, 8, 8, 16, 16, 2, 3, 1, 4,
	4, 4, 3, 5, 1, 1, 3, 1, 2, 4,
	6, 4, 3, 4, 4, 5, 4, 2, 3, 4,
	8, 2, 8, 5, 2, 6, 7, 2, 4, 4,
	11, 5, 6,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -11, 4, -3, -12, 6, 5,
	-6, -5, -14, -15, 12, 13, 8, 6, -13, -5,
	14, 7, 6, 9, 7, -7, -16, -17, -8, 11,
	17, -24, -20, 24, 18, 13, 11, 7, 10, -4,
	-18, 19, -6, -5, -34, 7, 7, 11, -6, -21,
	-25, 22, 24, 16, 21, 15, 13, 7, -6, 13,
	-9, -10, -8, -20, -31, -5, 13, 26, 16, 5,
	-22, 11, -26, -27, 16, 10, 8, 16, 11, 7,
	7, 11, -19, 5, 18, -30, 6, -28, 16, -32,
	17, 17, 8, 15, 6, 7, -10, -10, 18, 7,
	9, 4, 5, -20, 20, 13, -10, 16, -29, 13,
	15, -9, 16, 7, 16, 18, 19, 25, 17, 8,
	16, -21, 6, -10, 20, 15, -33, 14, 11, 6,
	4, 5, 13, 7, 10, -22, 15, 11, -32, 4,
	7, 23, 23, 6, 4, -23, 5, 7, 13, 6,
	18, 15, 7, -14, -15, 23, 23, 23, 18, 7,
	-18, -18, 23, 7, 21, 11, -6, -6, 7, -19,
	-19, 23, -20, -20, 5, -21, -21, 23, -22, -22,
	-23, -23, -15, -14,
}
var yyDef = [...]int{

	0, -2, 1, 2, 0, 0, 3, 0, 0, 0,
	0, 27, 5, 6, 0, 0, 0, 0, 0, 28,
	0, 0, 0, 0, 29, 10, 7, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 32, 0, 4,
	0, 0, 8, 9, 0, 0, 0, 0, 16, 0,
	0, 0, 0, 0, 0, 0, 31, 0, 11, 0,
	26, 24, 18, 0, 0, 25, 0, 0, 0, 0,
	17, 0, 0, 0, 0, 0, 0, 0, 0, 38,
	33, 30, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 34, 0, 0, 19, 20, 44, 0,
	0, 0, 0, 0, 37, 36, 22, 0, 0, 0,
	0, 0, 0, 0, 0, 35, 0, 0, 39, 0,
	0, 0, 0, 21, 47, 0, 0, 0, 0, 49,
	0, 0, 0, 43, 0, 0, 48, 0, 23, 0,
	0, 52, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 45, 42, 12, 13, 41, 0, 0, 0, 0,
	0, 0, 46, 51, 0, 40, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 50, 0, 0,
	0, 0, 14, 15,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 16, 10, 14, 15, 7,
	20, 19, 21, 5, 22, 3, 4, 9, 6, 13,
	17, 26, 11, 23, 18, 8, 24, 3, 3, 25,
	12,
}
var yyTok2 = [...]int{

	2, 3,
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

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:29
		{
			yyVAL.program = &Program{
				Lines: yyDollar[1].lines,
			}
			yylex.(*lex).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:37
		{
			yyVAL.lines = []*Line{yyDollar[1].line}
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:42
		{
			yyVAL.lines = append(yyDollar[1].lines, yyDollar[2].line)
		}
	case 4:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line syntax.y:47
		{
			yyVAL.line = yyDollar[6].line
			yyVAL.line.Num = yyDollar[3].num
			yyVAL.line.Stmt = yyDollar[5].stmt
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:54
		{
			yyVAL.bit = false
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:59
		{
			yyVAL.bit = true
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:64
		{
			yyVAL.stmt = &ReadStmt{}
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:69
		{
			// extension
			yyVAL.stmt = &ReadStmt{
				EOFLine: yyDollar[2].num,
			}
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:77
		{
			yyVAL.stmt = &PrintStmt{
				Bit: yyDollar[2].bit,
			}
		}
	case 10:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line syntax.y:84
		{
			yyVAL.line = &Line{}
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:89
		{
			yyVAL.line = &Line{
				Zero: yyDollar[2].num,
				One:  yyDollar[2].num,
			}
		}
	case 12:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line syntax.y:97
		{
			yyVAL.line = &Line{
				Zero: yyDollar[2].num,
			}
		}
	case 13:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line syntax.y:104
		{
			yyVAL.line = &Line{
				One: yyDollar[2].num,
			}
		}
	case 14:
		yyDollar = yyS[yypt-16 : yypt+1]
		//line syntax.y:111
		{
			yyVAL.line = &Line{
				Zero: yyDollar[2].num,
				One:  yyDollar[10].num,
			}
		}
	case 15:
		yyDollar = yyS[yypt-16 : yypt+1]
		//line syntax.y:119
		{
			yyVAL.line = &Line{
				Zero: yyDollar[10].num,
				One:  yyDollar[2].num,
			}
		}
	case 16:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:127
		{
			yyVAL.expr = &UnknownVariable{
				Num: yyDollar[2].num,
			}
		}
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:134
		{
			yyVAL.expr = &JumpRegister{}
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:139
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 19:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line syntax.y:144
		{
			if !yyDollar[4].expr.Pointer() {
				yylex.Error("not a pointer: " + yyDollar[4].expr.String())
			}
			yyVAL.expr = &ValueAt{
				Target: yyDollar[4].expr,
				Offset: 0,
			}
		}
	case 20:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line syntax.y:155
		{
			if !yyDollar[4].expr.Pointer() {
				yylex.Error("not a pointer: " + yyDollar[4].expr.String())
			}
			yyVAL.expr = &ValueAt{
				Target: yyDollar[4].expr,
				Offset: 1,
			}
		}
	case 21:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line syntax.y:166
		{
			if !yyDollar[4].expr.Addressable() {
				yylex.Error("not addressable: " + yyDollar[4].expr.String())
			}
			yyVAL.expr = &AddressOf{
				Variable: yyDollar[4].expr,
			}
		}
	case 22:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:176
		{
			if !yyDollar[1].expr.Value() {
				yylex.Error("not a value: " + yyDollar[1].expr.String())
			}
			if !yyDollar[3].expr.Value() {
				yylex.Error("not a value: " + yyDollar[3].expr.String())
			}
			yyVAL.expr = &Nand{
				Left:  yyDollar[1].expr,
				Right: yyDollar[3].expr,
			}
		}
	case 23:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line syntax.y:190
		{
			yyVAL.expr = yyDollar[3].expr
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:195
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:200
		{
			yyVAL.expr = &BitValue{
				Bit: yyDollar[1].bit,
			}
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:207
		{
			if (!yyDollar[1].expr.Pointer() || !yyDollar[3].expr.Pointer()) && (!yyDollar[1].expr.Value() || !yyDollar[3].expr.Value()) {
				yylex.Error("invalid assignment: " + yyDollar[1].expr.String() + " EQUALS " + yyDollar[3].expr.String())
			}
			yyVAL.stmt = &EqualsStmt{
				Left:  yyDollar[1].expr,
				Right: yyDollar[3].expr,
			}
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:218
		{
			yyVAL.num = &Number{}
			yyVAL.num.Append(yyDollar[1].bit)
		}
	case 28:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:224
		{
			yyVAL.num = yyDollar[1].num
			yyVAL.num.Append(yyDollar[2].bit)
		}
	}
	goto yystack /* stack new state and value */
}

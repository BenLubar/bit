//line syntax.y:2

//go:generate go tool yacc syntax.y

package bit

import __yyfmt__ "fmt"

//line syntax.y:4
import "fmt"

//line syntax.y:9
type yySymType struct {
	yys     int
	program *Program
	expr    Expr
	stmt    Stmt
	goto0   *uint64
	goto1   *uint64
	number  uint64
	line    struct {
		stmt   Stmt
		goto0  *uint64
		goto1  *uint64
		number uint64
	}
	numberbits struct {
		number uint64
		bits   uint8
	}
}

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"'Z'",
	"'E'",
	"'R'",
	"'O'",
	"'N'",
	"'G'",
	"'T'",
	"'L'",
	"'I'",
	"'U'",
	"'M'",
	"'B'",
	"'C'",
	"'D'",
	"'F'",
	"'H'",
	"'J'",
	"'P'",
	"'S'",
	"'V'",
	"'A'",
	"'Y'",
	"'Q'",
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

const yyNprod = 49
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 190

var yyAct = [...]int{

	71, 39, 49, 59, 43, 78, 12, 54, 83, 86,
	13, 16, 13, 45, 50, 22, 19, 33, 132, 140,
	16, 55, 22, 50, 66, 65, 33, 91, 44, 131,
	51, 167, 32, 66, 65, 99, 92, 40, 97, 51,
	55, 95, 85, 75, 67, 13, 165, 158, 64, 68,
	153, 13, 152, 150, 149, 73, 127, 72, 129, 26,
	160, 84, 56, 82, 111, 34, 16, 157, 148, 17,
	93, 22, 74, 98, 130, 60, 22, 120, 29, 135,
	117, 94, 76, 122, 53, 107, 30, 109, 16, 80,
	20, 17, 110, 166, 112, 113, 24, 116, 141, 125,
	105, 96, 11, 155, 142, 138, 118, 126, 4, 156,
	154, 136, 115, 33, 27, 128, 61, 8, 133, 60,
	151, 147, 137, 60, 16, 124, 119, 17, 28, 108,
	26, 21, 10, 146, 143, 88, 57, 52, 37, 18,
	164, 134, 123, 103, 101, 79, 69, 35, 163, 162,
	161, 159, 145, 144, 139, 61, 114, 106, 102, 100,
	81, 77, 70, 36, 31, 25, 2, 9, 5, 3,
	42, 41, 48, 121, 47, 90, 89, 63, 87, 62,
	58, 104, 23, 7, 15, 14, 46, 38, 6, 1,
}
var yyPact = [...]int{

	97, 97, 108, 124, 90, 108, -1000, 120, 132, 120,
	77, 123, 84, -1000, -1000, -1000, 160, 122, 104, 62,
	72, 159, -1000, 103, 47, 141, 158, 131, 7, 130,
	69, -1000, 20, 43, -1000, 129, -1000, -1000, -1000, 111,
	1, 120, 120, -1000, 140, 157, -1000, 36, 120, -1000,
	51, 19, 65, 156, 139, 76, 155, -1000, 16, 16,
	18, -17, 128, 12, 139, 64, 17, -1000, 120, 89,
	14, 16, 11, 120, 154, 138, 153, 137, 88, 152,
	71, -1000, -1000, 10, 115, 121, 74, 16, 46, 16,
	16, 151, 102, 150, 63, 95, 118, 60, 67, 136,
	117, 87, -1000, -1000, 120, 34, 106, 37, 57, 5,
	-1000, -1000, -1000, -1000, -7, -1000, 16, 135, 66, 101,
	-1000, 36, 94, 149, -1000, -5, -1000, -1000, 86, -1000,
	-1000, 93, 127, 115, 148, 147, -1000, -1000, 126, 113,
	53, 32, 31, 112, 30, -1000, 28, 100, 92, 99,
	-1000, 50, 25, 146, 41, 145, 144, -1000, -1000, -1000,
	143, -1000, 134, 24, -1000, 81, 9, -1000,
}
var yyPgo = [...]int{

	0, 189, 166, 188, 187, 1, 4, 186, 2, 6,
	185, 184, 183, 182, 8, 7, 5, 181, 180, 179,
	178, 177, 176, 175, 174, 0, 173, 172, 3, 171,
	170, 169, 167, 128,
}
var yyR1 = [...]int{

	0, 1, 1, 8, 8, 9, 9, 3, 3, 5,
	5, 6, 6, 6, 6, 7, 7, 7, 4, 4,
	4, 4, 4, 2, 2, 10, 11, 12, 31, 32,
	33, 13, 14, 15, 16, 17, 27, 21, 23, 22,
	19, 20, 18, 28, 24, 26, 25, 29, 30,
}
var yyR2 = [...]int{

	0, 1, 2, 1, 1, 1, 2, 2, 8, 3,
	1, 4, 4, 4, 1, 5, 2, 1, 3, 5,
	2, 1, 2, 5, 2, 4, 3, 4, 4, 6,
	4, 2, 3, 4, 8, 2, 8, 5, 2, 6,
	7, 2, 4, 6, 4, 5, 11, 5, 4,
}
var yyChk = [...]int{

	-1000, -1, -2, -31, 11, -2, -3, -12, 9, -32,
	8, 12, -9, -8, -10, -11, 4, 7, 7, -9,
	13, 8, -8, -13, 12, 5, 8, 10, -33, 16,
	14, 5, -14, 10, 18, 6, 5, 7, -4, -5,
	-14, -29, -30, -6, 21, 6, -7, -24, -27, -8,
	7, 23, 7, 15, -15, 20, 19, 7, -18, -28,
	8, 5, -19, -21, -15, 24, 23, -8, -9, 6,
	5, -25, 21, -9, 21, 24, 17, 5, -16, 6,
	13, 5, -6, -14, -5, 24, 26, -20, 7, -22,
	-23, 15, 24, -16, 17, 24, 12, 24, -5, 24,
	5, 6, 5, 6, -17, 12, 5, 14, 8, 13,
	-6, 18, -6, -6, 5, 10, -28, 17, 11, 8,
	17, -26, 16, 6, 8, 12, -8, 22, 9, 21,
	17, 24, 25, -5, 6, 13, 10, -25, 11, 5,
	24, 12, 11, 7, 5, 5, 7, 8, 15, 22,
	22, 8, 22, 22, 10, 11, 10, 17, 22, 5,
	19, 5, 5, 5, 6, 22, 12, 22,
}
var yyDef = [...]int{

	0, -2, 1, 0, 0, 2, 24, 0, 0, 0,
	0, 0, 7, 5, 3, 4, 0, 0, 0, 0,
	0, 0, 6, 0, 0, 0, 0, 0, 0, 0,
	0, 28, 0, 0, 31, 0, 26, 27, 23, 0,
	0, 0, 21, 10, 0, 0, 14, 0, 0, 17,
	0, 0, 0, 0, 0, 0, 0, 25, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 20, 22, 0,
	0, 0, 0, 16, 0, 0, 0, 0, 0, 0,
	0, 32, 9, 0, 18, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 30, 29, 0, 0, 0, 0, 0, 0,
	11, 41, 12, 13, 0, 38, 0, 0, 0, 0,
	48, 0, 0, 0, 44, 0, 8, 35, 0, 33,
	42, 0, 0, 19, 0, 0, 47, 15, 0, 0,
	0, 0, 0, 0, 0, 37, 0, 0, 0, 0,
	43, 0, 0, 0, 0, 0, 0, 39, 40, 45,
	0, 36, 0, 0, 34, 0, 0, 46,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 24, 15, 16, 17, 5,
	18, 9, 19, 12, 20, 3, 11, 14, 8, 7,
	21, 26, 6, 22, 10, 13, 23, 3, 3, 25,
	4,
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

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:39
		{
			yyVAL.program = new(Program)
			yyVAL.program.AddLine(yyDollar[1].line.number, yyDollar[1].line.stmt, yyDollar[1].line.goto0, yyDollar[1].line.goto1)
			yylex.(*lex).prog = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:45
		{
			yyVAL.program = yyDollar[1].program
			yyVAL.program.AddLine(yyDollar[2].line.number, yyDollar[2].line.stmt, yyDollar[2].line.goto0, yyDollar[2].line.goto1)
			yylex.(*lex).prog = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:54
		{
			yyVAL.number = 0
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:58
		{
			yyVAL.number = 1
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:65
		{
			yyVAL.numberbits.number = yyDollar[1].number
			yyVAL.numberbits.bits = 1
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:70
		{
			yyVAL.numberbits.number = yyDollar[1].numberbits.number<<1 | yyDollar[2].number
			yyVAL.numberbits.bits = yyDollar[1].numberbits.bits + 1
			if yyVAL.numberbits.bits == 64 {
				panic(fmt.Errorf("bit: integer overflow"))
			}
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:81
		{
			yyVAL.line.goto0 = new(uint64)
			yyVAL.line.goto1 = new(uint64)
			*yyVAL.line.goto0, *yyVAL.line.goto1 = yyDollar[2].numberbits.number, yyDollar[2].numberbits.number
		}
	case 8:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line syntax.y:87
		{
			if yyDollar[8].number == 0 {
				yyVAL.line.goto0 = new(uint64)
				*yyVAL.line.goto0 = yyDollar[2].numberbits.number
				yyVAL.line.goto1 = nil
			} else {
				yyVAL.line.goto1 = new(uint64)
				*yyVAL.line.goto1 = yyDollar[2].numberbits.number
				yyVAL.line.goto0 = nil
			}
		}
	case 9:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:102
		{
			yyVAL.expr = NandExpr{yyDollar[1].expr, yyDollar[3].expr}
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:106
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 11:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line syntax.y:113
		{
			yyVAL.expr = AddrExpr{yyDollar[4].expr}
		}
	case 12:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line syntax.y:117
		{
			yyVAL.expr = NextExpr{yyDollar[4].expr, 0}
		}
	case 13:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line syntax.y:121
		{
			yyVAL.expr = StarExpr{yyDollar[4].expr}
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:125
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 15:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line syntax.y:132
		{
			yyVAL.expr = yyDollar[3].expr
		}
	case 16:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:136
		{
			yyVAL.expr = VarExpr(yyDollar[2].numberbits.number)
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:140
		{
			if yyDollar[1].number == 0 {
				yyVAL.expr = BitExpr(false)
			} else {
				yyVAL.expr = BitExpr(true)
			}
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:151
		{
			yyVAL.stmt = AssignStmt{yyDollar[1].expr, yyDollar[3].expr}
		}
	case 19:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line syntax.y:155
		{
			yyVAL.stmt = JumpRegisterStmt{yyDollar[5].expr}
		}
	case 20:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:159
		{
			if yyDollar[2].number == 0 {
				yyVAL.stmt = PrintStmt(false)
			} else {
				yyVAL.stmt = PrintStmt(true)
			}
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:167
		{
			yyVAL.stmt = ReadStmt{}
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:171
		{
			pc := new(uint64)
			*pc = yyDollar[2].numberbits.number
			yyVAL.stmt = ReadStmt{
				pc: pc,
			}
		}
	case 23:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line syntax.y:182
		{
			yyVAL.line.number, yyVAL.line.stmt = yyDollar[3].numberbits.number, yyDollar[5].stmt
			yyVAL.line.goto0, yyVAL.line.goto1 = nil, nil
		}
	case 24:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:187
		{
			yyVAL.line.number, yyVAL.line.stmt = yyDollar[1].line.number, yyDollar[1].line.stmt
			yyVAL.line.goto0, yyVAL.line.goto1 = yyDollar[1].line.goto0, yyDollar[1].line.goto1

			if yyDollar[2].line.goto0 != nil {
				if yyVAL.line.goto0 != nil {
					panic(fmt.Errorf("bit: duplicate goto on line %v", yyVAL.line.number))
				}
				yyVAL.line.goto0 = yyDollar[2].line.goto0
			}

			if yyDollar[2].line.goto1 != nil {
				if yyVAL.line.goto1 != nil {
					panic(fmt.Errorf("bit: duplicate goto on line %v", yyVAL.line.number))
				}
				yyVAL.line.goto1 = yyDollar[2].line.goto1
			}
		}
	}
	goto yystack /* stack new state and value */
}

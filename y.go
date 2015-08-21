//line syntax.y:2

//go:generate go tool yacc syntax.y

package bit

import __yyfmt__ "fmt"

//line syntax.y:4
import "fmt"

//line syntax.y:9
type yySymType struct {
	yys     int
	program Program
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

const ZERO = 57346
const ONE = 57347
const GOTO = 57348
const LINE = 57349
const NUMBER = 57350
const CODE = 57351
const IF = 57352
const THE = 57353
const JUMP = 57354
const REGISTER = 57355
const IS = 57356
const VARIABLE = 57357
const VALUE = 57358
const AT = 57359
const BEYOND = 57360
const ADDRESS = 57361
const OF = 57362
const NAND = 57363
const EQUALS = 57364
const OPEN = 57365
const CLOSE = 57366
const PARENTHESIS = 57367
const PRINT = 57368
const READ = 57369

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"ZERO",
	"ONE",
	"GOTO",
	"LINE",
	"NUMBER",
	"CODE",
	"IF",
	"THE",
	"JUMP",
	"REGISTER",
	"IS",
	"VARIABLE",
	"VALUE",
	"AT",
	"BEYOND",
	"ADDRESS",
	"OF",
	"NAND",
	"EQUALS",
	"OPEN",
	"CLOSE",
	"PARENTHESIS",
	"PRINT",
	"READ",
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

const yyNprod = 24
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 71

var yyAct = [...]int{

	26, 22, 18, 10, 11, 10, 11, 9, 9, 13,
	38, 53, 19, 13, 25, 8, 25, 28, 34, 49,
	50, 33, 24, 12, 24, 28, 9, 20, 21, 40,
	37, 32, 39, 28, 29, 31, 13, 44, 30, 42,
	41, 35, 46, 47, 48, 31, 51, 45, 30, 43,
	36, 27, 52, 10, 11, 16, 10, 11, 7, 14,
	3, 15, 10, 11, 2, 6, 4, 23, 17, 5,
	1,
}
var yyPact = [...]int{

	53, 53, 59, 50, 59, -1000, 58, 58, 49, -1000,
	-1000, -1000, 52, -1000, 44, 1, 39, -1000, 12, 19,
	58, -1000, -1000, -1000, -7, 58, -1000, 37, -1, -1,
	9, 22, 36, -1000, -1, 58, 33, -1000, 29, 4,
	-1, -1, -1, -3, -4, 58, -1000, -1000, -1000, -1,
	-14, -1000, 4, -1000,
}
var yyPgo = [...]int{

	0, 70, 64, 69, 68, 2, 1, 67, 0, 15,
}
var yyR1 = [...]int{

	0, 1, 1, 8, 8, 9, 9, 3, 3, 5,
	5, 6, 6, 6, 6, 7, 7, 7, 4, 4,
	4, 4, 2, 2,
}
var yyR2 = [...]int{

	0, 1, 2, 1, 1, 1, 2, 2, 8, 3,
	1, 4, 4, 4, 1, 5, 2, 1, 3, 5,
	2, 1, 5, 2,
}
var yyChk = [...]int{

	-1000, -1, -2, 7, -2, -3, 6, 8, -9, -8,
	4, 5, -9, -8, 10, 9, 11, -4, -5, 11,
	26, 27, -6, -7, 23, 15, -8, 12, 21, 22,
	19, 16, 12, -8, 25, -9, 13, -6, 11, -5,
	20, 18, 17, 13, -5, 14, -6, -6, -6, 22,
	24, -8, -5, 25,
}
var yyDef = [...]int{

	0, -2, 1, 0, 2, 23, 0, 0, 7, 5,
	3, 4, 0, 6, 0, 0, 0, 22, 0, 0,
	0, 21, 10, 14, 0, 0, 17, 0, 0, 0,
	0, 0, 0, 20, 0, 16, 0, 9, 0, 18,
	0, 0, 0, 0, 0, 0, 11, 12, 13, 0,
	0, 8, 19, 15,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27,
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
		//line syntax.y:41
		{
			yyVAL.program = make(Program)
			if err := yyVAL.program.AddLine(yyDollar[1].line.number, yyDollar[1].line.stmt, yyDollar[1].line.goto0, yyDollar[1].line.goto1); err != nil {
				panic(err)
			}
			yylex.(*lex).prog = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:49
		{
			yyVAL.program = yyDollar[1].program
			if err := yyVAL.program.AddLine(yyDollar[2].line.number, yyDollar[2].line.stmt, yyDollar[2].line.goto0, yyDollar[2].line.goto1); err != nil {
				panic(err)
			}
			yylex.(*lex).prog = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:60
		{
			yyVAL.number = 0
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:64
		{
			yyVAL.number = 1
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:71
		{
			yyVAL.numberbits.number = yyDollar[1].number
			yyVAL.numberbits.bits = 1
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:76
		{
			yyVAL.numberbits.number = yyDollar[1].numberbits.number<<1 | yyDollar[2].number
			yyVAL.numberbits.bits = yyDollar[1].numberbits.bits + 1
			if yyVAL.numberbits.bits == 64 {
				panic(fmt.Errorf("bit: integer overflow"))
			}
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:87
		{
			yyVAL.line.goto0 = new(uint64)
			yyVAL.line.goto1 = new(uint64)
			*yyVAL.line.goto0, *yyVAL.line.goto1 = yyDollar[2].numberbits.number, yyDollar[2].numberbits.number
		}
	case 8:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line syntax.y:93
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
		//line syntax.y:108
		{
			yyVAL.expr = NandExpr{yyDollar[1].expr, yyDollar[3].expr}
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:112
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 11:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line syntax.y:119
		{
			yyVAL.expr = AddrExpr{yyDollar[4].expr}
		}
	case 12:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line syntax.y:123
		{
			yyVAL.expr = NextExpr{yyDollar[4].expr, 0}
		}
	case 13:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line syntax.y:127
		{
			yyVAL.expr = StarExpr{yyDollar[4].expr}
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:131
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 15:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line syntax.y:138
		{
			yyVAL.expr = yyDollar[3].expr
		}
	case 16:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:142
		{
			yyVAL.expr = VarExpr(yyDollar[2].numberbits.number)
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:146
		{
			if yyDollar[1].number == 0 {
				yyVAL.expr = BitExpr(false)
			} else {
				yyVAL.expr = BitExpr(true)
			}
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:157
		{
			yyVAL.stmt = AssignStmt{yyDollar[1].expr, yyDollar[3].expr}
		}
	case 19:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line syntax.y:161
		{
			yyVAL.stmt = JumpRegisterStmt{yyDollar[5].expr}
		}
	case 20:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:165
		{
			if yyDollar[2].number == 0 {
				yyVAL.stmt = PrintStmt(false)
			} else {
				yyVAL.stmt = PrintStmt(true)
			}
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:173
		{
			yyVAL.stmt = ReadStmt{}
		}
	case 22:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line syntax.y:180
		{
			yyVAL.line.number, yyVAL.line.stmt = yyDollar[3].numberbits.number, yyDollar[5].stmt
			yyVAL.line.goto0, yyVAL.line.goto1 = nil, nil
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:185
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

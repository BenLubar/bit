package bit

import (
	"errors"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

var tokens = map[int][]rune{
	ZERO:        []rune("ZERO"),
	ONE:         []rune("ONE"),
	GOTO:        []rune("GOTO"),
	LINE:        []rune("LINE"),
	NUMBER:      []rune("NUMBER"),
	CODE:        []rune("CODE"),
	IF:          []rune("IF"),
	THE:         []rune("THE"),
	JUMP:        []rune("JUMP"),
	REGISTER:    []rune("REGISTER"),
	IS:          []rune("IS"),
	VARIABLE:    []rune("VARIABLE"),
	VALUE:       []rune("VALUE"),
	AT:          []rune("AT"),
	BEYOND:      []rune("BEYOND"),
	ADDRESS:     []rune("ADDRESS"),
	OF:          []rune("OF"),
	NAND:        []rune("NAND"),
	EQUALS:      []rune("EQUALS"),
	OPEN:        []rune("OPEN"),
	CLOSE:       []rune("CLOSE"),
	PARENTHESIS: []rune("PARENTHESIS"),
	PRINT:       []rune("PRINT"),
	READ:        []rune("READ"),
}

var ErrInvalidUnicode = errors.New("bit: invalid unicode")

type lex struct {
	r    io.RuneReader
	line int
	col  int
	off  int
	buf  []rune
	prog Program
}

func (l *lex) Lex(lcal *yySymType) int {
	l.buf = l.buf[:0]

	for {
		r, size, err := l.r.ReadRune()
		if err == io.EOF {
			if len(l.buf) == 0 {
				return 0
			}
			panic(io.ErrUnexpectedEOF)
		}
		if err != nil {
			panic(err)
		}
		l.off += size
		l.col++
		if r == utf8.RuneError && size == 1 {
			panic(ErrInvalidUnicode)
		}

		if unicode.IsSpace(r) {
			if r == '\n' {
				l.line++
				l.col = 0
			}
			continue
		}

		l.buf = append(l.buf, r)
		any := false
	tokenLoop:
		for id, token := range tokens {
			if len(token) < len(l.buf) {
				continue
			}
			for i := range l.buf {
				if token[i] != l.buf[i] {
					continue tokenLoop
				}
			}
			any = true
			if len(token) == len(l.buf) {
				return id
			}
		}

		if !any {
			panic(fmt.Errorf("bit: unknown token %q", l.buf))
		}
	}
}

func (l *lex) Error(s string) {
	panic(errors.New(s))
}

type ParseError struct {
	Err    error
	Line   int
	Column int
	Offset int
}

func (err *ParseError) Error() string {
	return fmt.Sprintf("bit: on line %d column %d (offset %d): %v", err.Line, err.Column, err.Offset, err.Err)
}

func Parse(r io.RuneReader) (prog Program, err error) {
	l := &lex{r: r, line: 1}

	defer func() {
		/*if r := recover(); r != nil {
			err = r.(error)
			err = &ParseError{
				Err:    err,
				Line:   l.line,
				Column: l.col,
				Offset: l.off,
			}
		}*/
	}()

	yyParse(l)

	return l.prog.bake()
}

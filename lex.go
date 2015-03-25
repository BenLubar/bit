package bit

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

var tokens = map[string]int{
	"ZERO":        ZERO,
	"ONE":         ONE,
	"GOTO":        GOTO,
	"LINE":        LINE,
	"NUMBER":      NUMBER,
	"CODE":        CODE,
	"IF":          IF,
	"THE":         THE,
	"JUMP":        JUMP,
	"REGISTER":    REGISTER,
	"IS":          IS,
	"VARIABLE":    VARIABLE,
	"VALUE":       VALUE,
	"AT":          AT,
	"BEYOND":      BEYOND,
	"ADDRESS":     ADDRESS,
	"OF":          OF,
	"NAND":        NAND,
	"EQUALS":      EQUALS,
	"OPEN":        OPEN,
	"CLOSE":       CLOSE,
	"PARENTHESIS": PARENTHESIS,
	"PRINT":       PRINT,
	"READ":        READ,
}

var ErrInvalidUnicode = errors.New("bit: invalid unicode")

type lex struct {
	r    io.RuneReader
	line int
	col  int
	off  int
	prog Program
}

func (l *lex) Lex(lcal *yySymType) int {
	var s string

	for {
		r, size, err := l.r.ReadRune()
		if err == io.EOF {
			if s == "" {
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

		s += string(r)
		any := false
		for k, v := range tokens {
			if k == s {
				return v
			} else if strings.HasPrefix(k, s) {
				any = true
			}
		}

		if !any {
			panic(fmt.Errorf("bit: unknown token %q", s))
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
		if r := recover(); r != nil {
			err = r.(error)
			err = &ParseError{
				Err:    err,
				Line:   l.line,
				Column: l.col,
				Offset: l.off,
			}
		}
	}()

	yyParse(l)

	return l.prog, nil
}

package bit

import (
	"errors"
	"fmt"
	"io"
)

type lex struct {
	r    io.Reader
	read []byte
	buf  [4096]byte
	line int
	col  int
	off  int
	prog *Program
}

func (l *lex) Lex(lval *yySymType) int {
	for {
		for len(l.read) == 0 {
			n, err := l.r.Read(l.buf[:])
			if err != nil {
				if err == io.EOF {
					return 0
				}
				panic(err)
			}
			l.read = l.buf[:n]
		}

		var b byte
		b, l.read = l.read[0], l.read[1:]
		l.off++
		l.col++
		if b == '\n' {
			l.col = 0
			l.line++
			continue
		}
		if b == ' ' || b == '\t' {
			continue
		}
		return int(b)
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

func Parse(r io.Reader) (prog Program, err error) {
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

	err = l.prog.Init()
	if err != nil {
		return
	}
	return *l.prog, nil
}

package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/BenLubar/bit/cmd/brainfuckc/bf"
)

type lex struct {
	r    io.ByteReader
	line int
	col  int
	off  int
	prog []bf.Command
}

func (l *lex) Lex(lval *yySymType) int {
	for {
		b, err := l.r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return 0
			}
			panic(err)
		}
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
	return fmt.Sprintf("ook: on line %d column %d (offset %d): %v", err.Line, err.Column, err.Offset, err.Err)
}

func Parse(r io.ByteReader) (cmds []bf.Command, err error) {
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

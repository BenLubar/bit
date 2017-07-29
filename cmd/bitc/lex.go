package main

import (
	"bufio"
	"errors"
	"io"
	"unicode"
)

type lex struct {
	r       *bufio.Reader
	ext     bool
	program *Program
	line    int
	col     int
}

type lexError struct {
	Err  error
	Line int
	Col  int
}

func (l *lex) Lex(lval *yySymType) int {
	for {
		b, err := l.r.ReadByte()
		if err == io.EOF {
			return -1
		}
		if err != nil {
			panic(lexError{err, l.line + 1, l.col})
		}
		l.col++
		if !unicode.IsSpace(rune(b)) {
			return int(b)
		}
		if b == '\n' {
			l.line++
			l.col = 0
		}
	}
}

func (l *lex) Error(e string) {
	panic(lexError{errors.New(e), l.line + 1, l.col})
}

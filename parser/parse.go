//go:generate go get -u golang.org/x/tools/cmd/goyacc
//go:generate goyacc -o syntax.go syntax.y

package parser

import (
	"github.com/BenLubar/bit/ast"
	"github.com/BenLubar/bit/token"
)

func init() {
	yyErrorVerbose = true
}

type parseError struct {
	message string
}

func (err parseError) Error() string {
	return err.message
}

func Parse(tokens []token.Token) (prog *ast.Program, err error) {
	ch := make(chan token.Token, len(tokens))
	for _, t := range tokens {
		ch <- t
	}
	close(ch)

	p := &parser{tokens: ch, prog: &ast.Program{}}
	yyNewParser().Parse(p)

	return p.prog, nil
}

func ParseChan(tokens <-chan token.Token, lines chan<- *ast.Line) (err error) {
	defer close(lines)

	defer func() {
		if r := recover(); r != nil {
			err = r.(parseError)
		}
	}()

	p := &parser{tokens: tokens, lines: lines}
	yyNewParser().Parse(p)

	return nil
}

type parser struct {
	tokens <-chan token.Token
	lines  chan<- *ast.Line
	prog   *ast.Program
}

var tokenMap = [...]int{
	token.AddressOf:    tAddressOf,
	token.At:           tAt,
	token.Beyond:       tBeyond,
	token.Close:        tClose,
	token.Code:         tCode,
	token.Equals:       tEquals,
	token.Goto:         tGoto,
	token.If:           tIf,
	token.Is:           tIs,
	token.JumpRegister: tJumpRegister,
	token.LineNumber:   tLineNumber,
	token.Nand:         tNand,
	token.One:          tOne,
	token.Open:         tOpen,
	token.Parenthesis:  tParenthesis,
	token.Print:        tPrint,
	token.Read:         tRead,
	token.The:          tThe,
	token.Value:        tValue,
	token.Variable:     tVariable,
	token.Zero:         tZero,
}

func (p *parser) line(l *ast.Line) {
	if p.prog != nil {
		p.prog.Lines = append(p.prog.Lines, l)
	} else {
		p.lines <- l
	}
}

func (p *parser) Lex(*yySymType) int {
	if t, ok := <-p.tokens; ok {
		return tokenMap[t]
	}

	return -1
}

func (p *parser) Error(message string) {
	panic(parseError{message})
}

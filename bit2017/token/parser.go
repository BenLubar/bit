package token

import (
	"errors"
	"strconv"
	"unicode"
)

// Parser converts a stream of bytes or runes into a []Token.
//
// It is not safe for concurrent use across multiple goroutines; however,
// Parser does not modify any global or shared state, so multiple independent
// instances of Parser are safe to use concurrently.
//
// Methods that return error can return ErrSyntax or ErrInvalidChar. If
// ErrSyntax is returned by a Write* method, the current token is discarded
// to assist with tokenization error recovery.
type Parser struct {
	// Tokens are appended to the slice by writing to the parser. The
	// Parser does not change or access this slice in any other way.
	Tokens []Token

	state *state
	index int
}

// ErrSyntax is returned by Parser methods when a token is truncated or invalid.
var ErrSyntax = errors.New("bit2017/token: syntax error")

// ErrInvalidChar is returned by Parser methods when a character that is not an
// uppercase ASCII letter or a space is written.
type ErrInvalidChar struct {
	Char rune
}

// Error implements the error interface.
func (err ErrInvalidChar) Error() string {
	return "bit2017/token: invalid character: " + strconv.QuoteRune(err.Char)
}

// Write wraps the WriteByte method of Parser.
func (p *Parser) Write(b []byte) (n int, err error) {
	for i, c := range b {
		if err := p.WriteByte(c); err != nil {
			return i, err
		}
	}

	return len(b), nil
}

// WriteString wraps the WriteRune method of Parser.
func (p *Parser) WriteString(str string) (n int, err error) {
	for i, r := range str {
		if err := p.WriteRune(r); err != nil {
			return i, err
		}
	}

	return len(str), nil
}

func (p *Parser) write(c byte) error {
	if p.state == nil {
		p.state = &parser
	}

	if p.index < len(p.state.rest) {
		if p.state.rest[p.index] == c {
			p.index++
		} else {
			p.state = nil
			p.index = 0
			return ErrSyntax
		}

		for p.index < len(p.state.rest) && p.state.rest[p.index] == ' ' {
			p.index++
		}
	} else {
		if p.state.choice[c-'A'] == (state{}) {
			p.state = nil
			p.index = 0
			return ErrSyntax
		}

		p.state = &p.state.choice[c-'A']
		p.index = 0
	}

	if p.index == len(p.state.rest) {
		p.Tokens = append(p.Tokens, p.state.token)
		p.state = nil
		p.index = 0
	}

	return nil
}

// WriteByte adds a byte to the Parser input. If the byte is not an uppercase
// ASCII letter or a space character as defined by Unicode, this returns
// ErrInvalidChar. If an invalid token would be created by adding this byte,
// ErrSyntax is returned and the current token is discarded.
func (p *Parser) WriteByte(c byte) error {
	if unicode.IsSpace(rune(c)) {
		return nil
	}

	if c < 'A' || c > 'Z' {
		return ErrInvalidChar{Char: rune(c)}
	}

	return p.write(c)
}

// WriteRune supports Unicode but has the same restrictions and behavior as
// WriteByte. This method only differs in that it accepts non-ASCII whitespace
// and the ErrInvalidChar it returns will contain a full unicode code point.
func (p *Parser) WriteRune(c rune) error {
	if unicode.IsSpace(c) {
		return nil
	}

	if c < 'A' || c > 'Z' {
		return ErrInvalidChar{Char: c}
	}

	return p.write(byte(c))
}

// Done returns ErrSyntax if a partial token is buffered. Otherwise, it returns
// nil.
func (p *Parser) Done() error {
	if p.state == nil {
		return nil
	}

	return ErrSyntax
}

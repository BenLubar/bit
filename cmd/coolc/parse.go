//go:generate go tool yacc syntax.y

package main

import (
	"bytes"
	"fmt"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"unicode"
)

func init() {
	yyErrorVerbose = true
}

func (ast *AST) ParseFile(name string) (err error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return
	}

	if ast.FileSet == nil {
		ast.FileSet = token.NewFileSet()
	}
	file := ast.FileSet.AddFile(name, -1, len(b))
	file.SetLinesForContent(b)

	lex := &lexer{ast, file, bytes.NewReader(b), nil, -1}
	yyNewParser().Parse(lex)
	return lex.err
}

type lexer struct {
	ast    *AST
	file   *token.File
	r      *bytes.Reader
	err    error
	offset int64
}

var (
	errorSentinel = new(byte)
	eofSentinel   = new(byte)
	idTokens      = map[string]int{
		"case":     tokCASE,
		"class":    tokCLASS,
		"def":      tokDEF,
		"else":     tokELSE,
		"extends":  tokEXTENDS,
		"false":    tokFALSE,
		"if":       tokIF,
		"match":    tokMATCH,
		"new":      tokNEW,
		"null":     tokNULL,
		"override": tokOVERRIDE,
		"super":    tokSUPER,
		"this":     tokTHIS,
		"true":     tokTRUE,
		"var":      tokVAR,
		"while":    tokWHILE,

		"native": tokINVALID,

		"abstract":  tokINVALID,
		"catch":     tokINVALID,
		"do":        tokINVALID,
		"final":     tokINVALID,
		"finally":   tokINVALID,
		"for":       tokINVALID,
		"forSome":   tokINVALID,
		"implicit":  tokINVALID,
		"import":    tokINVALID,
		"lazy":      tokINVALID,
		"object":    tokINVALID,
		"package":   tokINVALID,
		"private":   tokINVALID,
		"protected": tokINVALID,
		"requires":  tokINVALID,
		"return":    tokINVALID,
		"sealed":    tokINVALID,
		"throw":     tokINVALID,
		"trait":     tokINVALID,
		"try":       tokINVALID,
		"type":      tokINVALID,
		"val":       tokINVALID,
		"with":      tokINVALID,
		"yield":     tokINVALID,
	}
)

func (l *lexer) Lex(lvalue *yySymType) (tok int) {
	defer func() {
		if r := recover(); r != nil {
			if r == errorSentinel {
				tok = tokINVALID
				return
			}
			if r == eofSentinel {
				tok = 0
				return
			}
			panic(r)
		}
	}()

	unexpected := false

	check := func(err error) {
		if err != nil {
			if err == io.EOF {
				if unexpected {
					err = io.ErrUnexpectedEOF
				} else {
					panic(eofSentinel)
				}
			}
			if l.err == nil {
				l.err = err
			}
			panic(errorSentinel)
		}
	}

	l.offset = -1
	offset, err := l.r.Seek(0, os.SEEK_CUR)
	check(err)
	l.offset = offset

	r, _, err := l.r.ReadRune()
	check(err)

	for {
		offset, err = l.r.Seek(0, os.SEEK_CUR)
		check(err)
		l.offset = offset

		if unicode.IsSpace(r) {
			r, _, err = l.r.ReadRune()
			check(err)
			continue
		}

		if r == '/' {
			r, _, err = l.r.ReadRune()
			switch r {
			case '/':
				for {
					r, _, err = l.r.ReadRune()
					check(err)
					if r == '\n' {
						break
					}
				}
				continue
			case '*':
				for {
					r, _, err = l.r.ReadRune()
					check(err)
					if r == '*' {
						r, _, err = l.r.ReadRune()
						check(err)
						if r == '/' {
							break
						}
						check(l.r.UnreadRune())
					}
				}
				continue
			default:
				r = '/'
				check(l.r.UnreadRune())
			}
		}

		break
	}

	unexpected = true

	validIdentifier := func(r rune) bool {
		return r == '_' || unicode.IsUpper(r) || unicode.IsLower(r) || (r >= '0' && r <= '9')
	}

	var buf []rune
	if unicode.IsUpper(r) {
		buf = append(buf, r)
		for {
			r, _, err = l.r.ReadRune()
			check(err)
			if validIdentifier(r) {
				buf = append(buf, r)
			} else {
				check(l.r.UnreadRune())
				lvalue.typ.Pos = l.file.Pos(int(offset))
				lvalue.typ.Name = string(buf)
				return tokTYPE
			}
		}
	}
	if unicode.IsLower(r) {
		buf = append(buf, r)
		for {
			r, _, err = l.r.ReadRune()
			check(err)
			if validIdentifier(r) {
				buf = append(buf, r)
			} else {
				check(l.r.UnreadRune())
				s := string(buf)
				if tok, ok := idTokens[s]; ok {
					return tok
				}
				lvalue.id.Pos = l.file.Pos(int(offset))
				lvalue.id.Name = s
				return tokID
			}
		}
	}
	if r == '0' {
		r, _, err = l.r.ReadRune()
		check(err)
		check(l.r.UnreadRune())

		if r >= '0' && r <= '9' {
			return tokINVALID
		}

		lvalue.n = 0
		return tokINTEGER
	}
	if r >= '1' && r <= '9' {
		buf = append(buf, r)

		for {
			r, _, err = l.r.ReadRune()
			check(err)
			if r >= '0' && r <= '9' {
				buf = append(buf, r)
			} else {
				check(l.r.UnreadRune())
				var n int64
				n, err = strconv.ParseInt(string(buf), 10, 32)
				check(err)
				lvalue.n = int32(n)
				return tokINTEGER
			}
		}
	}
	if r == '"' {
		r, _, err = l.r.ReadRune()
		check(err)
		if r == '"' {
			r, _, err = l.r.ReadRune()
			check(err)
			if r == '"' {
				for {
					r, _, err = l.r.ReadRune()
					check(err)
					buf = append(buf, r)

					if len(buf) >= 3 && buf[len(buf)-3] == '"' && buf[len(buf)-2] == '"' && buf[len(buf)-1] == '"' {
						lvalue.s = string(buf[:len(buf)-3])
						return tokSTRING
					}
				}
			}
			check(l.r.UnreadRune())
			lvalue.s = ""
			return tokSTRING
		}

		for {
			switch r {
			case '\n':
				return tokINVALID

			case '"':
				lvalue.s = string(buf)
				return tokSTRING

			case '\\':
				r, _, err = l.r.ReadRune()
				check(err)
				switch r {
				case '0':
					buf = append(buf, 0)
				case 'b':
					buf = append(buf, '\b')
				case 't':
					buf = append(buf, '\t')
				case 'n':
					buf = append(buf, '\n')
				case 'r':
					buf = append(buf, '\r')
				case 'f':
					buf = append(buf, '\f')
				case '"':
					buf = append(buf, '"')
				case '\\':
					buf = append(buf, '\\')
				default:
					return tokINVALID
				}

			default:
				buf = append(buf, r)
			}

			r, _, err = l.r.ReadRune()
			check(err)
		}
	}
	switch r {
	case '(':
		return tokLPAREN
	case ')':
		return tokRPAREN
	case ':':
		return tokCOLON
	case ',':
		return tokCOMMA
	case '{':
		return tokLBRACE
	case '}':
		return tokRBRACE
	case ';':
		return tokSEMICOLON
	case '.':
		return tokDOT
	case '!':
		return tokNEGATE
	case '*':
		return tokMULTIPLY
	case '/':
		return tokDIVIDE
	case '+':
		return tokPLUS
	case '-':
		return tokMINUS
	case '=':
		r, _, err = l.r.ReadRune()
		check(err)
		switch r {
		case '=':
			lvalue.id.Pos = l.file.Pos(int(offset))
			lvalue.id.Name = "equals"
			return tokEQUALEQUAL
		case '>':
			return tokARROW
		default:
			check(l.r.UnreadRune())
			return tokASSIGN
		}
	case '<':
		r, _, err = l.r.ReadRune()
		check(err)
		switch r {
		case '=':
			return tokLESSEQUAL
		default:
			check(l.r.UnreadRune())
			return tokLESSTHAN
		}
	default:
		return tokINVALID
	}
}

func (l *lexer) Error(s string) {
	if l.err == nil {
		if l.offset == -1 {
			l.err = fmt.Errorf("%s at %s:?", s, l.file.Name())
		} else {
			pos := l.file.Position(l.file.Pos(int(l.offset)))
			l.err = fmt.Errorf("%s at %v", s, pos)
		}
	}
}

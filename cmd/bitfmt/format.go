package main

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/BenLubar/bit/token"
)

func formatFile(name string) (string, bool) {
	fatal := func(err error) {
		if err != nil {
			log.Fatalf("%s: %v", name, err)
		}
	}

	f, err := os.Open(name)
	fatal(err)
	defer func() {
		fatal(f.Close())
	}()

	var (
		p   token.Parser
		br  = bufio.NewReader(f)
		buf []byte
	)

	var sinceLastToken []byte

	line := 1
	noWarnings := true
	var lastWarn error
	warn := func(err error) {
		if lastWarn == err {
			return
		}

		lastWarn = err
		noWarnings = false
		log.Printf("%s: %v on line %d", name, err, line)
	}

	for {
		b, err := br.ReadByte()
		if err == io.EOF {
			break
		}
		fatal(err)

		if b == '\n' {
			line++
		}

		if err = p.WriteByte(b); err != nil {
			if err == token.ErrSyntax {
				buf = append(buf, sinceLastToken...)
				buf = append(buf, b)
				sinceLastToken = sinceLastToken[:0]
			} else {
				buf = append(buf, b)
			}
			warn(err)
		} else if len(p.Tokens) != 0 {
			for _, t := range p.Tokens {
				if len(buf) != 0 {
					if t == token.LineNumber {
						buf = append(buf, '\n')
					} else {
						buf = append(buf, ' ')
					}
				}
				buf = append(buf, t.String()...)
			}
			p.Tokens = p.Tokens[:0]
			sinceLastToken = sinceLastToken[:0]
		} else {
			sinceLastToken = append(sinceLastToken, b)
		}
	}

	if err = p.Done(); err != nil {
		warn(err)
		buf = append(buf, sinceLastToken...)
	} else {
		buf = append(buf, '\n')
	}

	return string(buf), noWarnings
}

package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/BenLubar/bit/bitgen"
)

type Writer struct {
	*bitgen.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{bitgen.NewWriter(w)}
}

func (w *Writer) Program(s string) {
	var l bitgen.Line
	var m string

	write := func(v string) {
		if m == "" {
			m = v
		} else if v == "" {
			w.PrintString(l, m, 0)
		} else {
			next := w.ReserveLine()
			w.PrintString(l, m, next)
			l, m = next, v
		}
	}

	for _, r := range s {
		switch r {
		case 'h', 'H':
			write(Hello)

		case 'q', 'Q':
			write(s)

		case '9':
			write(Beer)
		}
	}
	write("")
}

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	w := NewWriter(os.Stdout)

	w.Program(string(b))

	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
}

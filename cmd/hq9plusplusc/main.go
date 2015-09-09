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

func (w *Writer) Program(s string) (n int64, err error) {
	var l bitgen.Line
	var m string

	write := func(v string) {
		if err != nil {
			return
		}

		var nn int64
		if m == "" {
			m = v
		} else if v == "" {
			nn, err = w.PrintString(l, m, 0)
		} else {
			next := w.ReserveLine()
			nn, err = w.PrintString(l, m, next)
			l, m = next, v
		}
		n += nn
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
	return
}

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	w := NewWriter(os.Stdout)

	_, err = w.Program(string(b))
	if err != nil {
		log.Fatal(err)
	}

	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
}

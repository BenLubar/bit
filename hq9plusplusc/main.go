package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Writer struct {
	w *bufio.Writer
	n uint64
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{bufio.NewWriter(w), 0}
}

func (w *Writer) Close() error {
	w.w.WriteByte('\n')
	return w.w.Flush()
}

func (w *Writer) Print(s string) {
	for i := range s {
		w.byte(s[i])
	}
}

func (w *Writer) byte(b byte) {
	w.bit(b&0x80 != 0)
	w.bit(b&0x40 != 0)
	w.bit(b&0x20 != 0)
	w.bit(b&0x10 != 0)
	w.bit(b&0x08 != 0)
	w.bit(b&0x04 != 0)
	w.bit(b&0x02 != 0)
	w.bit(b&0x01 != 0)
}

func (w *Writer) bit(b bool) {
	if w.n != 0 {
		w.w.WriteString(" GOTO")
		w.number(w.n)
		w.w.WriteByte('\n')
	}
	w.w.WriteString("LINE NUMBER")
	w.number(w.n)
	w.w.WriteString(" CODE PRINT")
	if b {
		w.w.WriteString(" ONE")
	} else {
		w.w.WriteString(" ZERO")
	}
	w.n++
}

func (w *Writer) number(n uint64) {
	if n == 0 {
		w.w.WriteString(" ZERO")
	} else {
		w.numberBits(n)
	}
}

func (w *Writer) numberBits(n uint64) {
	if n == 0 {
		return
	}

	w.numberBits(n >> 1)
	if n&1 == 0 {
		w.w.WriteString(" ZERO")
	} else {
		w.w.WriteString(" ONE")
	}
}

func (w *Writer) Program(s string) {
	for _, r := range s {
		switch r {
		case 'h', 'H':
			w.Print(Hello)

		case 'q', 'Q':
			w.Print(s)

		case '9':
			w.Print(Beer)
		}
	}
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

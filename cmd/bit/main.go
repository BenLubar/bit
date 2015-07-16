package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/BenLubar/bit"
)

func main() {
	if len(os.Args) != 2 {
		usage()
	}

	prog := parse()

	// we can't use bufio because it buffers at least one byte.
	r := &inefficient{rw: os.Stdin}
	w := &inefficient{rw: os.Stdout}

	err := prog.RunByte(r, w)
	if err != nil {
		handle(err)
		panic("unreachable")
	}
}

type inefficient struct {
	rw io.ReadWriter
	b  [1]byte
}

func (i *inefficient) ReadByte() (c byte, err error) {
	_, err = io.ReadFull(i.rw, i.b[:])
	return i.b[0], err
}

func (i *inefficient) WriteByte(c byte) error {
	i.b[0] = c
	n, err := i.rw.Write(i.b[:])
	if err != nil && n == 0 {
		err = io.ErrShortWrite
	}
	return err
}

func handle(err error) {
	fmt.Fprintf(os.Stderr, "%v\n\n", err)
	usage()
}

func parse() bit.Program {
	f, err := os.Open(os.Args[1])
	if err != nil {
		handle(err)
		panic("unreachable")
	}
	defer func() {
		if err := f.Close(); err != nil {
			handle(err)
			panic("unreachable")
		}
	}()

	prog, err := bit.Parse(bufio.NewReader(f))
	if err != nil {
		handle(err)
		panic("unreachable")
	}

	return prog
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %q:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s filename.bit\n", os.Args[0])
	os.Exit(2)
}
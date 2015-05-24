package main

import (
	"bufio"
	"bytes"
	"os"

	"github.com/BenLubar/bit"
)

func Example() {
	f, err := os.Open("hello.bf")
	if err != nil {
		panic(err)
	}

	commands, err := Parse(Tokenize(f))
	f.Close()
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.Program(commands)
	err = w.Close()
	if err != nil {
		panic(err)
	}

	prog, err := bit.Parse(&buf)
	if err != nil {
		panic(err)
	}

	bw := bufio.NewWriter(os.Stdout)
	defer bw.Flush()

	err = prog.RunByte(bufio.NewReader(os.Stdin), bw)
	if err != nil {
		panic(err)
	}

	// Output:
	// Hello World!
}

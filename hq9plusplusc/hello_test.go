package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"

	"github.com/BenLubar/bit"
)

func Example() {
	b, err := ioutil.ReadFile("hello.hq9")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.Program(string(b))
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
	// Hello, World!
}

package bf

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
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	w := NewWriter(&buf)
	_, err = w.Program(commands)
	if err != nil {
		panic(err)
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}

	prog, err := bit.Parse(&buf)
	if err != nil {
		panic(err)
	}

	bw := bufio.NewWriter(os.Stdout)

	err = prog.RunByte(bufio.NewReader(os.Stdin), bw)
	if err != nil {
		panic(err)
	}

	err = bw.Flush()
	if err != nil {
		panic(err)
	}

	// Output:
	// Hello World!
}

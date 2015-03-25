package bit

import (
	"bufio"
	"os"
)

func ExampleHello() {
	f, err := os.Open("hello.bit")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	prog, err := Parse(bufio.NewReader(f))
	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)

	err = prog.RunByte(r, w)
	if err != nil {
		panic(err)
	}

	err = w.Flush()
	if err != nil {
		panic(err)
	}

	// Output: Hello world!
}

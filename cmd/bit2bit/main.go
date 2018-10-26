package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var flagDecode = flag.Bool("d", false, "convert from ONE ZERO to bits. if not specified, converts from bits to ONE ZERO.")

func main() {
	flag.Parse()
	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(2)
	}

	var err error
	if *flagDecode {
		err = convertOneZero(os.Stdout, os.Stdin)
	} else {
		err = convertBit(os.Stdout, os.Stdin)
	}

	if err != io.EOF {
		fmt.Fprintln(os.Stderr, "Fatal error:", err)
	}
}

// Command bitfmt is a gofmt-like BIT prettifier.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var flagWrite = flag.Bool("w", false, "write changes back to .bit files instead of printing them on STDOUT.")

func main() {
	log.SetFlags(0)

	flag.Usage = func() {
		log.Println("Usage: bitfmt [flags] [input files]")
		log.Println()
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()

	if flag.NArg() == 0 {
		log.Println("Error: no input files specified")
		flag.Usage()
	}

	allOK := true

	for _, name := range flag.Args() {
		text, ok := formatFile(name)
		if !ok {
			allOK = false
		}

		if *flagWrite {
			ioutil.WriteFile(name, []byte(text), 0644)
		} else {
			fmt.Print(text)
		}
	}

	if !allOK {
		os.Exit(1)
	}
}

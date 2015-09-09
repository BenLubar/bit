package main

import (
	"bufio"
	"log"
	"os"

	"github.com/BenLubar/bit/cmd/brainfuckc/bf"
)

func main() {
	yyErrorVerbose = true

	list, err := Parse(bufio.NewReader(os.Stdin))
	if err != nil {
		log.Fatal(err)
	}

	w := bf.NewWriter(os.Stdout)

	_, err = w.Program(list)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
}

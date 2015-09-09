package main

import (
	"log"
	"os"

	"github.com/BenLubar/bit/cmd/brainfuckc/bf"
)

func main() {
	list, err := bf.Parse(bf.Tokenize(os.Stdin))
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

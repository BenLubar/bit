package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %q:\ncoolc [ -o fileout ] file1.cool file2.cool ... filen.cool\n", os.Args[0])
		flag.PrintDefaults()
	}
	output := flag.String("o", "", "output filename (defaults to first file name with a .bit extension)")

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
		return
	}

	var ast AST

	haveError := false
	for _, name := range flag.Args() {
		if err := ast.ParseFile(name); err != nil {
			fmt.Fprintf(os.Stderr, "%v: %v\n", name, err)
			haveError = true
		}
	}
	if haveError {
		os.Exit(1)
	}

	if err := ast.Semantic(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	name := *output
	if name == "" {
		name = strings.TrimSuffix(flag.Arg(0), ".cool") + ".bit"
	}
	f, err := os.Create(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", name, err)
		os.Exit(3)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			// don't exit - allow panics to go through
		}
	}()

	if err = ast.WriteTo(f); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(3)
	}
}

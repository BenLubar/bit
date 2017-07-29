package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %q:\nbitc [ -o output.s ] input.bit\n", os.Args[0])
		flag.PrintDefaults()
	}
	flagOutput := flag.String("o", "", "output filename. by default, this is the input filename with .s at the end.")
	flagExtensions := flag.Bool("ext", false, "allow some BIT extensions")
	flagOptimize := flag.Bool("opt", false, "optimize generated assembly code")
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
		return
	}

	if *flagOutput == "" {
		*flagOutput = strings.TrimSuffix(flag.Arg(0), ".bit") + ".s"
	}

	var prog *Program

	func() {
		defer func() {
			if r := recover(); r != nil {
				err := r.(lexError)
				fmt.Fprintln(os.Stderr, "parsing error:", err.Err)
				fmt.Fprintln(os.Stderr, "line", err.Line, "col", err.Col)
				os.Exit(1)
			}
		}()

		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
			return
		}
		defer f.Close()

		l := &lex{r: bufio.NewReader(f), ext: *flagExtensions}

		yyErrorVerbose = true
		if yyParse(l) != 0 {
			fmt.Fprintln(os.Stderr, "parsing failed: syntax error")
			os.Exit(1)
		}

		prog = l.program
	}()

	prog.CheckLineNumbers()

	prog.FindPointerVariables()

	if *flagOptimize {
		prog.Optimize()
	}

	f, err := os.Create(*flagOutput)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
		return
	}
	defer f.Close()

	err = prog.Compile(Linux64AssemblyWriter{f})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
		return
	}
}

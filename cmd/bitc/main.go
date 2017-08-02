package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

var osArch = map[[2]string]func(w io.Writer) AssemblyWriter{
	{"linux", "386"}: func(w io.Writer) AssemblyWriter {
		return Linux32AssemblyWriter{w}
	},
	{"linux", "amd64"}: func(w io.Writer) AssemblyWriter {
		return Linux64AssemblyWriter{w}
	},
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %q:\nbitc [ -o output.s ] input.bit\n", os.Args[0])
		flag.PrintDefaults()
	}
	flagOutput := flag.String("o", "", "output filename. by default, this is the input filename with .s at the end.")
	flagExtensions := flag.Bool("ext", false, "allow some BIT extensions")
	flagOptimize := flag.Bool("opt", false, "optimize generated assembly code")
	flagVerbose := flag.Int("v", 0, "logging verbosity")
	flagOS := flag.String("os", "linux", "the operating system to generate assembly for")
	flagArch := flag.String("arch", "amd64", "the processor architecture to generate assembly for")
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

		if *flagVerbose > 2 {
			yyDebug = *flagVerbose - 2
			realStdout := os.Stdout
			os.Stdout = os.Stderr
			defer func() {
				os.Stdout = realStdout
			}()
		}

		l := &lex{r: bufio.NewReader(f), ext: *flagExtensions}

		yyErrorVerbose = true
		if yyParse(l) != 0 {
			fmt.Fprintln(os.Stderr, "parsing failed: syntax error")
			os.Exit(1)
			return
		}

		prog = l.program
	}()

	if *flagVerbose >= 1 {
		fmt.Fprintln(os.Stderr, "parsed", len(prog.Lines), "lines")
	}

	prog.CheckLineNumbers(*flagVerbose)

	prog.FindPointerVariables(*flagVerbose)

	if *flagOptimize {
		prog.Optimize(*flagVerbose)
	}

	if a, ok := osArch[[2]string{*flagOS, *flagArch}]; ok {
		f, err := os.Create(*flagOutput)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
			return
		}
		defer f.Close()

		err = prog.Compile(a(f))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
			return
		}
	} else {
		fmt.Fprintln(os.Stderr, "unsupported operating system or architecture:", *flagOS+"/"+*flagArch)
		var supported []string
		for k := range osArch {
			supported = append(supported, k[0]+"/"+k[1])
		}
		sort.Strings(supported)
		for _, s := range supported {
			fmt.Fprintln(os.Stderr, "supported:", s)
		}
		os.Exit(1)
		return
	}
}

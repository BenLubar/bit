package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %q:\ncoolc [ -o fileout ] file1.cool file2.cool ... filen.cool\n", os.Args[0])
		flag.PrintDefaults()
	}
	output := flag.String("o", "", "output filename (defaults to first file name with a .bit extension)")
	cpuProfile := flag.String("cpuprofile", "", "Write a PPROF CPU profile")
	memProfile := flag.String("memprofile", "", "Write a PPROF heap profile")

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
		return
	}

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()

		err = pprof.StartCPUProfile(f)
		if err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
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

	if err := ast.Semantic(false); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	ast.Optimize()

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

	if *memProfile != "" {
		f, err := os.Create(*memProfile)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()

		err = pprof.WriteHeapProfile(f)
		if err != nil {
			panic(err)
		}
	}
}

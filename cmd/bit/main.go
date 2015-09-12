package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/pprof"

	"github.com/BenLubar/bit"
)

var flagCPUProfile = flag.String("cpuprofile", "", "Write a PPROF CPU profile")
var flagMemProfile = flag.String("memprofile", "", "Write a PPROF heap profile")
var flagNoOpt = flag.Bool("no-opt", false, "Don't build intrinsic versions of common patterns")

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
	}

	if *flagCPUProfile != "" {
		f, err := os.Create(*flagCPUProfile)
		if err != nil {
			handle(err)
			panic("unreachable")
		}
		defer func() {
			if err := f.Close(); err != nil {
				handle(err)
				panic("unreachable")
			}
		}()

		err = pprof.StartCPUProfile(f)
		if err != nil {
			handle(err)
			panic("unreachable")
		}
		defer pprof.StopCPUProfile()
	}

	prog := parse()

	// we can't use bufio because it buffers at least one byte.
	r := &inefficient{rw: os.Stdin}
	w := &inefficient{rw: os.Stdout}

	err := prog.RunByte(r, w)
	if err != nil {
		handle(err)
		panic("unreachable")
	}

	if *flagMemProfile != "" {
		f, err := os.Create(*flagMemProfile)
		if err != nil {
			handle(err)
			panic("unreachable")
		}
		defer func() {
			if err := f.Close(); err != nil {
				handle(err)
				panic("unreachable")
			}
		}()

		err = pprof.WriteHeapProfile(f)
		if err != nil {
			handle(err)
			panic("unreachable")
		}
	}
}

type inefficient struct {
	rw io.ReadWriter
	b  [1]byte
}

func (i *inefficient) ReadByte() (c byte, err error) {
	_, err = io.ReadFull(i.rw, i.b[:])
	return i.b[0], err
}

func (i *inefficient) WriteByte(c byte) error {
	i.b[0] = c
	n, err := i.rw.Write(i.b[:])
	if err != nil && n == 0 {
		err = io.ErrShortWrite
	}
	return err
}

func handle(err error) {
	fmt.Fprintf(os.Stderr, "%v\n\n", err)
	usage()
}

func parse() bit.Program {
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		handle(err)
		panic("unreachable")
	}
	defer func() {
		if err := f.Close(); err != nil {
			handle(err)
			panic("unreachable")
		}
	}()

	prog, err := bit.Parse(f)
	if err != nil {
		handle(err)
		panic("unreachable")
	}

	if !*flagNoOpt {
		prog.Optimize()
	}

	return prog
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %q:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s [options] filename.bit\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

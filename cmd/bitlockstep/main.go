package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/BenLubar/bit"
	"github.com/BenLubar/bit/bitio"
	"github.com/BenLubar/bit/internal/bitdebug"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
	}

	prog0, prog1 := parse()

	ch := make(chan interface{})

	var buf0, buf1 bytes.Buffer

	continueTo := make(chan uint64, 1)
	resulting := make(chan interface{})
	ready := make(chan struct{})

	p := func(tag string, line interface{}) {
		num, opt, stmt := bitdebug.LineNum(line), bitdebug.LineOpt(line), bitdebug.LineStmt(line)
		if opt == nil {
			fmt.Printf("%s%d = %#v\n", tag, num, stmt)
		} else {
			fmt.Printf("%s%d = %#v\n", tag, num, opt)
		}
	}

	compare := func(ctx0, ctx1 interface{}) {
		if !reflect.DeepEqual(ctx0, ctx1) || !bytes.Equal(buf0.Bytes(), buf1.Bytes()) {
			spew.Config.SortKeys = true
			spew.Dump(ctx0)
			spew.Dump(buf0.String())
			spew.Dump(ctx1)
			spew.Dump(buf1.String())
			panic("expected equal")
		}
	}

	step0 := func(line, ctx interface{}) {
		n := <-continueTo

		if n == bitdebug.LineNum(line) {
			resulting <- ctx

			<-ready
		} else {
			continueTo <- n
		}

		p("\t", line)
	}

	step1 := func(line, ctx1 interface{}) {
		continueTo <- bitdebug.LineNum(line)

		ctx0 := <-resulting

		compare(ctx0, ctx1)

		p("\n", line)

		ready <- struct{}{}
	}

	go trace(prog0, &buf0, step0, ch)
	go trace(prog1, &buf1, step1, ch)

	ctx1 := <-ch
	select {
	case continueTo <- 0:
	default:
		panic("optimized terminated second")
	}
	ctx0 := <-ch

	compare(ctx0, ctx1)
}

func trace(prog bit.Program, buf *bytes.Buffer, step func(interface{}, interface{}), ch chan<- interface{}) {
	ctx, err := bitdebug.RunTrace(prog, nil, bitio.NewWriter(buf), step)
	if err != nil {
		handle(err)
		panic("unreachable")
	}

	ch <- ctx
}

func handle(err error) {
	fmt.Fprintf(os.Stderr, "%v\n\n", err)
	usage()
}

func parse() (prog0, prog1 bit.Program) {
	ch := make(chan bit.Program)

	for i := 0; i < 2; i++ {
		go func() {
			f, err := os.Open(flag.Arg(0))
			if err != nil {
				handle(err)
				panic("unreachable")
			}

			prog, err := bit.Parse(bufio.NewReader(f))
			if err != nil {
				handle(err)
				panic("unreachable")
			}

			if err = f.Close(); err != nil {
				handle(err)
				panic("unreachable")
			}

			ch <- prog
		}()
	}

	prog0, prog1 = <-ch, <-ch

	prog1.Optimize()

	return
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %q:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s [options] filename.bit\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

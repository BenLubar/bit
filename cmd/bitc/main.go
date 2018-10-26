package main

import (
	"bufio"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/BenLubar/bit/ast"
	"github.com/BenLubar/bit/parser"
	"github.com/BenLubar/bit/token"
)

var flagOutput = flag.String("o", "", "output file (required)")

func main() {
	flag.Parse()
	if flag.NArg() != 0 || *flagOutput == "" {
		flag.Usage()
		os.Exit(2)
	}

	tokens := make(chan token.Token, 1000)
	go readTokens(tokens, bufio.NewReader(os.Stdin))
	lines := make(chan *ast.Line, 100)
	go parseTokens(tokens, lines)

	prog := &ast.Program{}

	var lineNum, varNum intern
	for line := range lines {
		prog.Lines = append(prog.Lines, line)
		if _, ok := lineNum.find(line.Num); ok {
			log.Fatalln("duplicate line number:", line.Num)
		}
		lineNum.intern(line.Num)
		if eq, ok := line.Stmt.(*ast.Equals); ok {
			findVars(&varNum, eq.Left)
			findVars(&varNum, eq.Right)
		}
	}

	semProg := typeCheck(prog, &lineNum, &varNum)
	ofile, err := ioutil.TempFile("", "bit*.o")
	if err != nil {
		panic(err)
	}
	defer os.Remove(ofile.Name())
	defer func() {
		if err := ofile.Close(); err != nil {
			panic(err)
		}
	}()
	as := exec.Command("as", "-32", "-L", "-g", "-o", "/dev/fd/3", "--")
	as.Stdout = os.Stdout
	as.Stderr = os.Stderr
	as.ExtraFiles = append(as.ExtraFiles, ofile)
	asmIn, err := as.StdinPipe()
	if err != nil {
		panic(err)
	}
	go func() {
		defer asmIn.Close()
		codeGen(asmIn, semProg)
	}()

	if err := as.Run(); err != nil {
		panic(err)
	}

	ld := exec.Command("ld", "--discard-none", "-melf_i386", "-o", *flagOutput, "/dev/fd/3")
	ld.Stdout = os.Stdout
	ld.Stderr = os.Stderr
	ld.ExtraFiles = append(ld.ExtraFiles, ofile)
	if err := ld.Run(); err != nil {
		panic(err)
	}
}

func readTokens(ch chan<- token.Token, br io.ByteReader) {
	var tok token.Parser

	for {
		b, err := br.ReadByte()
		if err == io.EOF {
			if err = tok.Done(); err != nil {
				panic(err)
			}
			close(ch)
			return
		}
		if err != nil {
			panic(err)
		}

		if err = tok.WriteByte(b); err != nil {
			panic(err)
		}

		if len(tok.Tokens) != 0 {
			for _, t := range tok.Tokens {
				ch <- t
			}
			tok.Tokens = tok.Tokens[:0]
		}
	}
}

func parseTokens(tokens <-chan token.Token, lines chan<- *ast.Line) {
	if err := parser.ParseChan(tokens, lines); err != nil {
		panic(err)
	}
}

func findVars(varNum *intern, expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.ValueAt:
		findVars(varNum, e.Ptr)
	case *ast.ValueBeyond:
		findVars(varNum, e.Ptr)
	case *ast.AddressOf:
		findVars(varNum, e.Val)
	case *ast.Nand:
		findVars(varNum, e.Left)
		findVars(varNum, e.Right)
	case *ast.JumpRegister:
		// no variables
	case *ast.Variable:
		varNum.intern(e.Num)
	case *ast.Constant:
		// no variables
	default:
		log.Panicf("unexpected AST type %T", e)
	}
}

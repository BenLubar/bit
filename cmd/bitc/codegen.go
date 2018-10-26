package main

import (
	"bufio"
	"fmt"
	"io"
	"log"

	"github.com/BenLubar/bit/ast"
)

func codeGen(w io.Writer, semProg *semanticProgram) {
	bw := bufio.NewWriter(w)
	defer func() {
		if err := bw.Flush(); err != nil {
			panic(err)
		}
	}()

	printf := func(format string, args ...interface{}) {
		if _, err := fmt.Fprintf(bw, format, args...); err != nil {
			panic(err)
		}
	}

	printf("%s", runtimeASM)

	for i := 1; i < len(semProg.bits); i++ {
		printf(".globl var_%d\nvar_%d:\n.long 0\n", i, i)
	}

	printf(`.text
.globl _start
_start:
	xor %%edx, %%edx
	jmp .Lline_%d
.Lline_0:
	jmp exit
`, semProg.start)

	for i, line := range semProg.lines {
		if i == 0 {
			continue
		}
		genLine(printf, semProg, i, line)
	}
}

func genLine(printf func(string, ...interface{}), semProg *semanticProgram, i int, line semanticLine) {
	printf("\n\nline_%d:\n.Lline_%d:\n", i, i)

	switch s := line.stmt.(type) {
	case *ast.Read:
		printf("\tcall read\n")
	case *ast.Print:
		if s.Val {
			printf("\tcall print1\n")
		} else {
			printf("\tcall print0\n")
		}
	case *ast.Equals:
		genExpr(printf, semProg, s.Right)
		if _, ok := s.Left.(*ast.JumpRegister); ok {
			printf("\tmovb (%%eax), %%dl\n")
			break
		}

		if v, ok := s.Left.(*ast.Variable); ok {
			vn, _ := semProg.varNum.find(v.Num)
			if semProg.ptrs[vn] {
				printf("\tmovl %%eax, var_%d\n", vn)
				break
			}
		}

		printf("\tmovb (%%eax), %%al\n")
		printf("\tpush %%eax\n")
		genExpr(printf, semProg, s.Left)
		printf("\tpop %%ebx\n")
		printf("\tmovb %%bl, (%%eax)\n")
	}

	if line.goto0 == line.goto1 {
		printf("\tjmp .Lline_%d\n", line.goto0)
	} else {
		printf("\ttest %%edx, %%edx\n\tjz .Lline_%d\n\tjmp .Lline_%d\n", line.goto0, line.goto1)
	}
}

func genExpr(printf func(string, ...interface{}), semProg *semanticProgram, expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.ValueAt:
		genExpr(printf, semProg, e.Ptr)
		// no-op
	case *ast.ValueBeyond:
		off := 1
		ptr := e.Ptr
		for {
			switch v := ptr.(type) {
			case *ast.ValueAt:
				ptr = v.Ptr
				continue
			case *ast.ValueBeyond:
				ptr = v.Ptr
				off++
				continue
			case *ast.AddressOf:
				ptr = v.Val
				continue
			}
			break
		}
		genExpr(printf, semProg, ptr)
		printf("\tmov $%d, %%ebx\n", off)
		printf("\tcall incr\n")
	case *ast.AddressOf:
		genExpr(printf, semProg, e.Val)
		// no-op
	case *ast.Nand:
		panic("TODO: NAND")
	case *ast.Constant:
		if e.Val {
			printf("\tleal one, %%eax\n")
		} else {
			printf("\tleal zero, %%eax\n")
		}
	case *ast.Variable:
		vn, _ := semProg.varNum.find(e.Num)
		if semProg.bits[vn] {
			printf("\tleal var_%d, %%eax\n", vn)
			printf("\tcall alloc\n")
			printf("\tmovl (%%eax), %%eax\n")
		} else {
			printf("\tmovl var_%d, %%eax\n", vn)
		}
	case *ast.JumpRegister:
		log.Panicln("unexpected THE JUMP REGISTER used as expression")
	default:
		log.Panicf("unexpected expression type %T", e)
	}
}

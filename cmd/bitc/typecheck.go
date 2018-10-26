package main

import (
	"log"

	"github.com/BenLubar/bit/ast"
)

func typeCheck(prog *ast.Program, lineNum, varNum *intern) *semanticProgram {
	semProg := &semanticProgram{
		lineNum: lineNum,
		varNum:  varNum,
		lines:   make([]semanticLine, lineNum.count()+1),
		bits:    make([]bool, varNum.count()+1),
		ptrs:    make([]bool, varNum.count()+1),
	}
	semProg.start, _ = lineNum.find(prog.Lines[0].Num)

	for _, line := range prog.Lines {
		num, _ := lineNum.find(line.Num)
		goto0, ok := lineNum.find(line.Goto0)
		if !ok && line.Goto0 != nil {
			log.Fatalln("Invalid goto line", line.Goto0, "on line", line.Num)
		}
		goto1, ok := lineNum.find(line.Goto1)
		if !ok && line.Goto1 != nil {
			log.Fatalln("Invalid goto line", line.Goto1, "on line", line.Num)
		}

		semProg.lines[num] = semanticLine{
			stmt:  line.Stmt,
			goto0: goto0,
			goto1: goto1,
		}
	}

	found := true
	for found {
		found = false

		for i, line := range semProg.lines[1:] {
			if eq, ok := line.stmt.(*ast.Equals); ok {
				isVal0, isPtr0, change0 := typeCheckVars(varNum, eq.Left, semProg.bits, semProg.ptrs)
				isVal1, isPtr1, change1 := typeCheckVars(varNum, eq.Right, semProg.bits, semProg.ptrs)
				if isVal0 && isPtr1 {
					log.Fatalln("Cannot store the value of", eq.Right, "(a pointer) in", eq.Left, "(a bit) on line", lineNum.name(i+1))
				}
				if isPtr0 && isVal1 {
					log.Fatalln("Cannot store the value of", eq.Right, "(a bit) in", eq.Left, "(a pointer) on line", lineNum.name(i+1))
				}
				var change2 bool
				if isVal0 && !isVal1 {
					change2 = setVarType(varNum, eq.Right, semProg.bits)
				} else if !isVal0 && isVal1 {
					change2 = setVarType(varNum, eq.Left, semProg.bits)
				} else if isPtr0 && !isPtr1 {
					change2 = setVarType(varNum, eq.Right, semProg.ptrs)
				} else if !isPtr0 && isPtr1 {
					change2 = setVarType(varNum, eq.Left, semProg.ptrs)
				}
				found = found || change0 || change1 || change2
			}
		}
	}

	for i := 1; i < len(semProg.bits); i++ {
		if !semProg.bits[i] && !semProg.ptrs[i] {
			log.Fatalln("Cannot determine type of variable", varNum.name(i))
		} else if semProg.bits[i] && semProg.ptrs[i] {
			log.Panicln("internal error:", varNum.name(i), "was determined to be both a pointer and a bit")
		}
	}

	return semProg
}

type semanticProgram struct {
	start   int
	lineNum *intern
	varNum  *intern
	lines   []semanticLine
	bits    []bool
	ptrs    []bool
}

type semanticLine struct {
	stmt  ast.Stmt
	goto0 int
	goto1 int
}

func setVarType(varNum *intern, expr ast.Expr, is []bool) bool {
	v, ok := expr.(*ast.Variable)
	if !ok {
		return false
	}

	if vn, _ := varNum.find(v.Num); !is[vn] {
		is[vn] = true
		return true
	}

	return false
}

func typeCheckVars(varNum *intern, expr ast.Expr, bits, ptrs []bool) (isVal, isPtr, changed bool) {
	switch e := expr.(type) {
	case *ast.ValueAt:
		isVal, _, changed = typeCheckVars(varNum, e.Ptr, bits, ptrs)
		if isVal {
			log.Fatalln("Cannot access", e, "(a bit)")
		}
		isVal = true
		changed = changed || setVarType(varNum, e.Ptr, ptrs)
	case *ast.ValueBeyond:
		isVal, _, changed = typeCheckVars(varNum, e.Ptr, bits, ptrs)
		if isVal {
			log.Fatalln("Cannot access", e, "(a bit)")
		}
		isVal = true
		changed = changed || setVarType(varNum, e.Ptr, ptrs)
	case *ast.AddressOf:
		_, isPtr, changed = typeCheckVars(varNum, e.Val, bits, ptrs)
		if isPtr {
			log.Fatalln("Cannot access", e, "(a pointer)")
		}
		isPtr = true
		changed = changed || setVarType(varNum, e.Val, bits)
	case *ast.Nand:
		_, lp, lc := typeCheckVars(varNum, e.Left, bits, ptrs)
		if lp {
			log.Fatalln("Cannot use", e.Left, "(a pointer) as the left-hand side of NAND")
		}
		_, rp, rc := typeCheckVars(varNum, e.Right, bits, ptrs)
		if rp {
			log.Fatalln("Cannot use", e.Right, "(a pointer) as the right-hand side of NAND")
		}
		lv := setVarType(varNum, e.Left, bits)
		rv := setVarType(varNum, e.Right, bits)
		return true, false, lc || rc || lv || rv
	case *ast.JumpRegister:
		return true, false, false
	case *ast.Variable:
		vn, _ := varNum.find(e.Num)
		return bits[vn], ptrs[vn], false
	case *ast.Constant:
		return true, false, false
	default:
		log.Panicf("unexpected AST type %T", e)
	}

	return
}

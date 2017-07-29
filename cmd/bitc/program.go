package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BenLubar/bit/bitnum"
)

type Program struct {
	Lines       []*Line
	Pointers    bitnum.Numbers
	notPointers bitnum.Numbers
}

func (p *Program) String() string {
	var buf []byte

	for _, l := range p.Lines {
		buf = append(buf, '\n')
		buf = append(buf, l.String()...)
	}

	return string(buf[1:])
}

func (p *Program) CheckLineNumbers(verbosity int) {
	var lines bitnum.NumberMap
	for _, l := range p.Lines {
		if o, ok := lines.Get(l.Num); ok {
			fmt.Fprintln(os.Stderr, "multiple lines with the same line number:")
			fmt.Fprintln(os.Stderr, o.(*Line))
			fmt.Fprintln(os.Stderr, l)
			os.Exit(1)
		}

		lines.Set(l.Num, l)
	}

	var seen bitnum.Numbers
	queue := []*Line{p.Lines[0]}

	for len(queue) != 0 {
		l := queue[0]
		queue = queue[1:]
		if !seen.Add(l.Num) {
			continue
		}

		if r, ok := l.Stmt.(*ReadStmt); ok && r.EOFLine != nil {
			if v, ok := lines.Get(r.EOFLine); !ok {
				fmt.Fprintln(os.Stderr, "undefined read line number:", r.EOFLine)
				fmt.Fprintln(os.Stderr, l)
				os.Exit(1)
			} else {
				eof := v.(*Line)
				r.gotoEOF = eof
				queue = append(queue, eof)
			}
		}

		if l.Zero != nil {
			if v, ok := lines.Get(l.Zero); !ok {
				fmt.Fprintln(os.Stderr, "undefined goto line number:", l.Zero)
				fmt.Fprintln(os.Stderr, l)
				os.Exit(1)
			} else {
				zero := v.(*Line)
				l.gotoZero = zero
				queue = append(queue, zero)
			}
		}

		if l.One != nil {
			if v, ok := lines.Get(l.One); !ok {
				fmt.Fprintln(os.Stderr, "undefined goto line number:", l.One)
				fmt.Fprintln(os.Stderr, l)
				os.Exit(1)
			} else {
				one := v.(*Line)
				l.gotoOne = one
				queue = append(queue, one)
			}
		}
	}

	if verbosity >= 2 && len(lines.Keys) != len(seen) {
		fmt.Fprintln(os.Stderr, "unreachable code:")
		for _, v := range lines.Values {
			line := v.(*Line)
			if !seen.Contains(line.Num) {
				fmt.Fprintln(os.Stderr, line)
			}
		}
	}
}

func (p *Program) FindPointerVariables(verbosity int) {
	anyIdentified, anyUnknown := true, true
	for anyIdentified {
		anyIdentified, anyUnknown = false, false
		for _, l := range p.Lines {
			if e, ok := l.Stmt.(*EqualsStmt); ok {
				i, u := p.findPointerVariables(&e.Left, e.Right.Pointer() && !e.Right.Value(), !e.Right.Pointer() && e.Right.Value(), false)
				anyIdentified = anyIdentified || i
				anyUnknown = anyUnknown || u

				i, u = p.findPointerVariables(&e.Right, e.Left.Pointer() && !e.Left.Value(), !e.Left.Pointer() && e.Left.Value(), false)
				anyIdentified = anyIdentified || i
				anyUnknown = anyUnknown || u
			}
		}
	}

	if anyUnknown {
		anyIdentified = true
		for anyIdentified {
			anyIdentified, anyUnknown = false, false
			for _, l := range p.Lines {
				if e, ok := l.Stmt.(*EqualsStmt); ok {
					i, u := p.findPointerVariables(&e.Left, e.Right.Pointer() && !e.Right.Value(), !e.Right.Pointer() && e.Right.Value(), true)
					anyIdentified = anyIdentified || i
					anyUnknown = anyUnknown || u

					i, u = p.findPointerVariables(&e.Right, e.Left.Pointer() && !e.Left.Value(), !e.Left.Pointer() && e.Left.Value(), true)
					anyIdentified = anyIdentified || i
					anyUnknown = anyUnknown || u
				}
			}
		}
	}

	if anyUnknown {
		log.Panicln("could not identify all variables")
	}

	if verbosity >= 1 {
		fmt.Fprintln(os.Stderr, "identified", len(p.Pointers), "pointer and", len(p.notPointers), "bit variables")
	}

	return
}

func (p *Program) findPointerVariables(expr *Expr, ptr, val, guess bool) (bool, bool) {
	switch e := (*expr).(type) {
	case *PointerVariable, *BitVariable, *JumpRegister, *BitValue:
		return false, false
	case *UnknownVariable:
		if !ptr && !val && guess {
			ptr = p.Pointers.Contains(e.Num)
			val = p.notPointers.Contains(e.Num)
			if ptr == val {
				return false, true
			}
		}
		if ptr {
			p.Pointers.Add(e.Num)
			*expr = &PointerVariable{
				Num: e.Num,
			}
			return true, false
		}
		if val {
			p.notPointers.Add(e.Num)
			*expr = &BitVariable{
				Num: e.Num,
			}
			return true, false
		}
		return false, true
	case *ValueAt:
		return p.findPointerVariables(&e.Target, true, false, guess)
	case *AddressOf:
		return p.findPointerVariables(&e.Variable, false, true, guess)
	case *Nand:
		lefti, leftu := p.findPointerVariables(&e.Left, false, true, guess)
		righti, rightu := p.findPointerVariables(&e.Right, false, true, guess)
		return lefti || righti, leftu || rightu
	case *Parenthesis:
		return p.findPointerVariables(&e.Inner, ptr, val, guess)
	default:
		log.Panicln("internal compiler error: unhandled expression:", e)
		panic("unreachable")
	}
}

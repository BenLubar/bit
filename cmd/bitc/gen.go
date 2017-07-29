package main

import "log"

type AssemblyWriter interface {
	Runtime() error
	DataSegment() error
	TextSegment() error
	DeclarePointer(n *Number) error
	Start() error
	DeclareLine(n *Number) error
	Goto(zero, one *Number) error
	Read(eof *Number) error
	Print() error
	SaveRegister(register int) error
	LoadRegister(register int) error
	PointerAddress(register int, n *Number) error
	BitAddress(register int, n *Number) error
	JumpAddress(register int) error
	PointerValue(register int, n *Number) error
	ReadBitPointer(dest, src int) error
	JumpValue(register int) error
	BitValue(register int, bit bool) error
	WritePointer(dest, src int) error
	WriteBit(dest, src int) error
	Advance(register, offset int) error
	NandBit(dest, src int) error
}

func (p *Program) Compile(w AssemblyWriter) error {
	if err := w.Runtime(); err != nil {
		return err
	}

	if len(p.Pointers) != 0 {
		if err := w.DataSegment(); err != nil {
			return err
		}
		for _, n := range p.Pointers {
			if err := w.DeclarePointer(n); err != nil {
				return err
			}
		}
	}

	if err := w.TextSegment(); err != nil {
		return err
	}
	if err := w.Start(); err != nil {
		return err
	}
	for _, l := range p.Lines {
		if err := w.DeclareLine(l.Num); err != nil {
			return err
		}
		if err := p.compileStmt(w, l.Stmt); err != nil {
			return err
		}
		if err := w.Goto(l.Zero, l.One); err != nil {
			return err
		}
	}

	return nil
}

func (p *Program) compileStmt(w AssemblyWriter, stmt Stmt) error {
	switch s := stmt.(type) {
	case *ReadStmt:
		if err := w.Read(s.EOFLine); err != nil {
			return err
		}
	case *PrintStmt:
		if err := w.BitValue(0, s.Bit); err != nil {
			return err
		}
		if err := w.Print(); err != nil {
			return err
		}
	case *EqualsStmt:
		if s.Left.Pointer() && s.Right.Pointer() && s.Left.Value() && s.Right.Value() {
			log.Panicln("fatal compiler error: cannot determine if assignment is pointer or value:", s)
		} else if s.Left.Pointer() && s.Right.Pointer() {
			if err := p.compileLValue(w, s.Left); err != nil {
				return err
			}
			if err := w.SaveRegister(0); err != nil {
				return err
			}
			if err := p.compileRValue(w, s.Right); err != nil {
				return err
			}
			if err := w.LoadRegister(1); err != nil {
				return err
			}
			if err := w.WritePointer(1, 0); err != nil {
				return err
			}
		} else if s.Left.Value() && s.Right.Value() {
			if err := p.compileLValue(w, s.Left); err != nil {
				return err
			}
			if err := w.SaveRegister(0); err != nil {
				return err
			}
			if err := p.compileRValue(w, s.Right); err != nil {
				return err
			}
			if err := w.LoadRegister(1); err != nil {
				return err
			}
			if err := w.WriteBit(1, 0); err != nil {
				return err
			}
		} else {
			log.Panicln("fatal compiler error: cannot determine assignment type:", s)
		}
	default:
		log.Panicln("fatal compiler error: unexpected statement:", s)
	}
	return nil
}

func (p *Program) compileLValue(w AssemblyWriter, expr Expr) error {
	switch e := expr.(type) {
	case *PointerVariable:
		if err := w.PointerAddress(0, e.Num); err != nil {
			return err
		}
	case *BitVariable:
		if err := w.BitAddress(0, e.Num); err != nil {
			return err
		}
	case *JumpRegister:
		if err := w.JumpAddress(0); err != nil {
			return err
		}
	case *ValueAt:
		if err := p.compileRValue(w, e.Target); err != nil {
			return err
		}
		if e.Offset != 0 {
			if err := w.Advance(0, e.Offset); err != nil {
				return err
			}
		}
	case *AddressOf:
		log.Panicln("internal compiler error: invalid lvalue:", e)
	case *BitValue:
		log.Panicln("internal compiler error: invalid lvalue:", e)
	case *Nand:
		log.Panicln("internal compiler error: invalid lvalue:", e)
	case *Parenthesis:
		return p.compileLValue(w, e.Inner)
	default:
		log.Panicln("internal compiler error: unrecognized lvalue type:", e)
	}
	return nil
}

func (p *Program) compileRValue(w AssemblyWriter, expr Expr) error {
	switch e := expr.(type) {
	case *PointerVariable:
		if err := w.PointerValue(0, e.Num); err != nil {
			return err
		}
	case *BitVariable:
		if err := w.BitAddress(0, e.Num); err != nil {
			return err
		}
		if err := w.ReadBitPointer(0, 0); err != nil {
			return err
		}
	case *JumpRegister:
		if err := w.JumpValue(0); err != nil {
			return err
		}
	case *ValueAt:
		if err := p.compileRValue(w, e.Target); err != nil {
			return err
		}
		if err := w.ReadBitPointer(0, 0); err != nil {
			return err
		}
	case *AddressOf:
		if err := p.compileLValue(w, e.Variable); err != nil {
			return err
		}
	case *BitValue:
		if err := w.BitValue(0, e.Bit); err != nil {
			return err
		}
	case *Nand:
		if err := p.compileRValue(w, e.Left); err != nil {
			return err
		}
		if err := w.SaveRegister(0); err != nil {
			return err
		}
		if err := p.compileRValue(w, e.Right); err != nil {
			return err
		}
		if err := w.LoadRegister(1); err != nil {
			return err
		}
		if err := w.NandBit(0, 1); err != nil {
			return err
		}
	case *Parenthesis:
		return p.compileRValue(w, e.Inner)
	default:
		log.Panicln("internal compiler error: unrecognized rvalue type:", e)
	}
	return nil
}

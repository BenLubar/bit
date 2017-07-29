package main

import "log"

// AssemblyWriter is an interface that is used to write compiler output.
type AssemblyWriter interface {
	// Runtime writes any helper code that is shared between programs.
	Runtime() error
	// DataSegment informs the writer that the following function calls
	// will define mutable data.
	DataSegment() error
	// TextSegment informs the writer that the following function calls
	// will define executable code.
	TextSegment() error
	// DeclarePointer creates an address-of-a-bit variable.
	DeclarePointer(n *Number) error
	// Start marks the location where code will begin executing.
	Start() error
	// DeclareLine marks a line of BIT code.
	DeclareLine(n *Number) error
	// Goto jumps to a line depending on the value of the jump register.
	// Jumping to a nil line number exits the program.
	Goto(zero, one *Number) error
	// Read reads one bit from the standard input into the jump register.
	// If an error occurs, execution will resume on the line number given.
	Read(eof *Number) error
	// Print writes the bit in register 0 to the standard output.
	Print() error
	// SaveRegister pushes the register onto the stack.
	SaveRegister(register int) error
	// LoadRegister pops the register from the stack.
	LoadRegister(register int) error
	// PointerAddress puts the address of an address-of-a-bit variable into
	// the register.
	PointerAddress(register int, n *Number) error
	// BitAddress puts the address of a bit variable into the register.
	BitAddress(register int, n *Number) error
	// JumpAddress puts the address of the jump register into the register.
	JumpAddress(register int) error
	// PointerValue puts the value of an address-of-a-bit variable into the
	// register.
	PointerValue(register int, n *Number) error
	// ReadBit sets the dest register to the bit at the address in the src
	// register. The registers may be the same.
	ReadBit(dest, src int) error
	// JumpValue puts the value of the jump register into the register.
	JumpValue(register int) error
	// BitValue puts a constant bit value into the register.
	BitValue(register int, bit bool) error
	// WritePointer puts the address-of-a-bit from src into the
	// address-of-a-bit variable at the address in dest.
	WritePointer(dest, src int) error
	// WriteBit puts the bit from src into the bit variable at the address
	// in dest.
	WriteBit(dest, src int) error
	// Advance increases the address of a bit variable in the register by
	// the offset.
	Advance(register, offset int) error
	// NandBit sets dest to (src NAND dest).
	NandBit(dest, src int) error
}

type compileError struct{ Err error }

func checkErr(err error) {
	if err != nil {
		panic(compileError{err})
	}
}

func (p *Program) Compile(w AssemblyWriter) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(compileError).Err
		}
	}()

	checkErr(w.Runtime())

	if len(p.Pointers) != 0 {
		checkErr(w.DataSegment())
		for _, n := range p.Pointers {
			checkErr(w.DeclarePointer(n))
		}
	}

	checkErr(w.TextSegment())
	checkErr(w.Start())
	for _, l := range p.Lines {
		checkErr(w.DeclareLine(l.Num))
		p.compileStmt(w, l.Stmt)
		checkErr(w.Goto(l.Zero, l.One))
	}

	return nil
}

func (p *Program) compileStmt(w AssemblyWriter, stmt Stmt) {
	switch s := stmt.(type) {
	case *ReadStmt:
		checkErr(w.Read(s.EOFLine))
	case *PrintStmt:
		checkErr(w.BitValue(0, s.Bit))
		checkErr(w.Print())
	case *EqualsStmt:
		p.compileLValue(w, s.Left)
		checkErr(w.SaveRegister(0))
		p.compileRValue(w, s.Right)
		checkErr(w.LoadRegister(1))
		if s.Left.Pointer() && s.Right.Pointer() && s.Left.Value() && s.Right.Value() {
			log.Panicln("fatal compiler error: cannot determine if assignment is pointer or value:", s)
		} else if s.Left.Pointer() && s.Right.Pointer() {
			checkErr(w.WritePointer(1, 0))
		} else if s.Left.Value() && s.Right.Value() {
			checkErr(w.WriteBit(1, 0))
		} else {
			log.Panicln("fatal compiler error: cannot determine assignment type:", s)
		}
	default:
		log.Panicln("fatal compiler error: unexpected statement:", s)
	}
}

func (p *Program) compileLValue(w AssemblyWriter, expr Expr) {
	switch e := expr.(type) {
	case *PointerVariable:
		checkErr(w.PointerAddress(0, e.Num))
	case *BitVariable:
		checkErr(w.BitAddress(0, e.Num))
	case *JumpRegister:
		checkErr(w.JumpAddress(0))
	case *ValueAt:
		p.compileRValue(w, e.Target)
		if e.Offset != 0 {
			checkErr(w.Advance(0, e.Offset))
		}
	case *AddressOf:
		log.Panicln("internal compiler error: invalid lvalue:", e)
	case *BitValue:
		log.Panicln("internal compiler error: invalid lvalue:", e)
	case *Nand:
		log.Panicln("internal compiler error: invalid lvalue:", e)
	case *Parenthesis:
		p.compileLValue(w, e.Inner)
	default:
		log.Panicln("internal compiler error: unrecognized lvalue type:", e)
	}
}

func (p *Program) compileRValue(w AssemblyWriter, expr Expr) {
	switch e := expr.(type) {
	case *PointerVariable:
		checkErr(w.PointerValue(0, e.Num))
	case *BitVariable:
		checkErr(w.BitAddress(0, e.Num))
		checkErr(w.ReadBit(0, 0))
	case *JumpRegister:
		checkErr(w.JumpValue(0))
	case *ValueAt:
		p.compileRValue(w, e.Target)
		checkErr(w.ReadBit(0, 0))
	case *AddressOf:
		p.compileLValue(w, e.Variable)
	case *BitValue:
		checkErr(w.BitValue(0, e.Bit))
	case *Nand:
		p.compileRValue(w, e.Left)
		checkErr(w.SaveRegister(0))
		p.compileRValue(w, e.Right)
		checkErr(w.LoadRegister(1))
		checkErr(w.NandBit(0, 1))
	case *Parenthesis:
		p.compileRValue(w, e.Inner)
	default:
		log.Panicln("internal compiler error: unrecognized rvalue type:", e)
	}
}

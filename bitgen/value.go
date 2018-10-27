package bitgen

type (
	// Pointer is a pointer variable.
	//
	// Pointer is Memory.
	Pointer uint64

	// Buffer is a bit variable.
	//
	// Buffer is Memory.
	Buffer uint64

	// Offset is a pointer to the bit Off bits after Base.
	//
	// Offset is Memory.
	Offset struct {
		Base Memory
		Off  uint64
	}
)

// Constant is a BIT value.
//
// Constant is Readable.
type Constant bool

const (
	// Zero is the Constant ZERO.
	Zero Constant = false
	// One is the Constant ONE.
	One Constant = true
)

type jumpReg struct{}

// JumpRegister is the jump register. Its value decides which path is taken
// when using Jump as a line ending.
//
// JumpRegister is Writable. To read from JumpRegister, use a Jump line ending.
var JumpRegister jumpReg

// Readable values can be read as bits.
type Readable interface {
	writeRead(w *Writer)
}

// Writable values can be written as bits.
type Writable interface {
	writeWrite(w *Writer)
}

// Memory values can be read as pointers.
//
// All Memory values are also Readable and Writable as bits.
type Memory interface {
	Readable
	Writable
	writeAddr(w *Writer)
}

// Nand represents the NAND binary operator. It is equivalent to ZERO if both
// Left and Right are ONE, and equivalent to ONE in any other case.
//
// Nand is Readable.
type Nand struct {
	Left  Readable
	Right Readable
}

func (p Pointer) writeAddr(w *Writer) {
	w.checkPtr(uint64(p), true)
	w.write(" VARIABLE")
	w.writeNumber(uint64(p))
}
func (b Buffer) writeAddr(w *Writer) {
	w.checkPtr(uint64(b), false)
	w.write(" THE ADDRESS OF VARIABLE")
	w.writeNumber(uint64(b))
}
func (o Offset) writeAddr(w *Writer) {
	for i := uint64(0); i < o.Off; i++ {
		w.write(" THE ADDRESS OF THE VALUE BEYOND")
	}
	o.Base.writeAddr(w)
}

func (j jumpReg) writeWrite(w *Writer) {
	w.write(" THE JUMP REGISTER")
}
func (p Pointer) writeWrite(w *Writer) {
	w.checkPtr(uint64(p), true)
	w.write(" THE VALUE AT VARIABLE")
	w.writeNumber(uint64(p))
}
func (b Buffer) writeWrite(w *Writer) {
	w.checkPtr(uint64(b), false)
	w.write(" VARIABLE")
	w.writeNumber(uint64(b))
}
func (o Offset) writeWrite(w *Writer) {
	if o.Off == 0 {
		o.Base.writeWrite(w)
		return
	}

	w.write(" THE VALUE BEYOND")
	for i := uint64(2); i < o.Off; i++ {
		w.write(" THE ADDRESS OF THE VALUE BEYOND")
	}
	o.Base.writeAddr(w)
}

func (c Constant) writeRead(w *Writer) {
	w.writeBit(bool(c))
}
func (p Pointer) writeRead(w *Writer) {
	p.writeWrite(w)
}
func (b Buffer) writeRead(w *Writer) {
	b.writeWrite(w)
}
func (o Offset) writeRead(w *Writer) {
	if o.Off == 0 {
		o.Base.writeRead(w)
		return
	}

	o.writeWrite(w)
}
func (n Nand) writeRead(w *Writer) {
	n.Left.writeRead(w)
	w.write(" NAND")
	if _, ok := n.Right.(binOp); ok {
		w.write(" OPEN PARENTHESIS")
		n.Right.writeRead(w)
		w.write(" CLOSE PARENTHESIS")
	} else {
		n.Right.writeRead(w)
	}
}

type binOp interface{ binOp() }

func (n Nand) binOp() {}

func (w *Writer) checkPtr(n uint64, want bool) {
	if n >= uint64(len(w.isPtr)) || w.isPtr[n] != want {
		panic("bitgen: variable is invalid or from a different Writer")
	}
}

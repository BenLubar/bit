package bitgen

// LineNumber is a line number in a BIT program.
type LineNumber uint64

// Exit is a reserved line number for exiting the program.
const Exit = LineNumber(0)

// Jump goes to a line number depending on the current value of the jump
// register.
type Jump struct {
	Zero LineNumber
	One  LineNumber
}

// Goto is either a LineNumber or a Jump.
type Goto interface {
	goto0() LineNumber
	goto1() LineNumber
}

func (n LineNumber) goto0() LineNumber { return n }
func (n LineNumber) goto1() LineNumber { return n }
func (j Jump) goto0() LineNumber       { return j.Zero }
func (j Jump) goto1() LineNumber       { return j.One }

func (w *Writer) startLine(n LineNumber) {
	if n >= LineNumber(len(w.wroteLine)) {
		panic("bitgen: line number is invalid or from a different Writer")
	}
	if w.wroteLine[n] {
		panic("bitgen: line number has already been used")
	}
	w.wroteLine[n] = true

	w.write("LINE NUMBER")
	w.writeNumber(uint64(n))
	w.write(" CODE")
}

func (w *Writer) useLine(n LineNumber) {
	if n >= LineNumber(len(w.refLine)) {
		panic("bitgen: line number is invalid or from a different Writer")
	}

	w.refLine[n] = true
}

func (w *Writer) endLine(target Goto) {
	defer w.write("\n")

	zero := target.goto0()
	one := target.goto1()
	w.useLine(zero)
	w.useLine(one)

	if zero == Exit && one == Exit {
		return
	}

	if zero == one {
		w.write(" GOTO")
		w.writeNumber(uint64(zero))
		return
	}

	if zero != Exit {
		w.write(" GOTO")
		w.writeNumber(uint64(zero))
		w.write(" IF THE JUMP REGISTER IS")
		w.writeBit(false)
	}

	if one != Exit {
		w.write(" GOTO")
		w.writeNumber(uint64(one))
		w.write(" IF THE JUMP REGISTER IS")
		w.writeBit(true)
	}
}

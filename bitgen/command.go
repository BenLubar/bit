package bitgen

// SetPtr stores a memory address in a pointer variable.
func (w *Writer) SetPtr(start LineNumber, end Goto, dest Pointer, src Memory) {
	w.startLine(start)
	dest.writeAddr(w)
	w.write(" EQUALS")
	src.writeAddr(w)
	w.endLine(end)
}

// SetVal stores a bit in a Writable destination.
func (w *Writer) SetVal(start LineNumber, end Goto, dest Writable, src Readable) {
	w.startLine(start)
	dest.writeWrite(w)
	w.write(" EQUALS")
	src.writeRead(w)
	w.endLine(end)
}

// Copy copies a sequence of count bits from src to dest.
//
// If src is after dest in the same buffer, set srcAfterDest to true.
// If dest is after src in the same buffer, set srcAfterDest to false.
// If src and dest are not the same buffer, srcAfterDest may have any value.
func (w *Writer) Copy(start LineNumber, end Goto, dest Memory, src Memory, count uint64, srcAfterDest bool) {
	off := func(base Memory, off uint64) Offset {
		if srcAfterDest {
			off = count - 1 - off
		}

		return Offset{
			Base: base,
			Off:  off,
		}
	}

	w.startLine(start)
	off(dest, 0).writeWrite(w)
	w.write(" EQUALS")
	if count == 0 {
		// make it a no-op
		off(dest, 0).writeRead(w)
	} else {
		off(src, 0).writeRead(w)
	}
	for i := uint64(1); i < count; i++ {
		next := w.ReserveLine()
		w.endLine(next)
		w.startLine(next)
		off(dest, i).writeWrite(w)
		w.write(" EQUALS")
		off(src, i).writeRead(w)
	}
	w.endLine(end)
}

// Print writes a single bit value to the standard output.
func (w *Writer) Print(start LineNumber, end Goto, value Constant) {
	w.startLine(start)
	w.write(" PRINT")
	w.writeBit(bool(value))
	w.endLine(end)
}

// Read reads a bit value from the standard input into the jump register.
func (w *Writer) Read(start LineNumber, end Goto) {
	w.startLine(start)
	w.write(" READ")
	w.endLine(end)
}

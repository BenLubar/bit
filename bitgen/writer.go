package bitgen

import (
	"bufio"
	"flag"
	"io"
)

func number(n uint64) string {
	if n == 0 {
		return "ZERO"
	}

	buf := []byte("ONE")
	first := true
	for i := 64 - 1; i >= 0; i-- {
		zero := n&(1<<uint(i)) == 0
		if first {
			first = zero
		} else {
			if zero {
				buf = append(buf, " ZERO"...)
			} else {
				buf = append(buf, " ONE"...)
			}
		}
	}
	return string(buf)
}

// Line represents a line number. It should be treated as an opaque type, except
// that 0 can be used as the first line number and the line number to exit on.
type Line uint64

// Writer generates BIT code.
type Writer struct {
	w     *bufio.Writer
	n     Line
	v     Variable
	heap  bool
	taken map[Line]bool
}

// NewWriter returns a *Writer. Remember to call Close when you're done!
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:     bufio.NewWriter(w),
		taken: make(map[Line]bool),
	}
}

// Close flushes the internal buffer.
func (w *Writer) Close() error {
	return w.w.Flush()
}

// ReserveLine returns a line number that has not yet been used.
func (w *Writer) ReserveLine() Line {
	w.n++
	return w.n
}

// ReserveVariable returns a Variable that has not yet been used.
func (w *Writer) ReserveVariable() Variable {
	if w.heap {
		panic("bitgen: cannot reserve after reserving heap")
	}
	v := w.v
	w.v++
	return v
}

// ReserveInteger returns an Integer covering memory locations that have not yet
// been used.
func (w *Writer) ReserveInteger(width uint) Integer {
	if w.heap {
		panic("bitgen: cannot reserve after reserving heap")
	}
	if width == 0 {
		panic("bitgen: illegal integer width: 0")
	}

	v := w.v
	w.v += Variable(width)
	return Integer{
		Start: v,
		Width: width,
	}
}

// ReserveHeap returns the first unused memory location after all reservations.
// No additional memory may be reserved after ReserveHeap is called.
func (w *Writer) ReserveHeap() Variable {
	w.heap = true
	return w.v
}

var flagCrashOnLine = flag.Uint64("crash-on-line", ^uint64(0), "debugging tool: crash when this line number is reserved")

func (w *Writer) line(n Line, line string, goto0, goto1 Line) {
	if w.taken[n] {
		panic("bitgen: duplicate line number")
	}
	w.taken[n] = true
	if n == Line(*flagCrashOnLine) {
		panic("DEBUG")
	}
	w.w.WriteString("LINE NUMBER ")
	w.w.WriteString(number(uint64(n)))
	w.w.WriteString(" CODE ")
	w.w.WriteString(line)

	if goto0 == goto1 {
		if goto0 != 0 {
			w.w.WriteString(" GOTO ")
			w.w.WriteString(number(uint64(goto0)))
		}
	} else {
		if goto0 != 0 {
			w.w.WriteString(" GOTO ")
			w.w.WriteString(number(uint64(goto0)))
			w.w.WriteString(" IF THE JUMP REGISTER IS ZERO")
		}
		if goto1 != 0 {
			w.w.WriteString(" GOTO ")
			w.w.WriteString(number(uint64(goto1)))
			w.w.WriteString(" IF THE JUMP REGISTER IS ONE")
		}
	}
	w.w.WriteString("\n")
}

// Jump to goto0 if value is ZERO or goto1 if value is ONE.
func (w *Writer) Jump(start Line, value Value, goto0, goto1 Line) {
	w.line(start, "THE JUMP REGISTER EQUALS "+value.simplify().String(), goto0, goto1)
}

// Assign sets left to right.
func (w *Writer) Assign(start Line, left, right Value, end Line) {
	w.line(start, left.simplify().String()+" EQUALS "+right.simplify().String(), end, end)
}

// PrintString writes the bytes in a string, most significant bit first.
func (w *Writer) PrintString(start Line, value string, end Line) {
	for i := range value {
		var next Line
		if i == len(value)-1 {
			next = end
		} else {
			next = w.ReserveLine()
		}
		w.Print(start, value[i], next)
		start = next
	}
}

// Print writes a byte, most significant bit first.
func (w *Writer) Print(start Line, value byte, end Line) {
	for i := uint(8) - 1; i < 8; i-- {
		var next Line
		if i == 0 {
			next = end
		} else {
			next = w.ReserveLine()
		}
		w.PrintBit(start, (value>>i)&1 == 1, next)
		start = next
	}
}

// PrintBit writes a bit to the standard output.
func (w *Writer) PrintBit(start Line, value Bit, end Line) {
	if value {
		w.line(start, "PRINT ONE", end, end)
	} else {
		w.line(start, "PRINT ZERO", end, end)
	}
}

// Output writes value, most significant bit first.
func (w *Writer) Output(start Line, value Integer, end Line) {
	for i := value.Width - 1; i < value.Width; i-- {
		var next Line
		if i == 0 {
			next = end
		} else {
			next = w.ReserveLine()
		}
		w.OutputBit(start, value.Bit(i), next)
		start = next
	}
}

// OutputBit writes value to the standard output.
func (w *Writer) OutputBit(start Line, value Value, end Line) {
	n0 := w.ReserveLine()
	n1 := w.ReserveLine()

	w.Jump(start, value, n0, n1)
	w.line(n0, "PRINT ZERO", end, end)
	w.line(n1, "PRINT ONE", end, end)
}

// Input reads value, most significant bit first.
func (w *Writer) Input(start Line, value Integer, end Line) {
	for i := value.Width - 1; i < value.Width; i-- {
		var next Line
		if i == 0 {
			next = end
		} else {
			next = w.ReserveLine()
		}
		w.InputBit(start, value.Bit(i), next)
		start = next
	}
}

// InputBit reads value from the standard input.
func (w *Writer) InputBit(start Line, value Value, end Line) {
	n0 := w.ReserveLine()
	n1 := w.ReserveLine()

	w.line(start, "READ", n0, n1)
	w.Assign(n0, value, Bit(false), end)
	w.Assign(n1, value, Bit(true), end)
}

// InputEOF reads value, most significant bit first.
func (w *Writer) InputEOF(start Line, value Integer, end, eof Line) {
	for i := value.Width - 1; i < value.Width; i-- {
		var next Line
		if i == 0 {
			next = end
		} else {
			next = w.ReserveLine()
		}

		if i == value.Width-1 {
			w.InputBitEOF(start, value.Bit(i), next, eof)
		} else {
			w.InputBit(start, value.Bit(i), next)
		}
		start = next
	}
}

// InputBitEOF reads value from the standard input.
func (w *Writer) InputBitEOF(start Line, value Value, end, eof Line) {
	n0 := w.ReserveLine()
	n1 := w.ReserveLine()

	w.line(start, "READ "+number(uint64(eof)), n0, n1)
	w.Assign(n0, value, Bit(false), end)
	w.Assign(n1, value, Bit(true), end)
}

// Cmp jumps to same if value == base or different otherwise.
func (w *Writer) Cmp(start Line, value Integer, base uint64, same, different Line) {
	if base&(1<<value.Width-1) != base {
		panic("bitgen: non-equal widths for Cmp")
	}

	for i := value.Width - 1; i < value.Width; i-- {
		var next Line
		if i == 0 {
			next = same
		} else {
			next = w.ReserveLine()
		}

		if base&(1<<i) == 0 {
			w.Jump(start, value.Bit(i), next, different)
		} else {
			w.Jump(start, value.Bit(i), different, next)
		}
		start = next
	}
}

func (w *Writer) AddInt(start Line, left, right Integer, end, overflow Line) {
	if left.Width != right.Width {
		panic("bitgen: non-equal widths for AddInt")
	}

	var carry Line

	for i := uint(0); i < left.Width; i++ {
		var next, nextCarry Line
		if i == left.Width-1 {
			next, nextCarry = end, overflow
		} else {
			next, nextCarry = w.ReserveLine(), w.ReserveLine()
		}

		one := w.ReserveLine()
		w.Jump(start, right.Bit(i), next, one)

		if carry != 0 {
			w.Jump(carry, right.Bit(i), one, nextCarry)
		}

		setOne, setTwo := w.ReserveLine(), w.ReserveLine()
		w.Jump(one, left.Bit(i), setOne, setTwo)
		w.Assign(setOne, left.Bit(i), Bit(true), next)
		w.Assign(setTwo, left.Bit(i), Bit(false), nextCarry)

		start, carry = next, nextCarry
	}
}

func (w *Writer) Add(start Line, left Integer, right uint64, end, overflow Line) {
	if right&(1<<left.Width-1) != right {
		panic("bitgen: non-equal widths for Add")
	}
	if right == 0 {
		panic("bitgen: cannot add 0")
	}

	for i := left.Width - 1; i < left.Width; i-- {
		var prev, carry Line
		if right&(1<<i-1) == 0 {
			prev = start
		} else {
			if right&(1<<i) != 0 {
				prev = w.ReserveLine()
			}
			carry = w.ReserveLine()
		}

		one := w.ReserveLine()
		w.Assign(one, left.Bit(i), Bit(true), end)

		two := w.ReserveLine()
		w.Assign(two, left.Bit(i), Bit(false), overflow)

		three := overflow

		if carry == 0 {
			w.Jump(prev, left.Bit(i), one, two)
			break
		}

		if prev == 0 {
			w.Jump(carry, left.Bit(i), one, two)
			overflow = carry
		} else {
			w.Jump(prev, left.Bit(i), one, two)
			w.Jump(carry, left.Bit(i), two, three)
			end, overflow = prev, carry
		}
	}
}

// Increment adds ONE to value, then jumps to end if successful or overflow if
// it needed to carry to a nonexistent bit.
func (w *Writer) Increment(start Line, value Integer, end, overflow Line) {
	w.Add(start, value, 1, end, overflow)
}

// Decrement subtracts ONE from value, then jumps to end if successful or
// underflow if it needed to borrow from a nonexistent bit.
func (w *Writer) Decrement(start Line, value Integer, end, underflow Line) {
	next := start

	for i := uint(0); i < value.Width; i++ {
		current := value.Bit(i)

		n1 := w.ReserveLine()
		n2 := w.ReserveLine()

		w.Jump(next, current, n1, n2)

		if i == value.Width-1 {
			next = underflow
		} else {
			next = w.ReserveLine()
		}
		w.Assign(n1, current, Bit(true), next)
		w.Assign(n2, current, Bit(false), end)
	}
}

// Copy sets left to right.
func (w *Writer) Copy(start Line, left, right Integer, end Line) {
	if left.Width != right.Width {
		panic("bitgen: cannot copy integers of varying width")
	}

	for i := left.Width - 1; i < left.Width; i-- {
		var next Line
		if i == 0 {
			next = end
		} else {
			next = w.ReserveLine()
		}
		w.Assign(start, left.Bit(i), right.Bit(i), next)
		start = next
	}
}

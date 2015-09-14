package bitgen

import (
	"bufio"
	"encoding/gob"
	"flag"
	"io"
	"os"
	"runtime"
)

func write(w io.Writer, s string, n *int64, err *error) {
	if *err != nil {
		return
	}
	var nn int
	nn, *err = io.WriteString(w, s)
	*n += int64(nn)
}

func writeTo(w io.Writer, wt io.WriterTo, n *int64, err *error) {
	if *err != nil {
		return
	}
	var nn int64
	nn, *err = wt.WriteTo(w)
	*n += nn
}

func number(w io.Writer, n uint64, nn *int64, err *error) {
	if *err != nil {
		return
	}

	if n == 0 {
		write(w, "ZERO", nn, err)
		return
	}

	write(w, "ONE", nn, err)
	first := true
	for i := 64 - 1; i >= 0; i-- {
		zero := n&(1<<uint(i)) == 0
		if first {
			first = zero
		} else {
			if zero {
				write(w, " ZERO", nn, err)
			} else {
				write(w, " ONE", nn, err)
			}
		}
	}
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
	taken map[Line][64]uint64
}

// NewWriter returns a *Writer. Remember to call Close when you're done!
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:     bufio.NewWriter(w),
		taken: make(map[Line][64]uint64),
	}
}

var flagMap = flag.String("map", "", "output mapping from BIT line number to stack trace to this file")

// Close flushes the internal buffer.
func (w *Writer) Close() error {
	if err := w.w.Flush(); err != nil {
		return err
	}

	if *flagMap != "" {
		f, err := os.Create(*flagMap)
		if err != nil {
			return err
		}

		err = gob.NewEncoder(f).Encode(w.taken)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
	}

	return nil
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

func (w *Writer) startLine(n Line, nn *int64, err *error) {
	if *err != nil {
		return
	}

	if _, ok := w.taken[n]; ok {
		panic("bitgen: duplicate line number")
	}
	var pc [64]uint64
	if *flagMap != "" {
		var pcnative [len(pc)]uintptr
		if i := runtime.Callers(1, pcnative[:]); i >= len(pc) {
			panic("call stack too deep for profile")
		}
		for i, n := range pcnative {
			pc[i] = uint64(n)
		}
	}
	w.taken[n] = pc
	if n == Line(*flagCrashOnLine) {
		panic("DEBUG")
	}

	write(w.w, "LINE NUMBER ", nn, err)
	number(w.w, uint64(n), nn, err)
	write(w.w, " CODE ", nn, err)
}

func (w *Writer) endLine(goto0, goto1 Line, n *int64, err *error) {
	if *err != nil {
		return
	}
	if goto0 == goto1 {
		if goto0 != 0 {
			write(w.w, " GOTO ", n, err)
			number(w.w, uint64(goto0), n, err)
		}
	} else {
		if goto0 != 0 {
			write(w.w, " GOTO ", n, err)
			number(w.w, uint64(goto0), n, err)
			write(w.w, " IF THE JUMP REGISTER IS ZERO", n, err)
		}
		if goto1 != 0 {
			write(w.w, " GOTO ", n, err)
			number(w.w, uint64(goto1), n, err)
			write(w.w, " IF THE JUMP REGISTER IS ONE", n, err)
		}
	}
	write(w.w, "\n", n, err)
}

// Jump to goto0 if value is ZERO or goto1 if value is ONE.
func (w *Writer) Jump(start Line, value Value, goto0, goto1 Line) (n int64, err error) {
	w.jump(start, value, goto0, goto1, &n, &err)
	return
}

func (w *Writer) jump(start Line, value Value, goto0, goto1 Line, n *int64, err *error) {
	if *err != nil {
		return
	}

	w.startLine(start, n, err)
	write(w.w, "THE JUMP REGISTER EQUALS ", n, err)
	writeTo(w.w, value.simplify(), n, err)
	w.endLine(goto0, goto1, n, err)
}

// Assign sets left to right.
func (w *Writer) Assign(start Line, left, right Value, end Line) (n int64, err error) {
	w.assign(start, left, right, end, &n, &err)
	return
}

func (w *Writer) assign(start Line, left, right Value, end Line, n *int64, err *error) {
	if *err != nil {
		return
	}

	w.startLine(start, n, err)
	writeTo(w.w, left.simplify(), n, err)
	write(w.w, " EQUALS ", n, err)
	writeTo(w.w, right.simplify(), n, err)
	w.endLine(end, end, n, err)
}

// PrintString writes the bytes in a string, most significant bit first.
func (w *Writer) PrintString(start Line, value string, end Line) (n int64, err error) {
	w.printString(start, value, end, &n, &err)
	return
}

func (w *Writer) printString(start Line, value string, end Line, n *int64, err *error) {
	if *err != nil {
		return
	}

	for i := range value {
		var next Line
		if i == len(value)-1 {
			next = end
		} else {
			next = w.ReserveLine()
		}
		w.print(start, value[i], next, n, err)
		start = next
	}
}

// Print writes a byte, most significant bit first.
func (w *Writer) Print(start Line, value byte, end Line) (n int64, err error) {
	w.print(start, value, end, &n, &err)
	return
}

func (w *Writer) print(start Line, value byte, end Line, n *int64, err *error) {
	if *err != nil {
		return
	}

	for i := uint(8) - 1; i < 8; i-- {
		var next Line
		if i == 0 {
			next = end
		} else {
			next = w.ReserveLine()
		}
		w.printBit(start, (value>>i)&1 == 1, next, n, err)
		start = next
	}
}

// PrintBit writes a bit to the standard output.
func (w *Writer) PrintBit(start Line, value Bit, end Line) (n int64, err error) {
	w.printBit(start, value, end, &n, &err)
	return
}

func (w *Writer) printBit(start Line, value Bit, end Line, n *int64, err *error) {
	if *err != nil {
		return
	}

	w.startLine(start, n, err)
	if value {
		write(w.w, "PRINT ONE", n, err)
	} else {
		write(w.w, "PRINT ZERO", n, err)
	}
	w.endLine(end, end, n, err)
}

// Output writes value, most significant bit first.
func (w *Writer) Output(start Line, value Integer, end Line) (n int64, err error) {
	w.output(start, value, end, &n, &err)
	return
}

func (w *Writer) output(start Line, value Integer, end Line, n *int64, err *error) {
	if *err != nil {
		return
	}

	for i := value.Width - 1; i < value.Width; i-- {
		var next Line
		if i == 0 {
			next = end
		} else {
			next = w.ReserveLine()
		}
		w.outputBit(start, value.Bit(i), next, n, err)
		start = next
	}
}

// OutputBit writes value to the standard output.
func (w *Writer) OutputBit(start Line, value Value, end Line) (n int64, err error) {
	w.outputBit(start, value, end, &n, &err)
	return
}

func (w *Writer) outputBit(start Line, value Value, end Line, n *int64, err *error) {
	if *err != nil {
		return
	}

	n0 := w.ReserveLine()
	n1 := w.ReserveLine()

	w.jump(start, value, n0, n1, n, err)

	w.startLine(n0, n, err)
	write(w.w, "PRINT ZERO", n, err)
	w.endLine(end, end, n, err)

	w.startLine(n1, n, err)
	write(w.w, "PRINT ONE", n, err)
	w.endLine(end, end, n, err)
}

// Input reads value, most significant bit first.
func (w *Writer) Input(start Line, value Integer, end Line) (n int64, err error) {
	w.input(start, value, end, &n, &err)
	return
}

func (w *Writer) input(start Line, value Integer, end Line, n *int64, err *error) {
	if *err != nil {
		return
	}

	for i := value.Width - 1; i < value.Width; i-- {
		var next Line
		if i == 0 {
			next = end
		} else {
			next = w.ReserveLine()
		}
		w.inputBit(start, value.Bit(i), next, n, err)
		start = next
	}
}

// InputBit reads value from the standard input.
func (w *Writer) InputBit(start Line, value Value, end Line) (n int64, err error) {
	w.inputBit(start, value, end, &n, &err)
	return
}

func (w *Writer) inputBit(start Line, value Value, end Line, n *int64, err *error) {
	if *err != nil {
		return
	}

	n0 := w.ReserveLine()
	n1 := w.ReserveLine()

	w.startLine(start, n, err)
	write(w.w, "READ", n, err)
	w.endLine(n0, n1, n, err)

	w.assign(n0, value, Bit(false), end, n, err)
	w.assign(n1, value, Bit(true), end, n, err)
}

// InputEOF reads value, most significant bit first.
func (w *Writer) InputEOF(start Line, value Integer, end, eof Line) (n int64, err error) {
	w.inputEOF(start, value, end, eof, &n, &err)
	return
}

func (w *Writer) inputEOF(start Line, value Integer, end, eof Line, n *int64, err *error) {
	if *err != nil {
		return
	}

	for i := value.Width - 1; i < value.Width; i-- {
		var next Line
		if i == 0 {
			next = end
		} else {
			next = w.ReserveLine()
		}

		if i == value.Width-1 {
			w.inputBitEOF(start, value.Bit(i), next, eof, n, err)
		} else {
			w.inputBit(start, value.Bit(i), next, n, err)
		}
		start = next
	}
}

// InputBitEOF reads value from the standard input.
func (w *Writer) InputBitEOF(start Line, value Value, end, eof Line) (n int64, err error) {
	w.inputBitEOF(start, value, end, eof, &n, &err)
	return
}

func (w *Writer) inputBitEOF(start Line, value Value, end, eof Line, n *int64, err *error) {
	if *err != nil {
		return
	}

	n0 := w.ReserveLine()
	n1 := w.ReserveLine()

	w.startLine(start, n, err)
	write(w.w, "READ ", n, err)
	number(w.w, uint64(eof), n, err)
	w.endLine(n0, n1, n, err)

	w.assign(n0, value, Bit(false), end, n, err)
	w.assign(n1, value, Bit(true), end, n, err)
}

// Cmp jumps to same if value == base or different otherwise.
func (w *Writer) Cmp(start Line, value Integer, base uint64, same, different Line) (n int64, err error) {
	w.cmp(start, value, base, same, different, &n, &err)
	return
}

func (w *Writer) cmp(start Line, value Integer, base uint64, same, different Line, n *int64, err *error) {
	if *err != nil {
		return
	}

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
			w.jump(start, value.Bit(i), next, different, n, err)
		} else {
			w.jump(start, value.Bit(i), different, next, n, err)
		}
		start = next
	}
}

func (w *Writer) AddInt(start Line, left, right Integer, end, overflow Line) (n int64, err error) {
	w.addInt(start, left, right, end, overflow, &n, &err)
	return
}

func (w *Writer) addInt(start Line, left, right Integer, end, overflow Line, n *int64, err *error) {
	if *err != nil {
		return
	}

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
		w.jump(start, right.Bit(i), next, one, n, err)

		if carry != 0 {
			w.jump(carry, right.Bit(i), one, nextCarry, n, err)
		}

		setOne, setTwo := w.ReserveLine(), w.ReserveLine()
		w.jump(one, left.Bit(i), setOne, setTwo, n, err)
		w.assign(setOne, left.Bit(i), Bit(true), next, n, err)
		w.assign(setTwo, left.Bit(i), Bit(false), nextCarry, n, err)

		start, carry = next, nextCarry
	}
}

func (w *Writer) Add(start Line, left Integer, right uint64, end, overflow Line) (n int64, err error) {
	w.add(start, left, right, end, overflow, &n, &err)
	return
}

func (w *Writer) add(start Line, left Integer, right uint64, end, overflow Line, n *int64, err *error) {
	if *err != nil {
		return
	}

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
		w.assign(one, left.Bit(i), Bit(true), end, n, err)

		two := w.ReserveLine()
		w.assign(two, left.Bit(i), Bit(false), overflow, n, err)

		three := overflow

		if carry == 0 {
			w.jump(prev, left.Bit(i), one, two, n, err)
			break
		}

		if prev == 0 {
			w.jump(carry, left.Bit(i), one, two, n, err)
			overflow = carry
		} else {
			w.jump(prev, left.Bit(i), one, two, n, err)
			w.jump(carry, left.Bit(i), two, three, n, err)
			end, overflow = prev, carry
		}
	}
}

// Increment adds ONE to value, then jumps to end if successful or overflow if
// it needed to carry to a nonexistent bit.
func (w *Writer) Increment(start Line, value Integer, end, overflow Line) (n int64, err error) {
	return w.Add(start, value, 1, end, overflow)
}

// Decrement subtracts ONE from value, then jumps to end if successful or
// underflow if it needed to borrow from a nonexistent bit.
func (w *Writer) Decrement(start Line, value Integer, end, underflow Line) (n int64, err error) {
	return w.Add(start, value, uint64(1)<<value.Width-1, underflow, end)
}

// Copy sets left to right.
func (w *Writer) Copy(start Line, left, right Integer, end Line) (n int64, err error) {
	w.copy(start, left, right, end, &n, &err)
	return
}

func (w *Writer) copy(start Line, left, right Integer, end Line, n *int64, err *error) {
	if *err != nil {
		return
	}

	if left.Width != right.Width {
		panic("bitgen: cannot copy integers of varying width")
	}

	for i := uint(0); i < left.Width; i++ {
		var next Line
		if i == left.Width-1 {
			next = end
		} else {
			next = w.ReserveLine()
		}
		w.assign(start, left.Bit(i), right.Bit(i), next, n, err)
		start = next
	}
}

func (w *Writer) Less(start Line, left, right Integer, less, equal, greater Line) (n int64, err error) {
	w.less(start, left, right, less, equal, greater, &n, &err)
	return
}

func (w *Writer) less(start Line, left, right Integer, less, equal, greater Line, n *int64, err *error) {
	if left.Width != right.Width {
		panic("bitgen: non-equal widths for Less")
	}

	for i := left.Width - 1; i < left.Width; i-- {
		zero, one := w.ReserveLine(), w.ReserveLine()
		w.jump(start, left.Bit(i), zero, one, n, err)

		var next Line
		if i == 0 {
			next = equal
		} else {
			next = w.ReserveLine()
		}

		w.jump(zero, right.Bit(i), next, less, n, err)
		w.jump(one, right.Bit(i), greater, next, n, err)

		start = next
	}
}

package bitgen

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"strings"
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

// Value represents an expression whose value is known at runtime.
type Value interface {
	fmt.Stringer
	simplify() Value
}

// Variable is a Value that represents a location in memory.
type Variable uint64

func (v Variable) String() string {
	return "VARIABLE " + number(uint64(v))
}
func (v Variable) simplify() Value { return v }

// ValueAt is a Value that represents the value at another Value. The inner
// Value must be an address-of-a-bit.
type ValueAt struct {
	Value
}

func (v ValueAt) String() string {
	return "THE VALUE AT " + v.Value.String()
}
func (v ValueAt) simplify() Value {
	v.Value = v.Value.simplify()
	if a, ok := v.Value.(AddressOf); ok {
		return a.Value
	}
	if a, ok := v.Value.(Offset); ok {
		if b, ok := a.Value.(AddressOf); ok {
			return offsetValue{b.Value, a.Offset}
		}
	}
	return v
}

// AddressOf is a Value that represents the address-of-a-bit of another Value.
// The inner Value must be a bit and must not be a constant.
type AddressOf struct {
	Value
}

func (v AddressOf) String() string {
	return "THE ADDRESS OF " + v.Value.String()
}
func (v AddressOf) simplify() Value {
	v.Value = v.Value.simplify()
	if a, ok := v.Value.(ValueAt); ok {
		return a.Value
	}
	return v
}

// Bit is a Value that represents a bit constant.
type Bit bool

func (v Bit) String() string {
	if v {
		return "ONE"
	}
	return "ZERO"
}
func (v Bit) simplify() Value {
	return v
}

// Offset is a Value that represents a chain of THE ADDRESS OF THE VALUE BEYOND.
// The inner Value must be an address-of-a-bit.
type Offset struct {
	Value
	Offset uint
}

func (v Offset) String() string {
	return strings.Repeat("THE ADDRESS OF THE VALUE BEYOND ", int(v.Offset)) + v.Value.String()
}
func (v Offset) simplify() Value {
	v.Value = v.Value.simplify()
	if v.Offset == 0 {
		return v.Value
	}
	return v
}

type offsetValue struct {
	Value
	Offset uint
}

func (v offsetValue) String() string {
	return strings.Repeat("THE VALUE BEYOND THE ADDRESS OF ", int(v.Offset)) + v.Value.String()
}
func (v offsetValue) simplify() Value {
	panic("bitgen: offsetValue.simplify should be unreachable")
}

// Integer represents an unsigned integer of a specified width.
type Integer struct {
	Start Value
	Width uint
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
	n1 := make([]Line, len(value))
	for i := range n1[:len(value)-1] {
		n1[i] = w.ReserveLine()
	}
	n1[len(value)-1] = end

	for i, next := range n1 {
		w.Print(start, value[i], next)
		start = next
	}
}

// Print writes a byte, most significant bit first.
func (w *Writer) Print(start Line, value byte, end Line) {
	n1 := w.ReserveLine()
	n2 := w.ReserveLine()
	n3 := w.ReserveLine()
	n4 := w.ReserveLine()
	n5 := w.ReserveLine()
	n6 := w.ReserveLine()
	n7 := w.ReserveLine()

	w.PrintBit(start, value&0x80 == 0x80, n1)
	w.PrintBit(n1, value&0x40 == 0x40, n2)
	w.PrintBit(n2, value&0x20 == 0x20, n3)
	w.PrintBit(n3, value&0x10 == 0x10, n4)
	w.PrintBit(n4, value&0x08 == 0x08, n5)
	w.PrintBit(n5, value&0x04 == 0x04, n6)
	w.PrintBit(n6, value&0x02 == 0x02, n7)
	w.PrintBit(n7, value&0x01 == 0x01, end)
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
	n1 := make([]Line, value.Width)
	for i := range n1[1:] {
		n1[i+1] = w.ReserveLine()
	}
	n1[0] = end

	for i := value.Width - 1; i < value.Width; i-- {
		w.OutputBit(start, ValueAt{Offset{AddressOf{value.Start}, i}}, n1[i])
		start = n1[i]
	}
}

// OutputBit writes value to the standard output.
func (w *Writer) OutputBit(start Line, value Value, end Line) {
	n1 := w.ReserveLine()
	n2 := w.ReserveLine()

	w.Jump(start, value, n1, n2)
	w.line(n1, "PRINT ZERO", end, end)
	w.line(n2, "PRINT ONE", end, end)
}

// Input reads value, most significant bit first.
func (w *Writer) Input(start Line, value Integer, end Line) {
	n1 := make([]Line, value.Width)
	for i := range n1[1:] {
		n1[i+1] = w.ReserveLine()
	}
	n1[0] = end

	for i := value.Width - 1; i < value.Width; i-- {
		w.InputBit(start, ValueAt{Offset{AddressOf{value.Start}, i}}, n1[i])
		start = n1[i]
	}
}

// InputBit reads value from the standard input.
func (w *Writer) InputBit(start Line, value Value, end Line) {
	n1 := w.ReserveLine()
	n2 := w.ReserveLine()

	w.line(start, "READ", n1, n2)
	w.Assign(n1, value, Bit(false), end)
	w.Assign(n2, value, Bit(true), end)
}

// InputEOF reads value, most significant bit first.
func (w *Writer) InputEOF(start Line, value Integer, end, eof Line) {
	n1 := make([]Line, value.Width)
	for i := range n1[1:] {
		n1[i+1] = w.ReserveLine()
	}
	n1[0] = end

	for i := value.Width - 1; i < value.Width; i-- {
		cur := ValueAt{Offset{AddressOf{value.Start}, i}}
		if i == value.Width-1 {
			w.InputBitEOF(start, cur, n1[i], eof)
		} else {
			w.InputBit(start, cur, n1[i])
		}
		start = n1[i]
	}
}

// InputBitEOF reads value from the standard input.
func (w *Writer) InputBitEOF(start Line, value Value, end, eof Line) {
	n1 := w.ReserveLine()
	n2 := w.ReserveLine()

	w.line(start, "READ "+number(uint64(eof)), n1, n2)
	w.Assign(n1, value, Bit(false), end)
	w.Assign(n2, value, Bit(true), end)
}

// Cmp jumps to same if value == base or different otherwise.
func (w *Writer) Cmp(start Line, value Integer, base uint64, same, different Line) {
	n1 := make([]Line, value.Width)
	for i := range n1[:value.Width-1] {
		n1[i] = w.ReserveLine()
	}
	n1[value.Width-1] = same

	for i, next := range n1 {
		if base&(1<<uint(i)) == 0 {
			w.Jump(start, ValueAt{Offset{AddressOf{value.Start}, uint(i)}}, next, different)
		} else {
			w.Jump(start, ValueAt{Offset{AddressOf{value.Start}, uint(i)}}, different, next)
		}
		start = next
	}
}

// Increment adds ONE to value, then jumps to end if successful or overflow if
// it needed to carry to a nonexistent bit.
func (w *Writer) Increment(start Line, value Integer, end, overflow Line) {
	next := start

	for i := uint(0); i < value.Width; i++ {
		current := ValueAt{Offset{AddressOf{value.Start}, i}}

		n1 := w.ReserveLine()
		n2 := w.ReserveLine()

		w.Jump(next, current, n1, n2)
		w.Assign(n1, current, Bit(true), end)

		if i == value.Width-1 {
			next = overflow
		} else {
			next = w.ReserveLine()
		}
		w.Assign(n2, current, Bit(false), next)
	}
}

// Decrement subtracts ONE from value, then jumps to end if successful or
// underflow if it needed to borrow from a nonexistent bit.
func (w *Writer) Decrement(start Line, value Integer, end, underflow Line) {
	next := start

	for i := uint(0); i < value.Width; i++ {
		current := ValueAt{Offset{AddressOf{value.Start}, i}}

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

	n1 := make([]Line, left.Width)
	for i := range n1[:left.Width-1] {
		n1[i] = w.ReserveLine()
	}
	n1[left.Width-1] = end

	for i, next := range n1 {
		w.Assign(start, ValueAt{Offset{AddressOf{left.Start}, uint(i)}}, ValueAt{Offset{AddressOf{right.Start}, uint(i)}}, next)
		start = next
	}
}

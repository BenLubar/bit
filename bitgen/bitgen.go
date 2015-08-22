package bitgen

import (
	"bufio"
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

type Value interface {
	fmt.Stringer
	simplify() Value
}

type Variable uint64

func (v Variable) String() string {
	return "VARIABLE " + number(uint64(v))
}
func (v Variable) simplify() Value { return v }

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

type Bit bool

func (v Bit) String() string {
	if v {
		return "ONE"
	} else {
		return "ZERO"
	}
}
func (v Bit) simplify() Value {
	return v
}

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

type Integer struct {
	Start Value
	Width uint
}

type Line uint64

type Writer struct {
	w     *bufio.Writer
	n     Line
	v     Variable
	heap  bool
	taken map[Line]bool
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:     bufio.NewWriter(w),
		taken: make(map[Line]bool),
	}
}

func (w *Writer) Close() error {
	return w.w.Flush()
}

func (w *Writer) ReserveLine() Line {
	w.n++
	return w.n
}

func (w *Writer) ReserveVariable() Variable {
	if w.heap {
		panic("bitgen: cannot reserve after reserving heap")
	}
	v := w.v
	w.v++
	return v
}

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

func (w *Writer) ReserveHeap() Variable {
	w.heap = true
	return w.v
}

func (w *Writer) line(n Line, line string, goto0, goto1 Line) {
	if w.taken[n] {
		panic("bitgen: duplicate line number")
	}
	w.taken[n] = true
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

func (w *Writer) Jump(start Line, value Value, goto0, goto1 Line) {
	w.line(start, "THE JUMP REGISTER EQUALS "+value.simplify().String(), goto0, goto1)
}

func (w *Writer) Assign(start Line, left, right Value, end Line) {
	w.line(start, left.simplify().String()+" EQUALS "+right.simplify().String(), end, end)
}

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

func (w *Writer) OutputBit(start Line, value Value, end Line) {
	n1 := w.ReserveLine()
	n2 := w.ReserveLine()

	w.Jump(start, value, n1, n2)
	w.line(n1, "PRINT ZERO", end, end)
	w.line(n2, "PRINT ONE", end, end)
}

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

func (w *Writer) InputBit(start Line, value Value, end Line) {
	n1 := w.ReserveLine()
	n2 := w.ReserveLine()

	w.line(start, "READ", n1, n2)
	w.Assign(n1, value, Bit(false), end)
	w.Assign(n2, value, Bit(true), end)
}

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

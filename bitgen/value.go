package bitgen

import "io"

// Value represents an expression whose value is known at runtime.
type Value interface {
	io.WriterTo
	simplify() Value
}

// Variable is a Value that represents a location in memory.
type Variable uint64

func (v Variable) WriteTo(w io.Writer) (n int64, err error) {
	write(w, "VARIABLE ", &n, &err)
	number(w, uint64(v), &n, &err)
	return
}
func (v Variable) simplify() Value { return v }

// ValueAt is a Value that represents the value at another Value. The inner
// Value must be an address-of-a-bit.
type ValueAt struct {
	Value
}

func (v ValueAt) WriteTo(w io.Writer) (n int64, err error) {
	write(w, "THE VALUE AT ", &n, &err)
	writeTo(w, v.Value, &n, &err)
	return
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

func (v AddressOf) WriteTo(w io.Writer) (n int64, err error) {
	write(w, "THE ADDRESS OF ", &n, &err)
	writeTo(w, v.Value, &n, &err)
	return
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

func (v Bit) WriteTo(w io.Writer) (n int64, err error) {
	if v {
		write(w, "ONE", &n, &err)
	} else {
		write(w, "ZERO", &n, &err)
	}
	return
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

func (v Offset) WriteTo(w io.Writer) (n int64, err error) {
	for i := uint(0); i < v.Offset; i++ {
		write(w, "THE ADDRESS OF THE VALUE BEYOND ", &n, &err)
	}
	writeTo(w, v.Value, &n, &err)
	return
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

func (v offsetValue) WriteTo(w io.Writer) (n int64, err error) {
	for i := uint(0); i < v.Offset; i++ {
		write(w, "THE VALUE BEYOND THE ADDRESS OF ", &n, &err)
	}
	writeTo(w, v.Value, &n, &err)
	return
}
func (v offsetValue) simplify() Value {
	panic("bitgen: offsetValue.simplify should be unreachable")
}

// Integer represents an unsigned integer of a specified width.
type Integer struct {
	Start Value
	Width uint
}

// Bit returns the Value of the nth bit of i, with the least significant bit
// being 0 and the most significant bit being i.Width - 1.
func (i Integer) Bit(n uint) Value {
	if i.Width <= n {
		panic("bitgen: Integer.Bit out of range")
	}
	return ValueAt{Offset{AddressOf{i.Start}, n}}
}

// Sub returns an Integer starting at the start'th bit and ending at end'th bit
// of i. The width of the returned Integer is end - start.
func (i Integer) Sub(start, end uint) Integer {
	if start >= end || end > i.Width {
		panic("bitgen: Integer.Sub out of range")
	}
	return Integer{i.Bit(start), end - start}
}

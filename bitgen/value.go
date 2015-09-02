package bitgen

import (
	"fmt"
	"strings"
)

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

func (i Integer) Bit(n uint) Value {
	if i.Width <= n {
		panic("bitgen: Integer.Bit out of range")
	}
	return ValueAt{Offset{AddressOf{i.Start}, n}}
}

func (i Integer) Sub(start, end uint) Integer {
	if start >= end || end > i.Width {
		panic("bitgen: Integer.Sub out of range")
	}
	return Integer{i.Bit(start), end - start}
}

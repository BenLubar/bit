// Package ast implements an abstract syntax tree for BIT.
package ast

import (
	"math/big"
)

// Program is a BIT program.
type Program struct {
	Lines []*Line
}

// Bit is a single ZERO or ONE.
type Bit bool

const (
	// Zero is the ZERO bit.
	Zero Bit = false
	// One is the ONE bit.
	One Bit = true
)

// String implements fmt.Stringer.
func (b Bit) String() string {
	if b {
		return "ONE"
	}
	return "ZERO"
}

// Bits is a series of ZERO and ONE.
type Bits []Bit

// String implements fmt.Stringer.
func (b Bits) String() string {
	var buf []byte
	for _, x := range b {
		if x {
			buf = append(buf, " ONE"...)
		} else {
			buf = append(buf, " ZERO"...)
		}
	}

	if len(buf) > 0 {
		return string(buf[1:])
	}
	return ""
}

// Set copies the numeric value of this Bits to a *big.Int. It returns the
// *big.Int for convenience.
func (b Bits) Set(z *big.Int) *big.Int {
	z.SetInt64(0)
	for i, v := range b {
		if v {
			z.SetBit(z, len(b)-i-1, 1)
		}
	}
	return z
}

// Line is a single line of code in a BIT program.
type Line struct {
	// Num is the line number.
	Num Bits
	// Stmt is the command on this line of code.
	Stmt Stmt
	// Goto0 is the line to go to next if the jump register is ZERO.
	// If Goto0 is nil, the program terminates.
	Goto0 Bits
	// Goto1 is the line to go to next if the jump register is ONE.
	// If Goto1 is nil, the program terminates.
	Goto1 Bits
}

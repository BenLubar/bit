package main

import "strings"

type Line struct {
	Num  *Number
	Stmt Stmt
	Zero *Number
	One  *Number

	gotoZero, gotoOne *Line
}

func (l *Line) String() string {
	var gotoString string

	if l.Zero != nil && l.One != nil {
		zero := l.Zero.String()
		one := l.One.String()
		if zero == one {
			gotoString = " GOTO " + zero
		} else {
			gotoString = " GOTO " + zero + " IF THE JUMP REGISTER IS ZERO GOTO " + one + " IF THE JUMP REGISTER IS ONE"
		}
	} else if l.Zero != nil {
		gotoString = " GOTO " + l.Zero.String() + " IF THE JUMP REGISTER IS ZERO"
	} else if l.One != nil {
		gotoString = " GOTO " + l.Zero.String() + " IF THE JUMP REGISTER IS ONE"
	}

	return "LINE NUMBER " + l.Num.String() + " CODE " + l.Stmt.String() + gotoString
}

type Stmt interface {
	String() string
}

type ReadStmt struct {
	EOFLine *Number // extension
	gotoEOF *Line
}

func (r *ReadStmt) String() string {
	if r.EOFLine != nil {
		return "READ " + r.EOFLine.String()
	}
	return "READ"
}

type PrintStmt struct {
	Bit bool
}

func (p *PrintStmt) String() string {
	if p.Bit {
		return "PRINT ONE"
	} else {
		return "PRINT ZERO"
	}
}

type EqualsStmt struct {
	Left, Right Expr
}

func (e *EqualsStmt) String() string {
	return e.Left.String() + " EQUALS " + e.Right.String()
}

type Expr interface {
	Pointer() bool
	Value() bool
	Addressable() bool
	String() string
}

type UnknownVariable struct {
	Num *Number
}

func (*UnknownVariable) Pointer() bool     { return true }
func (*UnknownVariable) Value() bool       { return true }
func (*UnknownVariable) Addressable() bool { return true }
func (v *UnknownVariable) String() string  { return "VARIABLE " + v.Num.String() }

type BitVariable struct {
	Num *Number
}

func (*BitVariable) Pointer() bool     { return false }
func (*BitVariable) Value() bool       { return true }
func (*BitVariable) Addressable() bool { return true }
func (v *BitVariable) String() string  { return "VARIABLE " + v.Num.String() }

type PointerVariable struct {
	Num *Number
}

func (*PointerVariable) Pointer() bool     { return true }
func (*PointerVariable) Value() bool       { return false }
func (*PointerVariable) Addressable() bool { return false }
func (v *PointerVariable) String() string  { return "VARIABLE " + v.Num.String() }

type JumpRegister struct {
}

func (*JumpRegister) Pointer() bool     { return false }
func (*JumpRegister) Value() bool       { return true }
func (*JumpRegister) Addressable() bool { return false }
func (*JumpRegister) String() string    { return "THE JUMP REGISTER" }

type ValueAt struct {
	Target Expr
	Offset int
}

func (*ValueAt) Pointer() bool     { return false }
func (*ValueAt) Value() bool       { return true }
func (*ValueAt) Addressable() bool { return true }
func (v *ValueAt) String() string {
	if v.Offset == 0 {
		return "THE VALUE AT " + v.Target.String()
	}
	return strings.Repeat("THE VALUE BEYOND THE ADDRESS OF ", v.Offset-1) + "THE VALUE BEYOND " + v.Target.String()
}

type AddressOf struct {
	Variable Expr
}

func (*AddressOf) Pointer() bool     { return true }
func (*AddressOf) Value() bool       { return false }
func (*AddressOf) Addressable() bool { return false }
func (a *AddressOf) String() string  { return "THE ADDRESS OF " + a.Variable.String() }

type BitValue struct {
	Bit bool
}

func (*BitValue) Pointer() bool     { return false }
func (*BitValue) Value() bool       { return true }
func (*BitValue) Addressable() bool { return false }
func (b *BitValue) String() string {
	if b.Bit {
		return "ONE"
	} else {
		return "ZERO"
	}
}

type Nand struct {
	Left, Right Expr
}

func (*Nand) Pointer() bool     { return false }
func (*Nand) Value() bool       { return true }
func (*Nand) Addressable() bool { return false }
func (n *Nand) String() string  { return n.Left.String() + " NAND " + n.Right.String() }

type Parenthesis struct {
	Inner Expr
}

func (p *Parenthesis) Pointer() bool     { return p.Inner.Pointer() }
func (p *Parenthesis) Value() bool       { return p.Inner.Value() }
func (p *Parenthesis) Addressable() bool { return p.Inner.Addressable() }
func (p *Parenthesis) String() string {
	return "LEFT PARENTHESIS " + p.Inner.String() + " RIGHT PARENTHESIS"
}

package ast

// Expr is the interface all expressions implement.
type Expr interface {
	canVal() bool
	canAddr() bool
	canDeref() bool
	canAssign() bool
	String() string
	expr()
}

// CanVal returns true if the expression can be read as BIT value.
func CanVal(expr Expr) bool { return expr.canVal() }

// CanAddr returns true if the expression is an addressable BIT value.
func CanAddr(expr Expr) bool { return expr.canAddr() }

// CanDeref returns true if the expression is a pointer value.
func CanDeref(expr Expr) bool { return expr.canDeref() }

// CanAssign returns true if the expression points to a mutable value.
func CanAssign(expr Expr) bool { return expr.canAssign() }

// ValueAt is the THE VALUE AT operator. It can be used on pointers,
// either stored in a variable or from the THE ADDRESS OF operator.
type ValueAt struct {
	Ptr Expr
}

// ValueBeyond is the THE VALUE BEYOND operator. It can be used on pointers,
// either stored in a variable or from the THE ADDRESS OF operator.
type ValueBeyond struct {
	Ptr Expr
}

// AddressOf is the THE ADDRESS OF operator. It can be used on variables and
// on the result of dereferencing a pointer using THE VALUE AT or THE VALUE
// BEYOND.
type AddressOf struct {
	Val Expr
}

// Nand is the NAND operator. It is left-associative and returns ONE if either
// expression is ZERO or ZERO if both expressions are ONE. The expressions
// must not result in pointer-to-a-bit values.
type Nand struct {
	Left, Right Expr
}

// Constant is a bit value constant.
type Constant struct {
	Val Bit
}

// Variable is a numbered variable slot.
//
// Each variable can either hold a bit or a pointer-to-a-bit. Changing the
// type of a variable within a program is an error. Bit variables are stored
// in isolated memory spaces, so there is no way to get from one bit variable
// to another with THE VALUE BEYOND or any other operator.
type Variable struct {
	Num Bits
}

// JumpRegister is the jump register. It can be written to, but not read from.
type JumpRegister struct {
}

func (*ValueAt) canVal() bool         { return true }
func (*ValueBeyond) canVal() bool     { return true }
func (*AddressOf) canVal() bool       { return false }
func (*Nand) canVal() bool            { return true }
func (*Constant) canVal() bool        { return true }
func (*Variable) canVal() bool        { return true }
func (*JumpRegister) canVal() bool    { return false }
func (*ValueAt) canAddr() bool        { return true }
func (*ValueBeyond) canAddr() bool    { return true }
func (*AddressOf) canAddr() bool      { return false }
func (*Nand) canAddr() bool           { return false }
func (*Constant) canAddr() bool       { return false }
func (*Variable) canAddr() bool       { return true }
func (*JumpRegister) canAddr() bool   { return false }
func (*ValueAt) canDeref() bool       { return false }
func (*ValueBeyond) canDeref() bool   { return false }
func (*AddressOf) canDeref() bool     { return true }
func (*Nand) canDeref() bool          { return false }
func (*Constant) canDeref() bool      { return false }
func (*Variable) canDeref() bool      { return true }
func (*JumpRegister) canDeref() bool  { return false }
func (*ValueAt) canAssign() bool      { return true }
func (*ValueBeyond) canAssign() bool  { return true }
func (*AddressOf) canAssign() bool    { return false }
func (*Nand) canAssign() bool         { return false }
func (*Constant) canAssign() bool     { return false }
func (*Variable) canAssign() bool     { return true }
func (*JumpRegister) canAssign() bool { return true }
func (v *ValueAt) String() string     { return "THE VALUE AT " + v.Ptr.String() }
func (v *ValueBeyond) String() string { return "THE VALUE BEYOND " + v.Ptr.String() }
func (a *AddressOf) String() string   { return "THE ADDRESS OF " + a.Val.String() }
func (n *Nand) String() string        { return n.Left.String() + " NAND " + n.Right.String() }
func (c *Constant) String() string    { return c.Val.String() }
func (v *Variable) String() string    { return "VARIABLE " + v.Num.String() }
func (*JumpRegister) String() string  { return "THE JUMP REGISTER" }
func (*ValueAt) expr()                {}
func (*ValueBeyond) expr()            {}
func (*AddressOf) expr()              {}
func (*Nand) expr()                   {}
func (*Constant) expr()               {}
func (*Variable) expr()               {}
func (*JumpRegister) expr()           {}

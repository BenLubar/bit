package ast

// Expr is the interface all expressions implement.
type Expr interface {
	CanVal() bool
	CanAddr() bool
	CanDeref() bool
	CanAssign() bool
	String() string
	expr()
}

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

func (*ValueAt) CanVal() bool         { return true }
func (*ValueBeyond) CanVal() bool     { return true }
func (*AddressOf) CanVal() bool       { return false }
func (*Nand) CanVal() bool            { return true }
func (*Constant) CanVal() bool        { return true }
func (*Variable) CanVal() bool        { return true }
func (*JumpRegister) CanVal() bool    { return false }
func (*ValueAt) CanAddr() bool        { return true }
func (*ValueBeyond) CanAddr() bool    { return true }
func (*AddressOf) CanAddr() bool      { return false }
func (*Nand) CanAddr() bool           { return false }
func (*Constant) CanAddr() bool       { return false }
func (*Variable) CanAddr() bool       { return true }
func (*JumpRegister) CanAddr() bool   { return false }
func (*ValueAt) CanDeref() bool       { return false }
func (*ValueBeyond) CanDeref() bool   { return false }
func (*AddressOf) CanDeref() bool     { return true }
func (*Nand) CanDeref() bool          { return false }
func (*Constant) CanDeref() bool      { return false }
func (*Variable) CanDeref() bool      { return true }
func (*JumpRegister) CanDeref() bool  { return false }
func (*ValueAt) CanAssign() bool      { return true }
func (*ValueBeyond) CanAssign() bool  { return true }
func (*AddressOf) CanAssign() bool    { return false }
func (*Nand) CanAssign() bool         { return false }
func (*Constant) CanAssign() bool     { return false }
func (*Variable) CanAssign() bool     { return true }
func (*JumpRegister) CanAssign() bool { return true }
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

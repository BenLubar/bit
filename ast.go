package bit

import (
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/BenLubar/bit/bitio"
)

var (
	ErrNotLValue          = errors.New("cannot assign to non-lvalue")
	ErrUnassignedVariable = errors.New("cannot read from unassigned variable")
	ErrAddressOfValue     = errors.New("cannot take address of value")
	ErrAddressOfAddress   = errors.New("cannot take address of address")
	ErrDereferenceValue   = errors.New("cannot dereference value")
	ErrValueOfAddress     = errors.New("cannot take value of address-of-bit variable")
)

type context struct {
	jump   bool
	memory *big.Int
	bVar   *big.Int
	aVar   map[uint64]uint64
}

type line struct {
	stmt  Stmt
	goto0 *uint64
	goto1 *uint64
}

type Program map[uint64]line

func (p Program) AddLine(n uint64, stmt Stmt, goto0, goto1 *uint64) error {
	if _, ok := p[n]; ok {
		return fmt.Errorf("bit: duplicate line number %v", n)
	}

	p[n] = line{stmt: stmt.simplify(), goto0: goto0, goto1: goto1}

	return nil
}

func (p Program) Start() uint64 {
	if len(p) == 0 {
		return 0
	}
	for i := uint64(0); ; i++ {
		if _, ok := p[i]; ok {
			return i
		}
	}
}

type ProgramError struct {
	Err  error
	Line uint64
}

func (err *ProgramError) Error() string {
	return fmt.Sprintf("bit: on line %d: %v", err.Line, err.Err)
}

func (p Program) Run(in bitio.BitReader, out bitio.BitWriter) error {
	ctx := &context{
		jump:   false,
		memory: new(big.Int),
		bVar:   new(big.Int),
		aVar:   make(map[uint64]uint64),
	}
	pc := p.Start()

	for {
		line := p[pc]
		err := line.stmt.run(in, out, ctx)
		if err != nil {
			return &ProgramError{Err: err, Line: pc}
		}
		if !ctx.jump {
			if line.goto0 != nil {
				pc = *line.goto0
			} else {
				return nil
			}
		} else {
			if line.goto1 != nil {
				pc = *line.goto1
			} else {
				return nil
			}
		}
	}
}

func (p Program) RunByte(r io.ByteReader, w io.ByteWriter) error {
	in := bitio.NewReader(r)
	out := bitio.NewWriter(w)

	if err := p.Run(in, out); err != nil {
		return err
	}

	// flush any remaining bits
	for i := 0; i < 7; i++ {
		if err := out.WriteBit(false); err != nil {
			return err
		}
	}
	return nil
}

type Stmt interface {
	run(in bitio.BitReader, out bitio.BitWriter, ctx *context) error
	simplify() Stmt
}

type AssignStmt struct {
	Left  Expr
	Right Expr
}

func (stmt AssignStmt) simplify() Stmt {
	return AssignStmt{
		Left:  stmt.Left.simplify(),
		Right: stmt.Right.simplify(),
	}
}

func (stmt AssignStmt) run(in bitio.BitReader, out bitio.BitWriter, ctx *context) error {
	left, err := stmt.Left.run(ctx)
	if err != nil {
		return err
	}
	right, err := stmt.Right.run(ctx)
	if err != nil {
		return err
	}
	if vv, ok := left.(varVal); ok {
		if vv.isAddr(ctx) {
			a, err := right.pointer(ctx)
			if err != nil {
				return err
			}
			ctx.aVar[uint64(vv)] = a
			return nil
		}
		if vv.isValue(ctx) {
			b, err := right.value(ctx)
			if err != nil {
				return err
			}
			var n uint
			if b {
				n = 1
			}
			ctx.memory.SetBit(ctx.memory, int(vv), n)
			return nil
		}
		if right.isAddr(ctx) {
			a, err := right.pointer(ctx)
			if err != nil {
				return err
			}
			ctx.aVar[uint64(vv)] = a
			return nil
		}
		if right.isValue(ctx) {
			b, err := right.value(ctx)
			if err != nil {
				return err
			}
			ctx.bVar.SetBit(ctx.bVar, int(vv), 1)
			var n uint
			if b {
				n = 1
			}
			ctx.memory.SetBit(ctx.memory, int(vv), n)
			return nil
		}
		return ErrUnassignedVariable
	}
	return ErrNotLValue
}

type JumpRegisterStmt struct {
	Right Expr
}

func (stmt JumpRegisterStmt) simplify() Stmt {
	return JumpRegisterStmt{
		Right: stmt.Right.simplify(),
	}
}

func (stmt JumpRegisterStmt) run(in bitio.BitReader, out bitio.BitWriter, ctx *context) error {
	right, err := stmt.Right.run(ctx)
	if err != nil {
		return err
	}

	v, err := right.value(ctx)
	if err != nil {
		return err
	}

	ctx.jump = v
	return nil
}

type PrintStmt bool

func (stmt PrintStmt) simplify() Stmt {
	return stmt
}

func (stmt PrintStmt) run(in bitio.BitReader, out bitio.BitWriter, ctx *context) error {
	return out.WriteBit(bool(stmt))
}

type ReadStmt struct{}

func (stmt ReadStmt) simplify() Stmt {
	return stmt
}

func (stmt ReadStmt) run(in bitio.BitReader, out bitio.BitWriter, ctx *context) error {
	c, err := in.ReadBit()
	if err != nil {
		return err
	}
	ctx.jump = c
	return nil
}

type Expr interface {
	run(*context) (Val, error)
	simplify() Expr
}

type NandExpr struct {
	Left  Expr
	Right Expr
}

func (expr NandExpr) simplify() Expr {
	return NandExpr{
		Left:  expr.Left.simplify(),
		Right: expr.Right.simplify(),
	}
}

func (expr NandExpr) run(ctx *context) (Val, error) {
	left, err := expr.Left.run(ctx)
	if err != nil {
		return nil, err
	}
	right, err := expr.Right.run(ctx)
	if err != nil {
		return nil, err
	}
	l, err := left.value(ctx)
	if err != nil {
		return nil, err
	}
	r, err := right.value(ctx)
	if err != nil {
		return nil, err
	}
	return !actualVal(l && r), nil
}

type VarExpr uint64

func (expr VarExpr) simplify() Expr {
	return expr
}

func (expr VarExpr) run(ctx *context) (Val, error) {
	return varVal(expr), nil
}

type AddrExpr struct {
	X Expr
}

func (expr AddrExpr) simplify() Expr {
	return AddrExpr{
		X: expr.X.simplify(),
	}
}

func (expr AddrExpr) run(ctx *context) (Val, error) {
	x, err := expr.X.run(ctx)
	if err != nil {
		return nil, err
	}
	return x.addr(ctx)
}

type NextExpr struct {
	X Expr

	additional uint64
}

func (expr NextExpr) simplify() Expr {
	x := expr.X.simplify()
	if addr, ok := x.(AddrExpr); ok {
		if next, ok := addr.X.(NextExpr); ok {
			return NextExpr{
				X: next.X,

				additional: expr.additional + next.additional + 1,
			}
		}
	}
	return NextExpr{
		X: x,

		additional: expr.additional,
	}
}

func (expr NextExpr) run(ctx *context) (Val, error) {
	x, err := expr.X.run(ctx)
	if err != nil {
		return nil, err
	}

	i, err := x.pointer(ctx)
	if err != nil {
		return nil, err
	}

	v := varVal(i + 1 + expr.additional)

	// remember that we're a bit
	_, err = v.value(ctx)

	return v, err
}

type StarExpr struct {
	X Expr
}

func (expr StarExpr) simplify() Expr {
	x := expr.X.simplify()

	if addr, ok := x.(AddrExpr); ok {
		if next, ok := addr.X.(NextExpr); ok {
			return next
		}
	}

	return StarExpr{
		X: x,
	}
}

func (expr StarExpr) run(ctx *context) (Val, error) {
	x, err := expr.X.run(ctx)
	if err != nil {
		return nil, err
	}

	i, err := x.pointer(ctx)
	if err != nil {
		return nil, err
	}

	v := varVal(i)
	_, err = v.value(ctx)
	return v, err
}

type BitExpr bool

func (expr BitExpr) simplify() Expr {
	return expr
}

func (expr BitExpr) run(ctx *context) (Val, error) {
	return actualVal(expr), nil
}

type Val interface {
	value(*context) (bool, error)
	pointer(*context) (uint64, error)
	addr(*context) (Val, error)
	isAddr(*context) bool
	isValue(*context) bool
}

type actualVal bool

func (v actualVal) value(*context) (bool, error)     { return bool(v), nil }
func (v actualVal) pointer(*context) (uint64, error) { return 0, ErrDereferenceValue }
func (v actualVal) addr(*context) (Val, error)       { return nil, ErrAddressOfValue }
func (v actualVal) isAddr(*context) bool             { return false }
func (v actualVal) isValue(*context) bool            { return true }

type addrVal uint64

func (v addrVal) value(*context) (bool, error)     { return false, ErrValueOfAddress }
func (v addrVal) pointer(*context) (uint64, error) { return uint64(v), nil }
func (v addrVal) addr(*context) (Val, error)       { return nil, ErrAddressOfAddress }
func (v addrVal) isAddr(*context) bool             { return true }
func (v addrVal) isValue(*context) bool            { return false }

type varVal uint64

func (v varVal) value(ctx *context) (bool, error) {
	if v.isAddr(ctx) {
		return false, ErrValueOfAddress
	}

	// remember that we used this variable as a bit.
	ctx.bVar.SetBit(ctx.bVar, int(v), 1)

	return ctx.memory.Bit(int(v)) != 0, nil
}
func (v varVal) pointer(ctx *context) (uint64, error) {
	if v.isValue(ctx) {
		return 0, ErrAddressOfValue
	}
	if a, ok := ctx.aVar[uint64(v)]; ok {
		return a, nil
	}

	return 0, ErrUnassignedVariable
}
func (v varVal) addr(ctx *context) (Val, error) {
	if v.isAddr(ctx) {
		return nil, ErrAddressOfAddress
	}

	// remember that we used this variable as a bit.
	ctx.bVar.SetBit(ctx.bVar, int(v), 1)

	return addrVal(v), nil
}
func (v varVal) isAddr(ctx *context) bool {
	_, ok := ctx.aVar[uint64(v)]
	return ok
}
func (v varVal) isValue(ctx *context) bool {
	return ctx.bVar.Bit(int(v)) != 0
}

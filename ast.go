package bit

import (
	"errors"
	"fmt"
	"io"
	"math/big"
	"sort"

	"github.com/BenLubar/bit/bitio"
)

var (
	ErrNotLValue          = errors.New("cannot assign to non-lvalue")
	ErrUnassignedVariable = errors.New("cannot read from unassigned variable")
	ErrAddressOfValue     = errors.New("cannot take address of value")
	ErrAddressOfAddress   = errors.New("cannot take address of address")
	ErrDereferenceValue   = errors.New("cannot dereference value")
	ErrValueOfAddress     = errors.New("cannot take value of address-of-bit variable")
	ErrMissingLine        = errors.New("missing line for goto")
)

type context struct {
	jump   bool
	memory *big.Int
	bVar   *big.Int
	aVar   map[uint64]uint64
	n0, n1 big.Int
}

type line struct {
	num   uint64
	stmt  Stmt
	goto0 *uint64
	goto1 *uint64
	line0 *line
	line1 *line
	opt   optimized
}

type Program []*line

func (p *Program) AddLine(n uint64, stmt Stmt, goto0, goto1 *uint64) {
	*p = append(*p, &line{
		num:   n,
		stmt:  stmt.simplify(),
		goto0: goto0,
		goto1: goto1,
	})
}

func (p Program) Init() error {
	sort.Sort(p)
	if len(p) != 1 {
		last := ^uint64(0)
		for _, l := range p {
			if l.num == last {
				return fmt.Errorf("bit: duplicate line number %v", l.num)
			}
			last = l.num
		}
	}
	return p.bake()
}

func (p Program) Len() int           { return len(p) }
func (p Program) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Program) Less(i, j int) bool { return p[i].num < p[j].num }
func (p Program) findLine(n uint64) (l *line, ok bool) {
	if i := sort.Search(len(p), func(i int) bool {
		return p[i].num >= n
	}); i < len(p) && p[i].num == n {
		return p[i], true
	}
	return nil, false
}

func (p Program) Start() uint64 {
	if len(p) == 0 {
		return 0
	}
	return p[0].num
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
	pc := new(uint64)
	*pc = p.Start()
	line := p[*pc]

	for line != nil {
		if line.opt != nil {
			newpc, newline, err := line.opt.run(in, out, ctx)
			if err != nil {
				return &ProgramError{Err: err, Line: *pc}
			}
			pc = newpc
			line = newline
		} else {
			err := line.stmt.run(in, out, ctx)
			if err != nil {
				if r, ok := line.stmt.(ReadStmt); ok && r.pc != nil {
					if l, ok := p.findLine(*r.pc); ok {
						pc = r.pc
						line = l
						continue
					}
				}
				return &ProgramError{Err: err, Line: *pc}
			}
			if !ctx.jump {
				pc = line.goto0
				line = line.line0
			} else {
				pc = line.goto1
				line = line.line1
			}
		}
	}
	return nil
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
	if !left.actual && !left.ptr {
		if left.isAddr(ctx) {
			a, err := right.pointer(ctx)
			if err != nil {
				return err
			}
			ctx.aVar[uint64(left.a)] = a
			return nil
		}
		if left.isValue(ctx) {
			b, err := right.value(ctx)
			if err != nil {
				return err
			}
			var n uint
			if b {
				n = 1
			}
			ctx.memory.SetBit(ctx.memory, int(left.a), n)
			return nil
		}
		if right.isAddr(ctx) {
			a, err := right.pointer(ctx)
			if err != nil {
				return err
			}
			ctx.aVar[left.a] = a
			return nil
		}
		if right.isValue(ctx) {
			b, err := right.value(ctx)
			if err != nil {
				return err
			}
			ctx.bVar.SetBit(ctx.bVar, int(left.a), 1)
			var n uint
			if b {
				n = 1
			}
			ctx.memory.SetBit(ctx.memory, int(left.a), n)
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

type ReadStmt struct {
	pc *uint64
}

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
		return Val{}, err
	}
	right, err := expr.Right.run(ctx)
	if err != nil {
		return Val{}, err
	}
	l, err := left.value(ctx)
	if err != nil {
		return Val{}, err
	}
	r, err := right.value(ctx)
	if err != nil {
		return Val{}, err
	}
	return actualVal(!(l && r)), nil
}

type VarExpr uint64

func (expr VarExpr) simplify() Expr {
	return expr
}

func (expr VarExpr) run(ctx *context) (Val, error) {
	return varVal(uint64(expr)), nil
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
		return Val{}, err
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
		return Val{}, err
	}

	i, err := x.pointer(ctx)
	if err != nil {
		return Val{}, err
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
		return Val{}, err
	}

	i, err := x.pointer(ctx)
	if err != nil {
		return Val{}, err
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
	return actualVal(bool(expr)), nil
}

type Val struct {
	a      uint64
	ptr    bool
	actual bool
}

func actualVal(b bool) Val { return Val{0, b, true} }
func addrVal(a uint64) Val { return Val{a, true, false} }
func varVal(a uint64) Val  { return Val{a, false, false} }

func (v Val) value(ctx *context) (bool, error) {
	if v.actual {
		return v.ptr, nil
	}
	if v.ptr {
		return false, ErrValueOfAddress
	}
	if v.isAddr(ctx) {
		return false, ErrValueOfAddress
	}

	// remember that we used this variable as a bit.
	ctx.bVar.SetBit(ctx.bVar, int(v.a), 1)

	return ctx.memory.Bit(int(v.a)) != 0, nil
}
func (v Val) pointer(ctx *context) (uint64, error) {
	if v.actual {
		return 0, ErrDereferenceValue
	}
	if v.ptr {
		return v.a, nil
	}
	if v.isValue(ctx) {
		return 0, ErrAddressOfValue
	}
	if a, ok := ctx.aVar[v.a]; ok {
		return a, nil
	}

	return 0, ErrUnassignedVariable
}
func (v Val) addr(ctx *context) (Val, error) {
	if v.actual {
		return Val{}, ErrAddressOfValue
	}
	if v.ptr {
		return Val{}, ErrAddressOfAddress
	}
	if v.isAddr(ctx) {
		return Val{}, ErrAddressOfAddress
	}

	// remember that we used this variable as a bit.
	ctx.bVar.SetBit(ctx.bVar, int(v.a), 1)

	return addrVal(v.a), nil
}
func (v Val) isAddr(ctx *context) bool {
	if v.actual {
		return false
	}
	if v.ptr {
		return true
	}
	_, ok := ctx.aVar[v.a]
	return ok
}
func (v Val) isValue(ctx *context) bool {
	if v.actual {
		return true
	}
	if v.ptr {
		return false
	}
	return ctx.bVar.Bit(int(v.a)) != 0
}

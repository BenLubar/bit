package bit

import (
	"errors"
	"fmt"
	"io"

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

	p[n] = line{stmt: stmt, goto0: goto0, goto1: goto1}

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
	var jump bool
	var memory []bool
	variables := make(map[uint64]*Val)
	pc := p.Start()

	for {
		line := p[pc]
		err := line.stmt.run(in, out, &jump, &memory, variables)
		if err != nil {
			return &ProgramError{Err: err, Line: pc}
		}
		if !jump {
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
	run(in bitio.BitReader, out bitio.BitWriter, jump *bool, memory *[]bool, variables map[uint64]*Val) error
}

type AssignStmt struct {
	Left  Expr
	Right Expr
}

func (stmt AssignStmt) run(in bitio.BitReader, out bitio.BitWriter, jump *bool, memory *[]bool, variables map[uint64]*Val) error {
	left, err := stmt.Left.run(*memory, variables)
	if err != nil {
		return err
	}
	right, err := stmt.Right.run(*memory, variables)
	if err != nil {
		return err
	}
	if !left.lvalue {
		return ErrNotLValue
	}
	if left.actual != nil {
		panic("internal error")
	}
	if right.unknown {
		return ErrUnassignedVariable
	}
	if left.unknown {
		left.unknown = false
		left.addr = right.addr
	}
	if left.addr {
		left.index, err = right.pointer()
		return err
	}
	v, err := right.value(*memory)
	if err != nil {
		return err
	}
	if left.index >= uint64(len(*memory)) {
		if !v {
			return nil
		}
		mem := make([]bool, left.index+1)
		copy(mem, *memory)
		*memory = mem
	}
	(*memory)[left.index] = v
	return nil
}

type JumpRegisterStmt struct {
	Right Expr
}

func (stmt JumpRegisterStmt) run(in bitio.BitReader, out bitio.BitWriter, jump *bool, memory *[]bool, variables map[uint64]*Val) error {
	right, err := stmt.Right.run(*memory, variables)
	if err != nil {
		return err
	}

	v, err := right.value(*memory)
	if err != nil {
		return err
	}

	*jump = v
	return nil
}

type PrintStmt bool

func (stmt PrintStmt) run(in bitio.BitReader, out bitio.BitWriter, jump *bool, memory *[]bool, variables map[uint64]*Val) error {
	return out.WriteBit(bool(stmt))
}

type ReadStmt struct{}

func (stmt ReadStmt) run(in bitio.BitReader, out bitio.BitWriter, jump *bool, memory *[]bool, variables map[uint64]*Val) error {
	c, err := in.ReadBit()
	if err != nil {
		return err
	}
	*jump = c
	return nil
}

type Expr interface {
	run(memory []bool, variables map[uint64]*Val) (*Val, error)
}

type NandExpr struct {
	Left  Expr
	Right Expr
}

func (expr NandExpr) run(memory []bool, variables map[uint64]*Val) (*Val, error) {
	left, err := expr.Left.run(memory, variables)
	if err != nil {
		return nil, err
	}
	right, err := expr.Right.run(memory, variables)
	if err != nil {
		return nil, err
	}
	l, err := left.value(memory)
	if err != nil {
		return nil, err
	}
	r, err := right.value(memory)
	if err != nil {
		return nil, err
	}
	v := !(l && r)
	return &Val{actual: &v}, nil
}

type VarExpr uint64

func (expr VarExpr) run(memory []bool, variables map[uint64]*Val) (*Val, error) {
	if v, ok := variables[uint64(expr)]; ok {
		return v, nil
	}
	v := &Val{unknown: true, lvalue: true, index: uint64(expr)}
	variables[uint64(expr)] = v
	return v, nil
}

type AddrExpr struct {
	X Expr
}

func (expr AddrExpr) run(memory []bool, variables map[uint64]*Val) (*Val, error) {
	x, err := expr.X.run(memory, variables)
	if err != nil {
		return nil, err
	}
	if x.actual != nil {
		return nil, ErrAddressOfValue
	}
	if x.unknown {
		x.unknown = false
	}
	if x.addr {
		return nil, ErrAddressOfAddress
	}
	return &Val{addr: true, index: x.index}, nil
}

type NextExpr struct {
	X Expr
}

func (expr NextExpr) run(memory []bool, variables map[uint64]*Val) (*Val, error) {
	x, err := expr.X.run(memory, variables)
	if err != nil {
		return nil, err
	}
	if x.unknown {
		return nil, ErrUnassignedVariable
	}
	if !x.addr {
		return nil, ErrDereferenceValue
	}
	return &Val{lvalue: true, index: x.index + 1}, nil
}

type StarExpr struct {
	X Expr
}

func (expr StarExpr) run(memory []bool, variables map[uint64]*Val) (*Val, error) {
	x, err := expr.X.run(memory, variables)
	if err != nil {
		return nil, err
	}
	if x.unknown {
		return nil, ErrUnassignedVariable
	}
	if !x.addr {
		return nil, ErrDereferenceValue
	}
	return &Val{lvalue: true, index: x.index}, nil
}

type BitExpr bool

func (expr BitExpr) run(memory []bool, variables map[uint64]*Val) (*Val, error) {
	v := bool(expr)
	return &Val{actual: &v}, nil
}

type Val struct {
	addr    bool
	unknown bool
	lvalue  bool
	index   uint64
	actual  *bool
}

func (v *Val) value(memory []bool) (bool, error) {
	if v.actual != nil {
		return *v.actual, nil
	}
	if v.unknown {
		return false, ErrUnassignedVariable
	}
	if v.addr {
		return false, ErrValueOfAddress
	}
	if uint64(len(memory)) < v.index {
		return false, nil
	}
	return memory[v.index], nil
}

func (v *Val) pointer() (uint64, error) {
	if v.actual != nil {
		return 0, ErrDereferenceValue
	}
	if v.unknown {
		return 0, ErrUnassignedVariable
	}
	if !v.addr {
		return 0, ErrDereferenceValue
	}
	return v.index, nil
}

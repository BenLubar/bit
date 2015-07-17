package bit

import (
	"math/big"

	"github.com/BenLubar/bit/bitio"
)

func (p Program) bake() (Program, error) {
	// precompute gotos
	var ok bool
	for pc, l := range p {
		if l.goto0 != nil {
			l.line0, ok = p[*l.goto0]
			if !ok {
				return nil, &ProgramError{ErrMissingLine, pc}
			}
		}
		if l.goto1 != nil {
			l.line1, ok = p[*l.goto1]
			if !ok {
				return nil, &ProgramError{ErrMissingLine, pc}
			}
		}
	}

	// intrinsic x++ and x--
	for _, l := range p {
		as, ok := l.stmt.(AssignStmt)
		if !ok {
			continue
		}
		ptr, ok := as.Left.(VarExpr)
		if !ok {
			continue
		}
		if l.line0 != l.line1 || l.line0 == nil {
			continue
		}
		j, ok := l.line0.stmt.(JumpRegisterStmt)
		if !ok {
			continue
		}
		star, ok := j.Right.(StarExpr)
		if !ok {
			continue
		}
		if star.X != ptr {
			continue
		}
		// it could be an increment, a decrement, or something else.
		if p.optimizeIncDec(l.line0, ptr, true) {
			continue
		}
		if p.optimizeIncDec(l.line0, ptr, false) {
			continue
		}
	}

	return p, nil
}

func (p Program) out(l *line, jr bool) *line {
	if jr {
		return l.line1
	}
	return l.line0
}

var one = big.NewInt(1)

func (p Program) optimizeIncDec(l *line, ptr VarExpr, inc bool) bool {
	done, flow, donepc, flowpc, bits, ok := p.verifyIncDec(l, ptr, 0, inc)
	if !ok {
		return false
	}

	var mask big.Int
	mask.SetBit(&mask, int(bits), 1)
	mask.Sub(&mask, one)

	l.opt = optIncDec{
		ptr:   ptr,
		bits:  bits,
		inc:   inc,
		doneg: donepc,
		donel: done,
		flowg: flowpc,
		flowl: flow,
		mask:  &mask,
	}
	return false
}

func (p Program) verifyIncDec(l *line, ptr VarExpr, offset uint64, inc bool) (done, flow *line, donepc, flowpc *uint64, bits uint64, ok bool) {
	var reg Expr
	if offset == 0 {
		reg = StarExpr{X: ptr}
	} else {
		reg = NextExpr{X: ptr, additional: offset - 1}
	}
	if l.stmt != (JumpRegisterStmt{Right: reg}) {
		return
	}
	l0, l1 := p.out(l, !inc), p.out(l, inc)
	if l0.stmt != (AssignStmt{Left: reg, Right: BitExpr(inc)}) {
		return
	}
	if l1.stmt != (AssignStmt{Left: reg, Right: BitExpr(!inc)}) {
		return
	}
	if l0.line0 != l0.line1 {
		return
	}
	done = l0.line0
	if l1.line0 != l1.line1 {
		return
	}
	if l1.line0 != nil {
		var check *line
		check, flow, donepc, flowpc, bits, ok = p.verifyIncDec(l1.line0, ptr, offset+1, inc)
		if done != check {
			ok = false
		}
	}
	if !ok {
		flow = l1.line0
		donepc = l0.goto0
		flowpc = l1.goto0
		bits = offset + 1
		ok = true
	}
	return
}

type optimized interface {
	run(bitio.BitReader, bitio.BitWriter, *context) (*uint64, *line, error)
}

type optIncDec struct {
	ptr   VarExpr
	bits  uint64
	inc   bool
	doneg *uint64
	donel *line
	flowg *uint64
	flowl *line
	mask  *big.Int
}

func (opt optIncDec) run(in bitio.BitReader, out bitio.BitWriter, ctx *context) (g *uint64, l *line, err error) {
	val, err := opt.ptr.run(ctx)
	if err != nil {
		return
	}
	ptr, err := val.pointer(ctx)
	if err != nil {
		return
	}

	ctx.n0.Rsh(ctx.memory, uint(ptr))
	ctx.n0.And(&ctx.n0, opt.mask)

	ctx.n1.Set(&ctx.n0)

	g, l = opt.doneg, opt.donel
	ctx.jump = !opt.inc

	if opt.inc {
		ctx.n0.Add(&ctx.n0, one)
		if ctx.n0.Cmp(opt.mask) > 0 {
			ctx.n0.SetInt64(0)
			g, l = opt.flowg, opt.flowl
			ctx.jump = true
		}
	} else {
		ctx.n0.Sub(&ctx.n0, one)
		if ctx.n0.Sign() < 0 {
			ctx.n0.Set(opt.mask)
			g, l = opt.flowg, opt.flowl
			ctx.jump = false
		}
	}

	ctx.n1.Xor(&ctx.n1, &ctx.n0)
	for i := opt.bits; i > 0; i-- {
		if j := i - 1; ctx.n1.Bit(int(j)) == 1 {
			_, err = varVal(ptr + j).addr(ctx)
			if err != nil {
				return
			}

			ctx.memory.SetBit(ctx.memory, int(ptr+j), ctx.n0.Bit(int(j)))
		}
	}

	return
}
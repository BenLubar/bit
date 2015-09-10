package bit

import "github.com/BenLubar/bit/bitio"

func (p Program) bake() error {
	// precompute gotos
	var ok bool
	for _, l := range p {
		if l.goto0 != nil {
			l.line0, ok = p.findLine(*l.goto0)
			if !ok {
				return &ProgramError{ErrMissingLine, l.num}
			}
		}
		if l.goto1 != nil {
			l.line1, ok = p.findLine(*l.goto1)
			if !ok {
				return &ProgramError{ErrMissingLine, l.num}
			}
		}
	}

	return nil
}

func (p Program) Optimize() {
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

	// intrinsic print(byte)
intrinsicPrint:
	for _, l := range p {
		var b byte
		ll := l
		var g *uint64

		for i := uint(8); i > 0; i-- {
			if ll == nil || ll.line0 != ll.line1 {
				continue intrinsicPrint
			}
			ps, ok := ll.stmt.(PrintStmt)
			if !ok {
				continue intrinsicPrint
			}
			if bool(ps) {
				b |= 1 << (i - 1)
			}
			g = ll.goto0
			ll = ll.line0
		}

		l.opt = optPrintByteConst{
			b: b,
			l: ll,
			g: g,
		}
	}

	// intrinsic print(*ptr)
	for _, l := range p {
		js, ok := l.stmt.(JumpRegisterStmt)
		if !ok {
			continue
		}
		next, ok := js.Right.(NextExpr)
		if !ok {
			continue
		}
		if next.additional != 8-2 {
			continue
		}
		ptr, ok := next.X.(VarExpr)
		if !ok {
			continue
		}
		done, donepc, ok := p.verifyPrint(l, ptr, 8-2)
		if !ok {
			continue
		}
		l.opt = optPrintByte{
			p: ptr,
			l: done,
			g: donepc,
		}
	}
}

func (p Program) out(l *line, jr bool) *line {
	if jr {
		return l.line1
	}
	return l.line0
}

func (p Program) optimizeIncDec(l *line, ptr VarExpr, inc bool) bool {
	done, flow, donepc, flowpc, bits, ok := p.verifyIncDec(l, ptr, 0, inc)
	if !ok || bits > 64 {
		return false
	}

	l.opt = optIncDec{
		ptr:   ptr,
		bits:  bits,
		inc:   inc,
		doneg: donepc,
		donel: done,
		flowg: flowpc,
		flowl: flow,
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

func (p Program) verifyPrint(l *line, ptr VarExpr, offset uint64) (done *line, donepc *uint64, ok bool) {
	if l == nil || l.line0 == nil || l.line1 == nil {
		return
	}
	if offset == 0 {
		if l.stmt != (JumpRegisterStmt{Right: StarExpr{X: ptr}}) {
			return
		}
	} else {
		if l.stmt != (JumpRegisterStmt{Right: NextExpr{X: ptr, additional: offset - 1}}) {
			return
		}
	}
	done, donepc = l.line0.line0, l.line0.goto0
	if done != l.line0.line1 || done != l.line1.line0 || done != l.line1.line1 {
		return
	}
	if l.line0.stmt != PrintStmt(false) || l.line1.stmt != PrintStmt(true) {
		return
	}
	if offset == 0 {
		ok = true
		return
	}
	return p.verifyPrint(done, ptr, offset-1)
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

	mask := uint64(1<<opt.bits - 1)
	n0 := ctx.memory.Uint64(ptr)
	n0 &= mask

	g, l = opt.doneg, opt.donel
	ctx.jump = !opt.inc

	var n1 uint64
	if opt.inc {
		n1 = (n0 + 1) & mask
		if n0 > n1 {
			n1 = 0
			g, l = opt.flowg, opt.flowl
			ctx.jump = true
		}
	} else {
		n1 = (n0 + 1) & mask
		if n0 < n1 {
			n1 = mask
			g, l = opt.flowg, opt.flowl
			ctx.jump = false
		}
	}

	for i := ptr; i < ptr+opt.bits; i++ {
		if _, err = varVal(i).addr(ctx); err != nil {
			return
		}

		ctx.memory.SetBit(i, n1&1 == 1)
		n1 >>= 1
	}

	return
}

type optPrintByteConst struct {
	b byte
	g *uint64
	l *line
}

func (opt optPrintByteConst) run(in bitio.BitReader, out bitio.BitWriter, ctx *context) (g *uint64, l *line, err error) {
	return opt.g, opt.l, bitio.WriteByte(out, opt.b)
}

type optPrintByte struct {
	p VarExpr
	g *uint64
	l *line
}

func (opt optPrintByte) run(in bitio.BitReader, out bitio.BitWriter, ctx *context) (g *uint64, l *line, err error) {
	n0 := ctx.bVar.Uint64(uint64(opt.p))
	for i := uint64(0); i < 8; i++ {
		if n0&1<<i == 0 {
			_, err = varVal(uint64(opt.p) + i).value(ctx)
			if err != nil {
				return
			}
		}
	}

	b := byte(ctx.memory.Uint64(uint64(opt.p)))

	ctx.jump = b&1 == 1

	return opt.g, opt.l, bitio.WriteByte(out, b)
}

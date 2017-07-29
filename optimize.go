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
	// intrinsic left += const right
	for _, l := range p {
		if j, ok := l.stmt.(JumpRegisterStmt); ok {
			p.optimizeAddConst(l, j.Right)
		}
	}

	// intrinsic while (x--) ptr++;
	for _, l := range p {
		if s, ok := l.stmt.(AssignStmt); ok && l.line0 == l.line1 && l.line0 != nil {
			if a, ok := s.Right.(AddrExpr); ok {
				if v, ok := a.X.(VarExpr); ok {
					if o, ok := l.line0.opt.(optAddConst); ok && o.right == 1<<o.width-1 && o.flowl != nil && o.flowl.line0 == o.flowl.line1 && o.flowl.line0 == l.line0 {
						if p, ok := o.flowl.stmt.(AssignStmt); ok && s.Left == p.Left {
							if pa, ok := p.Right.(AddrExpr); ok {
								if pn, ok := pa.X.(NextExpr); ok && pn.X == p.Left {
									l.opt = optPointerAdvance{
										base:  v,
										addr:  o.left,
										width: o.width,
										ptr:   s.Left,
										step:  pn.additional + 1,
										endl:  o.endl,
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// intrinsic x = y
	for _, l := range p {
		if a, ok := l.stmt.(AssignStmt); ok {
			if _, ok = a.Right.(BitExpr); ok {
				p.optimizeCopyConst(l, a.Left)
			} else {
				p.optimizeCopy(l, a.Left, a.Right)
			}
		}
	}

	// intrinsic print(byte)
intrinsicPrint:
	for _, l := range p {
		var b byte
		ll := l

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
			ll = ll.line0
		}

		l.opt = optPrintByteConst{
			b: b,
			l: ll,
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
		done, ok := p.verifyPrint(l, ptr, 8-2)
		if !ok {
			continue
		}
		l.opt = optPrintByte{
			p: ptr,
			l: done,
		}
	}

	// intrinsic x < y or x <= y
	for _, l := range p {
		if _, ok := l.stmt.(JumpRegisterStmt); ok && l.line0 != l.line1 && l.line0 != nil && l.line1 != nil {
			p.optimizeLess(l)
		}
	}
}

func (p Program) out(l *line, jr bool) *line {
	if jr {
		return l.line1
	}
	return l.line0
}

func (p Program) optimizeAddConst(l *line, v Expr) {
	var endl, flowl *line

	offset, one, two, ok := p.verifyAddConstJump(l, v)
	if !ok || offset != 0 {
		return
	}

	offset, endl, ok = p.verifyAddConstOne(one, v)
	if !ok || offset != 0 {
		return
	}

	offset, flowl, ok = p.verifyAddConstTwo(two, v)
	if !ok || offset != 0 {
		return
	}

	n := uint64(1)

	offset, endl, flowl, n = p.verifyAddConst(v, endl, flowl, n, 1)

	l.opt = optAddConst{
		width: offset,
		left:  v,
		right: n,
		endl:  endl,
		flowl: flowl,
	}
}

func (p Program) verifyAddConst(v Expr, prevl, carryl *line, n, i uint64) (offset uint64, endl, flowl *line, nout uint64) {
	var ok bool
	if i < 64 && carryl != nil {
		offset, endl, flowl, nout, ok = p.verifyAddConstFalse(v, prevl, carryl, n, i)

		if !ok && prevl != nil {
			offset, endl, flowl, nout, ok = p.verifyAddConstTrue(v, prevl, carryl, n, i)
		}
	}
	if !ok {
		return i, prevl, carryl, n
	}

	return p.verifyAddConst(v, endl, flowl, nout, i+1)
}

func (p Program) verifyAddConstOffset(v, base Expr) (offset uint64, ok bool) {
	if v == base {
		return 0, true
	}

	if n, ok := v.(NextExpr); ok {
		if n.X == (AddrExpr{base}) {
			return n.additional + 1, true
		}
		if (StarExpr{n.X}) == base {
			return n.additional + 1, true
		}
		if bn, ok := base.(NextExpr); ok && n.X == bn.X {
			return n.additional - bn.additional, true
		}
	}

	return 0, false
}

func (p Program) verifyAddConstFalse(v Expr, prevl, carryl *line, n, i uint64) (offset uint64, endl, flowl *line, nout uint64, ok bool) {
	nout = n

	offset, one, two, ok := p.verifyAddConstJump(carryl, v)
	if !ok || offset != i {
		ok = false
		return
	}

	offset, endl, ok = p.verifyAddConstOne(one, v)
	if !ok || offset != i {
		ok = false
		return
	}

	offset, flowl, ok = p.verifyAddConstTwo(two, v)
	if !ok || offset != i {
		ok = false
		return
	}

	return
}

func (p Program) verifyAddConstTrue(v Expr, prevl, carryl *line, n, i uint64) (offset uint64, endl, flowl *line, nout uint64, ok bool) {
	nout = n | (1 << i)

	offset, one, two, ok := p.verifyAddConstJump(prevl, v)
	if !ok || offset != i {
		ok = false
		return
	}

	offset, endl, ok = p.verifyAddConstOne(one, v)
	if !ok || offset != i {
		ok = false
		return
	}

	offset, flowl, ok = p.verifyAddConstTwo(two, v)
	if !ok || offset != i {
		ok = false
		return
	}

	offset, two_, three, ok := p.verifyAddConstJump(carryl, v)
	if !ok || offset != i || two_ != two || three != flowl {
		ok = false
		return
	}

	return
}

func (p Program) verifyAddConstJump(l *line, v Expr) (offset uint64, l0, l1 *line, ok bool) {
	if l.line0 == l.line1 || l.line0 == nil || l.line1 == nil {
		return
	}
	l0, l1 = l.line0, l.line1

	var j JumpRegisterStmt
	if j, ok = l.stmt.(JumpRegisterStmt); !ok {
		return
	}

	if offset, ok = p.verifyAddConstOffset(j.Right, v); !ok {
		return
	}

	return
}

func (p Program) verifyAddConstOne(l *line, v Expr) (offset uint64, endl *line, ok bool) {
	if l.line0 != l.line1 {
		return
	}
	endl = l.line0

	var a AssignStmt
	if a, ok = l.stmt.(AssignStmt); !ok {
		return
	}

	if offset, ok = p.verifyAddConstOffset(a.Left, v); !ok {
		return
	}

	var b BitExpr
	if b, ok = a.Right.(BitExpr); !ok {
		return
	}

	if !b {
		ok = false
		return
	}

	return
}

func (p Program) verifyAddConstTwo(l *line, v Expr) (offset uint64, flowl *line, ok bool) {
	if l.line0 != l.line1 {
		return
	}
	flowl = l.line0

	var a AssignStmt
	if a, ok = l.stmt.(AssignStmt); !ok {
		return
	}

	if offset, ok = p.verifyAddConstOffset(a.Left, v); !ok {
		return
	}

	var b BitExpr
	if b, ok = a.Right.(BitExpr); !ok {
		return
	}

	if b {
		ok = false
		return
	}

	return
}

func (p Program) optimizeCopyConst(l *line, v Expr) {
	width, n, end := p.verifyCopyConst(l, v, 0, 0)
	if width == 0 {
		return
	}

	l.opt = optCopyConst{
		width: width,
		left:  v,
		right: n,
		endl:  end,
	}
}

func (p Program) verifyCopyConst(l *line, v Expr, i, n uint64) (offset, nout uint64, end *line) {
	if i < 64 && l != nil && l.line0 == l.line1 {
		if a, ok := l.stmt.(AssignStmt); ok {
			if b, ok := a.Right.(BitExpr); ok {
				if offset, ok = p.verifyAddConstOffset(a.Left, v); ok && offset == i {
					if b {
						n |= 1 << offset
					}
					return p.verifyCopyConst(l.line0, v, i+1, n)
				}
			}
		}
	}

	return i, n, l
}

func (p Program) optimizeCopy(l *line, left, right Expr) {
	width, end := p.verifyCopy(l, left, right, 0)
	if width < 2 {
		return
	}

	l.opt = optCopy{
		width: width,
		left:  left,
		right: right,
		endl:  end,
	}
}

func (p Program) verifyCopy(l *line, left, right Expr, i uint64) (offset uint64, end *line) {
	if i < 64 && l != nil && l.line0 == l.line1 {
		if a, ok := l.stmt.(AssignStmt); ok {
			if offset, ok = p.verifyAddConstOffset(a.Left, left); ok && offset == i {
				if offset, ok = p.verifyAddConstOffset(a.Right, right); ok && offset == i {
					return p.verifyCopy(l.line0, left, right, i+1)
				}
			}
		}
	}

	return i, l
}

func (p Program) verifyPrint(l *line, ptr VarExpr, offset uint64) (done *line, ok bool) {
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
	done = l.line0.line0
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

func (p Program) optimizeLess(l *line) {
	less, equal, greater, left, right, width, ok := p.verifyLess(l)
	if !ok {
		return
	}

	l.opt = optLess{
		less:    less,
		equal:   equal,
		greater: greater,
		left:    left,
		right:   right,
		width:   width,
	}
}

func (p Program) verifyLess(l *line) (less, equal, greater *line, left, right Expr, width uint64, ok bool) {
	if l == nil {
		return
	}

	l0, l1 := l.line0, l.line1
	if l0 == nil || l1 == nil {
		return
	}

	j, ok := l.stmt.(JumpRegisterStmt)
	if !ok {
		return
	}
	j0, ok := l0.stmt.(JumpRegisterStmt)
	if !ok {
		return
	}
	j1, ok := l1.stmt.(JumpRegisterStmt)
	if !ok {
		return
	}
	if j0.Right != j1.Right {
		ok = false
		return
	}

	equal0, equal1 := l0.line0, l1.line1
	if equal0 != equal1 {
		ok = false
		return
	}

	less, equal, greater = l0.line1, equal0, l1.line0

	nextLess, nextEqual, nextGreater, nextLeft, nextRight, nextWidth, ok := p.verifyLess(equal)
	if !ok || nextWidth >= 64 {
		left, right, width, ok = j.Right, j0.Right, 1, true
		return
	}

	lo, ok := p.verifyAddConstOffset(left, nextLeft)
	if !ok || lo != nextWidth {
		ok = false
		return
	}

	ro, ok := p.verifyAddConstOffset(right, nextRight)
	if !ok || ro != nextWidth {
		ok = false
		return
	}

	less, equal, greater, left, right, width, ok = nextLess, nextEqual, nextGreater, nextLeft, nextRight, nextWidth+1, true
	return
}

type optimized interface {
	run(bitio.BitReader, bitio.BitWriter, *context) (*line, error)
}

type optAddConst struct {
	width uint64
	left  Expr
	right uint64
	endl  *line
	flowl *line
}

func (opt optAddConst) run(in bitio.BitReader, out bitio.BitWriter, ctx *context) (l *line, err error) {
	v, err := opt.left.run(ctx)
	if err != nil {
		return
	}
	v, err = v.addr(ctx)
	if err != nil {
		return
	}
	offset, err := v.pointer(ctx)
	if err != nil {
		return
	}

	n0 := ctx.memory.Uint64Offset(offset) & (1<<opt.width - 1)
	n1 := n0 + opt.right
	n2 := n1 & (1<<opt.width - 1)

	if n0 < n2 {
		l = opt.endl
		ctx.jump = false
	} else {
		n1 = 1<<opt.width - 1
		l = opt.flowl
		ctx.jump = true
	}

	for i := offset; i < offset+opt.width; i++ {
		if _, err = varVal(i).addr(ctx); err != nil {
			return
		}
		ctx.memory.SetBit(i, n2&1 == 1)
		n1 >>= 1
		n2 >>= 1
		if n1 == 0 && n2 == 0 {
			break
		}
	}

	return
}

type optPointerAdvance struct {
	base  VarExpr
	addr  Expr
	width uint64
	ptr   Expr
	step  uint64
	endl  *line
}

func (opt optPointerAdvance) run(in bitio.BitReader, out bitio.BitWriter, ctx *context) (l *line, err error) {
	addr, err := opt.addr.run(ctx)
	if err != nil {
		return
	}

	addraddr, err := addr.addr(ctx)
	if err != nil {
		return
	}

	offset, err := addraddr.pointer(ctx)
	if err != nil {
		return
	}

	err = (AssignStmt{
		Left:  opt.ptr,
		Right: AddrExpr{VarExpr((ctx.memory.Uint64Offset(offset)&(1<<opt.width-1))*opt.step) + opt.base},
	}).run(in, out, ctx)
	if err != nil {
		return
	}

	for i := offset; i < offset+opt.width; i++ {
		_, err = varVal(i).value(ctx)
		if err != nil {
			return
		}

		ctx.memory.SetBit(i, true)
	}

	ctx.jump = false
	return opt.endl, nil
}

type optCopyConst struct {
	width uint64
	left  Expr
	right uint64
	endl  *line
}

func (opt optCopyConst) run(in bitio.BitReader, out bitio.BitWriter, ctx *context) (l *line, err error) {
	left, err := opt.left.run(ctx)
	if err != nil {
		return
	}

	addr, err := left.addr(ctx)
	if err != nil {
		return
	}

	offset, err := addr.pointer(ctx)
	if err != nil {
		return
	}

	n := opt.right
	for i := offset; i < offset+opt.width; i++ {
		_, err = varVal(i).value(ctx)
		if err != nil {
			return
		}

		ctx.memory.SetBit(i, n&1 == 1)
		n >>= 1
	}

	return opt.endl, nil
}

type optCopy struct {
	width uint64
	left  Expr
	right Expr
	endl  *line
}

func (opt optCopy) run(in bitio.BitReader, out bitio.BitWriter, ctx *context) (l *line, err error) {
	left, err := opt.left.run(ctx)
	if err != nil {
		return
	}

	addr, err := left.addr(ctx)
	if err != nil {
		return
	}

	offsetL, err := addr.pointer(ctx)
	if err != nil {
		return
	}

	right, err := opt.right.run(ctx)
	if err != nil {
		return
	}

	addr, err = right.addr(ctx)
	if err != nil {
		return
	}

	offsetR, err := addr.pointer(ctx)
	if err != nil {
		return
	}

	for i, j := offsetL, offsetR; i < offsetL+opt.width; i, j = i+1, j+1 {
		_, err = varVal(i).value(ctx)
		if err != nil {
			return
		}

		var b bool
		b, err = varVal(j).value(ctx)
		if err != nil {
			return
		}

		ctx.memory.SetBit(i, b)
	}

	return opt.endl, nil
}

type optPrintByteConst struct {
	b byte
	l *line
}

func (opt optPrintByteConst) run(in bitio.BitReader, out bitio.BitWriter, ctx *context) (l *line, err error) {
	return opt.l, bitio.WriteByte(out, opt.b)
}

type optPrintByte struct {
	p VarExpr
	l *line
}

func (opt optPrintByte) run(in bitio.BitReader, out bitio.BitWriter, ctx *context) (l *line, err error) {
	n0 := ctx.bVar.Uint64Offset(uint64(opt.p))
	for i := uint64(0); i < 8; i++ {
		if n0&1<<i == 0 {
			_, err = varVal(uint64(opt.p) + i).value(ctx)
			if err != nil {
				return
			}
		}
	}

	b := byte(ctx.memory.Uint64Offset(uint64(opt.p)))

	ctx.jump = b&1 == 1

	return opt.l, bitio.WriteByte(out, b)
}

type optLess struct {
	less    *line
	equal   *line
	greater *line
	left    Expr
	right   Expr
	width   uint64
}

func (opt optLess) run(in bitio.BitReader, out bitio.BitWriter, ctx *context) (l *line, err error) {
	leftbit, err := opt.left.run(ctx)
	if err != nil {
		return
	}

	leftaddr, err := leftbit.addr(ctx)
	if err != nil {
		return
	}

	leftoffset, err := leftaddr.pointer(ctx)
	if err != nil {
		return
	}

	rightbit, err := opt.right.run(ctx)
	if err != nil {
		return
	}

	rightaddr, err := rightbit.addr(ctx)
	if err != nil {
		return
	}

	rightoffset, err := rightaddr.pointer(ctx)
	if err != nil {
		return
	}

	left := ctx.memory.Uint64Offset(leftoffset) & (1<<opt.width - 1)
	right := ctx.memory.Uint64Offset(rightoffset) & (1<<opt.width - 1)

	if left < right {
		ctx.jump = true
		return opt.less, nil
	}

	if left > right {
		ctx.jump = false
		return opt.greater, nil
	}

	ctx.jump = left&1 == 1
	return opt.equal, nil
}

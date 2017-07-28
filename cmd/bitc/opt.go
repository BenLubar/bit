package main

func (p *Program) Optimize() {
	p.optimizeOffsets()
}

func (p *Program) optimizeOffsets() {
	for _, l := range p.Lines {
		if e, ok := l.Stmt.(*EqualsStmt); ok {
			p.optimizeOffsetExpr(&e.Left)
			p.optimizeOffsetExpr(&e.Right)
		}
	}
}

func (p *Program) optimizeOffsetExpr(expr *Expr) {
	switch e := (*expr).(type) {
	case *ValueAt:
		p.optimizeOffsetExpr(&e.Target)
		if a, ok := e.Target.(*AddressOf); ok {
			if v, ok := a.Variable.(*ValueAt); ok {
				*expr = &ValueAt{
					Target: v.Target,
					Offset: e.Offset + v.Offset,
				}
			}
		}
	case *AddressOf:
		p.optimizeOffsetExpr(&e.Variable)
	case *Nand:
		p.optimizeOffsetExpr(&e.Left)
		p.optimizeOffsetExpr(&e.Right)
	case *Parenthesis:
		p.optimizeOffsetExpr(&e.Inner)
	}
}

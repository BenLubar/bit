package main

func (ast *AST) Optimize() {
	ast.usedTypes = make(map[*ClassDecl]bool)
	for _, t := range basicClasses {
		ast.usedTypes[t] = true
	}
	ast.findUsedTypes(ast.main)

	ast.overriddenMethods = make(map[*MethodFeature]bool)
	for t := range ast.usedTypes {
		ast.findOverriddenMethods(t)
	}
}

func (ast *AST) findUsedTypes(c *ClassDecl) {
	if ast.usedTypes[c] || c == basicDummyNull || c == basicDummyNothing {
		return
	}
	ast.usedTypes[c] = true

	ast.findUsedTypes(c.Extends.Type.target)
	for _, feature := range c.Body {
		switch f := feature.(type) {
		case *VarFeature:
			ast.findUsedTypes(f.Type.target)

		case *MethodFeature:
			for _, a := range f.Args {
				ast.findUsedTypes(a.Type.target)
			}
			ast.findUsedTypes(f.Return.target)
			ast.findUsedTypesRecurse(f.Body)

		case *BlockFeature:

		case *NativeFeature:

		default:
			panic(f)
		}
	}
}

func (ast *AST) findUsedTypesRecurse(expr Expr) {
	switch e := expr.(type) {
	case *ConstructorExpr:
		for _, a := range e.Args {
			ast.findUsedTypes(a.Type.target)
		}
		ast.findUsedTypesRecurse(e.Expr)

	case *AssignExpr:
		ast.findUsedTypesRecurse(e.Right)

	case *IfExpr:
		ast.findUsedTypesRecurse(e.Condition)
		ast.findUsedTypesRecurse(e.Then)
		ast.findUsedTypesRecurse(e.Else)

	case *WhileExpr:
		ast.findUsedTypesRecurse(e.Condition)
		ast.findUsedTypesRecurse(e.Do)

	case *MatchExpr:
		ast.findUsedTypesRecurse(e.Left)
		for _, a := range e.Cases {
			ast.findUsedTypes(a.Type.target)
			ast.findUsedTypesRecurse(a.Body)
		}

	case *CallExpr:
		ast.findUsedTypesRecurse(e.Left)
		for _, a := range e.Args {
			ast.findUsedTypesRecurse(a)
		}

	case *StaticCallExpr:
		for _, a := range e.Args {
			ast.findUsedTypesRecurse(a)
		}

	case *NewExpr:
		ast.findUsedTypes(e.Type.target)

	case *VarExpr:
		ast.findUsedTypesRecurse(e.Value)
		ast.findUsedTypesRecurse(e.Expr)

	case *ChainExpr:
		ast.findUsedTypesRecurse(e.Pre)
		ast.findUsedTypesRecurse(e.Expr)

	case *NullExpr, *UnitExpr, *NameExpr, *IntegerExpr, *StringExpr, *BooleanExpr, *ThisExpr:

	default:
		panic(e)
	}
}

func (ast *AST) findOverriddenMethods(c *ClassDecl) {
	for _, f := range c.Body {
		if m, ok := f.(*MethodFeature); ok {
			if m.Override {
				ast.overriddenMethods[c.Extends.Type.target.methods[m.Name.Name]] = true
			}
		}
	}
}

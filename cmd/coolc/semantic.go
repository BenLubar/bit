package main

import (
	"errors"
	"fmt"
	"go/token"
)

var extendsAny = &ExtendsDecl{
	Type: TYPE{
		Name:   "Any",
		target: basicAny,
	},
}

var basicDummyNull = &ClassDecl{Name: TYPE{Name: "Null"}}
var basicDummyNothing = &ClassDecl{Name: TYPE{Name: "Nothing"}}

func (ast *AST) Semantic(basic bool) (err error) {
	classes := make(map[string]*ClassDecl)

	for _, c := range basicClasses {
		classes[c.Name.Name] = c
	}

	for _, c := range ast.Classes {
		if c.Name.Name == "Nothing" || c.Name.Name == "Null" {
			cp := ast.FileSet.Position(c.Name.Pos)
			return fmt.Errorf("cannot define a class with name %s at %v", c.Name.Name, cp)
		}
		if o, ok := classes[c.Name.Name]; ok && c != o {
			op := ast.FileSet.Position(o.Name.Pos)
			cp := ast.FileSet.Position(c.Name.Pos)
			return fmt.Errorf("multiple classes with name %s at %v, %v", c.Name.Name, op, cp)
		}
		classes[c.Name.Name] = c
		if c.Extends == nil {
			c.Extends = extendsAny
		}
		if c.Name.Name == "Main" {
			ast.main = c
			if len(c.Args) != 0 {
				cp := ast.FileSet.Position(c.Name.Pos)
				return fmt.Errorf("class Main cannot have constructor arguments at %v", cp)
			}
		}
	}

	if ast.main == nil {
		return fmt.Errorf("missing Main class")
	}

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	for _, c := range ast.Classes {
		ast.recurse(classes, nil, &c.Extends.Type)
	}

	for _, c := range ast.Classes {
		ast.checkExtends([]*ClassDecl{c}, c.Extends.Type.target)
	}

	for _, c := range ast.Classes {
		ast.recurse(classes, nil, c)
	}

	if !basic {
		for _, c := range ast.Classes {
			ast.buildConstructor(c)
		}
	}

	for _, c := range ast.Classes {
		c.methods = make(map[string]*MethodFeature)
		ast.buildMethodTable(c.methods, c)
	}

	for _, c := range ast.Classes {
		ast.offset(c)
	}

	for _, c := range ast.Classes {
		ast.typecheck(c, c)
	}

	return nil
}

func (ast *AST) recurse(classes map[string]*ClassDecl, ns []*ID, value interface{}) {
	recurse := func(v interface{}) {
		ast.recurse(classes, ns, v)
	}
	addNSAllowShadow := func(target interface{}, name *ID) {
		ns = append(ns, name)
		name.target = target
	}
	addNS := func(target interface{}, name *ID) {
		for _, n := range ns {
			if n.Name == name.Name {
				nPos := ast.FileSet.Position(n.Pos)
				namePos := ast.FileSet.Position(name.Pos)
				panic(fmt.Errorf("shadowing is not allowed (%s) at %v, %v", name.Name, nPos, namePos))
			}
		}
		addNSAllowShadow(target, name)
	}
	switch v := value.(type) {
	case *ID:
		for i := len(ns) - 1; i >= 0; i-- {
			n := ns[i]
			if n.Name == v.Name {
				v.target = n.target
				return
			}
		}
		pos := ast.FileSet.Position(v.Pos)
		panic(fmt.Errorf("undeclared identifier (%s) at %v", v.Name, pos))

	case *TYPE:
		if v.Name == "Nothing" {
			v.target = basicDummyNothing
			return
		}
		if v.Name == "Null" {
			v.target = basicDummyNull
			return
		}
		if c, ok := classes[v.Name]; ok {
			v.target = c
			return
		}
		pos := ast.FileSet.Position(v.Pos)
		panic(fmt.Errorf("undeclared type (%s) at %v", v.Name, pos))

	case *ClassDecl:
		recurse(&v.Name)
		for c := v.Extends.Type.target; c != basicAny; c = c.Extends.Type.target {
			for _, a := range c.Args {
				addNSAllowShadow(a, &a.Name)
			}
			for _, f := range c.Body {
				if a, ok := f.(*VarFeature); ok {
					addNSAllowShadow(a, &a.Name)
				}
			}
		}
		for i, a := range v.Args {
			a.arg = uint(i)
			addNS(a, &a.Name)
			recurse(&a.Type)
		}
		recurse(v.Extends)
		for _, f := range v.Body {
			if a, ok := f.(*VarFeature); ok {
				addNS(a, &a.Name)
			}
		}
		for _, f := range v.Body {
			recurse(f)
		}

	case *ExtendsDecl:
		for _, e := range v.Args {
			recurse(e)
		}
		for _, f := range v.Type.target.Body {
			if _, ok := f.(*NativeFeature); ok {
				pos := ast.FileSet.Position(v.Type.Pos)
				panic(fmt.Errorf("cannot extend native class (%s) at %v", v.Type.Name, pos))
			}
		}

	case *VarFeature:
		recurse(&v.Type)
		recurse(v.Value)

	case *MethodFeature:
		recurse(&v.Return)
		if v.Return.Name != v.Name.Name {
			for i, a := range v.Args {
				a.arg = uint(i)
				addNS(a, &a.Name)
				recurse(&a.Type)
			}
		}
		recurse(v.Body)

	case *BlockFeature:
		recurse(v.Expr)

	case *NativeFeature:

	case *StaticCallExpr:
		for _, a := range v.Args {
			recurse(a)
		}

	case *CallExpr:
		recurse(v.Left)
		for _, a := range v.Args {
			recurse(a)
		}

	case *NewExpr:
		recurse(&v.Type)

	case *NameExpr:
		recurse(&v.Name)

	case *ThisExpr, *StringExpr, *IntegerExpr, *BooleanExpr, *NullExpr, *UnitExpr:

	case *IfExpr:
		recurse(v.Condition)
		recurse(v.Then)
		recurse(v.Else)

	case *WhileExpr:
		recurse(v.Condition)
		recurse(v.Do)

	case *MatchExpr:
		recurse(v.Left)
		for _, c := range v.Cases {
			c.match = v
			recurse(c)
		}

	case *Case:
		addNSAllowShadow(v, &v.Name)
		recurse(&v.Type)
		recurse(v.Body)

	case *AssignExpr:
		recurse(&v.Left)
		recurse(v.Right)

	case *VarExpr:
		recurse(v.Value)
		addNS(v, &v.Name)
		recurse(&v.Type)
		recurse(v.Expr)

	case *ChainExpr:
		recurse(v.Pre)
		recurse(v.Expr)

	case NativeExpr:

	default:
		panic(v)
	}
}

func (ast *AST) offset(c *ClassDecl) {
	if c.size != 0 {
		return
	}
	if c == basicAny {
		c.size = 32 / 8
	} else {
		ast.offset(c.Extends.Type.target)
		c.size = c.Extends.Type.target.size
	}
	for _, a := range c.Args {
		a.offset = c.size
		c.size += 32 / 8
	}
	for _, f := range c.Body {
		if a, ok := f.(*VarFeature); ok {
			a.offset = c.size
			c.size += 32 / 8
		}
	}
}

func (ast *AST) typecheck(this *ClassDecl, value interface{}) {
	switch v := value.(type) {
	case *ClassDecl:
		ast.typecheck(this, v.Extends)
		for _, f := range v.Body {
			ast.typecheck(this, f)
		}

	case *ExtendsDecl:
		if len(v.Args) != len(v.Type.target.Args) {
			pos := ast.FileSet.Position(v.Type.Pos)
			panic(fmt.Errorf("argument count mismatch (%d != %d) at %v", len(v.Args), len(v.Type.target.Args), pos))
		}
		for i, e := range v.Args {
			ast.checkType(ast.checkExpr(this, e), v.Type.target.Args[i].Type.target, v.Type.target.Args[i].Type.Pos)
		}

	case *VarFeature:
		ast.checkType(ast.checkExpr(this, v.Value), v.Type.target, v.Type.Pos)
		if v.Type.target == basicDummyNothing {
			pos := ast.FileSet.Position(v.Type.Pos)
			panic(fmt.Errorf("cannot use Nothing as the type of a variable at %v", pos))
		}

	case *MethodFeature:
		ast.checkType(ast.checkExpr(this, v.Body), v.Return.target, v.Return.Pos)

	case *BlockFeature:
		ast.checkExpr(this, v.Expr)

	case *NativeFeature:

	default:
		panic(v)
	}
}

func (ast *AST) checkExtends(path []*ClassDecl, parent *ClassDecl) {
	if parent == basicAny {
		for i, c := range path {
			c.depth = uint(len(path) - i)
		}
		return
	}

	var loopError []byte
	for _, c := range path {
		if c == parent {
			loopError = []byte("loop in class heirarchy:")
		}
		if loopError != nil {
			pos := ast.FileSet.Position(c.Extends.Type.Pos)
			loopError = append(loopError, "\n\t"...)
			loopError = append(loopError, c.Name.Name...)
			loopError = append(loopError, " extends "...)
			loopError = append(loopError, c.Extends.Type.Name...)
			loopError = append(loopError, " at "...)
			loopError = append(loopError, pos.String()...)
		}
	}
	if loopError != nil {
		panic(errors.New(string(loopError)))
	}

	invalid := parent == basicDummyNothing || parent == basicDummyNull
	if !invalid {
		for _, f := range parent.Body {
			if _, ok := f.(*NativeFeature); ok {
				invalid = true
				break
			}
		}
	}
	if invalid {
		t := path[len(path)-1].Extends.Type
		pos := ast.FileSet.Position(t.Pos)
		panic(fmt.Errorf("cannot extend %s at %v", t.Name, pos))
	}

	ast.checkExtends(append(path, parent), parent.Extends.Type.target)
}

func (ast *AST) buildConstructor(c *ClassDecl) {
	var expr Expr
	expr = &StaticCallExpr{
		Name: ID{
			Name: c.Extends.Type.Name,
			Pos:  c.Extends.Type.Pos,
		},
		Args: c.Extends.Args,
	}

	for _, f := range c.Body {
		switch v := f.(type) {
		case *VarFeature:
			expr = &ChainExpr{
				Pre: expr,
				Expr: &AssignExpr{
					Left:  v.Name,
					Right: v.Value,
				},
			}

		case *BlockFeature:
			expr = &ChainExpr{
				Pre:  expr,
				Expr: v.Expr,
			}
		}
	}

	c.Body = append(c.Body, &MethodFeature{
		Name: ID{
			Name: c.Name.Name,
			Pos:  c.Name.Pos,
		},
		Args:   c.Args,
		Return: c.Name,
		Body: &ConstructorExpr{
			Args: c.Args,
			Expr: expr,
		},
	})
}

func (ast *AST) buildMethodTable(methods map[string]*MethodFeature, c *ClassDecl) {
	if c != basicAny {
		ast.buildMethodTable(methods, c.Extends.Type.target)
	}

	for _, f := range c.Body {
		if m, ok := f.(*MethodFeature); ok {
			p, ok := methods[m.Name.Name]
			if ok && !m.Override {
				pos := ast.FileSet.Position(m.Name.Pos)
				panic(fmt.Errorf("cannot shadow method in parent class (%s) at %v", m.Name.Name, pos))
			}
			if !ok && m.Override {
				pos := ast.FileSet.Position(m.Name.Pos)
				panic(fmt.Errorf("missing parent for overridden function (%s) at %v", m.Name.Name, pos))
			}
			if ok {
				if len(m.Args) != len(p.Args) {
					pos := ast.FileSet.Position(m.Name.Pos)
					panic(fmt.Errorf("argument count mismatch (%d != %d) at %v", len(m.Args), len(p.Args), pos))
				}
				for i, a := range m.Args {
					if a.Type.target != p.Args[i].Type.target {
						pos := ast.FileSet.Position(a.Name.Pos)
						panic(fmt.Errorf("argument type mismatch (%s != %s) at %v", a.Type.Name, p.Args[i].Type.Name, pos))
					}
				}
				ast.checkType(m.Return.target, p.Return.target, m.Return.Pos)
				m.offset = p.offset
			} else {
				m.offset = uint(len(methods))
			}
			methods[m.Name.Name] = m
		}
	}
}

func (ast *AST) checkExpr(this *ClassDecl, value Expr) *ClassDecl {
	switch v := value.(type) {
	case NativeExpr:
		return basicDummyNothing

	case *ConstructorExpr:
		ast.checkExpr(this, v.Expr)
		return this

	case *IfExpr:
		if cond := ast.checkExpr(this, v.Condition); cond != basicBoolean {
			pos := ast.FileSet.Position(v.Pos)
			panic(fmt.Errorf("cannot use type %s as condition at %v", cond.Name.Name, pos))
		}
		return ast.leastUpperBound(ast.checkExpr(this, v.Then), ast.checkExpr(this, v.Else))

	case *WhileExpr:
		if cond := ast.checkExpr(this, v.Condition); cond != basicBoolean {
			pos := ast.FileSet.Position(v.Pos)
			panic(fmt.Errorf("cannot use type %s as condition at %v", cond.Name.Name, pos))
		}
		ast.checkExpr(this, v.Do)
		return basicUnit

	case *MatchExpr:
		left := ast.checkExpr(this, v.Left)

		var results []*ClassDecl
		handled := make(map[*ClassDecl]*Case)

		for _, c := range v.Cases {
			ast.checkCase(c, handled)
			if !ast.lessThan(left, c.Type.target) {
				ast.checkType(c.Type.target, left, c.Type.Pos)
			}
			results = append(results, ast.checkExpr(this, c.Body))
		}

		return ast.leastUpperBound(results...)

	case *VarExpr:
		ast.checkType(ast.checkExpr(this, v.Value), v.Type.target, v.Type.Pos)
		return ast.checkExpr(this, v.Expr)

	case *ChainExpr:
		ast.checkExpr(this, v.Pre)
		return ast.checkExpr(this, v.Expr)

	case *NameExpr:
		switch t := v.Name.target.(type) {
		case *VarDecl:
			return t.Type.target

		case *VarFeature:
			return t.Type.target

		case *VarExpr:
			return t.Type.target

		case *Case:
			return t.Type.target

		default:
			panic(t)
		}

	case *ThisExpr:
		return this

	case *NullExpr:
		return basicDummyNull

	case *StringExpr:
		return basicString

	case *IntegerExpr:
		return basicInt

	case *BooleanExpr:
		return basicBoolean

	case *UnitExpr:
		return basicUnit

	case *StaticCallExpr:
		left := this.Extends.Type.target
		return ast.checkCall(this, left, &v.Name, v.Args)

	case *CallExpr:
		left := ast.checkExpr(this, v.Left)
		return ast.checkCall(this, left, &v.Name, v.Args)

	case *NewExpr:
		if v.Type.target == basicDummyNull ||
			v.Type.target == basicDummyNothing ||
			v.Type.target == basicBoolean ||
			v.Type.target == basicInt ||
			v.Type.target == basicString ||
			v.Type.target == basicSymbol ||
			v.Type.target == basicUnit {
			pos := ast.FileSet.Position(v.Type.Pos)
			panic(fmt.Errorf("cannot instantiate type %s at %v", v.Type.Name, pos))
		}

		return v.Type.target

	case *AssignExpr:
		switch t := v.Left.target.(type) {
		case *VarDecl:
			ast.checkType(ast.checkExpr(this, v.Right), t.Type.target, v.Left.Pos)

		case *VarFeature:
			ast.checkType(ast.checkExpr(this, v.Right), t.Type.target, v.Left.Pos)

		case *VarExpr:
			ast.checkType(ast.checkExpr(this, v.Right), t.Type.target, v.Left.Pos)

		default:
			panic(t)
		}

		return basicUnit

	default:
		panic(v)
	}
}

func (ast *AST) checkCall(this, left *ClassDecl, name *ID, args []Expr) *ClassDecl {
	m, ok := left.methods[name.Name]
	if !ok {
		pos := ast.FileSet.Position(name.Pos)
		panic(fmt.Errorf("type %s has no %s method at %v", left.Name.Name, name.Name, pos))
	}

	name.target = m

	if len(args) != len(m.Args) {
		pos := ast.FileSet.Position(name.Pos)
		panic(fmt.Errorf("argument count mismatch (%d != %d) at %v", len(args), len(m.Args), pos))
	}

	for i := range args {
		ast.checkType(ast.checkExpr(this, args[i]), m.Args[i].Type.target, name.Pos)
	}

	return m.Return.target
}

func (ast *AST) checkType(left, right *ClassDecl, p token.Pos) {
	if !ast.lessThan(left, right) {
		pos := ast.FileSet.Position(p)
		panic(fmt.Errorf("cannot use type %s as %s at %v", left.Name.Name, right.Name.Name, pos))
	}
}

func (ast *AST) checkCase(c *Case, handled map[*ClassDecl]*Case) {
	any := false
	check := func(h *ClassDecl) {
		if _, ok := handled[h]; !ok && ast.lessThan(h, c.Type.target) {
			c.classes = append(c.classes, h)
			handled[h] = c
			any = true
		}
	}
	check(basicDummyNothing)
	check(basicDummyNull)
	for _, h := range basicClasses {
		check(h)
	}
	for _, h := range ast.Classes {
		check(h)
	}

	if !any {
		pos := ast.FileSet.Position(c.Type.Pos)
		panic(fmt.Errorf("inaccessible case for type %s at %v", c.Type.Name, pos))
	}
}

func (ast *AST) lessThan(left, right *ClassDecl) bool {
	// S-Nothing
	if left == basicDummyNothing {
		return true
	}
	// S-Self
	if left == right {
		return true
	}
	// S-Null
	if left == basicDummyNull {
		return right != basicDummyNothing && right != basicBoolean && right != basicInt && right != basicUnit
	}
	// S-Extends
	for c := left; ; c = c.Extends.Type.target {
		if c == right {
			return true
		}
		if c == basicAny {
			break
		}
	}
	return false
}

func (ast *AST) leastUpperBound(classes ...*ClassDecl) *ClassDecl {
	left := classes[0]
	for _, right := range classes[1:] {
		// G-Compare
		if ast.lessThan(left, right) {
			left = right
			continue
		}

		// G-Commute
		if ast.lessThan(right, left) {
			continue
		}

		// G-Extends
		if left.depth == 0 || right.depth == 0 {
			// null and a non-nullable type
			return basicAny
		}
		for left.depth > right.depth {
			left = left.Extends.Type.target
		}
		for right.depth > left.depth {
			right = right.Extends.Type.target
		}
		for left != right {
			left = left.Extends.Type.target
			right = right.Extends.Type.target
		}
	}
	return left
}

package main

import "fmt"

var extendsAny = &ExtendsDecl{
	Type: TYPE{
		Name:   "Any",
		target: basicAny,
	},
}

func (ast *AST) Semantic() (err error) {
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
		ast.recurse(classes, nil, c)
	}

	panic("unimplemented: typecheck")
}

func (ast *AST) recurse(classes map[string]*ClassDecl, ns []*ID, value interface{}) {
	recurse := func(v interface{}) {
		ast.recurse(classes, ns, v)
	}
	addNS := func(target interface{}, name *ID) {
		for _, n := range ns {
			if n.Name == name.Name {
				nPos := ast.FileSet.Position(n.Pos)
				namePos := ast.FileSet.Position(name.Pos)
				panic(fmt.Errorf("shadowing is not allowed (%s) at %v, %v", name.Name, nPos, namePos))
			}
		}
		ns = append(ns, name)
		name.target = target
	}
	switch v := value.(type) {
	case *ID:
		for _, n := range ns {
			if n.Name == v.Name {
				v.target = n.target
				return
			}
		}
		pos := ast.FileSet.Position(v.Pos)
		panic(fmt.Errorf("undeclared identifier (%s) at %v", v.Name, pos))

	case *TYPE:
		if v.Name == "Nothing" || v.Name == "Null" {
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
		for _, a := range v.Args {
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
		recurse(&v.Type)
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
		for _, a := range v.Args {
			addNS(a, &a.Name)
			recurse(&a.Type)
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
		for _, a := range v.Args {
			recurse(a)
		}

	case *NameExpr:
		recurse(&v.Name)

	case *ThisExpr, *StringExpr, *IntegerExpr, *BooleanExpr, *NullExpr:

	case *IfExpr:
		recurse(v.Condition)
		recurse(v.Then)
		recurse(v.Else)

	case *WhileExpr:
		recurse(v.Condition)
		recurse(v.Do)

	case *MatchExpr:
		recurse(v.Left)
		recurse(v.Cases)

	case *Cases:
		for _, c := range v.Cases {
			recurse(c)
		}
		if v.Null != nil {
			recurse(v.Null)
		}

	case *Case:
		addNS(v, &v.Name)
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

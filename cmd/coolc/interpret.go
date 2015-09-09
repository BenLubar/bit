package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type interpStop string

func (err interpStop) Error() string { return string(err) }

func (ast *AST) Interpret(in io.Reader, out io.Writer) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprint(out, r.(interpStop))
		}
	}()

	interp := &interpreter{
		In:  bufio.NewReader(in),
		Out: bufio.NewWriter(out),

		Symbol: make(map[string]*interpObject),
	}
	interp.New(ast.main).Call(ast.main.methods["Main"], nil)
}

type interpreter struct {
	In  *bufio.Reader
	Out *bufio.Writer

	Symbol map[string]*interpObject
}

type interpObject struct {
	Interp *interpreter
	Class  *ClassDecl

	// for user-defined types
	Attr map[interface{}]*interpObject

	// for basic.cool types
	ArrayAny []*interpObject
	String   string
	Int      int32
	Boolean  bool
}

func (interp *interpreter) New(c *ClassDecl) *interpObject {
	o := &interpObject{
		Interp: interp,
		Class:  c,
		Attr:   make(map[interface{}]*interpObject),
	}
	for p := c; p != basicAny; p = p.Extends.Type.target {
		for _, a := range p.Args {
			o.Attr[a] = interp.Default(a.Type.target)
		}
		for _, f := range p.Body {
			if v, ok := f.(*VarFeature); ok {
				o.Attr[v] = interp.Default(v.Type.target)
			}
		}
	}
	return o
}

func (interp *interpreter) Default(c *ClassDecl) *interpObject {
	switch c {
	case basicUnit, basicInt, basicBoolean:
		return &interpObject{
			Interp: interp,
			Class:  c,
		}

	default:
		return &interpObject{
			Interp: interp,
			Class:  basicDummyNull,
		}
	}
}

func (o *interpObject) Call(m *MethodFeature, args []*interpObject) *interpObject {
	if _, ok := m.Body.(NativeExpr); ok {
		switch m.Name.Name {
		case "out":
			_, err := o.Interp.Out.WriteString(args[0].String)
			if err != nil {
				panic(err)
			}
			err = o.Interp.Out.Flush()
			if err != nil {
				panic(err)
			}
			return o

		case "in":
			s, err := o.Interp.In.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					return &interpObject{
						Interp: o.Interp,
						Class:  basicDummyNull,
					}
				}
				panic(err)
			}

			return &interpObject{
				Interp: o.Interp,
				Class:  basicString,
				String: strings.TrimSuffix(s, "\n"),
			}

		case "abort":
			_, err := o.Interp.Out.WriteString("Abort: " + args[0].String)
			if err != nil {
				panic(err)
			}
			err = o.Interp.Out.Flush()
			if err != nil {
				panic(err)
			}
			panic(interpStop(""))

		case "symbol":
			if args[0].Class != basicString {
				panic(interpStop("null pointer dereference\n"))
			}
			if s, ok := o.Interp.Symbol[args[0].String]; ok {
				return s
			}
			s := &interpObject{
				Interp: o.Interp,
				Class:  basicSymbol,
				Attr: map[interface{}]*interpObject{
					basicSymbolHash: &interpObject{
						Interp: o.Interp,
						Class:  basicInt,
						Int:    int32(len(o.Interp.Symbol)),
					},
					basicSymbolName: args[0],
				},
			}
			o.Interp.Symbol[args[0].String] = s
			return s

		case "symbol_name":
			if args[0].Class != basicSymbol {
				panic(interpStop("null pointer dereference\n"))
			}
			return args[0].Attr[basicSymbolName]

		case "toString":
			var s string
			switch o.Class {
			case basicInt:
				s = strconv.FormatInt(int64(o.Int), 10)
			default:
				s = o.Class.Name.Name
			}
			return &interpObject{
				Interp: o.Interp,
				Class:  basicString,
				String: s,
			}

		case "equals":
			var equal bool
			switch o.Class {
			case basicInt:
				equal = args[0].Class == basicInt && o.Int == args[0].Int
			case basicString:
				equal = args[0].Class == basicString && o.String == args[0].String
			case basicBoolean:
				equal = args[0].Class == basicBoolean && o.Boolean == args[0].Boolean
			case basicUnit:
				equal = args[0].Class == basicUnit
			default:
				equal = o == args[0]
			}
			return &interpObject{
				Interp:  o.Interp,
				Class:   basicBoolean,
				Boolean: equal,
			}

		case "is_null":
			return &interpObject{
				Interp:  o.Interp,
				Class:   basicBoolean,
				Boolean: args[0].Class == basicDummyNull,
			}

		case "_lowest_bit":
			return &interpObject{
				Interp:  o.Interp,
				Class:   basicBoolean,
				Boolean: o.Int&1 == 1,
			}

		case "_lsh":
			return &interpObject{
				Interp: o.Interp,
				Class:  basicInt,
				Int:    int32(uint32(o.Int) << 1),
			}

		case "_rsh":
			return &interpObject{
				Interp: o.Interp,
				Class:  basicInt,
				Int:    int32(uint32(o.Int) >> 1),
			}

		case "_add":
			return &interpObject{
				Interp: o.Interp,
				Class:  basicInt,
				Int:    o.Int + args[0].Int,
			}

		case "_negative":
			return &interpObject{
				Interp: o.Interp,
				Class:  basicInt,
				Int:    -o.Int,
			}

		case "_less_unsigned":
			return &interpObject{
				Interp:  o.Interp,
				Class:   basicBoolean,
				Boolean: uint32(o.Int) < uint32(args[0].Int),
			}

		case "_less":
			return &interpObject{
				Interp:  o.Interp,
				Class:   basicBoolean,
				Boolean: o.Int < args[0].Int,
			}

		case "_less_equal":
			return &interpObject{
				Interp:  o.Interp,
				Class:   basicBoolean,
				Boolean: o.Int <= args[0].Int,
			}

		case "_check_divide_by_zero":
			if o.Int == 0 {
				panic(interpStop("division by zero\n"))
			}

			return &interpObject{
				Interp: o.Interp,
				Class:  basicUnit,
			}

		case "substring":
			return &interpObject{
				Interp: o.Interp,
				Class:  basicString,
				String: o.String[args[0].Int : args[1].Int-args[0].Int],
			}

		case "charAt":
			return &interpObject{
				Interp: o.Interp,
				Class:  basicInt,
				Int:    int32(o.String[args[0].Int]),
			}

		case "concat":
			if args[0].Class != basicString {
				panic(interpStop("null pointer dereference\n"))
			}
			return &interpObject{
				Interp: o.Interp,
				Class:  basicString,
				String: o.String + args[0].String,
			}

		case "ArrayAny":
			if args[0].Int < 0 {
				panic(interpStop("index out of range\n"))
			}
			o.ArrayAny = make([]*interpObject, args[0].Int)
			for i := range o.ArrayAny {
				o.ArrayAny[i] = &interpObject{
					Interp: o.Interp,
					Class:  basicDummyNull,
				}
			}
			o.Attr[basicArrayAnyLength] = args[0]
			return o

		case "get":
			if args[0].Int < 0 || len(o.ArrayAny) <= int(args[0].Int) {
				panic(interpStop("index out of range\n"))
			}
			return o.ArrayAny[args[0].Int]

		case "set":
			if args[0].Int < 0 || len(o.ArrayAny) <= int(args[0].Int) {
				panic(interpStop("index out of range\n"))
			}
			old := o.ArrayAny[args[0].Int]
			o.ArrayAny[args[0].Int] = args[1]
			return old

		default:
			panic(o.Class.Name.Name + "." + m.Name.Name)
		}
	}

	vars := make(map[interface{}]*interpObject)

	for i, a := range m.Args {
		vars[a] = args[i]
	}

	return o.Expr(vars, m.Body)
}

func (o *interpObject) Expr(vars map[interface{}]*interpObject, expr Expr) *interpObject {
	switch e := expr.(type) {
	case *ConstructorExpr:
		for _, a := range e.Args {
			o.Attr[a] = vars[a]
			delete(vars, a)
		}
		o.Expr(vars, e.Expr)
		return o

	case *NameExpr:
		if v, ok := vars[e.Name.target]; ok {
			return v
		} else if v, ok = o.Attr[e.Name.target]; ok {
			return v
		} else if e.Name.target == basicStringLength {
			return &interpObject{
				Interp: o.Interp,
				Class:  basicInt,
				Int:    int32(len(o.String)),
			}
		} else {
			panic("*NameExpr " + e.Name.Name)
		}

	case *ThisExpr:
		return o

	case *NullExpr:
		return &interpObject{
			Interp: o.Interp,
			Class:  basicDummyNull,
		}

	case *UnitExpr:
		return &interpObject{
			Interp: o.Interp,
			Class:  basicUnit,
		}

	case *BooleanExpr:
		return &interpObject{
			Interp:  o.Interp,
			Class:   basicBoolean,
			Boolean: e.B,
		}

	case *IntegerExpr:
		return &interpObject{
			Interp: o.Interp,
			Class:  basicInt,
			Int:    e.N,
		}

	case *StringExpr:
		return &interpObject{
			Interp: o.Interp,
			Class:  basicString,
			String: e.S,
		}

	case *NewExpr:
		return o.Interp.New(e.Type.target)

	case *AssignExpr:
		if _, ok := vars[e.Left.target]; ok {
			vars[e.Left.target] = o.Expr(vars, e.Right)
		} else if _, ok := o.Attr[e.Left.target]; ok {
			o.Attr[e.Left.target] = o.Expr(vars, e.Right)
		} else {
			panic("*AssignExpr " + e.Left.Name)
		}
		return &interpObject{
			Interp: o.Interp,
			Class:  basicUnit,
		}

	case *IfExpr:
		if o.Expr(vars, e.Condition).Boolean {
			return o.Expr(vars, e.Then)
		}
		return o.Expr(vars, e.Else)

	case *WhileExpr:
		for o.Expr(vars, e.Condition).Boolean {
			o.Expr(vars, e.Do)
		}
		return &interpObject{
			Interp: o.Interp,
			Class:  basicUnit,
		}

	case *VarExpr:
		vars[e] = o.Expr(vars, e.Value)
		defer delete(vars, e)
		return o.Expr(vars, e.Expr)

	case *ChainExpr:
		o.Expr(vars, e.Pre)
		return o.Expr(vars, e.Expr)

	case *CallExpr:
		this := o.Expr(vars, e.Left)
		if this.Class == basicDummyNull {
			panic(interpStop("null pointer dereference\n"))
		}
		m, ok := this.Class.methods[e.Name.Name]
		if !ok {
			panic("*CallExpr " + e.Name.Name)
		}
		return this.Call(m, o.Exprs(vars, e.Args))

	case *StaticCallExpr:
		return o.Call(e.Name.target.(*MethodFeature), o.Exprs(vars, e.Args))

	case *MatchExpr:
		left := o.Expr(vars, e.Left)
		for _, c := range e.Cases {
			for _, cl := range c.classes {
				if left.Class == cl {
					vars[c] = left
					defer delete(vars, c)
					return o.Expr(vars, c.Body)
				}
			}
		}
		panic(interpStop("no case for " + left.Class.Name.Name + "\n"))

	default:
		panic(e)
	}
}

func (o *interpObject) Exprs(vars map[interface{}]*interpObject, exprs []Expr) []*interpObject {
	result := make([]*interpObject, len(exprs))
	for i, e := range exprs {
		result[i] = o.Expr(vars, e)
	}
	return result
}

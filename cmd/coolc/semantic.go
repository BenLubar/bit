package main

import (
	"fmt"

	"github.com/BenLubar/bit/bitgen"
)

var extendsAny = &ExtendsDecl{
	Type: TYPE{
		Name:   "Any",
		target: basicAny,
	},
}

var basicAny = &ClassDecl{
	Name: TYPE{
		Name: "Any",
	},
	Body: []Feature{
		&MethodFeature{
			Name: ID{
				Name: "toString",
			},
			Return: TYPE{
				Name: "String",
			},
			Body: NativeExpr(func(w *writer, start, end bitgen.Line) {
				w.EndStack()

				next := w.ReserveLine()
				w.Load(start, w.This, w.This, 0, next)
				start = next

				next = w.ReserveLine()
				w.Load(start, w.Return, w.This, 32/8, next)
				start = next

				w.PopStack(start, end)
			}),
		},
		&MethodFeature{
			Name: ID{
				Name: "equals",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "x",
					},
					Type: TYPE{
						Name: "Any",
					},
				},
			},
			Return: TYPE{
				Name: "Boolean",
			},
			Body: NativeExpr(func(w *writer, start, end bitgen.Line) {
				w.EndStack()

				same, different := w.ReserveLine(), w.ReserveLine()
				w.CmpReg(start, w.This.Num, w.StackOffset(w.Arg(0)), same, different)

				next := w.ReserveLine()
				w.CopyReg(same, w.Return, w.True, next)
				w.CopyReg(same, w.Return, w.False, next)
				start = next

				w.PopStack(start, end)
			}),
		},
	},
}

var basicIO = &ClassDecl{
	Name: TYPE{
		Name: "IO",
	},
	Body: []Feature{
		&MethodFeature{
			Name: ID{
				Name: "abort",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "message",
					},
					Type: TYPE{
						Name: "String",
					},
				},
			},
			Return: TYPE{
				Name: "Nothing",
			},
			Body: NativeExpr(func(w *writer, start, end bitgen.Line) {
				w.EndStack()

				w.PrintStringArg(start, 0, 0)
			}),
		},
		&MethodFeature{
			Name: ID{
				Name: "out",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "arg",
					},
					Type: TYPE{
						Name: "String",
					},
				},
			},
			Return: TYPE{
				Name: "IO",
			},
			Body: NativeExpr(func(w *writer, start, end bitgen.Line) {
				w.EndStack()

				next := w.ReserveLine()
				w.CopyReg(start, w.Return, w.This, next)
				start = next

				next = w.ReserveLine()
				w.PrintStringArg(start, 0, next)
				start = next

				w.PopStack(start, end)
			}),
		},
		&MethodFeature{
			Name: ID{
				Name: "is_null",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "arg",
					},
					Type: TYPE{
						Name: "Any",
					},
				},
			},
			Return: TYPE{
				Name: "Boolean",
			},
			Body: &MatchExpr{
				Left: &NameExpr{
					Name: ID{
						Name: "arg",
					},
				},
				Cases: &Cases{
					Cases: []*Case{
						{
							Name: ID{
								Name: "x",
							},
							Type: TYPE{
								Name: "Any",
							},
							Body: &BooleanExpr{
								B: false,
							},
						},
					},
					Null: &BooleanExpr{
						B: true,
					},
				},
			},
		},
		&MethodFeature{
			Name: ID{
				Name: "out_any",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "arg",
					},
					Type: TYPE{
						Name: "Any",
					},
				},
			},
			Return: TYPE{
				Name: "IO",
			},
			Body: &SelfCallExpr{
				Name: ID{
					Name: "out",
				},
				Args: []Expr{
					&IfExpr{
						Condition: &SelfCallExpr{
							Name: ID{
								Name: "is_null",
							},
							Args: []Expr{
								&NameExpr{
									Name: ID{
										Name: "arg",
									},
								},
							},
						},
						Then: &StringExpr{
							S: "null",
						},
						Else: &CallExpr{
							Left: &NameExpr{
								Name: ID{
									Name: "arg",
								},
							},
							Name: ID{
								Name: "toString",
							},
						},
					},
				},
			},
		},
		&MethodFeature{
			Name: ID{
				Name: "in",
			},
			Return: TYPE{
				Name: "String",
			},
			Body: NativeExpr(func(w *writer, start, end bitgen.Line) {
				panic("unimplemented")
			}),
		},
		&MethodFeature{
			Name: ID{
				Name: "symbol",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "name",
					},
					Type: TYPE{
						Name: "String",
					},
				},
			},
			Return: TYPE{
				Name: "Symbol",
			},
			Body: NativeExpr(func(w *writer, start, end bitgen.Line) {
				panic("unimplemented")
			}),
		},
		&MethodFeature{
			Name: ID{
				Name: "symbol_name",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "sym",
					},
					Type: TYPE{
						Name: "Symbol",
					},
				},
			},
			Return: TYPE{
				Name: "String",
			},
			Body: NativeExpr(func(w *writer, start, end bitgen.Line) {
				panic("unimplemented")
			}),
		},
	},
}

var basicUnit = &ClassDecl{
	Name: TYPE{
		Name: "Unit",
	},
	Body: []Feature{
		&NativeFeature{},
	},
}

var basicInt = &ClassDecl{
	Name: TYPE{
		Name: "Int",
	},
	Body: []Feature{
		&NativeFeature{},
		&MethodFeature{
			Override: true,
			Name: ID{
				Name: "toString",
			},
			Return: TYPE{
				Name: "String",
			},
			Body: NativeExpr(func(w *writer, start, end bitgen.Line) {
				panic("unimplemented")
			}),
		},
		&MethodFeature{
			Override: true,
			Name: ID{
				Name: "equals",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "other",
					},
					Type: TYPE{
						Name: "Any",
					},
				},
			},
			Return: TYPE{
				Name: "Boolean",
			},
			Body: NativeExpr(func(w *writer, start, end bitgen.Line) {
				w.EndStack()

				next := w.ReserveLine()
				w.Load(start, w.General[0], w.Stack, w.Arg(0), next)
				start = next

				same, different := w.ReserveLine(), w.ReserveLine()

				next = w.ReserveLine()
				w.Cmp(start, w.General[0].Num, 0, different, next)
				start = next

				next = w.ReserveLine()
				w.CmpReg(start, bitgen.Integer{bitgen.ValueAt{w.This.Ptr}, 32}, bitgen.Integer{bitgen.ValueAt{w.General[0].Ptr}, 32}, next, different)
				start = next

				w.CmpReg(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.This.Ptr, 32}}, 32}, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[0].Ptr, 32}}, 32}, same, different)

				next = w.ReserveLine()
				w.CopyReg(same, w.Return, w.True, next)
				w.CopyReg(same, w.Return, w.False, next)
				start = next

				w.PopStack(start, end)
			}),
		},
	},
}

var basicBoolean = &ClassDecl{
	Name: TYPE{
		Name: "Boolean",
	},
	Body: []Feature{
		&NativeFeature{},
		&MethodFeature{
			Override: true,
			Name: ID{
				Name: "toString",
			},
			Return: TYPE{
				Name: "String",
			},
			Body: &IfExpr{
				Condition: &ThisExpr{},
				Then: &StringExpr{
					S: "true",
				},
				Else: &StringExpr{
					S: "false",
				},
			},
		},
	},
}

var basicStringLength = &VarFeature{
	VarDecl: VarDecl{
		Name: ID{
			Name: "length",
		},
		Type: TYPE{
			Name: "Int",
		},
	},
	Value: &IntegerExpr{
		N: 0,
	},
}

var basicString = &ClassDecl{
	Name: TYPE{
		Name: "String",
	},
	Body: []Feature{
		basicStringLength,
		&NativeFeature{},
		&MethodFeature{
			Override: true,
			Name: ID{
				Name: "toString",
			},
			Return: TYPE{
				Name: "String",
			},
			Body: &ThisExpr{},
		},
		&MethodFeature{
			Override: true,
			Name: ID{
				Name: "equals",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "other",
					},
					Type: TYPE{
						Name: "Any",
					},
				},
			},
			Return: TYPE{
				Name: "Boolean",
			},
			Body: NativeExpr(func(w *writer, start, end bitgen.Line) {
				panic("unimplemented")
			}),
		},
		&MethodFeature{
			Name: ID{
				Name: "length",
			},
			Return: TYPE{
				Name: "Int",
			},
			Body: &NameExpr{
				Name: ID{
					Name: "length",
				},
			},
		},
		&MethodFeature{
			Name: ID{
				Name: "concat",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "arg",
					},
					Type: TYPE{
						Name: "String",
					},
				},
			},
			Return: TYPE{
				Name: "String",
			},
			Body: NativeExpr(func(w *writer, start, end bitgen.Line) {
				panic("unimplemented")
			}),
		},
		&MethodFeature{
			Name: ID{
				Name: "substring",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "start",
					},
					Type: TYPE{
						Name: "Int",
					},
				},
				{
					Name: ID{
						Name: "end",
					},
					Type: TYPE{
						Name: "Int",
					},
				},
			},
			Return: TYPE{
				Name: "String",
			},
			Body: NativeExpr(func(w *writer, start, end bitgen.Line) {
				panic("unimplemented")
			}),
		},
		&MethodFeature{
			Name: ID{
				Name: "charAt",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "index",
					},
					Type: TYPE{
						Name: "Int",
					},
				},
			},
			Return: TYPE{
				Name: "Int",
			},
			Body: NativeExpr(func(w *writer, start, end bitgen.Line) {
				panic("unimplemented")
			}),
		},
		&MethodFeature{
			Name: ID{
				Name: "indexOf",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "sub",
					},
					Type: TYPE{
						Name: "String",
					},
				},
			},
			Return: TYPE{
				Name: "Int",
			},
			Body: &VarExpr{
				VarFeature: VarFeature{
					VarDecl: VarDecl{
						Name: ID{
							Name: "n",
						},
						Type: TYPE{
							Name: "Int",
						},
					},
					Value: &CallExpr{
						Left: &NameExpr{
							Name: ID{
								Name: "sub",
							},
						},
						Name: ID{
							Name: "length",
						},
					},
				},
				Expr: &VarExpr{
					VarFeature: VarFeature{
						VarDecl: VarDecl{
							Name: ID{
								Name: "diff",
							},
							Type: TYPE{
								Name: "Int",
							},
						},
						Value: &SubtractExpr{
							Left: &NameExpr{
								Name: ID{
									Name: "length",
								},
							},
							Right: &NameExpr{
								Name: ID{
									Name: "n",
								},
							},
						},
					},
					Expr: &VarExpr{
						VarFeature: VarFeature{
							VarDecl: VarDecl{
								Name: ID{
									Name: "i",
								},
								Type: TYPE{
									Name: "Int",
								},
							},
							Value: &IntegerExpr{
								N: 0,
							},
						},
						Expr: &VarExpr{
							VarFeature: VarFeature{
								VarDecl: VarDecl{
									Name: ID{
										Name: "result",
									},
									Type: TYPE{
										Name: "Int",
									},
								},
								Value: &NegativeExpr{
									Right: &IntegerExpr{
										N: 1,
									},
								},
							},
							Expr: &ChainExpr{
								Pre: &WhileExpr{
									Condition: &LessThanOrEqualExpr{
										Left: &NameExpr{
											Name: ID{
												Name: "i",
											},
										},
										Right: &NameExpr{
											Name: ID{
												Name: "diff",
											},
										},
									},
									Do: &IfExpr{
										Condition: &CallExpr{
											Left: &SelfCallExpr{
												Name: ID{
													Name: "substring",
												},
												Args: []Expr{
													&NameExpr{
														Name: ID{
															Name: "i",
														},
													},
													&AddExpr{
														Left: &NameExpr{
															Name: ID{
																Name: "i",
															},
														},
														Right: &NameExpr{
															Name: ID{
																Name: "n",
															},
														},
													},
												},
											},
											Name: ID{
												Name: "equals",
											},
											Args: []Expr{
												&NameExpr{
													Name: ID{
														Name: "sub",
													},
												},
											},
										},
										Then: &ChainExpr{
											Pre: &AssignExpr{
												Left: ID{
													Name: "result",
												},
												Right: &NameExpr{
													Name: ID{
														Name: "i",
													},
												},
											},
											Expr: &AssignExpr{
												Left: ID{
													Name: "i",
												},
												Right: &AddExpr{
													Left: &NameExpr{
														Name: ID{
															Name: "diff",
														},
													},
													Right: &IntegerExpr{
														N: 1,
													},
												},
											},
										},
										Else: &AssignExpr{
											Left: ID{
												Name: "i",
											},
											Right: &AddExpr{
												Left: &NameExpr{
													Name: ID{
														Name: "i",
													},
												},
												Right: &IntegerExpr{
													N: 1,
												},
											},
										},
									},
								},
								Expr: &NameExpr{
									Name: ID{
										Name: "result",
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

var basicSymbol = &ClassDecl{
	Name: TYPE{
		Name: "Symbol",
	},
	Body: []Feature{
		&VarFeature{
			VarDecl: VarDecl{
				Name: ID{
					Name: "name",
				},
				Type: TYPE{
					Name: "String",
				},
			},
			Value: &StringExpr{
				S: "",
			},
		},
		&VarFeature{
			VarDecl: VarDecl{
				Name: ID{
					Name: "hash",
				},
				Type: TYPE{
					Name: "Int",
				},
			},
			Value: &IntegerExpr{
				N: 0,
			},
		},
		&VarFeature{
			VarDecl: VarDecl{
				Name: ID{
					Name: "next",
				},
				Type: TYPE{
					Name: "Symbol",
				},
			},
			Value: &NullExpr{},
		},
		&NativeFeature{},
		&MethodFeature{
			Override: true,
			Name: ID{
				Name: "toString",
			},
			Return: TYPE{
				Name: "String",
			},
			Body: &CallExpr{
				Left: &StringExpr{
					S: "'",
				},
				Name: ID{
					Name: "concat",
				},
				Args: []Expr{
					&NameExpr{
						Name: ID{
							Name: "name",
						},
					},
				},
			},
		},
		&MethodFeature{
			Name: ID{
				Name: "hashCode",
			},
			Return: TYPE{
				Name: "Int",
			},
			Body: &NameExpr{
				Name: ID{
					Name: "hash",
				},
			},
		},
	},
}

var basicArrayAnyLength = &VarDecl{
	Name: ID{
		Name: "length",
	},
	Type: TYPE{
		Name: "Int",
	},
}

var basicArrayAny = &ClassDecl{
	Name: TYPE{
		Name: "ArrayAny",
	},
	Args: []*VarDecl{
		basicArrayAnyLength,
	},
	Body: []Feature{
		&NativeFeature{},
		&MethodFeature{
			Name: ID{
				Name: "length",
			},
			Return: TYPE{
				Name: "Int",
			},
			Body: &NameExpr{
				Name: ID{
					Name: "length",
				},
			},
		},
		&MethodFeature{
			Name: ID{
				Name: "resize",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "s",
					},
					Type: TYPE{
						Name: "Int",
					},
				},
			},
			Return: TYPE{
				Name: "ArrayAny",
			},
			Body: NativeExpr(func(w *writer, start, end bitgen.Line) {
				panic("unimplemented")
			}),
		},
		&MethodFeature{
			Name: ID{
				Name: "get",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "index",
					},
					Type: TYPE{
						Name: "Int",
					},
				},
			},
			Return: TYPE{
				Name: "Any",
			},
			Body: NativeExpr(func(w *writer, start, end bitgen.Line) {
				w.EndStack()

				next := w.ReserveLine()
				w.Load(start, w.General[0], w.Stack, w.Arg(0), next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, w.General[0].Num, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[0].Ptr, 32}}, 32}, next)
				start = next

				next = w.ReserveLine()
				w.Load(start, w.General[1], w.This, basicArrayAnyLength.offset, next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, w.General[1].Num, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[1].Ptr, 32}}, 32}, next)
				start = next

				next = w.ReserveLine()
				w.LessThanUnsigned(start, w.General[0].Num, w.General[1].Num, next, w.IndexRange, w.IndexRange)
				start = next

				next = w.ReserveLine()
				w.CopyReg(start, w.Return, w.This, next)
				start = next

				loop, done := start, w.ReserveLine()
				next = w.ReserveLine()
				w.Decrement(start, w.General[0].Num, next, done)
				start = next

				for i := 0; i < 32/8; i++ {
					next = w.ReserveLine()
					w.Increment(start, w.Return.Num, next, 0)
					start = next
				}

				w.Assign(start, w.Return.Ptr, bitgen.Offset{w.Return.Ptr, 32}, loop)

				next = w.ReserveLine()
				w.Load(done, w.Return, w.Return, basicArrayAnyLength.offset+32/8, next)
				start = next

				w.PopStack(start, end)
			}),
		},
		&MethodFeature{
			Name: ID{
				Name: "set",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "index",
					},
					Type: TYPE{
						Name: "Int",
					},
				},
				{
					Name: ID{
						Name: "obj",
					},
					Type: TYPE{
						Name: "Any",
					},
				},
			},
			Return: TYPE{
				Name: "Any",
			},
			Body: NativeExpr(func(w *writer, start, end bitgen.Line) {
				w.EndStack()

				next := w.ReserveLine()
				w.Load(start, w.General[0], w.Stack, w.Arg(0), next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, w.General[0].Num, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[0].Ptr, 32}}, 32}, next)
				start = next

				next = w.ReserveLine()
				w.Load(start, w.General[1], w.This, basicArrayAnyLength.offset, next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, w.General[1].Num, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[1].Ptr, 32}}, 32}, next)
				start = next

				next = w.ReserveLine()
				w.LessThanUnsigned(start, w.General[0].Num, w.General[1].Num, next, w.IndexRange, w.IndexRange)
				start = next

				next = w.ReserveLine()
				w.CopyReg(start, w.General[1], w.This, next)
				start = next

				loop, done := start, w.ReserveLine()
				next = w.ReserveLine()
				w.Decrement(start, w.General[0].Num, next, done)
				start = next

				for i := 0; i < 32/8; i++ {
					next = w.ReserveLine()
					w.Increment(start, w.General[1].Num, next, 0)
					start = next
				}

				w.Assign(start, w.General[1].Ptr, bitgen.Offset{w.General[1].Ptr, 32}, loop)

				next = w.ReserveLine()
				w.Load(done, w.Return, w.General[1], basicArrayAnyLength.offset+32/8, next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[1].Ptr, basicArrayAnyLength.offset + 32/8}}, 32}, w.StackOffset(w.Arg(1)), next)
				start = next

				w.PopStack(start, end)
			}),
		},
	},
}

var basicClasses = []*ClassDecl{
	basicAny,
	basicIO,
	basicUnit,
	basicInt,
	basicBoolean,
	basicString,
	basicSymbol,
	basicArrayAny,
}

func init() {
	basicAST := &AST{
		Classes: append([]*ClassDecl{
			&ClassDecl{
				Name: TYPE{
					Name: "Main",
				},
			},
		}, basicClasses...),
	}
	if err := basicAST.Semantic(); err != nil {
		panic(err)
	}
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
		}
	}

	if ast.main == nil {
		return fmt.Errorf("missing Main class")
	}

	var recurse func([]*ID, interface{})
	recurse = func(ns []*ID, value interface{}) {
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
			recurse(ns, &v.Name)
			for _, a := range v.Args {
				addNS(a, &a.Name)
				recurse(ns, &a.Type)
			}
			recurse(ns, v.Extends)
			for _, f := range v.Body {
				if a, ok := f.(*VarFeature); ok {
					addNS(a, &a.Name)
				}
			}
			for _, f := range v.Body {
				recurse(ns, f)
			}

		case *ExtendsDecl:
			recurse(ns, &v.Type)
			for _, e := range v.Args {
				recurse(ns, e)
			}

		case *VarFeature:
			recurse(ns, &v.Type)
			recurse(ns, v.Value)

		case *MethodFeature:
			recurse(ns, &v.Return)
			for _, a := range v.Args {
				addNS(a, &a.Name)
				recurse(ns, &a.Type)
			}
			recurse(ns, v.Body)

		case *BlockFeature:
			recurse(ns, v.Expr)

		case *NativeFeature:

		case *SelfCallExpr:
			for _, a := range v.Args {
				recurse(ns, a)
			}

		case *CallExpr:
			recurse(ns, v.Left)
			for _, a := range v.Args {
				recurse(ns, a)
			}

		case *NegativeExpr:
			recurse(ns, v.Right)

		case *NotExpr:
			recurse(ns, v.Right)

		case *AddExpr:
			recurse(ns, v.Left)
			recurse(ns, v.Right)

		case *SubtractExpr:
			recurse(ns, v.Left)
			recurse(ns, v.Right)

		case *MultiplyExpr:
			recurse(ns, v.Left)
			recurse(ns, v.Right)

		case *DivideExpr:
			recurse(ns, v.Left)
			recurse(ns, v.Right)

		case *LessThanExpr:
			recurse(ns, v.Left)
			recurse(ns, v.Right)

		case *LessThanOrEqualExpr:
			recurse(ns, v.Left)
			recurse(ns, v.Right)

		case *NameExpr:
			recurse(ns, &v.Name)

		case *ThisExpr:

		case *StringExpr, *IntegerExpr, *BooleanExpr, *NullExpr:

		case *IfExpr:
			recurse(ns, v.Condition)
			recurse(ns, v.Then)
			recurse(ns, v.Else)

		case *WhileExpr:
			recurse(ns, v.Condition)
			recurse(ns, v.Do)

		case *MatchExpr:
			recurse(ns, v.Left)
			recurse(ns, v.Cases)

		case *Cases:
			for _, c := range v.Cases {
				recurse(ns, c)
			}
			if v.Null != nil {
				recurse(ns, v.Null)
			}

		case *Case:
			addNS(v, &v.Name)
			recurse(ns, &v.Type)
			recurse(ns, v.Body)

		case *AssignExpr:
			recurse(ns, &v.Left)
			recurse(ns, v.Right)

		case *VarExpr:
			recurse(ns, v.Value)
			addNS(v, &v.Name)
			recurse(ns, &v.Type)
			recurse(ns, v.Expr)

		case *ChainExpr:
			recurse(ns, v.Pre)
			recurse(ns, v.Expr)

		case NativeExpr:

		default:
			panic(v)
		}
	}

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	for _, c := range ast.Classes {
		recurse(nil, c)
	}

	panic("unimplemented: typecheck")
}

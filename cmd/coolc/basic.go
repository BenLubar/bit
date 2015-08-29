package main

import "github.com/BenLubar/bit/bitgen"

var basicAny = &ClassDecl{
	Name: TYPE{
		Name: "Any",
	},
	Body: []Feature{
		&MethodFeature{
			Name: ID{
				Name: "Any",
			},
			Return: TYPE{
				Name: "Any",
			},
			Body: &ThisExpr{},
		},
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
				w.Load(start, w.Return, w.This, 0, next)
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
				w.CopyReg(different, w.Return, w.False, next)
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
				Name: "IO",
			},
			Return: TYPE{
				Name: "IO",
			},
			Body: &ThisExpr{},
		},
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
				Cases: []*Case{
					{
						Name: ID{
							Name: "_null",
						},
						Type: TYPE{
							Name: "Null",
						},
						Body: &BooleanExpr{
							B: true,
						},
					},
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
			Body: &CallExpr{
				Left: &ThisExpr{},
				Name: ID{
					Name: "out",
				},
				Args: []Expr{
					&IfExpr{
						Condition: &CallExpr{
							Left: &ThisExpr{},
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
				w.EndStack()

				next := w.ReserveLine()
				w.NewInt(start, w.General[0], 0, next)
				start = next

				next = w.ReserveLine()
				w.CopyReg(start, w.Return, w.Alloc, next)
				start = next

				for i := uint(0); i < basicString.size; i++ {
					next = w.ReserveLine()
					w.Increment(start, w.Alloc.Num, next, 0)
					start = next
				}

				next = w.ReserveLine()
				w.Assign(start, w.Alloc.Ptr, bitgen.Offset{w.Alloc.Ptr, basicString.size * 8}, next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, bitgen.Integer{bitgen.ValueAt{w.Return.Ptr}, 32}, w.Classes[basicString].Num, next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.Return.Ptr, basicStringLength.offset * 8}}, 32}, w.General[0].Num, next)
				start = next

				loop, done := start, w.ReserveLine()
				next = w.ReserveLine()
				w.InputEOF(loop, bitgen.Integer{bitgen.ValueAt{w.Alloc.Ptr}, 8}, next, done)
				start = next

				newline := w.ReserveLine()
				next = w.ReserveLine()
				w.Cmp(start, bitgen.Integer{bitgen.ValueAt{w.Alloc.Ptr}, 8}, '\n', newline, next)
				start = next

				next = w.ReserveLine()
				w.Increment(start, w.IntValue(w.General[0].Ptr), next, 0)
				start = next

				next = w.ReserveLine()
				w.Increment(start, w.Alloc.Num, next, 0)
				start = next

				w.Assign(start, w.Alloc.Ptr, bitgen.Offset{w.Alloc.Ptr, 8}, loop)

				next = w.ReserveLine()
				w.Increment(newline, w.Alloc.Num, next, 0)
				start = next

				w.Assign(start, w.Alloc.Ptr, bitgen.Offset{w.Alloc.Ptr, 8}, done)

				w.PopStack(done, end)
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
				next := w.ReserveLine()
				str, prevstr := w.StackAlloc(start, next)
				start = next

				w.EndStack()

				next = w.ReserveLine()
				w.CopyReg(start, w.General[0], w.Symbol, next)
				start = next

				next = w.ReserveLine()
				w.Cmp(start, w.StackOffset(w.Arg(0)), 0, w.Null, next)
				start = next

				loop, done, found := start, w.ReserveLine(), w.ReserveLine()
				next = w.ReserveLine()
				w.Cmp(start, w.General[0].Num, 0, done, next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, str, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[0].Ptr, basicSymbolName.offset * 8}}, 32}, next)
				start = next

				next = w.ReserveLine()
				w.BeginStack(start, next)
				start = next

				next = w.ReserveLine()
				arg, _ := w.StackAlloc(start, next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, arg, w.PrevStackOffset(w.Arg(0)), next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, w.This.Num, prevstr, next)
				start = next

				next = w.ReserveLine()
				w.Pointer(start, w.This.Ptr, w.This.Num, next)
				start = next

				next = w.ReserveLine()
				w.StaticCall(start, basicStringEquals, next)
				start = next

				w.CmpReg(start, w.True.Num, w.Return.Num, found, loop)

				next = w.ReserveLine()
				w.NewNative(done, w.General[0], basicSymbol, 0, next)
				start = next

				next = w.ReserveLine()
				w.NewInt(start, w.General[1], 0, next)
				start = next

				next, done = w.ReserveLine(), w.ReserveLine()
				w.Cmp(start, w.Symbol.Num, 0, done, next)
				start = next

				next = w.ReserveLine()
				w.Load(start, w.General[2], w.Symbol, basicSymbolHash.offset, next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, w.IntValue(w.General[1].Ptr), w.IntValue(w.General[2].Ptr), next)
				start = next

				w.Increment(start, w.IntValue(w.General[1].Ptr), done, 0)

				next = w.ReserveLine()
				w.Copy(done, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[0].Ptr, basicSymbolName.offset * 8}}, 32}, w.StackOffset(w.Arg(0)), next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[0].Ptr, basicSymbolHash.offset * 8}}, 32}, w.General[1].Num, next)
				start = next

				w.CopyReg(start, w.Symbol, w.General[0], found)

				next = w.ReserveLine()
				w.CopyReg(found, w.Return, w.General[0], next)
				start = next

				w.PopStack(start, end)
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
				w.EndStack()

				next := w.ReserveLine()
				w.Load(start, w.Return, w.Stack, w.Arg(0), next)
				start = next

				next = w.ReserveLine()
				w.Load(start, w.Return, w.Return, basicSymbolName.offset, next)
				start = next

				w.PopStack(start, end)
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

var basicIntSubtract = &MethodFeature{
	Name: ID{
		Name: "_subtract",
	},
	Args: []*VarDecl{
		{
			Name: ID{
				Name: "x",
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
		w.EndStack()

		next := w.ReserveLine()
		w.NewInt(start, w.Return, 0, next)
		start = next

		next = w.ReserveLine()
		w.Load(start, w.General[0], w.Stack, w.Arg(0), next)
		start = next

		for i := uint(0); i < 32; i++ {
			zero, one := w.ReserveLine(), w.ReserveLine()
			w.Jump(start, bitgen.ValueAt{bitgen.Offset{w.General[0].Ptr, 32 + i}}, zero, one)

			next = w.ReserveLine()
			w.Assign(zero, bitgen.ValueAt{bitgen.Offset{w.Return.Ptr, 32 + i}}, bitgen.Bit(true), next)
			w.Assign(one, bitgen.ValueAt{bitgen.Offset{w.Return.Ptr, 32 + i}}, bitgen.Bit(false), next)
			start = next
		}

		next = w.ReserveLine()
		w.Increment(start, w.IntValue(w.Return.Ptr), next, next)
		start = next

		next = w.ReserveLine()
		w.AddReg(start, w.IntValue(w.Return.Ptr), w.IntValue(w.This.Ptr), next)
		start = next

		w.PopStack(start, end)
	}),
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

				w.CmpReg(start, w.IntValue(w.This.Ptr), w.IntValue(w.General[0].Ptr), same, different)

				next = w.ReserveLine()
				w.CopyReg(same, w.Return, w.True, next)
				w.CopyReg(different, w.Return, w.False, next)
				start = next

				w.PopStack(start, end)
			}),
		},
		&MethodFeature{
			Name: ID{
				Name: "_negative",
			},
			Return: TYPE{
				Name: "Int",
			},
			Body: NativeExpr(func(w *writer, start, end bitgen.Line) {
				w.EndStack()

				next := w.ReserveLine()
				w.NewInt(start, w.Return, 0, next)
				start = next

				for i := uint(0); i < 32; i++ {
					zero, one := w.ReserveLine(), w.ReserveLine()
					w.Jump(start, bitgen.ValueAt{bitgen.Offset{w.This.Ptr, 32 + i}}, zero, one)

					next = w.ReserveLine()
					w.Assign(zero, bitgen.ValueAt{bitgen.Offset{w.Return.Ptr, 32 + i}}, bitgen.Bit(true), next)
					w.Assign(one, bitgen.ValueAt{bitgen.Offset{w.Return.Ptr, 32 + i}}, bitgen.Bit(false), next)
					start = next
				}

				next = w.ReserveLine()
				w.Increment(start, w.IntValue(w.Return.Ptr), next, next)
				start = next

				w.PopStack(start, end)
			}),
		},
		&MethodFeature{
			Name: ID{
				Name: "_add",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "x",
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
				w.EndStack()

				next := w.ReserveLine()
				w.NewInt(start, w.Return, 0, next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, w.IntValue(w.Return.Ptr), w.IntValue(w.This.Ptr), next)
				start = next

				next = w.ReserveLine()
				w.Load(start, w.General[0], w.Stack, w.Arg(0), next)
				start = next

				next = w.ReserveLine()
				w.AddReg(start, w.IntValue(w.Return.Ptr), w.IntValue(w.General[0].Ptr), next)
				start = next

				w.PopStack(start, end)
			}),
		},
		basicIntSubtract,
		&MethodFeature{
			Name: ID{
				Name: "_multiply",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "x",
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
				w.EndStack()

				next := w.ReserveLine()
				w.NewInt(start, w.Return, 0, next)
				start = next

				next = w.ReserveLine()
				w.Load(start, w.General[0], w.Stack, w.Arg(0), next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, w.This.Num, w.IntValue(w.This.Ptr), next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, w.General[0].Num, w.IntValue(w.General[0].Ptr), next)
				start = next

				loop, done := start, w.ReserveLine()
				next = w.ReserveLine()
				w.Cmp(loop, w.General[0].Num, 0, next, done)
				start = next

				next = w.ReserveLine()
				add := w.ReserveLine()
				w.Jump(start, w.General[0].Num.Start, add, next)

				w.AddReg(add, w.IntValue(w.Return.Ptr), w.This.Num, next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, bitgen.Integer{w.This.Num.Start, w.This.Num.Width - 1}, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{bitgen.AddressOf{w.This.Num.Start}, 1}}, w.This.Num.Width - 1}, next)
				start = next

				next = w.ReserveLine()
				w.Assign(start, w.This.Num.Start, bitgen.Bit(false), next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{bitgen.AddressOf{w.General[0].Num.Start}, 1}}, w.General[0].Num.Width - 1}, bitgen.Integer{w.General[0].Num.Start, w.General[0].Num.Width - 1}, next)
				start = next

				w.Assign(start, bitgen.ValueAt{bitgen.Offset{bitgen.AddressOf{w.General[0].Num.Start}, w.General[0].Num.Width - 1}}, bitgen.Bit(false), loop)

				w.PopStack(done, end)
			}),
		},
		&MethodFeature{
			Name: ID{
				Name: "_divide",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "x",
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
				Name: "_less",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "x",
					},
					Type: TYPE{
						Name: "Int",
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

				pos, neg := w.ReserveLine(), w.ReserveLine()
				w.Jump(start, bitgen.ValueAt{bitgen.Offset{w.This.Ptr, 32 + 32 - 1}}, pos, neg)

				yes, no := w.ReserveLine(), w.ReserveLine()

				next = w.ReserveLine()
				w.Jump(pos, bitgen.ValueAt{bitgen.Offset{w.General[0].Ptr, 32 + 32 - 1}}, next, no)
				w.Jump(pos, bitgen.ValueAt{bitgen.Offset{w.General[0].Ptr, 32 + 32 - 1}}, yes, next)
				start = next

				w.LessThanUnsigned(start, w.IntValue(w.This.Ptr), w.IntValue(w.General[0].Ptr), yes, no, no)

				next = w.ReserveLine()
				w.CopyReg(yes, w.Return, w.True, next)
				w.CopyReg(no, w.Return, w.False, next)
				start = next

				w.PopStack(start, end)
			}),
		},
		&MethodFeature{
			Name: ID{
				Name: "_less_equal",
			},
			Args: []*VarDecl{
				{
					Name: ID{
						Name: "x",
					},
					Type: TYPE{
						Name: "Int",
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

				pos, neg := w.ReserveLine(), w.ReserveLine()
				w.Jump(start, bitgen.ValueAt{bitgen.Offset{w.This.Ptr, 32 + 32 - 1}}, pos, neg)

				yes, no := w.ReserveLine(), w.ReserveLine()

				next = w.ReserveLine()
				w.Jump(pos, bitgen.ValueAt{bitgen.Offset{w.General[0].Ptr, 32 + 32 - 1}}, next, no)
				w.Jump(pos, bitgen.ValueAt{bitgen.Offset{w.General[0].Ptr, 32 + 32 - 1}}, yes, next)
				start = next

				w.LessThanUnsigned(start, w.IntValue(w.This.Ptr), w.IntValue(w.General[0].Ptr), yes, yes, no)

				next = w.ReserveLine()
				w.CopyReg(yes, w.Return, w.True, next)
				w.CopyReg(no, w.Return, w.False, next)
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
		&MethodFeature{
			Name: ID{
				Name: "_not",
			},
			Return: TYPE{
				Name: "Boolean",
			},
			Body: &IfExpr{
				Condition: &ThisExpr{},
				Then: &BooleanExpr{
					B: false,
				},
				Else: &BooleanExpr{
					B: true,
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

var basicStringEquals = &MethodFeature{
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

		next = w.ReserveLine()
		w.Load(start, w.General[1], w.General[0], basicStringLength.offset, next)
		start = next

		next = w.ReserveLine()
		w.Load(start, w.General[2], w.This, basicStringLength.offset, next)
		start = next

		next = w.ReserveLine()
		w.CmpReg(start, w.IntValue(w.General[1].Ptr), w.IntValue(w.General[2].Ptr), next, different)
		start = next

		next = w.ReserveLine()
		w.Copy(start, w.General[1].Num, w.IntValue(w.General[1].Ptr), next)
		start = next

		loop := start
		next = w.ReserveLine()
		w.Decrement(start, w.General[1].Num, next, same)
		start = next

		next = w.ReserveLine()
		w.CmpReg(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.This.Ptr, basicStringLength.offset*8 + 32}}, 8}, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[0].Ptr, basicStringLength.offset*8 + 32}}, 8}, next, different)
		start = next

		next = w.ReserveLine()
		w.Assign(start, w.This.Ptr, bitgen.Offset{w.This.Ptr, 8}, next)
		start = next

		w.Assign(start, w.General[0].Ptr, bitgen.Offset{w.General[0].Ptr, 8}, loop)

		next = w.ReserveLine()
		w.CopyReg(same, w.Return, w.True, next)
		w.CopyReg(different, w.Return, w.False, next)
		start = next

		w.PopStack(start, end)
	}),
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
		basicStringEquals,
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
				w.EndStack()

				next := w.ReserveLine()
				w.Load(start, w.General[0], w.Stack, w.Arg(0), next)
				start = next

				next = w.ReserveLine()
				w.Load(start, w.General[1], w.General[0], basicStringLength.offset, next)
				start = next

				next = w.ReserveLine()
				w.Load(start, w.General[2], w.This, basicStringLength.offset, next)
				start = next

				next = w.ReserveLine()
				w.NewInt(start, w.General[3], 0, next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, w.IntValue(w.General[3].Ptr), w.IntValue(w.General[1].Ptr), next)
				start = next

				next = w.ReserveLine()
				w.AddReg(start, w.IntValue(w.General[3].Ptr), w.IntValue(w.General[2].Ptr), next)
				start = next

				next = w.ReserveLine()
				w.NewNativeDynamic(start, w.Return, w.basicString, w.IntValue(w.General[3].Ptr), next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.Return.Ptr, basicStringLength.offset * 8}}, 32}, w.General[3].Num, next)
				start = next

				next = w.ReserveLine()
				w.CopyReg(start, w.General[3], w.Return, next)
				start = next

				appendString := func(str, length register) {
					next = w.ReserveLine()
					w.Copy(start, length.Num, w.IntValue(length.Ptr), next)
					start = next

					loop, done := start, w.ReserveLine()
					next = w.ReserveLine()
					w.Decrement(start, length.Num, next, done)
					start = next

					next = w.ReserveLine()
					w.Copy(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[3].Ptr, basicStringLength.offset*8 + 32}}, 8}, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{str.Ptr, basicStringLength.offset*8 + 32}}, 8}, next)
					start = next

					next = w.ReserveLine()
					w.Assign(start, str.Ptr, bitgen.Offset{str.Ptr, 8}, next)
					start = next

					w.Assign(start, w.General[3].Ptr, bitgen.Offset{w.General[3].Ptr, 8}, loop)

					start = done
				}

				appendString(w.This, w.General[2])
				appendString(w.General[0], w.General[1])

				w.PopStack(start, end)
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
				w.EndStack()

				next := w.ReserveLine()
				w.Load(start, w.General[0], w.Stack, w.Arg(0), next)
				start = next

				next = w.ReserveLine()
				w.Load(start, w.General[1], w.Stack, w.Arg(1), next)
				start = next

				next = w.ReserveLine()
				w.Load(start, w.General[2], w.This, basicStringLength.offset, next)
				start = next

				next = w.ReserveLine()
				w.LessThanUnsigned(start, w.IntValue(w.General[0].Ptr), w.IntValue(w.General[1].Ptr), next, next, w.IndexRange)
				start = next

				next = w.ReserveLine()
				w.LessThanUnsigned(start, w.IntValue(w.General[1].Ptr), w.IntValue(w.General[2].Ptr), next, w.IndexRange, w.IndexRange)
				start = next

				next = w.ReserveLine()
				w.BeginStack(start, next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, w.This.Num, w.PrevStackOffset(w.Arg(1)), next)
				start = next

				next = w.ReserveLine()
				w.Pointer(start, w.This.Ptr, w.This.Num, next)
				start = next

				arg0, _ := w.StackAlloc(start, next)
				next = w.ReserveLine()
				w.Copy(start, arg0, w.PrevStackOffset(w.Arg(0)), next)
				start = next

				next = w.ReserveLine()
				w.StaticCall(start, basicIntSubtract, next)
				start = next

				next = w.ReserveLine()
				w.CopyReg(start, w.General[2], w.Return, next)
				start = next

				next = w.ReserveLine()
				w.NewNativeDynamic(start, w.Return, w.basicString, w.IntValue(w.General[2].Ptr), next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.Return.Ptr, basicStringLength.offset * 8}}, 32}, w.General[2].Num, next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, w.General[0].Num, w.IntValue(w.General[0].Ptr), next)
				start = next

				loop, done := start, w.ReserveLine()
				next = w.ReserveLine()
				w.Decrement(loop, w.General[0].Num, next, done)
				start = next

				w.Assign(start, w.This.Ptr, bitgen.Offset{w.This.Ptr, 8}, loop)

				next = w.ReserveLine()
				w.Assign(done, w.General[0].Ptr, w.Return.Ptr, next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, w.General[2].Num, w.IntValue(w.General[2].Ptr), next)
				start = next

				loop, done = start, w.ReserveLine()
				next = w.ReserveLine()
				w.Decrement(loop, w.General[2].Num, next, done)
				start = next

				next = w.ReserveLine()
				w.Copy(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[0].Ptr, basicStringLength.offset*8 + 32}}, 8}, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.This.Ptr, basicStringLength.offset*8 + 32}}, 8}, next)
				start = next

				next = w.ReserveLine()
				w.Assign(start, w.This.Ptr, bitgen.Offset{w.This.Ptr, 8}, next)
				start = next

				w.Assign(start, w.General[0].Ptr, bitgen.Offset{w.General[0].Ptr, 8}, loop)

				w.PopStack(done, end)
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
				w.EndStack()

				next := w.ReserveLine()
				w.Load(start, w.General[0], w.Stack, w.Arg(0), next)
				start = next

				next = w.ReserveLine()
				w.Load(start, w.General[1], w.This, basicStringLength.offset, next)
				start = next

				next = w.ReserveLine()
				w.LessThanUnsigned(start, w.IntValue(w.General[0].Ptr), w.IntValue(w.General[1].Ptr), next, w.IndexRange, w.IndexRange)
				start = next

				next = w.ReserveLine()
				w.Copy(start, w.General[0].Num, w.IntValue(w.General[0].Ptr), next)
				start = next

				next = w.ReserveLine()
				w.NewInt(start, w.Return, 0, next)
				start = next

				loop, done := start, w.ReserveLine()
				next = w.ReserveLine()
				w.Decrement(start, w.General[0].Num, next, done)
				start = next

				w.Assign(start, w.This.Ptr, bitgen.Offset{w.This.Ptr, 8}, loop)

				next = w.ReserveLine()
				w.Copy(done, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.Return.Ptr, 32 + 32 - 8}}, 32}, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.This.Ptr, 32 + 32}}, 32}, next)
				start = next

				w.PopStack(start, end)
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
						Value: &CallExpr{
							Left: &NameExpr{
								Name: ID{
									Name: "length",
								},
							},
							Name: ID{
								Name: "_subtract",
							},
							Args: []Expr{
								&NameExpr{
									Name: ID{
										Name: "n",
									},
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
								Value: &CallExpr{
									Left: &IntegerExpr{
										N: 1,
									},
									Name: ID{
										Name: "_negative",
									},
								},
							},
							Expr: &ChainExpr{
								Pre: &WhileExpr{
									Condition: &CallExpr{
										Left: &NameExpr{
											Name: ID{
												Name: "i",
											},
										},
										Name: ID{
											Name: "_less_equal",
										},
										Args: []Expr{
											&NameExpr{
												Name: ID{
													Name: "diff",
												},
											},
										},
									},
									Do: &IfExpr{
										Condition: &CallExpr{
											Left: &CallExpr{
												Left: &ThisExpr{},
												Name: ID{
													Name: "substring",
												},
												Args: []Expr{
													&NameExpr{
														Name: ID{
															Name: "i",
														},
													},
													&CallExpr{
														Left: &NameExpr{
															Name: ID{
																Name: "i",
															},
														},
														Name: ID{
															Name: "_add",
														},
														Args: []Expr{
															&NameExpr{
																Name: ID{
																	Name: "n",
																},
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
												Right: &CallExpr{
													Left: &NameExpr{
														Name: ID{
															Name: "diff",
														},
													},
													Name: ID{
														Name: "_add",
													},
													Args: []Expr{
														&IntegerExpr{
															N: 1,
														},
													},
												},
											},
										},
										Else: &AssignExpr{
											Left: ID{
												Name: "i",
											},
											Right: &CallExpr{
												Left: &NameExpr{
													Name: ID{
														Name: "i",
													},
												},
												Name: ID{
													Name: "_add",
												},
												Args: []Expr{
													&IntegerExpr{
														N: 1,
													},
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

var basicSymbolName = &VarFeature{
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
}

var basicSymbolHash = &VarFeature{
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
}

var basicSymbolNext = &VarFeature{
	VarDecl: VarDecl{
		Name: ID{
			Name: "next",
		},
		Type: TYPE{
			Name: "Symbol",
		},
	},
	Value: &NullExpr{},
}

var basicSymbol = &ClassDecl{
	Name: TYPE{
		Name: "Symbol",
	},
	Body: []Feature{
		basicSymbolName,
		basicSymbolHash,
		basicSymbolNext,
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
			Body: &VarExpr{
				VarFeature: VarFeature{
					VarDecl: VarDecl{
						Name: ID{
							Name: "ret",
						},
						Type: TYPE{
							Name: "ArrayAny",
						},
					},
					Value: &NewExpr{
						Type: TYPE{
							Name: "ArrayAny",
						},
						Args: []Expr{
							&NameExpr{
								Name: ID{
									Name: "s",
								},
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
					Expr: &ChainExpr{
						Pre: &WhileExpr{
							Condition: &IfExpr{
								Condition: &CallExpr{
									Left: &NameExpr{
										Name: ID{
											Name: "i",
										},
									},
									Name: ID{
										Name: "_less",
									},
									Args: []Expr{
										&NameExpr{
											Name: ID{
												Name: "length",
											},
										},
									},
								},
								Then: &CallExpr{
									Left: &NameExpr{
										Name: ID{
											Name: "i",
										},
									},
									Name: ID{
										Name: "_less",
									},
									Args: []Expr{
										&NameExpr{
											Name: ID{
												Name: "length",
											},
										},
									},
								},
								Else: &BooleanExpr{
									B: false,
								},
							},
							Do: &ChainExpr{
								Pre: &CallExpr{
									Left: &NameExpr{
										Name: ID{
											Name: "ret",
										},
									},
									Name: ID{
										Name: "set",
									},
									Args: []Expr{
										&NameExpr{
											Name: ID{
												Name: "i",
											},
										},
										&CallExpr{
											Left: &ThisExpr{},
											Name: ID{
												Name: "get",
											},
											Args: []Expr{
												&NameExpr{
													Name: ID{
														Name: "i",
													},
												},
											},
										},
									},
								},
								Expr: &AssignExpr{
									Left: ID{
										Name: "i",
									},
									Right: &CallExpr{
										Left: &NameExpr{
											Name: ID{
												Name: "i",
											},
										},
										Name: ID{
											Name: "_add",
										},
										Args: []Expr{
											&IntegerExpr{
												N: 1,
											},
										},
									},
								},
							},
						},
						Expr: &NameExpr{
							Name: ID{
								Name: "ret",
							},
						},
					},
				},
			},
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
				w.Copy(start, w.General[0].Num, w.IntValue(w.General[0].Ptr), next)
				start = next

				next = w.ReserveLine()
				w.Load(start, w.General[1], w.This, basicArrayAnyLength.offset, next)
				start = next

				next = w.ReserveLine()
				w.LessThanUnsigned(start, w.General[0].Num, w.IntValue(w.General[1].Ptr), next, w.IndexRange, w.IndexRange)
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
				w.Copy(start, w.General[0].Num, w.IntValue(w.General[0].Ptr), next)
				start = next

				next = w.ReserveLine()
				w.Load(start, w.General[1], w.This, basicArrayAnyLength.offset, next)
				start = next

				next = w.ReserveLine()
				w.Copy(start, w.General[1].Num, w.IntValue(w.General[1].Ptr), next)
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
	if err := basicAST.Semantic(true); err != nil {
		panic(err)
	}
}

// vim: set ts=2 sw=2:

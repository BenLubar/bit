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
				w.Copy(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[1].Ptr, 32}}, 32}, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[2].Ptr, 32}}, 32}, next)
				start = next

				w.Increment(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[1].Ptr, 32}}, 32}, done, 0)

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
				panic("unimplemented")
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
				panic("unimplemented")
			}),
		},
		&MethodFeature{
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
				panic("unimplemented")
			}),
		},
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
				panic("unimplemented")
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
				Name: "Int",
			},
			Body: NativeExpr(func(w *writer, start, end bitgen.Line) {
				panic("unimplemented")
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
				Name: "Int",
			},
			Body: NativeExpr(func(w *writer, start, end bitgen.Line) {
				panic("unimplemented")
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
		w.CmpReg(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[1].Ptr, 32}}, 32}, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[2].Ptr, 32}}, 32}, next, different)
		start = next

		next = w.ReserveLine()
		w.Copy(start, w.General[1].Num, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[1].Ptr, 32}}, 32}, next)
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
				w.EndStack()

				next := w.ReserveLine()
				w.Load(start, w.General[0], w.Stack, w.Arg(0), next)
				start = next

				next = w.ReserveLine()
				w.Load(start, w.General[1], w.This, basicStringLength.offset, next)
				start = next

				next = w.ReserveLine()
				w.LessThanUnsigned(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[0].Ptr, 32}}, 32}, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[1].Ptr, 32}}, 32}, next, w.IndexRange, w.IndexRange)
				start = next

				next = w.ReserveLine()
				w.Copy(start, w.General[0].Num, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[0].Ptr, 32}}, 32}, next)
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
				w.LessThanUnsigned(start, w.General[0].Num, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[1].Ptr, 32}}, 32}, next, w.IndexRange, w.IndexRange)
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

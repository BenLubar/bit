package main

import (
	"go/token"

	"github.com/BenLubar/bit/bitgen"
)

type AST struct {
	FileSet *token.FileSet

	Classes []*ClassDecl
	main    *ClassDecl
}

type ID struct {
	Name   string
	Pos    token.Pos
	target interface{}
}

type TYPE struct {
	Name   string
	Pos    token.Pos
	target *ClassDecl
}

type ClassDecl struct {
	Name    TYPE
	Args    []*VarDecl
	Extends *ExtendsDecl
	Body    []Feature

	methods map[string]*MethodFeature
	size    uint
	depth   uint
}

type VarDecl struct {
	Name ID
	Type TYPE

	offset uint
	arg    uint
}

type ExtendsDecl struct {
	Type TYPE
	Args []Expr
}

type Feature interface{}

type MethodFeature struct {
	Override bool
	Name     ID
	Args     []*VarDecl
	Return   TYPE
	Body     Expr

	offset uint
}

type VarFeature struct {
	VarDecl
	Value Expr
}

type BlockFeature struct {
	Expr
}

type NativeFeature struct{}

type Expr interface {
	write(w *writer, start, end bitgen.Line)
	alloc(w *writer, start bitgen.Line) (next bitgen.Line)
}

type ConstructorExpr struct {
	Args []*VarDecl
	Expr Expr
}

func (e *ConstructorExpr) write(w *writer, start, end bitgen.Line) {
	for _, a := range e.Args {
		next := w.ReserveLine()
		w.Copy(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.This.Ptr, a.offset * 8}}, 32}, w.StackOffset(w.Arg(a.arg)), next)
		start = next
	}

	next := w.ReserveLine()
	e.Expr.write(w, start, next)
	start = next

	w.CopyReg(start, w.Return, w.This, end)
}

func (e *ConstructorExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	return e.Expr.alloc(w, start)
}

type AssignExpr struct {
	Left  ID
	Right Expr
}

func (e *AssignExpr) write(w *writer, start, end bitgen.Line) {
	next := w.ReserveLine()
	e.Right.write(w, start, next)
	start = next

	switch v := e.Left.target.(type) {
	case *VarDecl:
		if v.offset != 0 {
			w.Copy(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.This.Ptr, v.offset * 8}}, 32}, w.Return.Num, end)
		} else {
			w.Copy(start, w.StackOffset(w.Arg(v.arg)), w.Return.Num, end)
		}

	case *VarFeature:
		w.Copy(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.This.Ptr, v.offset * 8}}, 32}, w.Return.Num, end)

	case *VarExpr:
		w.Copy(start, v.slot, w.Return.Num, end)

	case *Case:
		w.Copy(start, v.match.slot, w.Return.Num, end)

	default:
		panic(v)
	}
}

func (e *AssignExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	return e.Right.alloc(w, start)
}

type IfExpr struct {
	Pos       token.Pos
	Condition Expr
	Then      Expr
	Else      Expr
}

func (e *IfExpr) write(w *writer, start, end bitgen.Line) {
	next := w.ReserveLine()
	e.Condition.write(w, start, next)
	start = next

	n0, n1 := w.ReserveLine(), w.ReserveLine()
	w.CmpReg(start, w.Return.Num, w.True.Num, n1, n0)

	e.Then.write(w, n1, end)
	e.Else.write(w, n0, end)
}

func (e *IfExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	next = e.Condition.alloc(w, start)
	next = e.Then.alloc(w, next)
	return e.Else.alloc(w, next)
}

type WhileExpr struct {
	Pos       token.Pos
	Condition Expr
	Do        Expr
}

func (e *WhileExpr) write(w *writer, start, end bitgen.Line) {
	loop, done := start, w.ReserveLine()
	next := w.ReserveLine()
	e.Condition.write(w, loop, next)
	start = next

	next = w.ReserveLine()
	w.CmpReg(start, w.Return.Num, w.True.Num, next, done)
	start = next

	e.Do.write(w, start, loop)

	w.CopyReg(done, w.Return, w.Unit, end)
}

func (e *WhileExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	next = e.Condition.alloc(w, start)
	return e.Do.alloc(w, next)
}

type MatchExpr struct {
	Left  Expr
	Cases []*Case

	slot bitgen.Integer
}

func (e *MatchExpr) write(w *writer, start, end bitgen.Line) {
	next := w.ReserveLine()
	e.Left.write(w, start, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, e.slot, w.Return.Num, next)
	start = next

	null := w.CaseNull
findNull:
	for _, c := range e.Cases {
		for _, h := range c.classes {
			if h == basicDummyNull {
				null = w.ReserveLine()
				break findNull
			}
		}
	}

	next = w.ReserveLine()
	w.Cmp(start, w.Return.Num, 0, null, next)
	start = next

	for _, c := range e.Cases {
		hasNull := false
		for _, h := range c.classes {
			if h == basicDummyNull {
				hasNull = true
				break
			}
		}
		target := null
		if !hasNull {
			target = w.ReserveLine()
		}

		for _, h := range c.classes {
			if h == basicDummyNull || h == basicDummyNothing {
				continue
			}

			next = w.ReserveLine()
			w.CmpReg(start, bitgen.Integer{bitgen.ValueAt{w.Return.Ptr}, 32}, w.Classes[h].Num, target, next)
			start = next
		}

		c.Body.write(w, target, end)
	}

	w.Jump(start, bitgen.Bit(false), w.NoCase, w.NoCase)
}

func (e *MatchExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	next = w.ReserveLine()
	e.slot, _ = w.StackAlloc(start, next)

	next = e.Left.alloc(w, next)
	for _, c := range e.Cases {
		next = c.Body.alloc(w, next)
	}
	return next
}

type CallExpr struct {
	Left Expr
	Name ID
	Args []Expr

	thisPre  bitgen.Integer
	thisPost bitgen.Integer
	argPre   []bitgen.Integer
	argPost  []bitgen.Integer
}

func (e *CallExpr) write(w *writer, start, end bitgen.Line) {
	next := w.ReserveLine()
	e.Left.write(w, start, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, e.thisPre, w.Return.Num, next)
	start = next

	for i, a := range e.Args {
		next = w.ReserveLine()
		a.write(w, start, next)
		start = next

		next = w.ReserveLine()
		w.Copy(start, e.argPre[i], w.Return.Num, next)
		start = next
	}

	next = w.ReserveLine()
	w.BeginStack(start, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, w.This.Num, e.thisPost, next)
	start = next

	next = w.ReserveLine()
	w.Pointer(start, w.This.Ptr, w.This.Num, next)
	start = next

	for _, a := range e.argPost {
		next = w.ReserveLine()
		arg, _ := w.StackAlloc(start, next)
		start = next

		next = w.ReserveLine()
		w.Copy(start, arg, a, next)
		start = next
	}

	w.DynamicCall(start, e, end)
}

func (e *CallExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	start = e.Left.alloc(w, start)

	for _, a := range e.Args {
		start = a.alloc(w, start)
	}

	next = w.ReserveLine()
	e.thisPre, e.thisPost = w.StackAlloc(start, next)
	start = next

	e.argPre = make([]bitgen.Integer, len(e.Args))
	e.argPost = make([]bitgen.Integer, len(e.Args))

	for i := range e.Args {
		next = w.ReserveLine()
		e.argPre[i], e.argPost[i] = w.StackAlloc(start, next)
		start = next
	}

	return start
}

type StaticCallExpr struct {
	Name ID
	Args []Expr

	thisPre  bitgen.Integer
	thisPost bitgen.Integer
	argPre   []bitgen.Integer
	argPost  []bitgen.Integer
}

func (e *StaticCallExpr) write(w *writer, start, end bitgen.Line) {
	next := w.ReserveLine()
	w.Copy(start, e.thisPre, w.This.Num, next)
	start = next

	for i, a := range e.Args {
		next = w.ReserveLine()
		a.write(w, start, next)
		start = next

		next = w.ReserveLine()
		w.Copy(start, e.argPre[i], w.Return.Num, next)
		start = next
	}

	next = w.ReserveLine()
	w.BeginStack(start, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, w.This.Num, e.thisPost, next)
	start = next

	next = w.ReserveLine()
	w.Pointer(start, w.This.Ptr, w.This.Num, next)
	start = next

	for _, a := range e.argPost {
		next = w.ReserveLine()
		arg, _ := w.StackAlloc(start, next)
		start = next

		next = w.ReserveLine()
		w.Copy(start, arg, a, next)
		start = next
	}

	w.StaticCall(start, e, end)
}

func (e *StaticCallExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	for _, a := range e.Args {
		start = a.alloc(w, start)
	}

	next = w.ReserveLine()
	e.thisPre, e.thisPost = w.StackAlloc(start, next)
	start = next

	e.argPre = make([]bitgen.Integer, len(e.Args))
	e.argPost = make([]bitgen.Integer, len(e.Args))

	for i := range e.Args {
		next = w.ReserveLine()
		e.argPre[i], e.argPost[i] = w.StackAlloc(start, next)
		start = next
	}

	return start
}

type NewExpr struct {
	Type TYPE
}

func (e *NewExpr) write(w *writer, start, end bitgen.Line) {
	w.New(start, w.Return, e.Type.target, end)
}

func (e *NewExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	return start
}

type VarExpr struct {
	VarFeature
	Expr Expr

	slot bitgen.Integer
}

func (e *VarExpr) write(w *writer, start, end bitgen.Line) {
	next := w.ReserveLine()
	e.Value.write(w, start, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, e.slot, w.Return.Num, next)
	start = next

	e.Expr.write(w, start, end)
}

func (e *VarExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	next = w.ReserveLine()
	e.slot, _ = w.StackAlloc(start, next)

	next = e.Value.alloc(w, next)
	return e.Expr.alloc(w, next)
}

type ChainExpr struct {
	Pre  Expr
	Expr Expr
}

func (e *ChainExpr) write(w *writer, start, end bitgen.Line) {
	next := w.ReserveLine()
	e.Pre.write(w, start, next)
	start = next

	e.Expr.write(w, start, end)
}

func (e *ChainExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	next = e.Pre.alloc(w, start)
	return e.Expr.alloc(w, next)
}

type NullExpr struct {
}

func (e *NullExpr) write(w *writer, start, end bitgen.Line) {
	for i := uint(0); i < w.Return.Num.Width; i++ {
		next := w.ReserveLine()
		w.Assign(start, w.Return.Num.Bit(i), bitgen.Bit(false), next)
		start = next
	}
	w.Assign(start, w.Return.Ptr, w.Heap, end)
}

func (e *NullExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	return start
}

type UnitExpr struct {
}

func (e *UnitExpr) write(w *writer, start, end bitgen.Line) {
	w.CopyReg(start, w.Return, w.Unit, end)
}

func (e *UnitExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	return start
}

type NameExpr struct {
	Name ID
}

func (e *NameExpr) write(w *writer, start, end bitgen.Line) {
	switch v := e.Name.target.(type) {
	case *VarDecl:
		if v.offset != 0 {
			w.Load(start, w.Return, w.This, v.offset, end)
		} else {
			w.Load(start, w.Return, w.Stack, w.Arg(v.arg), end)
		}

	case *VarFeature:
		w.Load(start, w.Return, w.This, v.offset, end)

	case *VarExpr:
		next := w.ReserveLine()
		w.Copy(start, w.Return.Num, v.slot, next)
		start = next

		w.Pointer(start, w.Return.Ptr, w.Return.Num, end)

	case *Case:
		next := w.ReserveLine()
		w.Copy(start, w.Return.Num, v.match.slot, next)
		start = next

		w.Pointer(start, w.Return.Ptr, w.Return.Num, end)

	default:
		panic(v)
	}
}

func (e *NameExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	return start
}

type IntegerExpr struct {
	N int32
}

func (e *IntegerExpr) write(w *writer, start, end bitgen.Line) {
	w.NewInt(start, w.Return, e.N, end)
}

func (e *IntegerExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	return start
}

type StringExpr struct {
	S string

	length bitgen.Integer
}

func (e *StringExpr) write(w *writer, start, end bitgen.Line) {
	next := w.ReserveLine()
	w.Copy(start, e.length, w.General[0].Num, next)
	start = next

	next = w.ReserveLine()
	w.NewString(start, w.Return, w.General[0], e.S, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, w.General[0].Num, e.length, next)
	start = next

	w.Pointer(start, w.General[0].Ptr, w.General[0].Num, end)
}

func (e *StringExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	next = w.ReserveLine()
	e.length, _ = w.StackAlloc(start, next)
	start = next

	return start
}

type BooleanExpr struct {
	B bool
}

func (e *BooleanExpr) write(w *writer, start, end bitgen.Line) {
	if e.B {
		w.CopyReg(start, w.Return, w.True, end)
	} else {
		w.CopyReg(start, w.Return, w.False, end)
	}
}

func (e *BooleanExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	return start
}

type ThisExpr struct {
}

func (e *ThisExpr) write(w *writer, start, end bitgen.Line) {
	w.CopyReg(start, w.Return, w.This, end)
}

func (e *ThisExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	return start
}

type Case struct {
	Name ID
	Type TYPE
	Body Expr

	match   *MatchExpr
	classes []*ClassDecl
}

type NativeExpr func(w *writer, start, end bitgen.Line)

func (e NativeExpr) write(w *writer, start, end bitgen.Line) {
	e(w, start, end)
}

func (e NativeExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	panic("cannot include NativeExpr inside another expression")
}

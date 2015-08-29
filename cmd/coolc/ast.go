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
	panic("unimplemented")
}

func (e *AssignExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	panic("unimplemented")
}

type IfExpr struct {
	Pos       token.Pos
	Condition Expr
	Then      Expr
	Else      Expr
}

func (e *IfExpr) write(w *writer, start, end bitgen.Line) {
	panic("unimplemented")
}

func (e *IfExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	panic("unimplemented")
}

type WhileExpr struct {
	Pos       token.Pos
	Condition Expr
	Do        Expr
}

func (e *WhileExpr) write(w *writer, start, end bitgen.Line) {
	panic("unimplemented")
}

func (e *WhileExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	panic("unimplemented")
}

type MatchExpr struct {
	Left  Expr
	Cases []*Case
}

func (e *MatchExpr) write(w *writer, start, end bitgen.Line) {
	panic("unimplemented")
}

func (e *MatchExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	panic("unimplemented")
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

	w.DynamicCall(start, e.Name.target.(*MethodFeature), end)
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

	w.StaticCall(start, e.Name.target.(*MethodFeature), end)
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
	panic("unimplemented")
}

func (e *NewExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	panic("unimplemented")
}

type VarExpr struct {
	VarFeature
	Expr Expr
}

func (e *VarExpr) write(w *writer, start, end bitgen.Line) {
	panic("unimplemented")
}

func (e *VarExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	panic("unimplemented")
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
	panic("unimplemented")
}

func (e *NullExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	panic("unimplemented")
}

type UnitExpr struct {
}

func (e *UnitExpr) write(w *writer, start, end bitgen.Line) {
	panic("unimplemented")
}

func (e *UnitExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	panic("unimplemented")
}

type NameExpr struct {
	Name ID
}

func (e *NameExpr) write(w *writer, start, end bitgen.Line) {
	panic("unimplemented")
}

func (e *NameExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	panic("unimplemented")
}

type IntegerExpr struct {
	N int32
}

func (e *IntegerExpr) write(w *writer, start, end bitgen.Line) {
	panic("unimplemented")
}

func (e *IntegerExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	panic("unimplemented")
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
	panic("unimplemented")
}

func (e *BooleanExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	panic("unimplemented")
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
}

type NativeExpr func(w *writer, start, end bitgen.Line)

func (e NativeExpr) write(w *writer, start, end bitgen.Line) {
	e(w, start, end)
}

func (e NativeExpr) alloc(w *writer, start bitgen.Line) (next bitgen.Line) {
	panic("cannot include NativeExpr inside another expression")
}

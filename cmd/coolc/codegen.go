package main

import (
	"io"

	"github.com/BenLubar/bit/bitgen"
)

func (ast *AST) WriteTo(out io.Writer) (err error) {
	w := &writer{Writer: bitgen.NewWriter(out)}
	defer func() {
		if e := w.Close(); err == nil {
			err = e
		}
	}()

	w.Init()

	panic("unimplemented")
}

type writer struct {
	*bitgen.Writer

	Ptr     bitgen.Integer
	Alloc   register
	Return  register
	Stack   register
	Offset  uint
	General [4]register
	Heap    bitgen.Variable
}

type register struct {
	Ptr bitgen.Variable
	Num bitgen.Integer
}

func (w *writer) Init() {
	reg := func(r *register) {
		r.Ptr = w.ReserveVariable()
		r.Num = w.ReserveInteger(32)
	}

	w.Ptr = w.ReserveInteger(32)
	reg(&w.Alloc)
	reg(&w.Return)
	reg(&w.Stack)
	for i := range w.General {
		reg(&w.General[i])
	}
	w.Heap = w.ReserveHeap()
}

func (w *writer) Pointer(start bitgen.Line, ptr bitgen.Variable, num bitgen.Integer, end bitgen.Line) {
	next := w.ReserveLine()
	w.Copy(start, w.Ptr, num, next)
	start = next

	next = w.ReserveLine()
	w.Assign(start, ptr, w.Heap, next)
	start = next

	loop := start
	next = w.ReserveLine()
	w.Decrement(start, w.Ptr, next, end)
	start = next

	w.Assign(start, ptr, bitgen.Offset{ptr, 8}, loop)
}

func (w *writer) Abort(start bitgen.Line, message string) {
	w.PrintString(start, message, 0)
}

func (w *writer) BeginStack(start, end bitgen.Line) {
	if w.Offset != 0 {
		panic("BeginStack while in stack")
	}

	next := w.ReserveLine()
	w.Copy(start, bitgen.Integer{bitgen.ValueAt{w.Alloc.Ptr}, 32}, w.Stack.Num, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, w.Stack.Num, w.Alloc.Num, next)
	start = next

	next = w.ReserveLine()
	w.Assign(start, w.Stack.Ptr, w.Alloc.Ptr, next)
	start = next

	w.Offset = uint(1 + len(w.General))
	for i := uint(0); i < 32/8*w.Offset; i++ {
		next = w.ReserveLine()
		w.Increment(start, w.Alloc.Num, next, 0)
		start = next
	}

	for i := range w.General {
		next = w.ReserveLine()
		w.Copy(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.Stack.Ptr, uint(32 + i*32)}}, 32}, w.General[i].Num, next)
		start = next

		for j := 0; j < 32; j++ {
			next = w.ReserveLine()
			w.Assign(start, bitgen.ValueAt{bitgen.Offset{bitgen.AddressOf{w.General[i].Num.Start}, uint(32 + i*32 + j)}}, bitgen.Bit(false), next)
			start = next
		}

		next = w.ReserveLine()
		w.Assign(start, w.General[i].Ptr, w.Heap, next)
		start = next
	}

	w.Assign(start, w.Alloc.Ptr, bitgen.Offset{w.Alloc.Ptr, w.Offset * 32}, end)
}

func (w *writer) StackAlloc(start, end bitgen.Line) bitgen.Integer {
	if w.Offset == 0 {
		panic("StackAlloc while not in stack")
	}

	ret := bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.Stack.Ptr, w.Offset * 32}}, 32}

	w.Offset++

	for i := 0; i < 32/8; i++ {
		next := w.ReserveLine()
		w.Increment(start, w.Alloc.Num, next, 0)
		start = next
	}

	w.Assign(start, w.Alloc.Ptr, bitgen.Offset{w.Alloc.Ptr, 32}, end)

	return ret
}

func (w *writer) EndStack(start, end bitgen.Line) {
	if w.Offset == 0 {
		panic("EndStack while not in stack")
	}

	w.Offset = 0
}

func (w *writer) PopStack(start, end bitgen.Line) {
	if w.Offset != 0 {
		panic("PopStack while in stack")
	}

	for i := range w.General {
		next := w.ReserveLine()
		w.Copy(start, w.General[i].Num, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.Stack.Ptr, uint(32 + i*32)}}, 32}, next)
		start = next

		next = w.ReserveLine()
		w.Pointer(start, w.General[i].Ptr, w.General[i].Num, next)
		start = next
	}

	next := w.ReserveLine()
	w.Copy(start, w.Stack.Num, bitgen.Integer{bitgen.ValueAt{w.Stack.Ptr}, 32}, next)
	start = next

	w.Pointer(start, w.Stack.Ptr, w.Stack.Num, end)
}

func (w *writer) StaticAlloc(start bitgen.Line, reg register, size uint, end bitgen.Line) {
	if w.Offset != 0 {
		panic("StaticAlloc while in stack")
	}

	next := w.ReserveLine()
	w.Copy(start, reg.Num, w.Alloc.Num, next)
	start = next

	next = w.ReserveLine()
	w.Assign(start, reg.Ptr, w.Alloc.Ptr, next)
	start = next

	for i := uint(0); i < size; i++ {
		next = w.ReserveLine()
		w.Increment(start, w.Alloc.Num, next, 0)
		start = next
	}

	w.Assign(start, w.Alloc.Ptr, bitgen.Offset{w.Alloc.Ptr, 8 * size}, end)
}

func (w *writer) DynamicAlloc(start bitgen.Line, reg register, size bitgen.Integer, end bitgen.Line) {
	if w.Offset != 0 {
		panic("DynamicAlloc while in stack")
	}

	next := w.ReserveLine()
	w.Copy(start, reg.Num, w.Alloc.Num, next)
	start = next

	next = w.ReserveLine()
	w.Assign(start, reg.Ptr, w.Alloc.Ptr, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, w.Ptr, size, next)
	start = next

	loop := start
	next = w.ReserveLine()
	w.Decrement(start, w.Ptr, next, end)
	start = next

	next = w.ReserveLine()
	w.Increment(start, w.Alloc.Num, next, 0)
	start = next

	w.Assign(start, w.Alloc.Ptr, bitgen.Offset{w.Alloc.Ptr, 8}, loop)
}

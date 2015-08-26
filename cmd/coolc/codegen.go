package main

import (
	"fmt"
	"io"

	"github.com/BenLubar/bit/bitgen"
)

func (ast *AST) WriteTo(out io.Writer) (err error) {
	w := &writer{Writer: bitgen.NewWriter(out), AST: ast}
	defer func() {
		if e := w.Close(); err == nil {
			err = e
		}
	}()

	start := w.Init()

	methods := w.MethodTable()

	_, _ = start, methods

	panic("unimplemented")
}

type writer struct {
	*bitgen.Writer

	AST *AST

	Ptr     bitgen.Integer   // used for counting pointers, internal
	Alloc   register         // points at next available memory, internal
	Unit    register         // points at the constant "()", internal
	True    register         // points at the constant "true", internal
	False   register         // points at the constant "false", internal
	Zero    register         // points at the constant "0", internal
	Symbol  register         // points at the last symbol created, internal
	Return  register         // points at return value, not saved
	This    register         // points at "this" value, saved
	Stack   register         // points at start of stack segment, saved
	Offset  uint             // current stack offset, internal
	Prev    bitgen.Variable  // points at previous stack segment, internal
	General [4]register      // general purpose registers, saved
	Heap    bitgen.AddressOf // heap start (also null pointer), internal

	Classes map[*ClassDecl]register // class definition pointers, internal

	basicInt      *ClassDecl // the same as the global basicInt
	basicString   *ClassDecl // the same as the global basicString
	basicArrayAny *ClassDecl // the same as the global basicArrayAny

	Null       bitgen.Line // null pointer dereference
	IndexRange bitgen.Line // ArrayAny index out of range
}

type register struct {
	Ptr bitgen.Variable
	Num bitgen.Integer
}

func (w *writer) Init() (start bitgen.Line) {
	w.basicInt = basicInt
	w.basicString = basicString
	w.basicArrayAny = basicArrayAny

	var registers []register
	reg := func(r *register) {
		r.Ptr = w.ReserveVariable()
		r.Num = w.ReserveInteger(32)
		registers = append(registers, *r)
	}

	w.Ptr = w.ReserveInteger(32)
	reg(&w.Alloc)
	reg(&w.Unit)
	reg(&w.True)
	reg(&w.False)
	reg(&w.Zero)
	reg(&w.Symbol)
	reg(&w.Return)
	reg(&w.This)
	reg(&w.Stack)
	w.Prev = w.ReserveVariable()
	for i := range w.General {
		reg(&w.General[i])
	}
	w.Classes = make(map[*ClassDecl]register)
	for _, c := range basicClasses {
		var r register
		reg(&r)
		w.Classes[c] = r
	}
	for _, c := range w.AST.Classes {
		var r register
		reg(&r)
		w.Classes[c] = r
	}
	w.Heap = bitgen.AddressOf{w.ReserveHeap()}

	for _, r := range registers {
		for i := uint(0); i < 32; i++ {
			next := w.ReserveLine()
			w.Assign(start, bitgen.ValueAt{bitgen.Offset{bitgen.AddressOf{r.Num.Start}, i}}, bitgen.Bit(false), next)
			start = next
		}

		next := w.ReserveLine()
		w.Assign(start, r.Ptr, w.Heap, next)
		start = next
	}

	next := w.ReserveLine()
	w.Increment(start, w.Alloc.Num, next, 0)
	start = next

	next = w.ReserveLine()
	w.Assign(start, w.Alloc.Ptr, bitgen.Offset{w.Alloc.Ptr, 8}, next)
	start = next

	w.Null = w.ReserveLine()
	w.Abort(w.Null, "runtime error: null pointer dereference")

	w.IndexRange = w.ReserveLine()
	w.Abort(w.IndexRange, "runtime error: index out of range")

	for _, c := range basicClasses {
		next = w.ReserveLine()
		w.ClassDecl(start, c, next)
		start = next
	}
	for _, c := range w.AST.Classes {
		next = w.ReserveLine()
		w.ClassDecl(start, c, next)
		start = next
	}

	for _, c := range basicClasses {
		next = w.ReserveLine()
		w.ClassDeclFixup(start, c, next)
		start = next
	}
	for _, c := range w.AST.Classes {
		next = w.ReserveLine()
		w.ClassDeclFixup(start, c, next)
		start = next
	}

	next = w.ReserveLine()
	w.NewNative(start, w.Unit, basicUnit, 0, next)
	start = next

	next = w.ReserveLine()
	w.NewNative(start, w.True, basicBoolean, 0, next)
	start = next

	next = w.ReserveLine()
	w.NewNative(start, w.False, basicBoolean, 0, next)
	start = next

	next = w.ReserveLine()
	w.NewInt(start, w.Zero, 0, next)
	start = next

	return
}

func (w *writer) ClassDecl(start bitgen.Line, c *ClassDecl, end bitgen.Line) {
	next := w.ReserveLine()
	w.NewString(start, w.General[0], w.General[1], c.Name.Name, next)
	start = next

	next = w.ReserveLine()
	w.StaticAlloc(start, w.Classes[c], (1+uint(len(c.methods)))+32/8, next)
	start = next

	w.Copy(start, bitgen.Integer{bitgen.ValueAt{w.Classes[c].Ptr}, 32}, w.General[0].Num, end)
}

func (w *writer) ClassDeclFixup(start bitgen.Line, c *ClassDecl, end bitgen.Line) {
	next := w.ReserveLine()
	w.Load(start, w.This, w.Classes[c], 0, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, bitgen.Integer{bitgen.ValueAt{w.This.Ptr}, 32}, w.Classes[basicString].Num, next)
	start = next

	next = w.ReserveLine()
	w.Load(start, w.This, w.This, basicStringLength.offset, next)
	start = next

	w.Copy(start, bitgen.Integer{bitgen.ValueAt{w.This.Ptr}, 32}, w.Classes[basicInt].Num, end)
}

func (w *writer) MethodTable() (entry bitgen.Line) {
	entry = w.ReserveLine()

	panic("unimplemented")

	return
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
	w.Assign(start, w.Prev, w.Stack.Ptr, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, bitgen.Integer{bitgen.ValueAt{w.Alloc.Ptr}, 32}, w.Stack.Num, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, w.Stack.Num, w.Alloc.Num, next)
	start = next

	next = w.ReserveLine()
	w.Assign(start, w.Stack.Ptr, w.Alloc.Ptr, next)
	start = next

	w.Offset = uint(2 + len(w.General))
	for i := uint(0); i < 32/8*w.Offset; i++ {
		next = w.ReserveLine()
		w.Increment(start, w.Alloc.Num, next, 0)
		start = next
	}

	next = w.ReserveLine()
	w.Copy(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.Stack.Ptr, 32}}, 32}, w.This.Num, next)
	start = next

	for i := 0; i < 32; i++ {
		next = w.ReserveLine()
		w.Assign(start, bitgen.ValueAt{bitgen.Offset{bitgen.AddressOf{w.This.Num.Start}, uint(32 + i)}}, bitgen.Bit(false), next)
		start = next
	}

	next = w.ReserveLine()
	w.Assign(start, w.This.Ptr, w.Heap, next)
	start = next

	for i := range w.General {
		next = w.ReserveLine()
		w.Copy(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.Stack.Ptr, uint(2*32 + i*32)}}, 32}, w.General[i].Num, next)
		start = next

		for j := 0; j < 32; j++ {
			next = w.ReserveLine()
			w.Assign(start, bitgen.ValueAt{bitgen.Offset{bitgen.AddressOf{w.General[i].Num.Start}, uint(2*32 + i*32 + j)}}, bitgen.Bit(false), next)
			start = next
		}

		next = w.ReserveLine()
		w.Assign(start, w.General[i].Ptr, w.Heap, next)
		start = next
	}

	w.Assign(start, w.Alloc.Ptr, bitgen.Offset{w.Alloc.Ptr, w.Offset * 32}, end)
}

func (w *writer) StackAlloc(start, end bitgen.Line) (cur, prev bitgen.Integer) {
	if w.Offset == 0 {
		panic("StackAlloc while not in stack")
	}

	cur, prev = w.StackOffset(w.Offset), w.PrevStackOffset(w.Offset)

	w.Offset++

	for i := 0; i < 32/8; i++ {
		next := w.ReserveLine()
		w.Increment(start, w.Alloc.Num, next, 0)
		start = next
	}

	w.Assign(start, w.Alloc.Ptr, bitgen.Offset{w.Alloc.Ptr, 32}, end)

	return
}

func (w *writer) EndStack() {
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
		w.Copy(start, w.General[i].Num, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.Stack.Ptr, uint(2*32 + i*32)}}, 32}, next)
		start = next

		next = w.ReserveLine()
		w.Pointer(start, w.General[i].Ptr, w.General[i].Num, next)
		start = next
	}

	next := w.ReserveLine()
	w.Copy(start, w.This.Num, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.Stack.Ptr, 32}}, 32}, next)
	start = next

	next = w.ReserveLine()
	w.Pointer(start, w.This.Ptr, w.This.Num, end)
	start = next

	next = w.ReserveLine()
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

func (w *writer) NewInt(start bitgen.Line, reg register, value int32, end bitgen.Line) {
	next := w.ReserveLine()
	w.NewNative(start, reg, w.basicInt, 32/8, next)
	start = next

	for i := uint(0); i < 32; i++ {
		if i == 32-1 {
			next = end
		} else {
			next = w.ReserveLine()
		}
		w.Assign(start, bitgen.ValueAt{bitgen.Offset{reg.Ptr, 32 + i}}, bitgen.Bit((uint32(value)<<uint(i))&0x80000000 == 0x80000000), next)
		start = next
	}
}

func (w *writer) NewString(start bitgen.Line, reg, length register, value string, end bitgen.Line) {
	next := w.ReserveLine()
	w.NewInt(start, length, int32(len(value)), next)
	start = next

	next = w.ReserveLine()
	w.NewNative(start, reg, basicString, uint(len(value)), next)
	start = next

	for i := range value {
		for j := 0; j < 8; j++ {
			next = w.ReserveLine()
			w.Assign(start, bitgen.ValueAt{bitgen.Offset{reg.Ptr, uint(32 + 32 + i*8 + j)}}, bitgen.Bit((value[i]<<uint(j))&0x80 == 0x80), next)
			start = next
		}
	}

	w.Copy(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{reg.Ptr, basicStringLength.offset * 8}}, 32}, length.Num, end)
}

func (w *writer) DynamicAlloc(start bitgen.Line, reg register, size bitgen.Integer, end bitgen.Line) {
	if w.Offset != 0 {
		panic("DynamicAlloc while in stack")
	}

	next := w.ReserveLine()
	w.Copy(start, w.Ptr, size, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, reg.Num, w.Alloc.Num, next)
	start = next

	next = w.ReserveLine()
	w.Assign(start, reg.Ptr, w.Alloc.Ptr, next)
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

func (w *writer) New(start bitgen.Line, reg register, c *ClassDecl, end bitgen.Line) {
	for _, f := range c.Body {
		if _, ok := f.(*NativeFeature); ok {
			panic(fmt.Errorf("cannot construct native class %s", c.Name.Name))
		}
	}

	next := end

	for _, f := range c.Body {
		if v, ok := f.(*VarFeature); ok {
			var val bitgen.Integer
			switch v.Type.target {
			case basicBoolean:
				val = w.False.Num
			case w.basicInt:
				val = w.Zero.Num
			case basicUnit:
				val = w.Unit.Num
			default:
				continue
			}
			prev := w.ReserveLine()
			w.Copy(prev, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{reg.Ptr, v.offset * 8}}, 32}, val, next)
			next = prev
		}
	}

	w.NewNative(start, reg, c, 0, next)
}

func (w *writer) NewNative(start bitgen.Line, reg register, c *ClassDecl, additional uint, end bitgen.Line) {
	next := w.ReserveLine()
	w.StaticAlloc(start, reg, c.size+additional, next)
	start = next

	w.Copy(start, bitgen.Integer{bitgen.ValueAt{reg.Ptr}, 32}, w.Classes[c].Num, end)
}

func (w *writer) NewNativeDynamic(start bitgen.Line, reg register, c *ClassDecl, additional bitgen.Integer, end bitgen.Line) {
	if reg.Num != additional {
		next := w.ReserveLine()
		w.Copy(start, reg.Num, additional, next)
		start = next
	}

	for i := uint(0); i < c.size; i++ {
		next := w.ReserveLine()
		w.Increment(start, reg.Num, next, 0)
		start = next
	}

	next := w.ReserveLine()
	w.DynamicAlloc(start, reg, reg.Num, next)
	start = next

	w.Copy(start, bitgen.Integer{bitgen.ValueAt{reg.Ptr}, 32}, w.Classes[c].Num, end)
}

func (w *writer) NewArrayAny(start bitgen.Line, end bitgen.Line) {
	w.EndStack()

	next := w.ReserveLine()
	w.Load(start, w.General[0], w.Stack, w.Arg(0), next)
	start = next

	// 1<<2 == 32/8

	for i := uint(0); i <= 2; i++ {
		next = w.ReserveLine()
		w.Jump(start, bitgen.ValueAt{bitgen.Offset{w.General[0].Ptr, 32 + 32 - 1 - i}}, next, w.IndexRange)
		start = next
	}

	// multiply array size by 32/8 (elements -> bytes)
	next = w.ReserveLine()
	w.Copy(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{bitgen.AddressOf{w.Return.Num.Start}, 2}}, w.Return.Num.Width - 2}, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[0].Ptr, 32}}, 32 - 2}, next)
	start = next

	// clear the bottom 2 bits
	for i := uint(0); i < 2; i++ {
		next = w.ReserveLine()
		w.Assign(start, bitgen.ValueAt{bitgen.Offset{bitgen.AddressOf{w.Return.Num.Start}, i}}, bitgen.Bit(false), next)
		start = next
	}

	next = w.ReserveLine()
	w.NewNativeDynamic(start, w.Return, w.basicArrayAny, w.Return.Num, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.Return.Ptr, basicArrayAnyLength.offset * 8}}, 32}, w.General[0].Num, next)
	start = next

	w.PopStack(start, end)
}

// Load puts the pointer [offset] bytes after [right] in [left]. If [right] is
// null, a runtime error occurs.
func (w *writer) Load(start bitgen.Line, left, right register, offset uint, end bitgen.Line) {
	next := w.ReserveLine()
	w.Cmp(start, right.Num, 0, w.Null, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, left.Num, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{bitgen.AddressOf{right.Ptr}, 8 * offset}}, 32}, next)
	start = next

	w.Pointer(start, left.Ptr, left.Num, end)
}

// Arg returns the stack offset of the i'th argument to the current function.
// Example usage: w.Load(start, reg, w.Stack, w.Arg(i), end)
// Example usage: w.StackOffset(w.Arg(i))
func (w *writer) Arg(i uint) uint {
	return (2 + uint(len(w.General)) + i) * 32 / 8
}

func (w *writer) StackOffset(offset uint) bitgen.Integer {
	return bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.Stack.Ptr, offset * 8}}, 32}
}

func (w *writer) PrevStackOffset(offset uint) bitgen.Integer {
	return bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.Prev, offset * 8}}, 32}
}

func (w *writer) CmpReg(start bitgen.Line, left, right bitgen.Integer, same, different bitgen.Line) {
	if left.Width != right.Width {
		panic("non-equal widths for CmpReg")
	}

	for i := uint(0); i < left.Width; i++ {
		zero, one := w.ReserveLine(), w.ReserveLine()
		w.Jump(start, bitgen.ValueAt{bitgen.Offset{bitgen.AddressOf{left.Start}, i}}, zero, one)

		var next bitgen.Line
		if i == left.Width-1 {
			next = same
		} else {
			next = w.ReserveLine()
		}

		w.Jump(zero, bitgen.ValueAt{bitgen.Offset{bitgen.AddressOf{right.Start}, i}}, next, different)
		w.Jump(one, bitgen.ValueAt{bitgen.Offset{bitgen.AddressOf{right.Start}, i}}, different, next)

		start = next
	}
}

func (w *writer) LessThanUnsigned(start bitgen.Line, left, right bitgen.Integer, less, equal, greater bitgen.Line) {
	if left.Width != right.Width {
		panic("non-equal widths for CmpReg")
	}

	for i := left.Width - 1; i < left.Width; i-- {
		zero, one := w.ReserveLine(), w.ReserveLine()
		w.Jump(start, bitgen.ValueAt{bitgen.Offset{bitgen.AddressOf{left.Start}, i}}, zero, one)

		var next bitgen.Line
		if i == 0 {
			next = equal
		} else {
			next = w.ReserveLine()
		}

		w.Jump(zero, bitgen.ValueAt{bitgen.Offset{bitgen.AddressOf{right.Start}, i}}, next, greater)
		w.Jump(one, bitgen.ValueAt{bitgen.Offset{bitgen.AddressOf{right.Start}, i}}, less, next)

		start = next
	}
}

func (w *writer) IntValue(ptr bitgen.Value) bitgen.Integer {
	return bitgen.Integer{bitgen.ValueAt{bitgen.Offset{ptr, 32}}, 32}
}

func (w *writer) AddReg(start bitgen.Line, left, right bitgen.Integer, end bitgen.Line) {
	if left.Width != right.Width {
		panic("non-equal widths for AddReg")
	}

	var carry bitgen.Line

	for i := uint(0); i < left.Width; i++ {
		curL := bitgen.ValueAt{bitgen.Offset{bitgen.AddressOf{left.Start}, i}}
		curR := bitgen.ValueAt{bitgen.Offset{bitgen.AddressOf{right.Start}, i}}

		var next, nextCarry bitgen.Line
		if i == left.Width-1 {
			next, nextCarry = end, end
		} else {
			next, nextCarry = w.ReserveLine(), w.ReserveLine()
		}

		one := w.ReserveLine()
		w.Jump(start, curR, next, one)

		if carry != 0 {
			w.Jump(carry, curR, one, nextCarry)
		}

		setOne, setTwo := w.ReserveLine(), w.ReserveLine()
		w.Jump(one, curL, setOne, setTwo)
		w.Assign(setOne, curL, bitgen.Bit(true), next)
		w.Assign(setTwo, curL, bitgen.Bit(false), nextCarry)

		start, carry = next, nextCarry
	}
}

func (w *writer) CopyReg(start bitgen.Line, left, right register, end bitgen.Line) {
	next := w.ReserveLine()
	w.Assign(start, left.Ptr, right.Ptr, next)
	start = next

	w.Copy(start, left.Num, right.Num, end)
}

func (w *writer) PrintStringArg(start bitgen.Line, arg uint, end bitgen.Line) {
	next := w.ReserveLine()
	w.Load(start, w.General[0], w.Stack, w.Arg(arg), next)
	start = next

	next = w.ReserveLine()
	w.Load(start, w.General[1], w.General[0], basicStringLength.offset, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, w.General[1].Num, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[1].Ptr, 32}}, 32}, next)
	start = next

	next = w.ReserveLine()
	w.Assign(start, w.General[0].Ptr, bitgen.Offset{w.General[0].Ptr, basicStringLength.offset*8 + 32}, next)
	start = next

	loop := start
	next = w.ReserveLine()
	w.Decrement(start, w.General[1].Num, next, end)
	start = next

	next = w.ReserveLine()
	w.Output(start, bitgen.Integer{bitgen.ValueAt{w.General[0].Ptr}, 8}, next)
	start = next

	w.Assign(start, w.General[0].Ptr, bitgen.Offset{w.General[0].Ptr, 8}, loop)
}

func (w *writer) StaticCall(start bitgen.Line, m *MethodFeature, end bitgen.Line) {
	panic("unimplemented")
}

func (w *writer) DynamicCall(start bitgen.Line, m *MethodFeature, end bitgen.Line) {
	panic("unimplemented")
}

package main

import (
	"flag"
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

	w.InitJumpTable()

	next := w.ReserveLine()
	w.MethodTables(start, next)
	start = next

	next = w.ReserveLine()
	w.New(start, w.Return, ast.main, next)
	start = next

	next = w.ReserveLine()
	w.BeginStack(start, next)
	start = next

	next = w.ReserveLine()
	w.CopyReg(start, w.This, w.Return, next)
	start = next

	w.StaticCall(start, &StaticCallExpr{
		Name: ID{
			Name:   "Main",
			target: ast.main.methods["Main"],
		},
	}, 0)

	return nil
}

type writer struct {
	*bitgen.Writer

	AST *AST

	Goto    bitgen.Integer    // a position in the jump table, saved
	Next    bitgen.Integer    // the value for Goto after the jump, internal
	Ptr     bitgen.Integer    // used for counting pointers, internal
	Alloc   register          // points at next available memory, internal
	Unit    register          // points at the constant "()", internal
	True    register          // points at the constant "true", internal
	False   register          // points at the constant "false", internal
	Zero    register          // points at the constant "0", internal
	Symbol  register          // points at the last symbol created, internal
	Return  register          // points at return value, not saved
	This    register          // points at "this" value, saved
	Stack   register          // points at start of stack segment, saved
	Offset  uint              // current stack offset, internal
	Prev    bitgen.Variable   // points at previous stack segment, internal
	General [4]register       // general purpose registers, saved
	Save    [4]bitgen.Integer // saved General registers, internal
	Heap    bitgen.AddressOf  // heap start (also null pointer), internal

	Classes map[*ClassDecl]register // class definition pointers, internal

	basicInt      *ClassDecl   // the same as the global basicInt
	basicString   *ClassDecl   // the same as the global basicString
	basicArrayAny *ClassDecl   // the same as the global basicArrayAny
	basicClasses  []*ClassDecl // the same as the global basicClasses

	Panic      bitgen.Line // AAAAAAAAAAAAAAAAAA
	Null       bitgen.Line // null pointer dereference
	IndexRange bitgen.Line // ArrayAny index out of range
	CaseNull   bitgen.Line // special case of NoCase for Null
	NoCase     bitgen.Line // missing case in match expression
	DivZero    bitgen.Line // divide by zero
	HeapRange  bitgen.Line // pointer outside of heap

	JumpTableEntry bitgen.Line
	MethodStarts   map[*MethodFeature]uint32
	DynamicCalls   map[*CallExpr]uint32
	StaticCalls    map[*StaticCallExpr]uint32
	Jumps          []bitgen.Line
}

type register struct {
	Ptr bitgen.Variable
	Num bitgen.Integer
}

func (w *writer) Init() (start bitgen.Line) {
	w.basicInt = basicInt
	w.basicString = basicString
	w.basicArrayAny = basicArrayAny
	w.basicClasses = basicClasses

	var registers []register
	reg := func(r *register) {
		r.Ptr = w.ReserveVariable()
		r.Num = w.ReserveInteger(32)
		registers = append(registers, *r)
	}

	w.Goto = w.ReserveInteger(32)
	w.Next = w.ReserveInteger(32)
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
		if w.AST.usedTypes[c] {
			var r register
			reg(&r)
			w.Classes[c] = r
		}
	}
	w.Heap = bitgen.AddressOf{w.ReserveHeap()}

	for _, r := range registers {
		for i := uint(0); i < 32; i++ {
			next := w.ReserveLine()
			w.Assign(start, r.Num.Bit(i), bitgen.Bit(false), next)
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

	w.Panic = w.ReserveLine()
	w.Dump(w.Panic, 0)

	w.Null = w.Abort("runtime error: null pointer dereference\n")

	w.IndexRange = w.Abort("runtime error: index out of range\n")

	w.CaseNull = w.ReserveLine()
	w.PrintString(w.CaseNull, "runtime error: missing case for Null\n", 0)

	w.NoCase = w.ReserveLine()
	noCaseB := w.ReserveLine()
	w.PrintString(w.NoCase, "runtime error: missing case for ", noCaseB)
	noCaseA := noCaseB

	noCaseB = w.ReserveLine()
	w.Copy(noCaseA, w.StackOffset(w.Arg(0)), bitgen.Integer{bitgen.ValueAt{w.Return.Ptr}, 32}, noCaseB)
	noCaseA = noCaseB

	noCaseB = w.ReserveLine()
	w.PrintStringArg(noCaseA, 0, noCaseB)
	noCaseA = noCaseB

	w.Print(noCaseA, '\n', 0)

	w.DivZero = w.Abort("runtime error: division by zero\n")

	w.HeapRange = w.Abort("runtime error: pointer outside of heap\n")

	for _, c := range basicClasses {
		next = w.ReserveLine()
		w.ClassDecl(start, c, next)
		start = next
	}
	for _, c := range w.AST.Classes {
		if w.AST.usedTypes[c] {
			next = w.ReserveLine()
			w.ClassDecl(start, c, next)
			start = next
		}
	}

	for _, c := range basicClasses {
		next = w.ReserveLine()
		w.ClassDeclFixup(start, c, next)
		start = next
	}
	for _, c := range w.AST.Classes {
		if w.AST.usedTypes[c] {
			next = w.ReserveLine()
			w.ClassDeclFixup(start, c, next)
			start = next
		}
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

	next = w.ReserveLine()
	w.CopyReg(start, w.General[0], w.General[2], next)
	start = next

	next = w.ReserveLine()
	w.CopyReg(start, w.General[1], w.General[2], next)
	start = next

	return
}

func (w *writer) Abort(message string) bitgen.Line {
	start := w.ReserveLine()
	w.PrintString(start, message, w.Panic)
	return start
}

func (w *writer) ClassDecl(start bitgen.Line, c *ClassDecl, end bitgen.Line) {
	next := w.ReserveLine()
	w.NewString(start, w.General[0], w.General[1], c.Name.Name, next)
	start = next

	next = w.ReserveLine()
	w.StaticAlloc(start, w.Classes[c], (1+uint(len(c.methods)))*32/8, next)
	start = next

	w.Copy(start, bitgen.Integer{bitgen.ValueAt{w.Classes[c].Ptr}, 32}, w.General[0].Num, end)
}

func (w *writer) ClassDeclFixup(start bitgen.Line, c *ClassDecl, end bitgen.Line) {
	next := w.ReserveLine()
	w.Load(start, w.General[0], w.Classes[c], 0, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, bitgen.Integer{bitgen.ValueAt{w.General[0].Ptr}, 32}, w.Classes[basicString].Num, next)
	start = next

	next = w.ReserveLine()
	w.Load(start, w.General[0], w.General[0], basicStringLength.offset, next)
	start = next

	w.Copy(start, bitgen.Integer{bitgen.ValueAt{w.General[0].Ptr}, 32}, w.Classes[basicInt].Num, end)
}

var flagPrintJumpTable = flag.Bool("print-jump-table", false, "debugging tool: print jump table")

func (w *writer) InitJumpTable() {
	w.MethodStarts = make(map[*MethodFeature]uint32)
	w.DynamicCalls = make(map[*CallExpr]uint32)
	w.StaticCalls = make(map[*StaticCallExpr]uint32)
	w.Jumps = append(w.Jumps, 0)

	var recurse func(Expr)
	recurse = func(value Expr) {
		switch v := value.(type) {
		case *ConstructorExpr:
			recurse(v.Expr)

		case *AssignExpr:
			recurse(v.Right)

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
				recurse(c.Body)
			}

		case *CallExpr:
			recurse(v.Left)
			for _, a := range v.Args {
				recurse(a)
			}

			l := w.ReserveLine()
			j := uint32(len(w.Jumps))
			w.DynamicCalls[v] = j
			w.Jumps = append(w.Jumps, l)
			if *flagPrintJumpTable {
				pos := w.AST.FileSet.Position(v.Name.Pos)
				fmt.Printf("%08X\t%d\t%v\n", j, l, pos)
			}

		case *StaticCallExpr:
			for _, a := range v.Args {
				recurse(a)
			}

			l := w.ReserveLine()
			j := uint32(len(w.Jumps))
			w.StaticCalls[v] = j
			w.Jumps = append(w.Jumps, l)
			if *flagPrintJumpTable {
				pos := w.AST.FileSet.Position(v.Name.Pos)
				fmt.Printf("%08X\t%d\t%v\n", j, l, pos)
			}

		case *NewExpr:

		case *VarExpr:
			recurse(v.Value)
			recurse(v.Expr)

		case *ChainExpr:
			recurse(v.Pre)
			recurse(v.Expr)

		case *NullExpr:

		case *UnitExpr:

		case *NameExpr:

		case *IntegerExpr:

		case *StringExpr:

		case *BooleanExpr:

		case *ThisExpr:

		case NativeExpr:

		default:
			panic(v)
		}
	}

	method := func(m *MethodFeature) {
		l := w.ReserveLine()
		j := uint32(len(w.Jumps))
		w.MethodStarts[m] = j
		w.Jumps = append(w.Jumps, l)
		if *flagPrintJumpTable {
			pos := w.AST.FileSet.Position(m.Name.Pos)
			fmt.Printf("%08X\t%d\t%v\n", j, l, pos)
		}

		recurse(m.Body)
	}

	for _, m := range basicDummyCalls {
		recurse(m)
	}
	for _, c := range basicClasses {
		for _, f := range c.Body {
			if m, ok := f.(*MethodFeature); ok {
				method(m)
			}
		}
	}
	for _, c := range w.AST.Classes {
		if w.AST.usedTypes[c] {
			for _, f := range c.Body {
				if m, ok := f.(*MethodFeature); ok {
					method(m)
				}
			}
		}
	}

	w.JumpTableEntry = w.ReserveLine()

	jump := func(n uint32) bitgen.Line {
		if n == uint32(len(w.Jumps)) {
			return w.Panic
		}

		end := w.Jumps[n]

		prev := w.ReserveLine()
		w.Copy(prev, w.Goto, w.Next, end)
		end = prev

		return end
	}

	var binary func(bit uint, min, max uint32) bitgen.Line
	binary = func(bit uint, min, max uint32) bitgen.Line {
		if min >= uint32(len(w.Jumps)) {
			return w.Panic
		}

		var n0, n1 bitgen.Line
		if bit == 0 {
			n0 = jump(min)
			n1 = jump(max)
		} else {
			mid := (min + max) >> 1
			n0 = binary(bit-1, min, mid)
			n1 = binary(bit-1, mid+1, max)
		}

		start := w.ReserveLine()
		w.Jump(start, w.Goto.Bit(bit), n0, n1)

		return start
	}

	n0 := binary(32-2, 0, 1<<(32-1)-1)
	n1 := binary(32-2, 1<<(32-1), 1<<32-1)
	w.Jump(w.JumpTableEntry, w.Goto.Bit(32-1), n0, n1)

	pop := w.ReserveLine()
	w.SimplePopStack(pop, w.JumpTableEntry)

	method = func(m *MethodFeature) {
		start := w.Jumps[w.MethodStarts[m]]

		w.Offset = 3 + uint(len(m.Args))

		if _, ok := m.Body.(NativeExpr); ok {
			next := w.ReserveLine()
			w.SaveRegisters(start, next)
			start = next

			m.Body.write(w, start, pop)
		} else {
			start = m.Body.alloc(w, start)
			w.EndStack()

			m.Body.write(w, start, pop)
		}
	}

	for _, c := range basicClasses {
		for _, f := range c.Body {
			if m, ok := f.(*MethodFeature); ok {
				method(m)
			}
		}
	}
	for _, c := range w.AST.Classes {
		if w.AST.usedTypes[c] {
			for _, f := range c.Body {
				if m, ok := f.(*MethodFeature); ok {
					method(m)
				}
			}
		}
	}
}

func (w *writer) MethodTables(start, end bitgen.Line) {
	methods := func(start bitgen.Line, c *ClassDecl, end bitgen.Line) {
		cr := w.Classes[c].Ptr
		for i := uint(0); i < uint(len(c.methods)); i++ {
			for _, m := range c.methods {
				if m.offset != i {
					continue
				}
				jump := w.MethodStarts[m]
				for j := uint(32 - 1); j < 32; j-- {
					if (jump>>j)&1 == 0 {
						continue
					}
					var next bitgen.Line
					if i == uint(len(c.methods))-1 && jump&(1<<j-1) == 0 {
						next = end
					} else {
						next = w.ReserveLine()
					}

					w.Assign(start, bitgen.ValueAt{bitgen.Offset{cr, 32 + 32*i + j}}, bitgen.Bit(true), next)

					start = next
				}
			}
		}
	}

	var next bitgen.Line
	for _, c := range basicClasses {
		next = w.ReserveLine()
		methods(start, c, next)
		start = next
	}

	var usedClasses []*ClassDecl
	for _, c := range w.AST.Classes {
		if w.AST.usedTypes[c] {
			usedClasses = append(usedClasses, c)
		}
	}

	for i, c := range usedClasses {
		if i == len(usedClasses)-1 {
			next = end
		} else {
			next = w.ReserveLine()
		}
		methods(start, c, next)
		start = next
	}
}

func (w *writer) Pointer(start bitgen.Line, ptr bitgen.Variable, num bitgen.Integer, end bitgen.Line) {
	next := w.ReserveLine()
	w.LessThanUnsigned(start, num, w.Alloc.Num, next, w.HeapRange, w.HeapRange)
	start = next

	next = w.ReserveLine()
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

	w.Offset = uint(3)
	next = w.ReserveLine()
	w.Add(start, w.Alloc.Num, uint64(w.Offset*32/8), next, 0)
	start = next

	next = w.ReserveLine()
	w.Assign(start, w.Alloc.Ptr, bitgen.Offset{w.Alloc.Ptr, w.Offset * 32}, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, w.StackOffset(32/8), w.This.Num, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, w.StackOffset(32/8+32/8), w.Goto, next)
	start = next

	for i := 0; i < 32; i++ {
		next = w.ReserveLine()
		w.Assign(start, w.This.Num.Bit(uint(i)), bitgen.Bit(false), next)
		start = next

		next = w.ReserveLine()
		w.Assign(start, w.Goto.Bit(uint(i)), bitgen.Bit(false), next)
		start = next
	}

	w.Assign(start, w.This.Ptr, w.Heap, end)
}

func (w *writer) SaveRegisters(start, end bitgen.Line) {
	if w.Offset == 0 {
		panic("SaveRegisters while not in stack")
	}

	for i := range w.General {
		next := w.ReserveLine()
		w.Save[i], _ = w.StackAlloc(start, next)
		start = next

		next = w.ReserveLine()
		w.Copy(start, w.Save[i], w.General[i].Num, next)
		start = next

		for j := uint(0); j < 32; j++ {
			next = w.ReserveLine()
			w.Assign(start, w.General[i].Num.Bit(j), bitgen.Bit(false), next)
			start = next
		}

		if i == len(w.General)-1 {
			next = end
		} else {
			next = w.ReserveLine()
		}
		w.Assign(start, w.General[i].Ptr, w.Heap, next)
		start = next
	}
}

func (w *writer) StackAlloc(start, end bitgen.Line) (cur, prev bitgen.Integer) {
	if w.Offset == 0 {
		panic("StackAlloc while not in stack")
	}

	cur, prev = w.StackOffset(w.Offset*32/8), w.PrevStackOffset(w.Offset*32/8)

	w.Offset++

	next := w.ReserveLine()
	w.Add(start, w.Alloc.Num, 32/8, next, 0)
	start = next

	w.Assign(start, w.Alloc.Ptr, bitgen.Offset{w.Alloc.Ptr, 32}, end)

	return
}

func (w *writer) EndStack() {
	if w.Offset == 0 {
		panic("EndStack while not in stack")
	}

	w.Offset = 0
}

func (w *writer) SimplePopStack(start, end bitgen.Line) {
	if w.Offset != 0 {
		panic("SimplePopStack while in stack")
	}

	next := w.ReserveLine()
	w.Copy(start, w.Next, w.StackOffset(32/8+32/8), next)
	start = next

	next = w.ReserveLine()
	w.Load(start, w.This, w.Stack, 32/8, next)
	start = next

	w.Load(start, w.Stack, w.Stack, 0, end)
}

func (w *writer) PopStack(start, end bitgen.Line) {
	if w.Offset != 0 {
		panic("PopStack while in stack")
	}

	for i := range w.General {
		next := w.ReserveLine()
		w.Copy(start, w.General[i].Num, w.Save[i], next)
		start = next

		if i == len(w.General)-1 {
			next = end
		} else {
			next = w.ReserveLine()
		}
		w.Pointer(start, w.General[i].Ptr, w.General[i].Num, next)
		start = next
	}
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

	next = w.ReserveLine()
	w.Add(start, w.Alloc.Num, uint64(size), next, 0)
	start = next

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
		w.Assign(start, bitgen.ValueAt{bitgen.Offset{reg.Ptr, 32 + i}}, bitgen.Bit((uint32(value)>>uint(i))&1 == 1), next)
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
			w.Assign(start, bitgen.ValueAt{bitgen.Offset{reg.Ptr, basicStringLength.offset*8 + 32 + uint(i*8+j)}}, bitgen.Bit((value[i]>>uint(j))&1 == 1), next)
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
	if c != basicArrayAny {
		for _, f := range c.Body {
			if _, ok := f.(*NativeFeature); ok {
				panic(fmt.Errorf("cannot construct native class %s", c.Name.Name))
			}
		}
	}

	next := end

	for p := c; p != basicAny; p = p.Extends.Type.target {
		for _, a := range p.Args {
			var val bitgen.Integer
			switch a.Type.target {
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
			w.Copy(prev, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{reg.Ptr, a.offset * 8}}, 32}, val, next)
			next = prev
		}
		for _, f := range p.Body {
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

	next := w.ReserveLine()
	w.Add(start, reg.Num, uint64(c.size), next, 0)
	start = next

	next = w.ReserveLine()
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
	w.Copy(start, w.Return.Num.Sub(2, 32), w.IntValue(w.General[0].Ptr).Sub(0, 32-2), next)
	start = next

	// clear the bottom 2 bits
	for i := uint(0); i < 2; i++ {
		next = w.ReserveLine()
		w.Assign(start, w.Return.Num.Bit(i), bitgen.Bit(false), next)
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
	w.Copy(start, left.Num, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{right.Ptr, 8 * offset}}, 32}, next)
	start = next

	w.Pointer(start, left.Ptr, left.Num, end)
}

// Arg returns the stack offset of the i'th argument to the current function.
// Example usage: w.Load(start, reg, w.Stack, w.Arg(i), end)
// Example usage: w.StackOffset(w.Arg(i))
func (w *writer) Arg(i uint) uint {
	return (3 + i) * 32 / 8
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
		w.Jump(start, left.Bit(i), zero, one)

		var next bitgen.Line
		if i == left.Width-1 {
			next = same
		} else {
			next = w.ReserveLine()
		}

		w.Jump(zero, right.Bit(i), next, different)
		w.Jump(one, right.Bit(i), different, next)

		start = next
	}
}

func (w *writer) LessThanUnsigned(start bitgen.Line, left, right bitgen.Integer, less, equal, greater bitgen.Line) {
	if left.Width != right.Width {
		panic("non-equal widths for LessThanUnsigned")
	}

	for i := left.Width - 1; i < left.Width; i-- {
		zero, one := w.ReserveLine(), w.ReserveLine()
		w.Jump(start, left.Bit(i), zero, one)

		var next bitgen.Line
		if i == 0 {
			next = equal
		} else {
			next = w.ReserveLine()
		}

		w.Jump(zero, right.Bit(i), next, less)
		w.Jump(one, right.Bit(i), greater, next)

		start = next
	}
}

func (w *writer) IntValue(ptr bitgen.Value) bitgen.Integer {
	return bitgen.Integer{bitgen.ValueAt{bitgen.Offset{ptr, 32}}, 32}
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
	w.Copy(start, w.General[1].Num, w.IntValue(w.General[1].Ptr), next)
	start = next

	loop := start
	next = w.ReserveLine()
	w.Decrement(start, w.General[1].Num, next, end)
	start = next

	next = w.ReserveLine()
	w.Output(start, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.General[0].Ptr, basicStringLength.offset*8 + 32}}, 8}, next)
	start = next

	w.Assign(start, w.General[0].Ptr, bitgen.Offset{w.General[0].Ptr, 8}, loop)
}

func (w *writer) gotoNext(start bitgen.Line, nextVal uint32, end bitgen.Line) {
	for i := uint(0); i < 32; i++ {
		var next bitgen.Line
		if i == 32-1 {
			next = w.JumpTableEntry
		} else {
			next = w.ReserveLine()
		}
		w.Assign(start, w.Next.Bit(i), bitgen.Bit((nextVal>>i)&1 == 1), next)
		start = next
	}

	if w.Jumps[nextVal] != 0 || end != 0 {
		w.Jump(w.Jumps[nextVal], bitgen.Bit(false), end, end)
	}
}

func (w *writer) StaticCall(start bitgen.Line, m *StaticCallExpr, end bitgen.Line) {
	w.EndStack()

	nextVal := w.StaticCalls[m]
	gotoVal := w.MethodStarts[m.Name.target.(*MethodFeature)]

	for i := uint(0); i < 32; i++ {
		next := w.ReserveLine()
		w.Assign(start, w.Goto.Bit(i), bitgen.Bit((gotoVal>>i)&1 == 1), next)
		start = next
	}

	w.gotoNext(start, nextVal, end)
}

func (w *writer) DynamicCall(start bitgen.Line, m *CallExpr, end bitgen.Line) {
	w.EndStack()

	nextVal := w.DynamicCalls[m]

	next := w.ReserveLine()
	w.Load(start, w.Return, w.This, 0, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, w.Goto, bitgen.Integer{bitgen.ValueAt{bitgen.Offset{w.Return.Ptr, 32 + m.Name.target.(*MethodFeature).offset*32}}, 32}, next)
	start = next

	w.gotoNext(start, nextVal, end)
}

func (w *writer) hexDigit(start bitgen.Line, n bitgen.Integer, end bitgen.Line) {
	n0, n1 := w.ReserveLine(), w.ReserveLine()
	w.Jump(start, n.Bit(3), n0, n1)

	n00, n01 := w.ReserveLine(), w.ReserveLine()
	w.Jump(n0, n.Bit(2), n00, n01)

	n10, n11 := w.ReserveLine(), w.ReserveLine()
	w.Jump(n1, n.Bit(2), n10, n11)

	n000, n001 := w.ReserveLine(), w.ReserveLine()
	w.Jump(n00, n.Bit(1), n000, n001)

	n010, n011 := w.ReserveLine(), w.ReserveLine()
	w.Jump(n01, n.Bit(1), n010, n011)

	n100, n101 := w.ReserveLine(), w.ReserveLine()
	w.Jump(n10, n.Bit(1), n100, n101)

	n110, n111 := w.ReserveLine(), w.ReserveLine()
	w.Jump(n11, n.Bit(1), n110, n111)

	n0000, n0001 := w.ReserveLine(), w.ReserveLine()
	w.Jump(n000, n.Bit(0), n0000, n0001)

	n0010, n0011 := w.ReserveLine(), w.ReserveLine()
	w.Jump(n001, n.Bit(0), n0010, n0011)

	n0100, n0101 := w.ReserveLine(), w.ReserveLine()
	w.Jump(n010, n.Bit(0), n0100, n0101)

	n0110, n0111 := w.ReserveLine(), w.ReserveLine()
	w.Jump(n011, n.Bit(0), n0110, n0111)

	n1000, n1001 := w.ReserveLine(), w.ReserveLine()
	w.Jump(n100, n.Bit(0), n1000, n1001)

	n1010, n1011 := w.ReserveLine(), w.ReserveLine()
	w.Jump(n101, n.Bit(0), n1010, n1011)

	n1100, n1101 := w.ReserveLine(), w.ReserveLine()
	w.Jump(n110, n.Bit(0), n1100, n1101)

	n1110, n1111 := w.ReserveLine(), w.ReserveLine()
	w.Jump(n111, n.Bit(0), n1110, n1111)

	w.Print(n0000, '0', end)
	w.Print(n0001, '1', end)
	w.Print(n0010, '2', end)
	w.Print(n0011, '3', end)
	w.Print(n0100, '4', end)
	w.Print(n0101, '5', end)
	w.Print(n0110, '6', end)
	w.Print(n0111, '7', end)
	w.Print(n1000, '8', end)
	w.Print(n1001, '9', end)
	w.Print(n1010, 'A', end)
	w.Print(n1011, 'B', end)
	w.Print(n1100, 'C', end)
	w.Print(n1101, 'D', end)
	w.Print(n1110, 'E', end)
	w.Print(n1111, 'F', end)
}

func (w *writer) hex(start bitgen.Line, n bitgen.Integer, end bitgen.Line) {
	if n.Width&(4-1) != 0 {
		panic(fmt.Sprintf("invalid hex width: %d", n.Width))
	}

	for i := n.Width - 4; i < n.Width; i -= 4 {
		var next bitgen.Line
		if i == 0 {
			next = end
		} else {
			next = w.ReserveLine()
		}

		w.hexDigit(start, n.Sub(i, i+4), next)

		start = next
	}
}

func (w *writer) Dump(start, end bitgen.Line) {
	for _, r := range [...]struct {
		name string
		reg  bitgen.Integer
	}{
		{"Goto", w.Goto},
		{"Next", w.Next},
		{"Alloc", w.Alloc.Num},
		{"Unit", w.Unit.Num},
		{"True", w.True.Num},
		{"False", w.False.Num},
		{"Zero", w.Zero.Num},
		{"Symbol", w.Symbol.Num},
		{"Return", w.Return.Num},
		{"This", w.This.Num},
		{"Stack", w.Stack.Num},
		{"Gen-0", w.General[0].Num},
		{"Gen-1", w.General[1].Num},
		{"Gen-2", w.General[2].Num},
		{"Gen-3", w.General[3].Num},
	} {
		next := w.ReserveLine()
		w.PrintString(start, "\n"+r.name+"\t", next)
		start = next

		next = w.ReserveLine()
		w.hex(start, r.reg, next)
		start = next
	}

	next := w.ReserveLine()
	w.PrintString(start, "\n\nClasses:", next)
	start = next

	for _, c := range w.basicClasses {
		next = w.ReserveLine()
		w.PrintString(start, "\n"+c.Name.Name+"\t", next)
		start = next

		next = w.ReserveLine()
		w.hex(start, w.Classes[c].Num, next)
		start = next
	}
	for _, c := range w.AST.Classes {
		next = w.ReserveLine()
		w.PrintString(start, "\n"+c.Name.Name+"\t", next)
		start = next

		if w.AST.usedTypes[c] {
			next = w.ReserveLine()
			w.hex(start, w.Classes[c].Num, next)
			start = next
		} else {
			next = w.ReserveLine()
			w.PrintString(start, "OMITTED", next)
			start = next
		}
	}

	next = w.ReserveLine()
	w.PrintString(start, "\nthis:\t", next)
	start = next

	next = w.ReserveLine()
	w.hex(start, bitgen.Integer{bitgen.ValueAt{w.This.Ptr}, 32}, next)
	start = next

	next = w.ReserveLine()
	w.PrintString(start, "\n\nHeap:\n", next)
	start = next

	next = w.ReserveLine()
	w.Assign(start, w.Alloc.Ptr, w.Heap, next)
	start = next

	next = w.ReserveLine()
	w.Copy(start, w.Ptr, w.Alloc.Num, next)
	start = next

	loop, done := start, w.ReserveLine()
	next = w.ReserveLine()
	w.Decrement(loop, w.Ptr, next, done)
	start = next

	next = w.ReserveLine()
	w.hex(start, bitgen.Integer{bitgen.ValueAt{w.Alloc.Ptr}, 8}, next)
	start = next

	w.Assign(start, w.Alloc.Ptr, bitgen.Offset{w.Alloc.Ptr, 8}, loop)

	w.PrintString(done, "\n", end)
}

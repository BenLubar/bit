package bit

import (
	"bufio"
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/BenLubar/bit/bitgen"
	"github.com/BenLubar/bit/bitio"
	"github.com/BenLubar/bit/cmd/brainfuckc/bf"
)

func testOpt(t *testing.T, bfprog []bf.Command) {
	r, w := io.Pipe()

	errs := make(chan error, 2)
	ch := make(chan Program, 1)

	go func() {
		bfw := bf.NewWriter(w)
		_, err := bfw.Program(bfprog)
		if err != nil {
			t.Error(err)
			errs <- err
			return
		}
		err = bfw.Close()
		if err != nil {
			t.Error(err)
			errs <- err
			return
		}
		err = w.Close()
		if err != nil {
			t.Error(err)
			errs <- err
			return
		}
	}()

	go func() {
		prog, err := Parse(bufio.NewReader(r))
		if err != nil {
			t.Error(err)
			errs <- err
			return
		}
		ch <- prog
	}()

	select {
	case err := <-errs:
		panic(err)

	case prog := <-ch:
		ctx0, err := prog.run(nil, nil, nil)
		if err != nil {
			panic(err)
		}

		prog.Optimize()

		trace := func(l *line, c *context) {
			if l.opt != nil {
				t.Logf("%d: %#v", l.num, l.opt)
			} else {
				t.Logf("%d: %#v", l.num, l.stmt)
			}
		}
		ctx1, err := prog.run(nil, nil, trace)
		if err != nil {
			panic(err)
		}

		if !reflect.DeepEqual(ctx0, ctx1) {
			t.Errorf("expected equal:\n%#v\n%#v", ctx0, ctx1)
		}
	}
}

func TestOptIncrement(t *testing.T) {
	testOpt(t, []bf.Command{
		{Token: bf.Increment},
	})
}

func TestOptDecrement(t *testing.T) {
	testOpt(t, []bf.Command{
		{Token: bf.Decrement},
	})
}

func TestOptIncrementDecrement(t *testing.T) {
	testOpt(t, []bf.Command{
		{Token: bf.Increment},
		{Token: bf.Decrement},
	})
}

func TestOptDecrementIncrement(t *testing.T) {
	testOpt(t, []bf.Command{
		{Token: bf.Decrement},
		{Token: bf.Increment},
	})
}

func TestOptIncrementIncrement(t *testing.T) {
	testOpt(t, []bf.Command{
		{Token: bf.Increment},
		{Token: bf.Increment},
	})
}

func TestOptDecrementDecrement(t *testing.T) {
	testOpt(t, []bf.Command{
		{Token: bf.Decrement},
		{Token: bf.Decrement},
	})
}

func TestOptRight(t *testing.T) {
	testOpt(t, []bf.Command{
		{Token: bf.Right},
	})
}

func TestOptLeft(t *testing.T) {
	testOpt(t, []bf.Command{
		{Token: bf.Left},
	})
}

func TestOptRightLeft(t *testing.T) {
	testOpt(t, []bf.Command{
		{Token: bf.Right},
		{Token: bf.Left},
	})
}

func testOptRaw(t *testing.T, f func(w *bitgen.Writer) (int64, error)) {
	r, w := io.Pipe()

	errs := make(chan error, 2)
	ch := make(chan Program, 1)

	go func() {
		bw := bitgen.NewWriter(w)

		_, err := f(bw)
		if err != nil {
			t.Error(err)
			errs <- err
			return
		}

		err = bw.Close()
		if err != nil {
			t.Error(err)
			errs <- err
			return
		}

		err = w.Close()
		if err != nil {
			t.Error(err)
			errs <- err
			return
		}
	}()

	go func() {
		prog, err := Parse(bufio.NewReader(r))
		if err != nil {
			t.Error(err)
			errs <- err
			return
		}
		ch <- prog
	}()

	select {
	case err := <-errs:
		panic(err)

	case prog := <-ch:
		var buf0, buf1 bytes.Buffer

		ctx0, err := prog.run(nil, bitio.NewWriter(&buf0), nil)
		if err != nil {
			panic(err)
		}

		prog.Optimize()

		trace := func(l *line, c *context) {
			if l.opt != nil {
				t.Logf("%d: %#v", l.num, l.opt)
			} else {
				t.Logf("%d: %#v", l.num, l.stmt)
			}
		}
		ctx1, err := prog.run(nil, bitio.NewWriter(&buf1), trace)
		if err != nil {
			panic(err)
		}

		if !reflect.DeepEqual(ctx0, ctx1) || !bytes.Equal(buf0.Bytes(), buf1.Bytes()) {
			t.Errorf("expected equal:\n%#v\n%q\n%#v\n%q", ctx0, buf0.Bytes(), ctx1, buf1.Bytes())
		}
	}
}

func TestOptZero32(t *testing.T) {
	testOptRaw(t, func(w *bitgen.Writer) (n int64, err error) {
		num := w.ReserveInteger(32)

		var start bitgen.Line
		for i := uint(0); i < num.Width; i++ {
			var next bitgen.Line
			if i != num.Width-1 {
				next = w.ReserveLine()
			}
			var nn int64
			nn, err = w.Assign(start, num.Bit(i), bitgen.Bit(false), next)
			n += nn
			if err != nil {
				return
			}
			start = next
		}
		return
	})
}

func TestOptAlloc64(t *testing.T) {
	testOptRaw(t, func(w *bitgen.Writer) (n int64, err error) {
		type register struct {
			Ptr bitgen.Variable
			Num bitgen.Integer
		}
		reg := register{
			Ptr: w.ReserveVariable(),
			Num: w.ReserveInteger(32),
		}
		alloc := register{
			Ptr: w.ReserveVariable(),
			Num: w.ReserveInteger(32),
		}
		heapStart := bitgen.AddressOf{w.ReserveHeap()}

		var start bitgen.Line
		for _, r := range []register{reg, alloc} {
			for i := uint(0); i < r.Num.Width; i++ {
				next := w.ReserveLine()
				w.Assign(start, r.Num.Bit(i), bitgen.Bit(false), next)
				start = next
			}

			next := w.ReserveLine()
			w.Assign(start, r.Ptr, heapStart, next)
			start = next
		}

		next := w.ReserveLine()
		w.Increment(start, alloc.Num, next, 0)
		start = next

		next = w.ReserveLine()
		w.Assign(start, alloc.Ptr, bitgen.Offset{alloc.Ptr, 8}, next)
		start = next

		next = w.ReserveLine()
		w.Copy(start, reg.Num, alloc.Num, next)
		start = next

		next = w.ReserveLine()
		w.Assign(start, reg.Ptr, alloc.Ptr, next)
		start = next

		const size = 64 / 8

		next = w.ReserveLine()
		w.Add(start, alloc.Num, size, next, 0)
		start = next

		w.Assign(start, alloc.Ptr, bitgen.Offset{alloc.Ptr, 8 * size}, 0)

		return
	})
}

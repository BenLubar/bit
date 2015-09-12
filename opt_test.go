package bit

import (
	"bufio"
	"io"
	"reflect"
	"testing"

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

		trace := func(l *line) {
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

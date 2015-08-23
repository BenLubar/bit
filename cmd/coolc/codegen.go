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

	panic("unimplemented")
}

type writer struct {
	*bitgen.Writer
}

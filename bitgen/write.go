package bitgen

import (
	"io"
	"math/bits"
)

func (w *Writer) write(str string) {
	if w.err != nil {
		return
	}

	_, w.err = io.WriteString(w.writer, str)
}

func (w *Writer) writeBit(b bool) {
	if b {
		w.write(" ONE")
	} else {
		w.write(" ZERO")
	}
}

func (w *Writer) writeNumber(n uint64) {
	if n == 0 {
		w.writeBit(false)
		return
	}

	for i := uint(64 - bits.LeadingZeros64(n)); i > 0; i-- {
		w.writeBit((n>>(i-1))&1 == 1)
	}
}

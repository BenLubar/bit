// Package bitgen implements a code generator backend for BIT.
package bitgen

import (
	"fmt"
	"io"
)

// Writer implements the BIT code generator.
type Writer struct {
	writer    io.Writer
	err       error
	wroteLine []bool
	refLine   []bool
	isPtr     []bool
}

// NewWriter returns a new Writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
		wroteLine: []bool{
			Exit: true,
		},
		refLine: []bool{
			Exit: false,
		},
	}
}

// Close returns the first error encountered while writing. If no write error
// occurred but a line number has been referenced without being written, the
// returned error will be non-nil. If no lines have been written, the error
// will be non-nil. In any other case, the error is nil.
//
// Close does not invalidate the Writer. However, if Close returns nil,
// additional lines added to the program are guaranteed to be dead code.
func (w *Writer) Close() error {
	if w.err != nil {
		return w.err
	}

	for i, b := range w.refLine {
		if b && !w.wroteLine[i] {
			return fmt.Errorf("bitgen: undefined reference to line %v", LineNumber(i))
		}
	}

	for _, b := range w.wroteLine[1:] {
		if b {
			return nil
		}
	}

	return fmt.Errorf("bitgen: no lines written")
}

// ReserveLine returns a new LineNumber unique within this Writer.
func (w *Writer) ReserveLine() LineNumber {
	num := LineNumber(len(w.wroteLine))
	w.wroteLine = append(w.wroteLine, false)
	w.refLine = append(w.refLine, false)
	return num
}

// NewPointer returns a new Pointer unique within this Writer.
func (w *Writer) NewPointer() Pointer {
	ptr := Pointer(len(w.isPtr))
	w.isPtr = append(w.isPtr, true)
	return ptr
}

// NewBuffer returns a new Buffer unique within this Writer.
func (w *Writer) NewBuffer() Buffer {
	buf := Buffer(len(w.isPtr))
	w.isPtr = append(w.isPtr, false)
	return buf
}

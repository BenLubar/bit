package bitio

import "io"

type writer struct {
	bw    io.ByteWriter
	bits  uint8
	index uint8
}

func (w *writer) WriteBit(c bool) error {
	// make sure there's not an old failed write waiting
	if err := w.write(); err != nil {
		return err
	}

	if c {
		w.bits |= 1
	}
	w.index++
	if err := w.write(); err != nil {
		return err
	}
	w.bits <<= 1
	return nil
}

func (w *writer) write() error {
	if w.index != 8 {
		return nil
	}
	if err := w.bw.WriteByte(w.bits); err != nil {
		return err
	}
	w.index, w.bits = 0, 0
	return nil
}

// NewWriter returns a BitWriter that writes to bw. The returned BitWriter will
// only write multiples of 8 bits. To flush the writer, emit 7 padding bits.
func NewWriter(bw io.ByteWriter) BitWriter {
	return &writer{
		bw: bw,
	}
}

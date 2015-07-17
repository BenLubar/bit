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

func (w *writer) WriteByte(c byte) error {
	// make sure there's not an old failed write waiting
	if err := w.write(); err != nil {
		return err
	}

	if w.index == 0 {
		// fast path: write the entire byte at once
		return w.bw.WriteByte(c)
	}

	for i := 0; i < 8; i++ {
		w.bits |= c >> uint(7-i) & 1
		w.index++
		if err := w.write(); err != nil {
			return err
		}
		w.bits <<= 1
	}
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

package bitio

import "io"

type pipe chan bool

func (p pipe) ReadBit() (c bool, err error) {
	c, ok := <-p
	if !ok {
		err = io.EOF
	}
	return
}

func (p pipe) WriteBit(c bool) error {
	p <- c
	return nil
}

// Pipe creates a synchronous in-memory pipe. It can be used to connect code
// expecting a BitReader with code expecting a BitWriter. Reads on one end are
// matched with writes on the other, copying data directly between the two;
// there is no internal buffering. It is safe to call Read and Write in parallel
// with each other. Parallel calls to Read, and parallel calls to Write, are
// also safe: the individual calls will be gated sequentially.
func Pipe() (BitReader, BitWriter) {
	p := make(pipe)
	return p, p
}

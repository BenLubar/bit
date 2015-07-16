package bitio

import "io"

// Null is a BitReader, BitScanner, and BitWriter that emulates /dev/null.
// It is a pointer to avoid a pointer dereference when calling it in interface
// form.
var Null *null
var (
	_ BitReader  = Null
	_ BitScanner = Null
	_ BitWriter  = Null
)

type null struct{}

func (*null) ReadBit() (bool, error) { return false, io.EOF }
func (*null) UnreadBit() error       { return nil }
func (*null) WriteBit(bool) error    { return nil }

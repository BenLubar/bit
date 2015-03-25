// Package bitio provides basic interfaces to bit I/O primitives.
package bitio

// BitReader is the interface that wraps the ReadBit method.
type BitReader interface {
	// ReadBit reads and returns the next bit from the input. If no bit is
	// available, err will be set.
	ReadBit() (c bool, err error)
}

// BitScanner is the interface that adds the UnreadBit method to the basic
// ReadBit method.
type BitScanner interface {
	BitReader

	// UnreadBit causes the next call to ReadByte to return the same bit as
	// the previous call to ReadBit. It may be an error to call UnreadBit
	// twice without an intervening call to ReadBit.
	UnreadBit() error
}

// BitWriter is the interface that wraps the WriteBit method.
type BitWriter interface {
	WriteBit(c bool) error
}

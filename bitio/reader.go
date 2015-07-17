package bitio

import (
	"errors"
	"io"
)

type reader struct {
	br     io.ByteReader
	bits   uint16
	remain uint8
	unread bool
}

var ErrDoubleUnread = errors.New("bitio: UnreadBit called without a call to ReadBit")

func (r *reader) ReadBit() (c bool, err error) {
	if r.remain == 0 {
		b, err := r.br.ReadByte()
		if err != nil {
			return false, err
		}
		r.remain = 8
		r.bits <<= 8
		r.bits |= uint16(b)
	}

	r.remain--
	r.unread = false
	return (r.bits>>r.remain)&1 == 1, nil
}

func (r *reader) ReadByte() (c byte, err error) {
	if r.remain < 8 {
		b, err := r.br.ReadByte()
		if err != nil {
			return 0, err
		}
		r.remain += 8
		r.bits <<= 8
		r.bits |= uint16(b)
	}
	r.remain -= 8
	r.unread = true // don't allow UnreadBit after ReadByte
	return byte(r.bits >> r.remain), nil
}

func (r *reader) UnreadBit() error {
	if r.unread {
		return ErrDoubleUnread
	}
	r.unread = true
	r.remain++
	return nil
}

// NewReader returns a BitScanner that reads from br. Calling UnreadBit without
// an immediately preceeding call to ReadBit is an error.
func NewReader(br io.ByteReader) BitScanner {
	return &reader{
		br:     br,
		unread: true,
	}
}

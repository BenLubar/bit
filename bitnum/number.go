package bitnum

type Number struct {
	Bits []uint64
	Size uint64
}

func (n *Number) Bit(i uint64) bool {
	if n.Size <= i {
		return false
	}

	return n.Bits[i>>6]&(1<<(i&63)) != 0
}

func (n *Number) Uint64() uint64 {
	if n.Size > 64 {
		panic("number too big: " + n.String())
	}
	return n.Uint64Offset(0)
}

func (n *Number) Uint64Offset(offset uint64) uint64 {
	if n.Size <= offset {
		return 0
	}
	words, bits := offset/64, offset%64
	var high, low uint64
	if words < uint64(len(n.Bits)) {
		low = n.Bits[offset/64]
	}
	if bits != 0 && words+1 < uint64(len(n.Bits)) {
		high = n.Bits[offset/64+1]
	}
	return high<<(64-bits) | low>>bits
}

func (n *Number) SetBit(i uint64, v bool) {
	if n.Size <= i {
		for uint64(len(n.Bits)) <= i>>6 {
			n.Bits = append(n.Bits, 0)
		}
		n.Size = i + 1
	}

	if v {
		n.Bits[i>>6] |= 1 << (i & 63)
	} else {
		n.Bits[i>>6] &^= 1 << (i & 63)
	}
}

func (n *Number) Append(v bool) {
	for uint64(len(n.Bits)) <= n.Size>>6 {
		n.Bits = append(n.Bits, 0)
	}
	n.Size++

	var carry uint64
	if v {
		carry = 1
	}

	for i := 0; i < len(n.Bits); i++ {
		next := n.Bits[i] >> 63
		n.Bits[i] <<= 1
		n.Bits[i] |= carry
		carry = next
	}
}

func (n *Number) String() string {
	var buf []byte

	for i := n.Size; i > 0; i-- {
		if n.Bit(i - 1) {
			buf = append(buf, " ONE"...)
		} else {
			buf = append(buf, " ZERO"...)
		}
	}

	return string(buf[1:])
}

func (n *Number) ShortString() string {
	var buf []byte

	for i := n.Size; i > 0; i-- {
		if n.Bit(i - 1) {
			buf = append(buf, 'o')
		} else {
			buf = append(buf, 'z')
		}
	}

	return string(buf)
}

func (n *Number) Equal(o *Number) bool {
	s := n.Size
	if so := o.Size; s < so {
		s = so
	}

	for i := s; i > 0; i-- {
		if n.Bit(i-1) != o.Bit(i-1) {
			return false
		}
	}

	return true
}

func (n *Number) Less(o *Number) bool {
	s := n.Size
	if so := o.Size; s < so {
		s = so
	}

	for i := s; i > 0; i-- {
		if n.Bit(i-1) && !o.Bit(i-1) {
			return false
		} else if !n.Bit(i-1) && o.Bit(i-1) {
			return true
		}
	}

	return false
}

package main

import "sort"

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
	return n.Bits[0]
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

func (n *Number) shortString() string {
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

type Numbers []*Number

func (s Numbers) Len() int           { return len(s) }
func (s Numbers) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Numbers) Less(i, j int) bool { return s[i].Less(s[j]) }

func (s Numbers) Search(n *Number) int {
	return sort.Search(len(s), func(i int) bool {
		return !s[i].Less(n)
	})
}

func (s Numbers) Contains(n *Number) bool {
	i := s.Search(n)
	return i < len(s) && s[i].Equal(n)
}

func (s *Numbers) Add(n *Number) bool {
	i := s.Search(n)
	if i < len(*s) && (*s)[i].Equal(n) {
		return false
	}

	*s = append(*s, nil)
	copy((*s)[i+1:], (*s)[i:])
	(*s)[i] = n

	return true
}

type NumberMap struct {
	Keys   Numbers
	Values []interface{}
}

func (m *NumberMap) Get(k *Number) (v interface{}, ok bool) {
	i := m.Keys.Search(k)
	if i < len(m.Keys) && m.Keys[i].Equal(k) {
		return m.Values[i], true
	}
	return nil, false
}

func (m *NumberMap) Set(k *Number, v interface{}) {
	i := m.Keys.Search(k)
	if i >= len(m.Keys) || !m.Keys[i].Equal(k) {
		m.Keys = append(m.Keys, nil)
		m.Values = append(m.Values, nil)
		copy(m.Keys[i+1:], m.Keys[i:])
		copy(m.Values[i+1:], m.Values[i:])
		m.Keys[i] = k
	}
	m.Values[i] = v
}

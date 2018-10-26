package main

import (
	"math/big"

	"github.com/BenLubar/bit/ast"
)

type intern struct {
	bits      []ast.Bits
	lookup    map[uint64]int
	lookupBig map[string]int
}

func (i *intern) intern(b ast.Bits) int {
	if len(b) == 0 {
		return 0
	}

	var z big.Int
	if b.Set(&z).IsUint64() {
		if x, ok := i.lookup[z.Uint64()]; ok {
			return x
		}

		if i.lookup == nil {
			i.lookup = make(map[uint64]int)
		}

		bc := make(ast.Bits, len(b))
		copy(bc, b)

		i.bits = append(i.bits, bc)
		i.lookup[z.Uint64()] = len(i.bits)
		return len(i.bits)
	}

	zb := string(z.Bytes())
	if x, ok := i.lookupBig[zb]; ok {
		return x
	}

	if i.lookupBig == nil {
		i.lookupBig = make(map[string]int)
	}

	bc := make(ast.Bits, len(b))
	copy(bc, b)

	i.bits = append(i.bits, bc)
	i.lookupBig[zb] = len(i.bits)
	return len(i.bits)
}

func (i *intern) find(b ast.Bits) (int, bool) {
	if len(b) == 0 {
		return 0, false
	}

	var z big.Int
	if b.Set(&z).IsUint64() {
		x, ok := i.lookup[z.Uint64()]
		return x, ok
	}

	zb := string(z.Bytes())
	x, ok := i.lookupBig[zb]
	return x, ok
}

func (i *intern) name(x int) string {
	if x == 0 {
		return ""
	}

	return i.bits[x-1].String()
}

func (i *intern) count() int {
	return len(i.bits)
}

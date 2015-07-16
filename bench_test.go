package bit

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/BenLubar/bit/bitio"
)

func BenchmarkParse(b *testing.B) {
	buf, err := ioutil.ReadFile("hello.bit")
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			prog, err := Parse(bytes.NewReader(buf))
			if err != nil {
				panic(err)
			}
			_ = prog
		}
	})
}

func BenchmarkHello(b *testing.B) {
	buf, err := ioutil.ReadFile("hello.bit")
	if err != nil {
		panic(err)
	}

	prog, err := Parse(bytes.NewReader(buf))
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := prog.Run(bitio.Null, bitio.Null)
			if err != nil {
				panic(err)
			}
		}
	})
}

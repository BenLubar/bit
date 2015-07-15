package bit

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
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
			err := prog.Run(&discard{}, &discard{})
			if err != nil {
				panic(err)
			}
		}
	})
}

type discard struct{}

func (*discard) ReadBit() (bool, error) { return false, io.EOF }
func (*discard) WriteBit(bool) error    { return nil }

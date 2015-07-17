package bit

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/BenLubar/bit/bitio"
)

func BenchmarkBIT_Parse(b *testing.B) {
	buf, err := ioutil.ReadFile("hello.bit")
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		prog, err := Parse(bytes.NewReader(buf))
		if err != nil {
			panic(err)
		}
		_ = prog
	}
}

func BenchmarkBIT_Optimize(b *testing.B) {
	buf, err := ioutil.ReadFile("hello.bit")
	if err != nil {
		panic(err)
	}

	prog, err := Parse(bytes.NewReader(buf))
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		prog.Optimize()
	}
}

func BenchmarkBIT_Hello(b *testing.B) {
	buf, err := ioutil.ReadFile("hello.bit")
	if err != nil {
		panic(err)
	}

	prog, err := Parse(bytes.NewReader(buf))
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := prog.Run(bitio.Null, bitio.Null)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkBIT_HelloOptimized(b *testing.B) {
	buf, err := ioutil.ReadFile("hello.bit")
	if err != nil {
		panic(err)
	}

	prog, err := Parse(bytes.NewReader(buf))
	if err != nil {
		panic(err)
	}

	prog.Optimize()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := prog.Run(bitio.Null, bitio.Null)
		if err != nil {
			panic(err)
		}
	}
}

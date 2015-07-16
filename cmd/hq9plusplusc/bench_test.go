package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/BenLubar/bit"
	"github.com/BenLubar/bit/bitio"
)

func BenchmarkWrite(b *testing.B) {
	hq9, err := ioutil.ReadFile("hello.hq9")
	if err != nil {
		panic(err)
	}

	s := string(hq9)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := NewWriter(ioutil.Discard)
			w.Program(s)
			err = w.Close()
			if err != nil {
				panic(err)
			}
		}
	})
}

func BenchmarkBITParse(b *testing.B) {
	hq9, err := ioutil.ReadFile("hello.hq9")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.Program(string(hq9))
	err = w.Close()
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			prog, err := bit.Parse(bytes.NewReader(buf.Bytes()))
			if err != nil {
				panic(err)
			}
			_ = prog
		}
	})
}

func BenchmarkHello(b *testing.B) {
	hq9, err := ioutil.ReadFile("hello.hq9")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.Program(string(hq9))
	err = w.Close()
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	prog, err := bit.Parse(&buf)
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

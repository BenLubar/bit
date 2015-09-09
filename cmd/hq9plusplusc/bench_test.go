package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/BenLubar/bit"
	"github.com/BenLubar/bit/bitio"
)

func BenchmarkHQ9_Write(b *testing.B) {
	hq9, err := ioutil.ReadFile("hello.hq9")
	if err != nil {
		panic(err)
	}

	s := string(hq9)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w := NewWriter(ioutil.Discard)
		_, err = w.Program(s)
		if err != nil {
			panic(err)
		}
		err = w.Close()
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkHQ9_Parse(b *testing.B) {
	hq9, err := ioutil.ReadFile("hello.hq9")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	w := NewWriter(&buf)
	_, err = w.Program(string(hq9))
	if err != nil {
		panic(err)
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		prog, err := bit.Parse(bytes.NewReader(buf.Bytes()))
		if err != nil {
			panic(err)
		}
		_ = prog
	}
}

func BenchmarkHQ9_Optimize(b *testing.B) {
	hq9, err := ioutil.ReadFile("hello.hq9")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	w := NewWriter(&buf)
	_, err = w.Program(string(hq9))
	if err != nil {
		panic(err)
	}
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

	for i := 0; i < b.N; i++ {
		prog.Optimize()
	}
}

func BenchmarkHQ9_Hello(b *testing.B) {
	hq9, err := ioutil.ReadFile("hello.hq9")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	w := NewWriter(&buf)
	_, err = w.Program(string(hq9))
	if err != nil {
		panic(err)
	}
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

	for i := 0; i < b.N; i++ {
		err := prog.Run(bitio.Null, bitio.Null)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkHQ9_HelloOptimized(b *testing.B) {
	hq9, err := ioutil.ReadFile("hello.hq9")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	w := NewWriter(&buf)
	_, err = w.Program(string(hq9))
	if err != nil {
		panic(err)
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	prog, err := bit.Parse(&buf)
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

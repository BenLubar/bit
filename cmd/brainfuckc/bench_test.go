package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/BenLubar/bit"
	"github.com/BenLubar/bit/bitio"
)

func BenchmarkBF_BFParse(b *testing.B) {
	bf, err := ioutil.ReadFile("hello.bf")
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		commands, err := Parse(Tokenize(bytes.NewReader(bf)))
		if err != nil {
			panic(err)
		}
		_ = commands
	}
}

func BenchmarkBF_Write(b *testing.B) {
	bf, err := ioutil.ReadFile("hello.bf")
	if err != nil {
		panic(err)
	}

	commands, err := Parse(Tokenize(bytes.NewReader(bf)))
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w := NewWriter(ioutil.Discard)
		w.Program(commands)
		err = w.Close()
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkBF_Parse(b *testing.B) {
	bf, err := ioutil.ReadFile("hello.bf")
	if err != nil {
		panic(err)
	}

	commands, err := Parse(Tokenize(bytes.NewReader(bf)))
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.Program(commands)
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

func BenchmarkBF_Optimize(b *testing.B) {
	bf, err := ioutil.ReadFile("hello.bf")
	if err != nil {
		panic(err)
	}

	commands, err := Parse(Tokenize(bytes.NewReader(bf)))
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.Program(commands)
	err = w.Close()
	if err != nil {
		panic(err)
	}

	prog, err := bit.Parse(&buf)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		prog.Optimize()
	}
}

func BenchmarkBF_Hello(b *testing.B) {
	bf, err := ioutil.ReadFile("hello.bf")
	if err != nil {
		panic(err)
	}

	commands, err := Parse(Tokenize(bytes.NewReader(bf)))
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.Program(commands)
	err = w.Close()
	if err != nil {
		panic(err)
	}

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

func BenchmarkBF_HelloOptimized(b *testing.B) {
	bf, err := ioutil.ReadFile("hello.bf")
	if err != nil {
		panic(err)
	}

	commands, err := Parse(Tokenize(bytes.NewReader(bf)))
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.Program(commands)
	err = w.Close()
	if err != nil {
		panic(err)
	}

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

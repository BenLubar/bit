package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/BenLubar/bit"
	"github.com/BenLubar/bit/bitio"
)

func BenchmarkBFParse(b *testing.B) {
	bf, err := ioutil.ReadFile("hello.bf")
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			commands, err := Parse(Tokenize(bytes.NewReader(bf)))
			if err != nil {
				panic(err)
			}
			_ = commands
		}
	})
}

func BenchmarkWrite(b *testing.B) {
	bf, err := ioutil.ReadFile("hello.bf")
	if err != nil {
		panic(err)
	}

	commands, err := Parse(Tokenize(bytes.NewReader(bf)))
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := NewWriter(ioutil.Discard)
			w.Program(commands)
			err = w.Close()
			if err != nil {
				panic(err)
			}
		}
	})
}

func BenchmarkBITParse(b *testing.B) {
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

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := prog.Run(bitio.Null, bitio.Null)
			if err != nil {
				panic(err)
			}
		}
	})
}

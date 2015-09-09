package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func BenchmarkOok_OokParse(b *testing.B) {
	ook, err := ioutil.ReadFile("hello.ook")
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		commands, err := Parse(bytes.NewReader(ook))
		if err != nil {
			panic(err)
		}
		_ = commands
	}
}

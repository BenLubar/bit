package token_test

import (
	"io/ioutil"
	"testing"

	"github.com/BenLubar/bit/token"
)

func BenchmarkHello(b *testing.B) {
	src, err := ioutil.ReadFile("../hello.bit")
	if err != nil {
		b.Fatalf("could not read hello.bit: %v", err)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var p token.Parser
			_, err := p.Write(src)
			if err == nil {
				err = p.Done()
			}
			if err != nil {
				b.Errorf("failed to tokenize hello.bit: %v", err)
			}
		}
	})
}

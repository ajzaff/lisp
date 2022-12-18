package hash

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/fuzzutil"
)

var (
	strDB     = make(map[string]struct{})
	valDB     = make(map[lisp.Val]struct{})
	maphashDB = make(map[uint64]struct{})
)

func BenchmarkMapHashMap(b *testing.B) {
	g := fuzzutil.NewGenerator(rand.New(rand.NewSource(1337)))
	g.ExprWeight = 0
	var h MapHash

	for i := 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			v := g.Next()
			h.Reset()
			h.WriteVal(v)
			maphashDB[h.Sum64()] = struct{}{}
		}
	}
}
func BenchmarkValMap(b *testing.B) {
	g := fuzzutil.NewGenerator(rand.New(rand.NewSource(1337)))
	g.ExprWeight = 0

	for i := 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			v := g.Next()
			valDB[v] = struct{}{}
		}
	}
}

func BenchmarkBaselineStringHash(b *testing.B) {
	g := fuzzutil.NewGenerator(rand.New(rand.NewSource(1337)))
	g.ExprWeight = 0

	for i := 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			v := g.Next()

			var buf bytes.Buffer
			lisp.StdPrinter(&buf).Print(v)
			strDB[buf.String()] = struct{}{}
		}
	}
}

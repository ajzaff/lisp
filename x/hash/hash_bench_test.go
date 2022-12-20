package hash

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/x/fuzzutil"
)

var res int

func BenchmarkMapHashMap(b *testing.B) {
	g := fuzzutil.NewGenerator(rand.New(rand.NewSource(1337)))
	var h MapHash
	maphashDB := make(map[uint64]struct{})

	i := 0
	for i = 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			v := g.Next()
			h.Reset()
			h.WriteVal(v)
			maphashDB[h.Sum64()] = struct{}{}
		}
	}
	res = i
}

func BenchmarkValMap(b *testing.B) {
	g := fuzzutil.NewGenerator(rand.New(rand.NewSource(1337)))
	valDB := make(map[lisp.Val]struct{})

	i := 0
	for i = 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			v := g.Next()
			valDB[v] = struct{}{}
		}
	}
	res = i
}

func BenchmarkBaselineStringHash(b *testing.B) {
	g := fuzzutil.NewGenerator(rand.New(rand.NewSource(1337)))
	strDB := make(map[string]struct{})

	i := 0
	for i = 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			v := g.Next()

			var buf bytes.Buffer
			lisp.StdPrinter(&buf).Print(v)
			strDB[buf.String()] = struct{}{}
		}
	}
	res = i
}

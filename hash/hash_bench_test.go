package hash

import (
	"bytes"
	"testing"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/lisputil"
)

var (
	strDB     = make(map[string]struct{})
	valDB     = make(map[lisp.Val]struct{})
	maphashDB = make(map[uint64]struct{})
)

func BenchmarkMapHashMap(b *testing.B) {
	var g lisputil.NodeGenerator
	var h MapHash

	for i := 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			v := g.NextValNoExpr()
			h.Reset()
			h.WriteVal(v)
			maphashDB[h.Sum64()] = struct{}{}
		}
	}
}
func BenchmarkValMap(b *testing.B) {
	var g lisputil.NodeGenerator

	for i := 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			v := g.NextValNoExpr()
			valDB[v] = struct{}{}
		}
	}
}

func BenchmarkBaselineStringHash(b *testing.B) {
	var g lisputil.NodeGenerator

	for i := 0; i < b.N; i++ {
		for j := 0; j < 256; j++ {
			v := g.NextValNoExpr()

			var buf bytes.Buffer
			lisp.StdPrinter(&buf).Print(v)
			strDB[buf.String()] = struct{}{}
		}
	}
}

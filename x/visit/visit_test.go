package visit

import (
	"math/rand"
	"testing"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/visit"
	"github.com/ajzaff/lisp/x/fuzzutil"
)

var res int

func BenchmarkVisitBaseline(b *testing.B) {
	g := fuzzutil.NewGenerator(rand.New(rand.NewSource(1337)))
	g.ConsWeight = 3
	g.ConsMaxDepth = 10

	var v visit.Visitor
	v.SetValVisitor(func(lisp.Val) {})
	v.SetLitVisitor(func(lisp.Lit) {})
	v.SetBeforeConsVisitor(func(*lisp.Cons) {})
	v.SetConsVisitor(func(*lisp.Cons) {})
	v.SetAfterConsVisitor(func(*lisp.Cons) {})

	var r int
	for i := 0; i < b.N; i++ {
		r++
		n := g.Next()

		v.Visit(n)
	}
	res = r
}

func BenchmarkVisitExperimental(b *testing.B) {
	g := fuzzutil.NewGenerator(rand.New(rand.NewSource(1337)))
	g.ConsWeight = 3
	g.ConsMaxDepth = 10

	visitFn := func(lisp.Val) {}

	var r int
	for i := 0; i < b.N; i++ {
		r++
		n := g.Next()

		Visit(n, visitFn)
	}
	res = r
}

func BenchmarkVisitExperimentalStack(b *testing.B) {
	g := fuzzutil.NewGenerator(rand.New(rand.NewSource(1337)))
	g.ConsWeight = 3
	g.ConsMaxDepth = 10

	queue := make([]lisp.Val, 0, 128)
	visitFn := func(lisp.Val) {}

	var r int
	for i := 0; i < b.N; i++ {
		r++
		n := g.Next()

		VisitStack(n, queue, visitFn)
	}
	res = r
}

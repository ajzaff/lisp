package visit

import (
	"math/rand"
	"testing"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/fuzzutil"
)

var res int

func BenchmarkVisitBaseline(b *testing.B) {
	g := fuzzutil.NewGenerator(rand.New(rand.NewSource(1337)))
	g.ExprWeight = 3
	g.ExprMaxDepth = 10

	var v lisp.Visitor
	v.SetBeforeValVisitor(func(lisp.Val) {})
	v.SetAfterValVisitor(func(lisp.Val) {})
	v.SetBeforeExprVisitor(func(lisp.Expr) {})
	v.SetBeforeExprVisitor(func(lisp.Expr) {})
	v.SetLitVisitor(func(lisp.Lit) {})

	var r int
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r++
		n := g.Next()
		b.StartTimer()

		v.Visit(n)
	}
	res = r
}

func BenchmarkVisitExperimental(b *testing.B) {
	g := fuzzutil.NewGenerator(rand.New(rand.NewSource(1337)))
	g.ExprWeight = 3
	g.ExprMaxDepth = 10

	visitFn := func(lisp.Val) {}

	var r int
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r++
		n := g.Next()
		b.StartTimer()

		Visit(n, visitFn)
	}
	res = r
}

func BenchmarkVisitExperimentalStack(b *testing.B) {
	g := fuzzutil.NewGenerator(rand.New(rand.NewSource(1337)))
	g.ExprWeight = 3
	g.ExprMaxDepth = 10

	queue := make([]lisp.Val, 0, 128)
	visitFn := func(lisp.Val) {}

	var r int
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r++
		n := g.Next()
		b.StartTimer()

		VisitStack(n, queue, visitFn)
	}
	res = r
}

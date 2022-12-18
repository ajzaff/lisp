package bench

import (
	"math/rand"
	"testing"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/fuzzutil"
	"github.com/ajzaff/lisp/x/visit"
)

var res int

func BenchmarkVisit(b *testing.B) {
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

func BenchmarkVisitQueue(b *testing.B) {
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

		visit.VisitStack(n, queue, visitFn)
	}
	res = r
}

func BenchmarkVisitQueueVisit(b *testing.B) {
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

		visit.Visit(n, visitFn)
	}
	res = r
}
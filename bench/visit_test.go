package bench

import (
	"math/rand"
	"testing"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/fuzzutil"
	"github.com/ajzaff/lisp/x/visitqueue"
)

var res int

func BenchmarkVisit(b *testing.B) {
	g := fuzzutil.NewGenerator(rand.New(rand.NewSource(1337)))
	g.ExprWeight = 3

	var v lisp.Visitor
	v.SetBeforeValVisitor(func(lisp.Val) {})
	v.SetAfterValVisitor(func(lisp.Val) {})
	v.SetBeforeExprVisitor(func(lisp.Expr) {})
	v.SetBeforeExprVisitor(func(lisp.Expr) {})
	v.SetLitVisitor(func(lisp.Lit) {})

	var r int
	for i := 0; i < b.N; i++ {
		n := g.Next()

		v.Visit(n.Val)
		r++
	}
	res = r
}

func BenchmarkVisitQueue(b *testing.B) {
	g := fuzzutil.NewGenerator(rand.New(rand.NewSource(1337)))
	g.ExprWeight = 3

	buf := make([]lisp.Val, 0, 16)

	var r int
	for i := 0; i < b.N; i++ {
		n := g.Next()

		visitqueue.VisitQueue(n.Val, buf, func(lisp.Val) {})
		r++
	}
	res = r
}

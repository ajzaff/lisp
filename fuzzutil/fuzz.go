package fuzzutil

import (
	"fmt"
	"math"
	"strconv"

	"github.com/ajzaff/lisp"
)

type Rand interface {
	Seed(int64)
	Intn(n int) int
	Uint64() uint64
	Float64() float64
	ExpFloat64() float64
}

type Generator struct {
	IdWeight     int
	IntWeight    int
	ExprWeight   int
	ExprMaxDepth int

	termFn func(depth int) int

	r Rand
}

func NewGenerator(r Rand) *Generator {
	g := &Generator{
		IdWeight:     1,
		IntWeight:    1,
		ExprWeight:   1,
		ExprMaxDepth: 3,
		r:            r,
	}
	g.termFn = g.expTermFn
	return g
}

func (g *Generator) expTermFn(depth int) int {
	return int(math.Max(float64(g.ExprMaxDepth-depth), 1) * g.r.ExpFloat64())
}

func (g *Generator) Seed(seed int64) {
	g.r.Seed(seed)
}

func (g *Generator) weight() int {
	return g.IdWeight + g.IntWeight + g.ExprWeight
}

func (g *Generator) Token() lisp.Token {
	// Shuffle order to make equal weights fair.
	// FIXME: can we do better? :)
	tok := [3]lisp.Token{lisp.Id, lisp.Int, lisp.LParen}
	w := [3]int{g.IdWeight, g.IntWeight, g.ExprWeight}
	i := g.r.Intn(3)
	tok[2], w[2], tok[i], w[i] = tok[i], w[i], tok[2], w[2]
	i = g.r.Intn(2)
	tok[1], w[1], tok[i], w[i] = tok[i], w[i], tok[1], w[1]

	v := g.r.Intn(g.weight())

	if w[0] != 0 && v <= w[0] {
		return tok[0]
	}
	v -= w[0]
	if w[1] != 0 && v <= w[1] {
		return tok[1]
	}
	return tok[2]
}

func (g *Generator) Next() lisp.Val {
	return g.nextDepth(0)
}

func (g *Generator) nextDepth(depth int) lisp.Val {
	switch g.Token() {
	case lisp.Id:
		return g.NextId()
	case lisp.Int:
		return g.NextInt()
	default: // Expr
		return g.nextExprDepth(depth)
	}
}

func (g *Generator) NextId() lisp.Val {
	return lisp.Lit{Token: lisp.Id, Text: fmt.Sprintf("a%d", g.r.Uint64())}
}

func (g *Generator) NextInt() lisp.Val {
	return lisp.Lit{Token: lisp.Int, Text: strconv.FormatUint(g.r.Uint64(), 10)}
}

func (g *Generator) NextExpr() lisp.Val {
	return g.nextExprDepth(0)
}

func (g *Generator) nextExprDepth(depth int) lisp.Val {
	n := g.termFn(depth)

	p := g.ExprMaxDepth - depth
	if p < 0 {
		p = 0
	}
	expr := make(lisp.Expr, 0, p)
	for i := 0; i < n; i++ {
		expr = append(expr, lisp.Node{Val: g.nextDepth(depth + 1)})
	}
	return expr
}

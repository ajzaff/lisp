package fuzzutil

import (
	"fmt"
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
	IdWeight   int
	IntWeight  int
	ExprWeight int

	termFn func() int

	r Rand
}

func NewGenerator(r Rand) *Generator {
	g := &Generator{
		IdWeight:   1,
		IntWeight:  1,
		ExprWeight: 1,
		r:          r,
	}
	g.termFn = g.expTermFn
	return g
}

func (g *Generator) expTermFn() int {
	return int(g.r.ExpFloat64())
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

func (g *Generator) Next() lisp.Node {
	switch g.Token() {
	case lisp.Id:
		return g.NextId()
	case lisp.Int:
		return g.NextInt()
	default: // Expr
		return g.NextExpr()
	}
}

func (g *Generator) NextId() lisp.Node {
	return lisp.Node{Val: lisp.Lit{Token: lisp.Id, Text: fmt.Sprintf("a%d", g.r.Uint64())}}
}

func (g *Generator) NextInt() lisp.Node {
	return lisp.Node{Val: lisp.Lit{Token: lisp.Int, Text: strconv.FormatUint(g.r.Uint64(), 10)}}
}

func (g *Generator) NextExpr() lisp.Node {
	var n int
	for n <= 0 {
		n = g.termFn()
	}
	var xs lisp.Expr
	for i := 0; i < n; i++ {
		xs = append(xs, g.Next())
	}
	return lisp.Node{Val: xs}
}

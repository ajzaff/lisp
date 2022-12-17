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
	IdWeight     int
	NumberWeight int
	ExprWeight   int

	termFn func() int

	r Rand
}

func NewGenerator(r Rand) *Generator {
	g := &Generator{
		IdWeight:     1,
		NumberWeight: 1,
		ExprWeight:   1,
		r:            r,
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
	return g.IdWeight + g.NumberWeight + g.ExprWeight
}

func (g *Generator) Token() lisp.Token {
	v := g.r.Intn(g.weight())
	if g.IdWeight != 0 && v <= 0 {
		return lisp.Id
	}
	v -= g.IdWeight
	if g.NumberWeight != 0 && v <= 0 {
		return lisp.Number
	}
	v -= g.NumberWeight
	if g.ExprWeight != 0 {
		return lisp.LParen // Expr
	}
	panic("Generator.Token: invalid weights resulted in no Token being emitted")
}

func (g *Generator) Next() lisp.Node {
	switch g.Token() {
	case lisp.Id:
		return g.NextId()
	case lisp.Number:
		return g.NextNumber()
	case lisp.LParen: // Expr
		return g.NextExpr()
	default:
		panic("Generator.Val: unreachable")
	}
}

func (g *Generator) NextId() lisp.Node {
	return lisp.Node{Val: lisp.Lit{Token: lisp.Id, Text: fmt.Sprintf("a%d", g.r.Uint64())}}
}

func (g *Generator) NextNumber() lisp.Node {
	if g.r.Intn(2) == 0 {
		return g.NextInt()
	} else { // 1
		return g.NextFloat()
	}
}

func (g *Generator) NextInt() lisp.Node {
	return lisp.Node{Val: lisp.Lit{Token: lisp.Number, Text: strconv.FormatInt(int64(g.r.Uint64()), 10)}}
}

func (g *Generator) NextFloat() lisp.Node {
	return lisp.Node{Val: lisp.Lit{Token: lisp.Number, Text: strconv.FormatFloat(g.r.Float64(), 'f', -1, 64)}}
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

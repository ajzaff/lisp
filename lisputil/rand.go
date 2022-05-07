package lisputil

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/ajzaff/lisp"
)

type NodeGenerator struct {
	r    *rand.Rand
	once sync.Once
}

func (g *NodeGenerator) init() {
	g.r = rand.New(rand.NewSource(1337))
}

func (g *NodeGenerator) SetSeed(seed int64) {
	g.once.Do(g.init)
	g.r.Seed(seed)
}

func (g *NodeGenerator) NextNodeNoExpr() lisp.Node {
	g.once.Do(g.init)
	return g.nextNode(true)
}

func (g *NodeGenerator) NextNode() lisp.Node {
	g.once.Do(g.init)
	return g.nextNode(false)
}

func (g *NodeGenerator) nextNode(noExpr bool) lisp.Node {
	n := 5
	if noExpr {
		n = 4
	}
	var l lisp.Lit
	switch g.r.Intn(n) {
	case 0: // Id
		l = g.nextId()
	case 1: // Int
		l = g.nextInt()
	case 2: // Float
		l = g.nextFloat()
	case 3: // String
		l = g.nextString()
	case 4: // Expr
		return &lisp.ExprNode{Expr: g.nextExpr()}
	default:
		panic("unreachable")
	}
	return &lisp.LitNode{Lit: l}
}

func (g *NodeGenerator) NextId() lisp.IdLit {
	g.once.Do(g.init)
	return g.nextId()
}

func (g *NodeGenerator) nextId() lisp.IdLit {
	return lisp.IdLit(fmt.Sprintf("a%d", g.r.Int63()))
}

func (g *NodeGenerator) NextString() lisp.Lit {
	g.once.Do(g.init)
	return g.nextString()
}

func (g *NodeGenerator) nextString() lisp.Lit {
	return lisp.StringLit(fmt.Sprintf("a%d", g.r.Int63()))
}

func (g *NodeGenerator) NextInt() lisp.Lit {
	g.once.Do(g.init)
	return g.nextInt()
}

func (g *NodeGenerator) nextInt() lisp.Lit {
	return lisp.IntLit(uint64(g.r.Int63()))
}

func (g *NodeGenerator) NextFloat() lisp.Lit {
	g.once.Do(g.init)
	return g.nextFloat()
}

func (g *NodeGenerator) nextFloat() lisp.Lit {
	return lisp.FloatLit(g.r.Float64())
}

func (g *NodeGenerator) NextExpr() lisp.Expr {
	g.once.Do(g.init)
	return g.nextExpr()
}

func (g *NodeGenerator) nextExpr() lisp.Expr {
	var n int
	for n <= 0 {
		n = int(2*g.r.NormFloat64() + 5)
	}
	var xs []lisp.Node
	for i := 0; i < n; i++ {
		xs = append(xs, g.nextNode(false))
	}
	return lisp.Expr(xs)
}

func (g *NodeGenerator) NextValNoExpr() lisp.Val {
	g.once.Do(g.init)
	return g.nextVal(true)
}

func (g *NodeGenerator) NextVal() lisp.Val {
	g.once.Do(g.init)
	return g.nextVal(false)
}

func (g *NodeGenerator) nextVal(noExpr bool) lisp.Val {
	return g.nextNode(noExpr).Val()
}

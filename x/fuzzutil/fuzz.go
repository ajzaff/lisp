package fuzzutil

import (
	"math"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ajzaff/lisp"
	"golang.org/x/text/unicode/rangetable"
)

type Rand interface {
	Seed(int64)
	Intn(n int) int
	Uint64() uint64
	Float64() float64
	ExpFloat64() float64
}

type Generator struct {
	IdWeight      int
	IntWeight     int
	GroupWeight   int
	GroupMeanLen  int
	GroupMaxDepth int

	// Generates number of links in Group.
	termFn func() int

	// Generates lengths of Ids.
	idFn func() int

	r Rand
}

func NewGenerator(r Rand) *Generator {
	g := &Generator{
		IdWeight:      1,
		IntWeight:     1,
		GroupWeight:   1,
		GroupMeanLen:  3,
		GroupMaxDepth: 3,
		r:             r,
	}
	g.termFn = g.expTermFn
	g.idFn = g.expIdFn
	return g
}

func (g *Generator) expTermFn() int {
	return int(float64(g.GroupMeanLen) * g.r.ExpFloat64())
}

func (g *Generator) expIdFn() int {
	// Generate the approx. length of an ID in bytes (most runes are 4 bytes long).
	return int(math.Ceil(40 * g.r.ExpFloat64()))
}

func (g *Generator) Seed(seed int64) {
	g.r.Seed(seed)
}

func (g *Generator) weight() int {
	return g.IdWeight + g.IntWeight + g.GroupWeight
}

func (g *Generator) Token() lisp.Token {
	return g.tokenDepth(0)
}

func (g *Generator) tokenDepth(depth int) lisp.Token {
	tok := []lisp.Token{lisp.Id, lisp.Nat, lisp.LParen}
	w := []int{g.IdWeight, g.IntWeight, g.GroupWeight}
	weightMax := g.weight()
	if g.GroupMaxDepth <= depth {
		// No more Group.
		tok = tok[:2]
		w = w[:2]
		weightMax -= g.GroupWeight
	} else {
		// Swap once.
		i := g.r.Intn(3)
		tok[2], w[2], tok[i], w[i] = tok[i], w[i], tok[2], w[2]
	}
	// Swap again.
	i := g.r.Intn(2)
	tok[1], w[1], tok[i], w[i] = tok[i], w[i], tok[1], w[1]
	// tok, w are shuffled.

	// Use weighted selection.
	v := g.r.Intn(weightMax)
	if w[0] != 0 && v <= w[0] {
		return tok[0]
	}
	if len(tok) == 2 {
		// When no Group we can return early.
		return tok[1]
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
	switch g.tokenDepth(depth) {
	case lisp.Id:
		return g.NextId()
	case lisp.Nat:
		return g.NextInt()
	default: // Group
		return g.nextGroupDepth(depth)
	}
}

var idTab = make([]rune, 0, 131241) // len(unicode.L)

func init() {
	rangetable.Visit(unicode.L, func(r rune) { idTab = append(idTab, r) })
}

func (g *Generator) NextId() lisp.Val {
	n := g.expIdFn()
	var sb strings.Builder
	sb.Grow(n)
	for i := 0; i < n; i++ {
		r := idTab[g.r.Intn(len(idTab))]
		size := utf8.RuneLen(r)
		sb.WriteRune(r)
		i += size
	}
	return lisp.Lit{Token: lisp.Id, Text: sb.String()}
}

func (g *Generator) NextInt() lisp.Val {
	return lisp.Lit{Token: lisp.Nat, Text: strconv.FormatUint(g.r.Uint64(), 10)}
}

func (g *Generator) NextGroup() lisp.Val {
	return g.nextGroupDepth(0)
}

func (g *Generator) nextGroupDepth(depth int) lisp.Val {
	n := g.termFn()
	group := lisp.Group{}
	for i := 0; i < n; i++ {
		group = append(group, g.nextDepth(depth+1))
	}
	return group
}

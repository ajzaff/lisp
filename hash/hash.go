package hash

import (
	"hash/maphash"

	"github.com/ajzaff/innit"
)

// Hash the node using the given maphash.
// Must call h.Reset() after each call.
func Hash(h *maphash.Hash, n innit.Node) uint64 {
	var v innit.Visitor
	first := true
	v.SetBeforeExprVisitor(func(*innit.Expr) {
		first = true
		h.WriteByte('(')
	})
	v.SetAfterExprVisitor(func(*innit.Expr) { h.WriteByte(')') })
	v.SetLitVisitor(func(e *innit.Lit) {
		if !first {
			h.WriteByte(' ')
		}
		first = false
		h.WriteString(e.Value)
	})
	v.Visit(n)
	return h.Sum64()
}

// Expr hashes e with the given maphash.
// Can be more efficient when the type of node is known.
func Expr(h *maphash.Hash, e *innit.Expr) uint64 {
	h.WriteByte('(')
	Hash(h, e.X)
	h.WriteByte(')')
	return h.Sum64()
}

// Lit hashes e with the given maphash.
// Can be more efficient when the type of node is known.
func Lit(h *maphash.Hash, e *innit.Lit) uint64 {
	h.WriteString(e.Value)
	return h.Sum64()
}

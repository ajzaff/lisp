package hash

import (
	"hash/maphash"

	"github.com/ajzaff/innit"
)

// Node hashes the node using the given maphash.
func Node(h *maphash.Hash, n innit.Node) {
	var v innit.Visitor
	first := true
	v.SetBeforeExprVisitor(func(*innit.Expr) { first = true; h.WriteByte('(') })
	v.SetAfterExprVisitor(func(*innit.Expr) { first = true; h.WriteByte(')') })
	v.SetLitVisitor(func(e *innit.Lit) {
		switch {
		case first:
			first = false
			h.WriteString(e.Value)
		default:
			h.WriteByte(' ')
			h.WriteString(e.Value)
		}
	})
	v.Visit(n)
}

// Expr hashes e with the given maphash.
// Can be more efficient when the type of node is known.
func Expr(h *maphash.Hash, e *innit.Expr) {
	h.WriteByte('(')
	Node(h, e.X)
	h.WriteByte(')')
}

// Lit hashes e with the given maphash.
// Can be more efficient when the type of node is known.
func Lit(h *maphash.Hash, e *innit.Lit) {
	h.WriteString(e.Value)
}

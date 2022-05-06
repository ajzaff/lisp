package hash

import (
	"hash/maphash"

	"github.com/ajzaff/lisp"
)

// Val hashes the Val using the given maphash.
func Val(h *maphash.Hash, v lisp.Val) {
	switch v := v.(type) {
	case lisp.Expr:
		Expr(h, v)
	case lisp.Lit:
		Lit(h, v)
	}
}

// Expr hashes e with the given maphash.
func Expr(h *maphash.Hash, es lisp.Expr) {
	h.WriteByte('(')
	var lastLit bool
	for _, e := range es {
		if lit, ok := e.(lisp.Lit); ok {
			if lastLit {
				h.WriteByte(' ')
			}
			lastLit = true
			Lit(h, lit)
			continue
		}
		lastLit = false
		Val(h, e.Val())
	}
	h.WriteByte(')')
}

// Lit hashes e with the given maphash.
func Lit(h *maphash.Hash, e lisp.Lit) { h.WriteString(e.String()) }

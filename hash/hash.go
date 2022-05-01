package hash

import (
	"hash/maphash"

	"github.com/ajzaff/innit"
)

// Val hashes the Val using the given maphash.
func Val(h *maphash.Hash, v innit.Val) {
	switch v := v.(type) {
	case innit.Expr:
		Expr(h, v)
	case innit.Lit:
		Lit(h, v)
	}
}

// Expr hashes e with the given maphash.
func Expr(h *maphash.Hash, es innit.Expr) {
	h.WriteByte('(')
	var lastLit bool
	for _, e := range es {
		if lit, ok := e.(innit.Lit); ok {
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
func Lit(h *maphash.Hash, e innit.Lit) { h.WriteString(e.String()) }

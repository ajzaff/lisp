package lisputil

import (
	"strings"

	"github.com/ajzaff/lisp"
)

// FIXME
func Compare(a, b lisp.Val) int {
	switch a := a.(type) {
	case lisp.Expr:
		switch b := b.(type) {
		case lisp.Expr:
			if len(a) != len(b) {
				if len(a) < len(b) {
					return -1
				}
				return 1
			}
			for i := range a {
				if v := Compare(a[i].Val, b[i].Val); v != 0 {
					return v
				}
			}
			return 0
		default: // Lit, Expr.
			return 1
		}
	case lisp.Lit:
		switch b := b.(type) {
		case lisp.Lit:
			return strings.Compare(a.String(), b.String())
		default: // Expr
			return -1
		}
	default:
		panic("compare not supported")
	}
}

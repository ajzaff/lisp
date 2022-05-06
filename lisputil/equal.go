package lisputil

import (
	"github.com/ajzaff/lisp"
)

// FIXME
func Equal(a, b lisp.Val) bool {
	switch a := a.(type) {
	case lisp.Expr:
		switch b := b.(type) {
		case lisp.Expr:
			if len(a) != len(b) {
				return false
			}
			for i := range a {
				if !Equal(a[i].Val(), b[i].Val()) {
					return false
				}
			}
			return true
		default:
			return false
		}
	case lisp.Lit:
		switch b := b.(type) {
		case lisp.Lit:
			return a.String() == b.String()
		default:
			return false
		}
	default: // other
		panic("equal not supported")
	}
}

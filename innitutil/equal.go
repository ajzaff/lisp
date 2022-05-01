package innitutil

import (
	"github.com/ajzaff/innit"
)

// FIXME
func Equal(a, b innit.Val) bool {
	switch a := a.(type) {
	case innit.Expr:
		switch b := b.(type) {
		case innit.Expr:
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
	case innit.Lit:
		switch b := b.(type) {
		case innit.Lit:
			return a.String() == b.String()
		default:
			return false
		}
	default: // other
		panic("equal not supported")
	}
}

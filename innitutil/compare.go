package innitutil

import (
	"strings"

	"github.com/ajzaff/innit"
)

// FIXME
func Compare(a, b innit.Val) int {
	switch a := a.(type) {
	case innit.Expr:
		switch b := b.(type) {
		case innit.Expr:
			if len(a) != len(b) {
				if len(a) < len(b) {
					return -1
				}
				return 1
			}
			for i := range a {
				if v := Compare(a[i].Val(), b[i].Val()); v != 0 {
					return v
				}
			}
			return 0
		default: // Lit, Expr.
			return 1
		}
	case innit.Lit:
		switch b := b.(type) {
		case innit.Lit:
			return strings.Compare(a.String(), b.String())
		default: // Expr
			return -1
		}
	default:
		panic("compare not supported")
	}
}

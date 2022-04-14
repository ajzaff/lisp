package innitutil

import (
	"github.com/ajzaff/innit"
)

func Equal(a, b innit.Node) bool {
	switch aa := a.(type) {
	case *innit.Expr:
		switch bb := b.(type) {
		case *innit.Expr:
			if len(aa.X) != len(bb.X) {
				return false
			}
			for i := range aa.X {
				if !Equal(aa.X[i], bb.X[i]) {
					return false
				}
			}
			return true
		default:
			return false
		}
	case *innit.Lit:
		switch bb := b.(type) {
		case *innit.Lit:
			return aa.Tok == bb.Tok && aa.Value == bb.Value
		default:
			return false
		}
	case innit.NodeList:
		switch bb := b.(type) {
		case innit.NodeList:
			if len(aa) != len(bb) {
				return false
			}
			for i := range aa {
				if !Equal(aa[i], bb[i]) {
					return false
				}
			}
			return true
		default:
			return false
		}
	default: // other
		switch b.(type) {
		case *innit.Lit, innit.NodeList, *innit.Expr:
			return false
		default: // other
			return true
		}
	}
}

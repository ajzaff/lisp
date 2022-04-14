package innitutil

import (
	"strings"

	"github.com/ajzaff/innit"
)

func Compare(a, b innit.Node) int {
	switch aa := a.(type) {
	case *innit.Expr:
		switch bb := b.(type) {
		case *innit.Expr:
			if len(aa.X) != len(bb.X) {
				if len(aa.X) < len(bb.X) {
					return -1
				}
				return 1
			}
			for i := range aa.X {
				if v := Compare(aa.X[i], bb.X[i]); v != 0 {
					return v
				}
			}
			return 0
		default: // *Lit, *NodeList, other.
			return 1
		}
	case *innit.Lit:
		switch bb := b.(type) {
		case *innit.Lit:
			if aa.Tok != bb.Tok {
				if aa.Tok < bb.Tok {
					return -1
				}
				return 1
			}
			return strings.Compare(aa.Value, bb.Value)
		default: // *Lit, NodeList, other, *Expr
			return -1
		}
	case innit.NodeList:
		switch bb := b.(type) {
		case innit.NodeList:
			if len(aa) != len(bb) {
				if len(aa) < len(bb) {
					return -1
				}
				return 1
			}
			for i := range aa {
				if v := Compare(aa[i], bb[i]); v != 0 {
					return v
				}
			}
			return 0
		case *innit.Lit:
			return 1
		default: // *Lit, other, *Expr
			return -1
		}
	default: // other
		switch b.(type) {
		case *innit.Expr:
			return -1
		case *innit.Lit, innit.NodeList:
			return 1
		default: // other
			return 0
		}
	}
}

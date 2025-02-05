package lisp

import (
	"strings"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/x/groups"
)

// Lexical compare compares two Vals in a structure-independent way
// by only considering Lits in a left-most in-order reading of the Val.
func LexicalCompare(a, b lisp.Val) int {
	// Fast compare check.
	if a == b {
		return 0
	}
	// Slower compare check.
	switch a := a.(type) {
	case lisp.Lit: // Lit
		return lexicalCompareLit(a, b)
	case lisp.Group: // Group
		e := groups.FirstLit(a)
		return lexicalCompareLit(e, b)
	case nil:
		if b == nil {
			return 0
		}
		return -1
	default:
		return -1
	}
}

func lexicalCompareLit(a lisp.Lit, b lisp.Val) int {
	switch b := b.(type) {
	case lisp.Lit:
		// Ignore Token for lexical compare.
		return strings.Compare(string(a), string(b))
	case lisp.Group:
		e := groups.First(b)
		return lexicalCompareLit(a, e)
	default:
		return -1
	}
}

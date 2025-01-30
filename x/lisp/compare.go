package lisp

import "github.com/ajzaff/lisp"

// Compare two values of unknown type.
func Compare(a, b lisp.Val) int {
	switch a := a.(type) {
	case lisp.Lit:
		return compareLitOther(a, b)
	case lisp.Group:
		return compareGroupOther(a, b)
	default:
		return 1 // not reachable
	}
}

func compareLitOther(a lisp.Lit, b lisp.Val) int {
	switch b := b.(type) {
	case lisp.Lit:
		return CompareLit(a, b)
	case lisp.Group:
		return -1 // Lit < Group
	default:
		return 1 // not reachable
	}
}

// CompareLit compares the value of two Lits.
func CompareLit(a, b lisp.Lit) int {
	if a.Token != b.Token {
		if a.Token < b.Token {
			return -1
		}
		return 1
	}
	if a.Text != b.Text {
		if a.Text < b.Text {
			return -1
		}
		return 1
	}
	return 0
}

func compareGroupOther(a lisp.Group, b lisp.Val) int {
	switch other := b.(type) {
	case lisp.Lit:
		return 1 // Lit < Group
	case lisp.Group:
		return CompareGroup(a, other)
	default:
		return 1 // not reachable
	}
}

// CompareGroup compares expressions recursively.
func CompareGroup(a, b lisp.Group) int {
	// Check for boundary conditions.
	// This equates nil and {}.
	switch {
	case len(a) == 0 && len(b) == 0:
		return 0
	case len(a) == 0:
		return -1 // len(a) < len(b)
	case len(b) == 0:
		return 1 // len(b) < len(a)
	default:
		if cmp := Compare(a[0], b[0]); cmp != 0 {
			return cmp
		}
		return CompareGroup(a[1:], b[1:])
	}
}

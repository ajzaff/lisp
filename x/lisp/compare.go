package lisp

import "github.com/ajzaff/lisp"

// Compare two values of unknown type.
func Compare(a, b lisp.Val) int {
	switch a := a.(type) {
	case lisp.Lit:
		return compareLitOther(a, b)
	case *lisp.Cons:
		return compareConsOther(a, b)
	default:
		return 1 // not reachable
	}
}

func compareLitOther(a lisp.Lit, b lisp.Val) int {
	switch b := b.(type) {
	case lisp.Lit:
		return CompareLit(a, b)
	case *lisp.Cons:
		return -1 // Lit < Cons
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

func compareConsOther(a *lisp.Cons, b lisp.Val) int {
	switch other := b.(type) {
	case lisp.Lit:
		return 1 // Lit < Cons
	case *lisp.Cons:
		return CompareCons(a, other)
	default:
		return 1 // not reachable
	}
}

// CompareCons compares expressions recursively.
func CompareCons(a, b *lisp.Cons) int {
	// Check for boundary conditions.
	switch {
	case a == nil && b == nil:
		return 0
	case a == nil:
		return -1 // len(a) < len(b)
	case b == nil:
		return 1 // len(b) < len(a)
	}
	// a != nil && b != nil:
	// Compare Vals.
	if w := Compare(a.Val, b.Val); w != 0 {
		return w
	}
	// Compare Cons recursively.
	return CompareCons(a.Cons, b.Cons)
}

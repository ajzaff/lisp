package lisputil

import "github.com/ajzaff/lisp"

// Equal returns whether two values are syntactically equivalent.
func Equal(a, b lisp.Val) bool {
	switch first := a.(type) {
	case lisp.Lit:
		return equalLitOther(first, b)
	case *lisp.Cons:
		return equalConsOther(first, b)
	default:
		return false // not reachable
	}
}

func equalLitOther(a lisp.Lit, b lisp.Val) bool {
	switch second := b.(type) {
	case lisp.Lit:
		return EqualLit(a, second)
	default:
		return false
	}
}

// EqualLit returns whether a and b are syntactically equivalent.
func EqualLit(a, b lisp.Lit) bool {
	return a.Token == b.Token && a.Text == b.Text
}

func equalConsOther(a *lisp.Cons, b lisp.Val) bool {
	switch second := b.(type) {
	case *lisp.Cons:
		return EqualCons(a, second)
	default:
		return false
	}
}

// EqualCons returns whether two expressions are syntactically equivalent by equating elements recursively.
func EqualCons(a, b *lisp.Cons) bool {
	// Check for boundary conditions.
	if a == nil {
		return b == nil
	}
	// a != nil.
	if b == nil {
		return false
	}
	// a != nil && b != nil:
	// Compare Vals.
	if !Equal(a.Val, b.Val) {
		return false
	}
	// Equate Cons recursively.
	return EqualCons(a.Cons, b.Cons)
}

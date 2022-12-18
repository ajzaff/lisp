package lisputil

import (
	"github.com/ajzaff/lisp"
)

// Equal returns whether two values are syntactically equivalent.
func Equal(a, b lisp.Val) bool {
	switch first := a.(type) {
	case lisp.Lit:
		return equalLitOther(first, b)
	case lisp.Expr:
		return equalExprOther(first, b)
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

func equalExprOther(a lisp.Expr, b lisp.Val) bool {
	switch second := b.(type) {
	case lisp.Expr:
		return EqualExpr(a, second)
	default:
		return false
	}
}

// EqualExpr returns whether two expressions are syntactically equivalent by equating elements recursively.
func EqualExpr(a, b lisp.Expr) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !Equal(a[i].Val, b[i].Val) {
			return false
		}
	}
	return true
}

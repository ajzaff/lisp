package lisputil

import (
	"github.com/ajzaff/lisp"
)

// Compare two values of unknown type.
func Compare(a, b lisp.Val) int {
	switch first := a.(type) {
	case lisp.Lit:
		return compareLitOther(first, b)
	case lisp.Expr:
		return compareExprOther(first, b)
	default:
		return 1 // not reachable
	}
}

func compareLitOther(a lisp.Lit, b lisp.Val) int {
	switch other := b.(type) {
	case lisp.Lit:
		return CompareLit(a, other)
	case lisp.Expr:
		return -1 // Lit < Expr
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

func compareExprOther(a lisp.Expr, b lisp.Val) int {
	switch other := b.(type) {
	case lisp.Lit:
		return 1 // Lit < Expr
	case lisp.Expr:
		return CompareExpr(a, other)
	default:
		return 1 // not reachable
	}
}

// CompareExpr compares expressions recursively.
func CompareExpr(a, b lisp.Expr) int {
	if len(a) != len(b) {
		if len(a) < len(b) {
			return -1
		}
		return 1
	}
	for i := range a {
		r := Compare(a[i].Val, b[i].Val)
		if r != 0 {
			return r
		}
	}
	return 0
}

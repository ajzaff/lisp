package lisp

import "github.com/ajzaff/lisp"

// Equal returns whether two values are syntactically equivalent.
func Equal(a, b lisp.Val) bool {
	switch first := a.(type) {
	case lisp.Lit:
		return equalLitOther(first, b)
	case lisp.Group:
		return equalGroupOther(first, b)
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

func equalGroupOther(a lisp.Group, b lisp.Val) bool {
	switch second := b.(type) {
	case lisp.Group:
		return EqualGroup(a, second)
	default:
		return false
	}
}

// EqualGroup returns whether two expressions are syntactically equivalent by equating elements recursively.
func EqualGroup(a, b lisp.Group) bool {
	switch {
	case len(a) == len(b):
		for i := range a {
			if !Equal(a[i], b[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

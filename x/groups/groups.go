package groups

import "github.com/ajzaff/lisp"

// First returns the first value in the group or nil.
func First(group lisp.Group) lisp.Val {
	if len(group) == 0 {
		return nil
	}
	return group[0]
}

// First returns the left-most in-order Lit in the group.
//
// A group with no Lits will return a zero value which is not a valid Lit.
func FirstLit(group lisp.Group) lisp.Lit {
	for _, e := range group {
		switch e := e.(type) {
		case lisp.Lit:
			return e
		case lisp.Group:
			if v := FirstLit(e); v != "" {
				return v
			}
		}
	}
	return ""
}

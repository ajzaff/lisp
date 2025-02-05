package lisp

import "github.com/ajzaff/lisp"

// Len returns the length of the Group.
//
// Len is a O(1)-time operation.
func Len(x lisp.Group) (n int) { return len(x) }

// Head returns the first Val in an group or nil.
func Head(x lisp.Group) lisp.Val {
	if len(x) > 0 {
		return x[0]
	}
	return nil
}

// Tail returns the first linked group or nil.
func Tail(x lisp.Group) lisp.Group {
	if len(x) > 0 {
		return x[1:]
	}
	return nil
}

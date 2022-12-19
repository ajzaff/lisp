package lisputil

import "github.com/ajzaff/lisp"

// Cons constructs a cons by linking the nodes together.
func Cons(ns ...lisp.Node) *lisp.Cons {
	root := &lisp.Cons{}
	e := root
	for i, n := range ns {
		e.Node = n
		if i+1 < len(ns) {
			e.Cons = &lisp.Cons{}
			e = e.Cons
		}
	}
	return root
}

// Len returns the length of the Cons.
//
// Len is a linear-time operation.
func Len(x *lisp.Cons) (n int) {
	for x := x; x != nil; n, x = n+1, x.Cons {
	}
	return
}

// Head returns the first Val in an cons or nil.
func Head(x *lisp.Cons) lisp.Val {
	if x != nil {
		return x.Val
	}
	return nil
}

// Tail returns the first linked cons or nil.
func Tail(x *lisp.Cons) *lisp.Cons {
	if x != nil {
		return x.Cons
	}
	return nil
}

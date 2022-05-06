package lisputil

import "github.com/ajzaff/lisp"

// Head returns the first element in an expression list or nil.
func Head(v lisp.Val) lisp.Val {
	if x, ok := v.(lisp.Expr); ok {
		return x[0].Val()
	}
	return nil
}

// Tail returns all but the first element of the expression list or nil.
func Tail(v lisp.Val) lisp.Val {
	if x, ok := v.(lisp.Expr); ok && len(x) > 0 {
		return x[1:]
	}
	return nil
}

package visit

import "github.com/ajzaff/lisp"

// Visit visits elements of root in order.
func Visit(root lisp.Val, visitFn func(v lisp.Val)) {
	switch root := root.(type) {
	case *lisp.Cons:
		VisitCons(root, visitFn)
	default:
		visitFn(root)
	}
}

// VisitCons visits the elements of the Cons recursively.
func VisitCons(root *lisp.Cons, visitFn func(v lisp.Val)) {
	if root != nil {
		Visit(root.Val, visitFn)
		VisitCons(root.Cons, visitFn)
	}
}

// VisitStack visits the elements of root in order.
//
// It is a bit less featureful than the full Visitor but may be better suited to larger cons.
func VisitStack(root lisp.Val, stack []lisp.Val, visitFn func(v lisp.Val)) {
	if root == nil {
		return
	}
	stack = stack[:0]
	stack = append(stack, root)
	for len(stack) > 0 {
		n := len(stack) - 1
		x := stack[n]
		stack = stack[:n]
		visitFn(x)
		if e, ok := x.(*lisp.Cons); ok {
			// Use closure to trigger defers.
			func() {
				for e := e; e != nil; e = e.Cons {
					// Use defer to reverse the linked list.
					// FIXME: we should benchmark how a non-defer solution fares here.
					defer func(e *lisp.Cons) { stack = append(stack, e.Val) }(e)
				}
			}()
		}
	}
}

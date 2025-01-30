package visit

import "github.com/ajzaff/lisp"

// Visit visits elements of root in order.
//
// It is a bit less featureful than the full Visitor but may be better suited to optimized workloads.
func Visit(root lisp.Val, visitFn func(v lisp.Val)) {
	switch root := root.(type) {
	case lisp.Group:
		VisitGroup(root, visitFn)
	default:
		visitFn(root)
	}
}

// VisitGroup visits the elements of the Group recursively.
func VisitGroup(root lisp.Group, visitFn func(v lisp.Val)) {
	for _, e := range root {
		Visit(e, visitFn)
	}
}

func VisitStack(root lisp.Val, stack []lisp.Val, visitFn func(v lisp.Val)) {
	if root == nil {
		return
	}
	stack = stack[:0:cap(stack)]
	stack = append(stack, root)
	for len(stack) > 0 {
		n := len(stack) - 1
		x := stack[n]
		stack = stack[:n]
		visitFn(x)
		if e, ok := x.(lisp.Group); ok {
			for _, x := range e {
				stack = append(stack, x)
			}
		}
	}
}

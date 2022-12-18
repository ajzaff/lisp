package visit

import "github.com/ajzaff/lisp"

// Visit visits the elements of root in order.
func Visit(root lisp.Val, visitFn func(v lisp.Val)) {
	if root == nil {
		return
	}
	visitFn(root)
	if expr, ok := root.(lisp.Expr); ok {
		for _, e := range expr {
			Visit(e.Val, visitFn)
		}
	}
}

// VisitStack visits the elements of root in order.
//
// It is a bit less featureful than the full Visitor but may be better suited to larger expressions.
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
		if expr, ok := x.(lisp.Expr); ok {
			for i := len(expr) - 1; i >= 0; i-- {
				stack = append(stack, expr[i].Val)
			}
		}
	}
}

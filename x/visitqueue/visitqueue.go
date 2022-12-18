package visitqueue

import "github.com/ajzaff/lisp"

// VisitQueue visits the elements of root in FIFO order.
//
// It is a bit less featureful than the full Visitor but may be better suited to large expressions.
func VisitQueue(root lisp.Val, queue []lisp.Val, visitFn func(v lisp.Val)) {
	if root == nil {
		return
	}
	queue = queue[:0]
	queue = append(queue, root)
	for len(queue) > 0 {
		x := queue[0]
		queue = queue[1:]
		visitFn(x)
		if expr, ok := x.(lisp.Expr); ok {
			for _, e := range expr {
				queue = append(queue, e.Val)
			}
		}
	}
}

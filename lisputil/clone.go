package lisputil

import "github.com/ajzaff/lisp"

func Clone(n lisp.Node) lisp.Node {
	if n == nil {
		return nil
	}
	var out lisp.Node
	switch x := n.(type) {
	case *lisp.LitNode:
		out = CloneLitNode(x)
	case *lisp.ExprNode:
		out = CloneExprNode(x)
	default:
		panic("clone not supported")
	}
	return out
}

func CloneExprNode(e *lisp.ExprNode) *lisp.ExprNode {
	if e == nil {
		return nil
	}
	out := new(lisp.ExprNode)
	*out = *e
	out.Expr = CloneExpr(e.Expr)
	return out
}

func CloneExpr(e lisp.Expr) lisp.Expr {
	if e == nil {
		return nil
	}
	out := make(lisp.Expr, len(e))
	for i := range out {
		out[i] = Clone(e[i])
	}
	return out
}

func CloneLitNode(e *lisp.LitNode) *lisp.LitNode {
	if e == nil {
		return nil
	}
	out := new(lisp.LitNode)
	*out = *e
	out.Lit = CloneLit(e.Lit)
	return out
}

func CloneLit(e lisp.Lit) lisp.Lit {
	switch e := e.(type) {
	case lisp.IdLit:
		return e[:]
	case lisp.StringLit:
		return e[:]
	case lisp.NumberLit:
		return e[:]
	default:
		panic("clone not supported")
	}
}

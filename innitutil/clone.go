package innitutil

import "github.com/ajzaff/innit"

func Clone(n innit.Node) innit.Node {
	if n == nil {
		return nil
	}
	var out innit.Node
	switch x := n.(type) {
	case *innit.LitNode:
		out = CloneLitNode(x)
	case *innit.ExprNode:
		out = CloneExprNode(x)
	default:
		panic("clone not supported")
	}
	return out
}

func CloneExprNode(e *innit.ExprNode) *innit.ExprNode {
	if e == nil {
		return nil
	}
	out := new(innit.ExprNode)
	*out = *e
	out.Expr = CloneExpr(e.Expr)
	return out
}

func CloneExpr(e innit.Expr) innit.Expr {
	if e == nil {
		return nil
	}
	out := make(innit.Expr, len(e))
	for i := range out {
		out[i] = Clone(e[i])
	}
	return out
}

func CloneLitNode(e *innit.LitNode) *innit.LitNode {
	if e == nil {
		return nil
	}
	out := new(innit.LitNode)
	*out = *e
	out.Lit = CloneLit(e.Lit)
	return out
}

func CloneLit(e innit.Lit) innit.Lit {
	switch e := e.(type) {
	case innit.IdLit:
		return e[:]
	case innit.StringLit:
		return e[:]
	case innit.IntLit, innit.FloatLit:
		return e
	default:
		panic("clone not supported")
	}
}

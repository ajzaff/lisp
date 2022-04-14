package innitutil

import "github.com/ajzaff/innit"

func Clone(n innit.Node) innit.Node {
	if n == nil {
		return nil
	}
	var out innit.Node
	switch x := n.(type) {
	case *innit.Lit:
		out = CloneLit(x)
	case *innit.Expr:
		out = CloneExpr(x)
	case innit.NodeList:
		out = CloneNodeList(x)
	default:
		panic("clone not supported")
	}
	return out
}

func CloneExpr(e *innit.Expr) *innit.Expr {
	if e == nil {
		return nil
	}
	out := new(innit.Expr)
	*out = *e
	out.X = CloneNodeList(e.X)
	return out
}

func CloneNodeList(e innit.NodeList) innit.NodeList {
	if e == nil {
		return nil
	}
	out := make(innit.NodeList, len(e))
	for i := range out {
		out[i] = Clone(e[i])
	}
	return out
}

func CloneLit(e *innit.Lit) *innit.Lit {
	if e == nil {
		return nil
	}
	out := new(innit.Lit)
	*out = *e
	return out
}

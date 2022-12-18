package lisp

import "github.com/ajzaff/lisp"

func FromId(v lisp.Val) string {
	return string(v.(lisp.Lit).Text)
}

func FromString(v lisp.Val) string {
	return string(v.(lisp.Lit).Text)
}

func IdTuple(v lisp.Val) []string {
	expr := v.(lisp.Expr)
	res := make([]string, 0, len(expr))
	for _, x := range expr {
		res = append(res, FromId(x.Val))
	}
	return res
}

func IdSet(v lisp.Val) map[string]struct{} {
	expr := v.(lisp.Expr)
	if FromId(expr[0].Val) != "set" {
		panic("IdSet: set should have Val marker")
	}
	m := make(map[string]struct{}, len(expr)-1)
	for _, x := range expr[1:] {
		m[FromId(x.Val)] = struct{}{}
	}
	return m
}

func StringSet(v lisp.Val) map[string]struct{} {
	expr := v.(lisp.Expr)
	if FromId(expr[0].Val) != "set" {
		panic("IdSet: set should have Val marker")
	}
	m := make(map[string]struct{}, len(expr)-1)
	for _, x := range expr[1:] {
		m[FromString(x.Val)] = struct{}{}
	}
	return m
}

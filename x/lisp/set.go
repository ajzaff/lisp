package lisp

import "github.com/ajzaff/lisp"

func FromId(v lisp.Val) string {
	return string(v.(lisp.Lit).Text)
}

func IdTuple(v lisp.Val) []string {
	group := v.(lisp.Group)
	var res []string
	for _, e := range group {
		res = append(res, FromId(e))
	}
	return res
}

func IdSet(v lisp.Val) map[string]struct{} {
	group := v.(lisp.Group)
	if FromId(group[0]) != "set" {
		panic("IdSet: set should have Val marker")
	}
	m := map[string]struct{}{}
	for _, e := range group[1:] {
		m[FromId(e)] = struct{}{}
	}
	return m
}

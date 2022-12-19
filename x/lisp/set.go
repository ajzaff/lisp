package lisp

import (
	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/lisputil"
)

func FromId(v lisp.Val) string {
	return string(v.(lisp.Lit).Text)
}

func IdTuple(v lisp.Val) []string {
	e := v.(*lisp.Cons)
	var res []string
	for e := e; e != nil; e = e.Cons {
		res = append(res, FromId(e.Val))
	}
	return res
}

func IdSet(v lisp.Val) map[string]struct{} {
	cons := v.(*lisp.Cons)
	if FromId(cons.Val) != "set" {
		panic("IdSet: set should have Val marker")
	}
	m := map[string]struct{}{}
	for e := lisputil.Tail(cons); e != nil; e = e.Cons {
		m[FromId(e.Val)] = struct{}{}
	}
	return m
}

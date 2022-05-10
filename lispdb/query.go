package lispdb

import (
	"strings"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/hash"
	"github.com/ajzaff/lisp/lisputil"
)

type QueryInterface interface {
	LispDB
	LoadInterface
	EachRef(ID, func(ID) bool)
	EachInverseRef(ID, func(ID) bool)
}

func QueryOneID(db QueryInterface, id ID) (lisp.Val, float64) {
	v, w := db.Load(id)
	if v != nil {
		return v, w
	}
	n, w := queryOneNode(db, id)
	return n.Val(), w
}

func queryOneNode(db QueryInterface, id ID) (lisp.Node, float64) {
	v, w := db.Load(id)
	if v != nil {
		return &lisp.LitNode{Lit: v}, w
	}
	var x lisp.Expr
	db.EachRef(id, func(i ID) bool {
		e, _ := queryOneNode(db, i)
		x = append(x, e)
		return true
	})
	return &lisp.ExprNode{Expr: x}, w
}

func EachTransRef(db QueryInterface, root ID, fn func(ID) bool) {
	stack := []ID{root}
	for len(stack) > 0 {
		e := stack[0]
		stack = stack[1:]
		if !fn(e) {
			return
		}
		db.EachRef(e, func(i ID) bool {
			stack = append(stack, i)
			return true
		})
	}
}

func EachTransInverseRef(db QueryInterface, root ID, fn func(ID) bool) {
	stack := []ID{root}
	for len(stack) > 0 {
		e := stack[0]
		stack = stack[1:]
		if !fn(e) {
			return
		}
		db.EachInverseRef(e, func(i ID) bool {
			stack = append(stack, i)
			return true
		})
	}
}

type QueryResult struct {
	matches [][]ID
	elems   []string
	err     error
}

func (r *QueryResult) Err() error {
	return r.err
}

func (r *QueryResult) Elements() []string {
	eCopy := make([]string, len(r.elems))
	copy(eCopy, r.elems)
	return eCopy
}

func (r *QueryResult) EachMatch(fn func(id []ID) bool) {
	for _, m := range r.matches {
		if !fn(m) {
			return
		}
	}
}

// Query uses q to query for matching elements in db.
//
// Matching elements are prefixed with "?".
//
// Example:
//	r := Query(db, "(?who is-on first)")
//	r.EachMatch(...)
//	// []ID{834583485} // "who"
func Query(db QueryInterface, q string) *QueryResult {
	var r QueryResult
	var qn []lisp.Val
	sc := lisp.NewNodeScanner(lisp.NewTokenScanner(strings.NewReader(q)))
	for sc.Scan() {
		qn = append(qn, sc.Node().Val())
	}
	if err := sc.Err(); err != nil {
		r.err = err
		return &r
	}
	if len(qn) != 1 {
		panic("union of multiple query expressions is not yet supported")
	}
	var h hash.MapHash
	h.SetSeed(db.Seed())
	qh := h.Sum64()
	h.WriteVal(qn[0])
	if _, w := db.Load(qh); w > 0 {
		// Exact match.
		r.matches = [][]ID{{qh}}
		return &r
	}
	r.elems = queryElements(qn[0])
	panic("not implemented")
}

func queryElements(q lisp.Val) (elems []string) {
	var v lisp.Visitor
	v.SetBeforeExprVisitor(func(e lisp.Expr) {
		if x := lisputil.Head(e); x != nil {
			if lisputil.Equal(x, lisp.IdLit("?")) {
				name := lisputil.Head(e[1:])
				if name != nil {
					if x, ok := name.(lisp.IdLit); ok {
						elems = append(elems, x.String())
					}
				}
			}
		}
	})
	v.Visit(q)
	return
}

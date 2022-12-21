package lispdb

import (
	"strings"

	"github.com/ajzaff/lisp"
	lisputil "github.com/ajzaff/lisp/lisp"
	"github.com/ajzaff/lisp/scan"
	"github.com/ajzaff/lisp/visit"
	"github.com/ajzaff/lisp/x/hash"
)

type QueryInterface interface {
	LispDB
	LoadInterface
	EachRef(ID, func(ID) bool)
	EachInverseRef(ID, func(ID) bool)
}

func QueryOneID(db QueryInterface, id ID) (lisp.Val, float64) {
	v, w := db.Load(id)
	if v.Token != lisp.Invalid {
		return v, w
	}
	return queryOneVal(db, id)
}

func queryOneVal(db QueryInterface, id ID) (lisp.Val, float64) {
	v, w := db.Load(id)
	if v.Token != lisp.Invalid {
		return v, w
	}
	x := &lisp.Cons{}
	db.EachRef(id, func(i ID) bool {
		e, _ := queryOneVal(db, i)
		x.Val = e
		// FIXME: This produces an incorrect Cons.
		// FIXME: Use Cons Builder.
		x.Cons = &lisp.Cons{}
		x = x.Cons
		return true
	})
	return x, w
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
//
//	r := Query(db, "(?who is-on first)")
//	r.EachMatch(...)
//	// []ID{834583485} // "who"
func Query(db QueryInterface, q string) *QueryResult {
	var r QueryResult
	var qv []lisp.Val
	var s scan.TokenScanner
	s.Reset(strings.NewReader(q))
	var sc scan.NodeScanner
	sc.Reset(&s)
	for sc.Scan() {
		_, _, v := sc.Node()
		qv = append(qv, v)
	}
	if err := sc.Err(); err != nil {
		r.err = err
		return &r
	}
	if len(qv) != 1 {
		panic("union of multiple query expressions is not yet supported")
	}
	var h hash.MapHash
	h.SetSeed(db.Seed())
	qh := h.Sum64()
	h.WriteVal(qv[0])
	if _, w := db.Load(qh); w > 0 {
		// Exact match.
		r.matches = [][]ID{{qh}}
		return &r
	}
	r.elems = queryElements(qv[0])
	panic("not implemented")
}

func queryElements(q lisp.Val) (elems []string) {
	var v visit.Visitor
	v.SetBeforeConsVisitor(func(e *lisp.Cons) {
		if x := lisputil.Head(e); x != nil {
			if lisputil.Equal(x, lisp.Lit{Token: lisp.Id, Text: "q"}) {
				name := lisputil.Head(e.Cons)
				if name != nil {
					if x, ok := name.(lisp.Lit); ok {
						elems = append(elems, x.String())
					}
				}
			}
		}
	})
	v.Visit(q)
	return
}

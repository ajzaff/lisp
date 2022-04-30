package innitdb

import (
	"hash/maphash"
	"strings"

	"github.com/ajzaff/innit"
	"github.com/ajzaff/innit/hash"
)

type QueryInterface interface {
	InnitDB
	LoadInterface
	EachRef(ID, func(ID) bool)
	EachInverseRef(ID, func(ID) bool)
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
	qn, err := innit.Parse(q)
	if err != nil {
		r.err = err
		return &r
	}
	var h maphash.Hash
	h.SetSeed(db.Seed())
	qh := h.Sum64()
	hash.Node(&h, qn)
	if v, _ := db.Load(qh); v != nil {
		// Exact match.
		r.matches = [][]ID{{qh}}
		return &r
	}
	r.elems = queryElements(qn)
	panic("not implemented")
}

func queryElements(q innit.Node) (elems []string) {
	var v innit.Visitor
	v.SetLitVisitor(func(e *innit.Lit) {
		if e.Tok == innit.Id && strings.HasPrefix(e.Value, "?") {
			elems = append(elems, e.Value)
		}
	})
	v.Visit(q)
	return
}

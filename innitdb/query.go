package innitdb

import (
	"fmt"
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
	if len(qn) != 1 {
		r.err = fmt.Errorf("expected exactly 1 Val in query expression, got %d", len(qn))
		return &r
	}
	var h maphash.Hash
	h.SetSeed(db.Seed())
	qh := h.Sum64()
	hash.Val(&h, qn[0].Val())
	if v, _ := db.Load(qh); v != nil {
		// Exact match.
		r.matches = [][]ID{{qh}}
		return &r
	}
	r.elems = queryElements(qn[0].Val())
	panic("not implemented")
}

func queryElements(q innit.Val) (elems []string) {
	var v innit.Visitor
	v.SetLitVisitor(func(e innit.Lit) {
		if id, ok := e.(innit.IdLit); ok {
			if strings.HasPrefix(string(id), "?") {
				elems = append(elems, string(id[1:]))
			}
		}
	})
	v.Visit(q)
	return
}

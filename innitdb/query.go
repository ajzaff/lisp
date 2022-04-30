package innitdb

import (
	"hash/maphash"

	"github.com/ajzaff/innit"
	"github.com/ajzaff/innit/hash"
)

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
// Matching elements are prefixed with ":".
//
// Example:
//	r := Query(db, "(:who ?is-on first)")
//	r.EachMatch(...)
//	// []ID{834583485} // "who"
func Query(db InnitDB, q string) *QueryResult {
	var r QueryResult
	qn, err := innit.Parse(q)
	if err != nil {
		r.err = err
		return &r
	}
	var h maphash.Hash
	h.SetSeed(db.Seed())
	qh := hash.Hash(&h, qn)
	if v, _ := db.Load(qh); v != nil {
		// Exact match.
		r.matches = [][]ID{{qh}}
		return &r
	}
	panic("not implemented")
}

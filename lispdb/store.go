package lispdb

import (
	"hash/maphash"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/hash"
)

type StoreInterface interface {
	InnitDB
	Store([]*TVal, float64) error
}

type TVal struct {
	ID
	lisp.Val
	Refs        []uint64
	InverseRefs []uint64
}

func Store(s StoreInterface, n lisp.Val, w float64) error {
	var (
		stack []*TVal
		t     []*TVal
		h     maphash.Hash
		v     lisp.Visitor
	)
	h.SetSeed(s.Seed())

	v.SetBeforeExprVisitor(func(e lisp.Expr) {
		h.Reset()
		hash.Expr(&h, e)
		id := h.Sum64()
		entry := &TVal{ID: id, Val: e}
		for _, parent := range stack {
			parent.Refs = append(parent.Refs, id)
			entry.InverseRefs = append(entry.InverseRefs, parent.ID)
		}
		stack = append(stack, entry)
	})
	v.SetAfterExprVisitor(func(e lisp.Expr) {
		entry := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		t = append(t, entry)
	})
	v.SetLitVisitor(func(e lisp.Lit) {
		h.Reset()
		hash.Lit(&h, e)
		id := h.Sum64()
		entry := &TVal{ID: id, Val: e}
		for _, parent := range stack {
			parent.Refs = append(parent.Refs, id)
			entry.InverseRefs = append(entry.InverseRefs, parent.ID)
		}
		t = append(t, entry)
	})
	v.Visit(n)
	return s.Store(t, w)
}

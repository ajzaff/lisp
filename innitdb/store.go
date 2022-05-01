package innitdb

import (
	"hash/maphash"

	"github.com/ajzaff/innit"
	"github.com/ajzaff/innit/hash"
)

type StoreInterface interface {
	InnitDB
	Store([]*TVal, float64) error
}

type TVal struct {
	ID
	innit.Val
	Refs        []uint64
	InverseRefs []uint64
}

func Store(s StoreInterface, n innit.Val, w float64) error {
	var (
		stack []*TVal
		t     []*TVal
		h     maphash.Hash
		v     innit.Visitor
	)
	h.SetSeed(s.Seed())

	v.SetBeforeExprVisitor(func(e innit.Expr) {
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
	v.SetAfterExprVisitor(func(e innit.Expr) {
		entry := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		t = append(t, entry)
	})
	v.SetLitVisitor(func(e innit.Lit) {
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

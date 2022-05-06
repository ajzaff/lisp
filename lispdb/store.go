package lispdb

import (
	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/hash"
)

type StoreInterface interface {
	LispDB
	Store([]*TVal, float64) error
}

type TVal struct {
	ID
	lisp.Lit
	Refs        []uint64
	InverseRefs []uint64
}

func Store(s StoreInterface, n lisp.Val, w float64) error {
	var (
		stack []*TVal
		t     []*TVal
		h     hash.MapHash
		v     lisp.Visitor
	)
	h.SetSeed(s.Seed())

	v.SetBeforeExprVisitor(func(e lisp.Expr) {
		h.Reset()
		h.WriteExpr(e)
		id := h.Sum64()
		entry := &TVal{ID: id}
		if len(stack) > 0 {
			stack[len(stack)-1].Refs = append(stack[len(stack)-1].Refs, id)
			entry.InverseRefs = append(entry.InverseRefs, stack[len(stack)-1].ID)
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
		h.WriteLit(e)
		id := h.Sum64()
		entry := &TVal{ID: id, Lit: e}
		if len(stack) > 0 {
			stack[len(stack)-1].Refs = append(stack[len(stack)-1].Refs, id)
			entry.InverseRefs = append(entry.InverseRefs, stack[len(stack)-1].ID)
		}
		t = append(t, entry)
	})
	v.Visit(n)
	return s.Store(t, w)
}

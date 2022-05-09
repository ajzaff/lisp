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

// Store each value in s each with weight w.
//
// Store writes to s in a single transaction built in memory.
func Store(s StoreInterface, vals []lisp.Val, w float64) error {
	var (
		stack []*TVal
		t     []*TVal
		h     hash.MapHash
		v     lisp.Visitor
	)
	h.SetSeed(s.Seed())

	v.SetBeforeExprVisitor(func(e lisp.Expr) {
		h.Reset()
		h.WriteVal(e)
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
		h.WriteVal(e)
		id := h.Sum64()
		entry := &TVal{ID: id, Lit: e}
		if len(stack) > 0 {
			stack[len(stack)-1].Refs = append(stack[len(stack)-1].Refs, id)
			entry.InverseRefs = append(entry.InverseRefs, stack[len(stack)-1].ID)
		}
		t = append(t, entry)
	})
	for _, val := range vals {
		v.Visit(val)
		if len(stack) > 0 {
			panic("internal error: stack not empty after Visit")
		}
	}
	return s.Store(t, w)
}

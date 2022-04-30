package innitdb

import (
	"hash/maphash"

	"github.com/ajzaff/innit"
	"github.com/ajzaff/innit/hash"
)

type StoreInterface interface {
	Seed() maphash.Seed
	Store(Transaction) error
}

type Transaction struct {
	Fc    int
	Nodes []TransactionNode
}

type TransactionNode struct {
	Id    uint64
	Node  innit.Node
	Refs  []uint64
	IRefs []uint64
}

func Store(s StoreInterface, n innit.Node, fc int) (rootId uint64, err error) {
	var (
		stack []TransactionNode
		t     Transaction
		h     maphash.Hash
		v     innit.Visitor
		first = true
	)
	t.Fc = fc
	h.SetSeed(s.Seed())

	v.SetBeforeExprVisitor(func(e *innit.Expr) {
		hash.Expr(&h, e)
		id := h.Sum64()
		h.Reset()
		if first {
			rootId = id
			first = false
		}
		entry := TransactionNode{Id: id, Node: e}
		for _, parent := range stack {
			parent.Refs = append(parent.Refs, id)
			entry.IRefs = append(entry.IRefs, parent.Id)
		}
		stack = append(stack, entry)
	})
	v.SetAfterExprVisitor(func(e *innit.Expr) {
		entry := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		t.Nodes = append(t.Nodes, entry)
	})
	v.SetLitVisitor(func(e *innit.Lit) {
		hash.Lit(&h, e)
		id := h.Sum64()
		h.Reset()
		if first {
			rootId = id
			first = false
		}
		entry := TransactionNode{Id: id, Node: e}
		for _, parent := range stack {
			parent.Refs = append(parent.Refs, id)
			entry.IRefs = append(entry.IRefs, parent.Id)
		}
		t.Nodes = append(t.Nodes, entry)
	})
	v.Visit(n)
	return rootId, s.Store(t)
}

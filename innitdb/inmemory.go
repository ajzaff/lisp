package innitdb

import (
	"hash/maphash"
	"sync"

	"github.com/ajzaff/innit"
	"github.com/ajzaff/innit/hash"
)

type InMemory struct {
	nodes       map[uint64]innit.Node // node hash       => node
	refs        map[uint64][]uint64   // expr hash       => child nodes
	inverseRefs map[uint64]uint64     // child node hash => parent expr hash

	h  maphash.Hash
	rw sync.RWMutex // guards struct
}

func NewInMemory() *InMemory {
	return &InMemory{
		nodes:       make(map[uint64]innit.Node),
		refs:        make(map[uint64][]uint64),
		inverseRefs: make(map[uint64]uint64),
	}
}

func (m *InMemory) Store(value innit.Node) (rootId uint64) {
	m.rw.Lock()
	defer m.rw.Unlock()

	var stack []uint64
	first := true

	var v innit.Visitor
	v.SetBeforeExprVisitor(func(e *innit.Expr) {
		id := hash.Expr(&m.h, e)
		if first {
			rootId = id
			first = false
		}
		m.nodes[id] = e
		for _, parentId := range stack {
			m.refs[parentId] = append(m.refs[parentId], id)
			m.inverseRefs[id] = parentId
		}
		stack = append(stack, id)
	})
	v.SetAfterExprVisitor(func(e *innit.Expr) {
		stack = stack[:len(stack)-1]
	})
	v.SetLitVisitor(func(e *innit.Lit) {
		id := hash.Lit(&m.h, e)
		if first {
			rootId = id
			first = false
		}
		m.nodes[id] = e
		for _, parentId := range stack {
			m.refs[parentId] = append(m.refs[parentId], id)
			m.inverseRefs[id] = parentId
		}
	})
	v.Visit(value)
	return rootId
}

package innitdb

import (
	"hash/maphash"
	"sync"

	"github.com/ajzaff/innit"
	"github.com/ajzaff/innit/hash"
)

type InMemory struct {
	nodes       map[ID]innit.Node // node hash       => node
	refs        map[ID][]ID       // expr hash       => child nodes
	inverseRefs map[ID]ID         // child node hash => parent expr hash

	hs maphash.Seed
	rw sync.RWMutex // guards struct
}

func NewInMemory() *InMemory {
	return &InMemory{
		nodes:       make(map[ID]innit.Node),
		refs:        make(map[ID][]ID),
		inverseRefs: make(map[ID]ID),
	}
}

func (m *InMemory) Seed() maphash.Seed { return m.hs }

func (m *InMemory) Load(id ID) innit.Node {
	m.rw.RLock()
	defer m.rw.RUnlock()

	return m.nodes[id]
}

func (m *InMemory) Store(value innit.Node) (rootId ID) {
	var stack []ID
	first := true

	var h maphash.Hash
	h.SetSeed(m.Seed())

	m.rw.Lock()
	defer m.rw.Unlock()

	var v innit.Visitor
	v.SetBeforeExprVisitor(func(e *innit.Expr) {
		id := hash.Expr(&h, e)
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
		id := hash.Lit(&h, e)
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

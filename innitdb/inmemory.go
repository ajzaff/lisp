package innitdb

import (
	"hash/maphash"
	"sync"

	"github.com/ajzaff/innit"
	"github.com/ajzaff/innit/hash"
)

type inMemoryEntry struct {
	innit.Node
	Weight float64
}

type InMemory struct {
	entries     map[ID]*inMemoryEntry // node hash       => node
	refs        map[ID][]ID           // expr hash       => child nodes
	inverseRefs map[ID][]ID           // child node hash => parent expr hash

	hs maphash.Seed
	rw sync.RWMutex // guards struct
}

func NewInMemory() *InMemory {
	return &InMemory{
		entries:     make(map[ID]*inMemoryEntry),
		refs:        make(map[ID][]ID),
		inverseRefs: make(map[ID][]ID),
		hs:          maphash.MakeSeed(),
	}
}

func (m *InMemory) Seed() maphash.Seed { return m.hs }

func (m *InMemory) Load(id ID) (innit.Node, float64) {
	m.rw.RLock()
	defer m.rw.RUnlock()

	e := m.entries[id]
	return e.Node, e.Weight
}

func (m *InMemory) Store(value innit.Node, weight float64) (rootId ID) {
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
		if entry, ok := m.entries[id]; ok {
			entry.Weight += weight
		} else {
			m.entries[id] = &inMemoryEntry{Node: e, Weight: weight}
		}
		for _, parentId := range stack {
			m.refs[parentId] = append(m.refs[parentId], id)
			m.inverseRefs[id] = append(m.inverseRefs[id], parentId)
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
		if entry, ok := m.entries[id]; ok {
			entry.Weight += weight
		} else {
			m.entries[id] = &inMemoryEntry{Node: e, Weight: weight}
		}
		for _, parentId := range stack {
			m.refs[parentId] = append(m.refs[parentId], id)
			m.inverseRefs[id] = append(m.inverseRefs[id], parentId)
		}
	})
	v.Visit(value)
	return rootId
}

func (m *InMemory) EachRef(root ID, fn func(ID) bool) {
	for _, r := range m.refs[root] {
		if !fn(r) {
			return
		}
	}
}

func (m *InMemory) EachInverseRef(root ID, fn func(ID) bool) {
	for _, r := range m.inverseRefs[root] {
		if !fn(r) {
			return
		}
	}
}

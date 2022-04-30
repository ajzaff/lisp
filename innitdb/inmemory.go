package innitdb

import (
	"hash/maphash"
	"sync"

	"github.com/ajzaff/innit"
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
	if e == nil {
		return nil, 0
	}
	return e.Node, e.Weight
}

func (m *InMemory) Store(t []*TNode, w float64) error {
	m.rw.Lock()
	defer m.rw.Unlock()
	for _, te := range t {
		if e, ok := m.entries[te.ID]; ok {
			e.Weight += w
		} else {
			m.entries[te.ID] = &inMemoryEntry{Node: te.Node, Weight: w}
		}
		m.refs[te.ID] = append(m.refs[te.ID], te.Refs...)
		m.inverseRefs[te.ID] = append(m.inverseRefs[te.ID], te.InverseRefs...)
	}
	return nil
}

func (m *InMemory) EachRef(root ID, fn func(ID) bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	for _, r := range m.refs[root] {
		if !fn(r) {
			return
		}
	}
}

func (m *InMemory) EachInverseRef(root ID, fn func(ID) bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	for _, r := range m.inverseRefs[root] {
		if !fn(r) {
			return
		}
	}
}

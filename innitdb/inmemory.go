package innitdb

import (
	"hash/maphash"
	"sync"

	"github.com/ajzaff/innit"
)

type inMemoryEntry struct {
	innit.Val
	Weight float64
}

type InMemory struct {
	entries     map[ID]*inMemoryEntry // Val hash       => Val entry
	refs        map[ID][]ID           // expr hash       => child Vals
	inverseRefs map[ID][]ID           // child Val hash => parent expr hash

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

func (m *InMemory) Load(id ID) (innit.Val, float64) {
	m.rw.RLock()
	defer m.rw.RUnlock()

	e := m.entries[id]
	if e == nil {
		return nil, 0
	}
	return e.Val, e.Weight
}

func (m *InMemory) Store(t []*TVal, w float64) error {
	m.rw.Lock()
	defer m.rw.Unlock()
	for _, te := range t {
		if e, ok := m.entries[te.ID]; ok {
			e.Weight += w
		} else {
			m.entries[te.ID] = &inMemoryEntry{Val: te.Val, Weight: w}
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

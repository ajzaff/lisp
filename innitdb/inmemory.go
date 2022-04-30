package innitdb

import (
	"hash/maphash"
	"sync"

	"github.com/ajzaff/innit"
)

type InMemory struct {
	fd    map[uint64]int                 // node hash => frequency count
	nodes map[uint64]innit.Node          // node hash       => node
	refs  map[uint64]map[uint64]struct{} // expr hash       => child nodes
	irefs map[uint64]map[uint64]struct{} // child node hash => parent exprs

	seed maphash.Seed
	rw   sync.RWMutex // guards struct
}

func NewInMemory() *InMemory {
	return &InMemory{
		fd:    make(map[uint64]int),
		nodes: make(map[uint64]innit.Node),
		refs:  make(map[uint64]map[uint64]struct{}),
		irefs: make(map[uint64]map[uint64]struct{}),
		seed:  maphash.MakeSeed(),
	}
}

func (m *InMemory) Seed() maphash.Seed {
	return m.seed
}

func (m *InMemory) Store(t Transaction) error {
	m.rw.Lock()
	defer m.rw.Unlock()

	for _, n := range t.Nodes {
		m.fd[n.Id] += t.Fc
		m.nodes[n.Id] = n.Node
		refs, ok := m.refs[n.Id]
		if !ok {
			refs = make(map[uint64]struct{})
			m.refs[n.Id] = refs
		}
		for _, ref := range n.Refs {
			refs[ref] = struct{}{}
		}
		irefs, ok := m.irefs[n.Id]
		if !ok {
			irefs = make(map[uint64]struct{})
			m.irefs[n.Id] = irefs
		}
		for _, iref := range n.IRefs {
			irefs[iref] = struct{}{}
		}
	}
	return nil
}

func (m *InMemory) Load(id uint64) (fc int, refs, irefs []uint64) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	for ref := range m.refs[id] {
		refs = append(refs, ref)
	}
	for iref := range m.irefs[id] {
		irefs = append(irefs, iref)
	}
	return m.fd[id], refs, irefs
}

func (m *InMemory) EachInverseRef(id uint64, fn func(uint64)) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	for id := range m.irefs[id] {
		fn(id)
	}
}

func (m *InMemory) EachRef(id uint64, fn func(uint64)) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	for id := range m.refs[id] {
		fn(id)
	}
}

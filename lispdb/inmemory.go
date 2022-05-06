package lispdb

import (
	"hash/maphash"
	"sync"

	"github.com/ajzaff/lisp"
)

type entryInMemory interface {
	entry()
	Weight() float64
	AddWeight(float64)
}

type litEntryInMemory struct {
	lisp.Lit
	weight float64
}

func (e *litEntryInMemory) entry()              {}
func (e *litEntryInMemory) Weight() float64     { return e.weight }
func (e *litEntryInMemory) AddWeight(w float64) { e.weight += w }

type exprEntryInMemory struct {
	refs        []ID
	inverseRefs []ID
	weight      float64
}

func (e *exprEntryInMemory) entry()              {}
func (e *exprEntryInMemory) Weight() float64     { return e.weight }
func (e *exprEntryInMemory) AddWeight(w float64) { e.weight += w }

type InMemory struct {
	entries map[ID]entryInMemory // hash ID => entry

	hs maphash.Seed
	rw sync.RWMutex // guards InMemory
}

func NewInMemory() *InMemory {
	return &InMemory{
		entries: make(map[ID]entryInMemory),
		hs:      maphash.MakeSeed(),
	}
}

func (m *InMemory) Seed() maphash.Seed { return m.hs }

func (m *InMemory) Load(id ID) (lit lisp.Lit, w float64) {
	m.rw.RLock()
	defer m.rw.RUnlock()

	e, ok := m.entries[id]
	if !ok {
		return nil, 0
	}
	w = e.Weight()
	if e, ok := e.(*litEntryInMemory); ok {
		lit = e.Lit
	}
	return
}

func (m *InMemory) Store(t []*TVal, w float64) error {
	m.rw.Lock()
	defer m.rw.Unlock()
	for _, te := range t {
		if e, ok := m.entries[te.ID]; ok {
			e.AddWeight(w)
		} else {
			switch v := te.Val.(type) {
			case lisp.Lit:
				m.entries[te.ID] = &litEntryInMemory{Lit: v, weight: w}
			case lisp.Expr:
				m.entries[te.ID] = &exprEntryInMemory{refs: te.Refs, inverseRefs: te.InverseRefs, weight: w}
			default:
				panic("unreachable")
			}
		}
	}
	return nil
}

func (m *InMemory) EachRef(root ID, fn func(ID) bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()

	e, ok := m.entries[root]
	if !ok {
		return
	}
	if e, ok := e.(*exprEntryInMemory); ok {
		for _, r := range e.refs {
			if !fn(r) {
				return
			}
		}
	}
}

func (m *InMemory) EachInverseRef(root ID, fn func(ID) bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()

	e, ok := m.entries[root]
	if !ok {
		return
	}
	if e, ok := e.(*exprEntryInMemory); ok {
		for _, r := range e.inverseRefs {
			if !fn(r) {
				return
			}
		}
	}
}

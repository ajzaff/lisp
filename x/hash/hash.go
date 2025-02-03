package hash

import (
	"hash/maphash"
	"sync"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/visit"
)

// MapHasher wraps a maphash for writing Lisp Values.
type MapHash struct {
	maphash.Hash

	v    visit.Visitor // init by once
	once sync.Once
}

func (h *MapHash) initVisitor() {
	var delim bool
	h.v.SetLitVisitor(func(x lisp.Lit) {
		if delim {
			h.WriteByte(' ')
		}
		h.WriteString(string(x))
		delim = true
	})
	h.v.SetBeforeGroupVisitor(func(lisp.Group) { h.WriteByte('('); delim = false })
	h.v.SetAfterGroupVisitor(func(lisp.Group) { h.WriteByte(')'); delim = false })
}

// WriteValue hashes the Val into the MapHash.
func (h *MapHash) WriteVal(v lisp.Val) {
	h.once.Do(h.initVisitor)
	h.v.Visit(v)
}

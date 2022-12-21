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
		h.WriteString(x.Text)
		delim = true
	})
	h.v.SetBeforeConsVisitor(func(*lisp.Cons) { h.WriteByte('('); delim = false })
	h.v.SetAfterConsVisitor(func(*lisp.Cons) { h.WriteByte(')'); delim = false })
}

// WriteValue hashes the Val into the MapHash.
func (h *MapHash) WriteVal(v lisp.Val) {
	h.once.Do(h.initVisitor)
	h.v.Visit(v)
}

package hash

import (
	"hash/maphash"
	"sync"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/print"
	"github.com/ajzaff/lisp/visit"
)

// MapHasher wraps a maphash for writing Lisp Values.
type MapHash struct {
	maphash.Hash

	v    visit.Visitor // init by once
	once sync.Once
}

// WriteValue hashes the Val into the MapHash.
func (h *MapHash) WriteVal(v lisp.Val) {
	print.StdPrinter(&h.Hash).Print(v)
}

var mapHashVisitor visit.Visitor

func (h *MapHash) initVisitor() {
	mapHashVisitor.SetLitVisitor(func(x lisp.Lit) { h.WriteString(x.Text) })
	mapHashVisitor.SetBeforeConsVisitor(func(*lisp.Cons) { h.WriteByte('(') })
	mapHashVisitor.SetAfterConsVisitor(func(*lisp.Cons) { h.WriteByte(')') })
}

func (h *MapHash) WriteVisitedVal(root lisp.Val) {
	h.once.Do(h.initVisitor)
	h.v.Visit(root)
}

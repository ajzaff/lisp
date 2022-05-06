package hash

import (
	"fmt"
	"hash/maphash"
	"sync"

	"github.com/ajzaff/lisp"
)

// MapHasher wraps a maphash for writing Lisp Values.
type MapHash struct {
	maphash.Hash
	delim string
	once  sync.Once
}

func (h *MapHash) init() {
	h.delim = " "
}

// WriteValue hashes the Val into the MapHash.
func (h *MapHash) WriteVal(v lisp.Val) {
	h.once.Do(h.init)
	h.writeValDelim(v, "")
}

func (h *MapHash) writeValDelim(v lisp.Val, delim string) {
	switch v := v.(type) {
	case lisp.Expr:
		h.writeExprDelim(v, delim)
	case lisp.Lit:
		h.writeLitDelim(v, delim)
	}
}

// WriteExpr hashes v into the MapHash.
func (h *MapHash) WriteExpr(v lisp.Expr) {
	h.once.Do(h.init)
	h.writeExprDelim(v, h.delim)
}

func (h *MapHash) writeExprDelim(v lisp.Expr, delim string) {
	h.WriteByte('(')
	for i, e := range v {
		switch e := e.Val().(type) {
		case lisp.Expr:
			h.writeExprDelim(e, delim)
		case lisp.Lit:
			lastElem := i+1 == len(v)
			if lastElem {
				delim = ""
			}
			h.writeLitDelim(e, delim)
		}
	}
	h.WriteByte(')')
}

// WriteLit hashes v into the MapHash.
func (h *MapHash) WriteLit(v lisp.Lit) {
	h.once.Do(h.init)
	h.writeLitDelim(v, h.delim)
}

func (h *MapHash) writeLitDelim(v lisp.Lit, delim string) {
	fmt.Fprint(h, v.String(), delim)
}

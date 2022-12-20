package hash

import (
	"hash/maphash"

	"github.com/ajzaff/lisp"
)

// MapHasher wraps a maphash for writing Lisp Values.
type MapHash struct {
	maphash.Hash
}

// WriteValue hashes the Val into the MapHash.
func (h *MapHash) WriteVal(v lisp.Val) {
	lisp.StdPrinter(&h.Hash).Print(v)
}

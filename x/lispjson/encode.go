package lispjson

import (
	"io"

	"github.com/ajzaff/lisp"
)

type Encoder struct{ w io.Writer }

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

func (e *Encoder) Encode(v lisp.Val) {
	switch v := v.(type) {
	case lisp.Lit:
		// FIXME: No need to Quote the Id if its valid.
		// (Number | Letter) & ~Print = {}.
		e.w.Write([]byte(v.String()))
	case *lisp.Cons:
		e.w.Write([]byte{'['})
		i := 0
		for v := v; v != nil; v = v.Cons {
			if i > 0 {
				e.w.Write([]byte{','})
			}
			e.Encode(v.Val)
			i++
		}
		e.w.Write([]byte{']'})
	}
}

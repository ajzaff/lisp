package lispjson

import (
	"io"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/x/stringer"
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
		e.w.Write([]byte(stringer.Lit(v)))
	case lisp.Group:
		e.w.Write([]byte{'['})
		for i, x := range v {
			if i > 0 {
				e.w.Write([]byte{','})
			}
			e.Encode(x)
		}
		e.w.Write([]byte{']'})
	}
}

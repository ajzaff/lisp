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
	case lisp.Expr:
		e.w.Write([]byte{'['})
		for i, elem := range v {
			e.Encode(elem.Val)
			if i+1 < len(v) {
				e.w.Write([]byte{','})
			}
		}
		e.w.Write([]byte{']'})
	}
}

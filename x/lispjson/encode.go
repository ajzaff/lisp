package lispjson

import (
	"io"
	"strconv"

	"github.com/ajzaff/lisp"
)

type Encoder struct{ w io.Writer }

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

func (e *Encoder) Encode(v lisp.Val) {
	switch v := v.(type) {
	case lisp.Lit:
		if v.Token == lisp.Id {
			e.w.Write([]byte(strconv.Quote(v.String())))
		}
		// Number
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

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
		// IntLit, FloatLit, StringLit
		e.w.Write([]byte(strconv.Quote(v.String())))
	case lisp.Expr:
		e.w.Write([]byte{'['})
		for i, elem := range v {
			e.Encode(elem.Val())
			if i+1 < len(v) {
				e.w.Write([]byte{','})
			}
		}
		e.w.Write([]byte{']'})
	}
}

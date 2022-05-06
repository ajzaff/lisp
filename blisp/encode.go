package blisp

import (
	"encoding/binary"
	"io"

	"github.com/ajzaff/lisp"
)

const magic = "lisp\n"

const (
	lit  = 0
	expr = 1
)

type Encoder struct {
	io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

func (e *Encoder) Encode(v lisp.Val) error {
	b := make([]byte, EncodedLen(v))
	encode(v, b)
	_, err := e.Writer.Write(b)
	return err
}

func encode(v lisp.Val, b []byte) int {
	switch v := v.(type) {
	case lisp.Lit:
		b[0] = lit
		var i int
		switch v := v.(type) {
		case lisp.IdLit:
			b[1] = byte(lisp.Id)
		case lisp.IntLit:
			b[1] = byte(lisp.Int)
		case lisp.FloatLit:
			b[1] = byte(lisp.Float)
		case lisp.StringLit:
			b[1] = byte(lisp.String)
			s := v.String()
			i = copy(b[2:], s[1:len(s)-1])
			return 2 + i
		}
		i = copy(b[2:], []byte(v.String()))
		return 2 + i
	case lisp.Expr:
		b[0] = expr
		size := 0
		for _, e := range v {
			size += EncodedLen(e.Val())
		}
		i := binary.PutUvarint(b, uint64(size))
		for _, e := range v {
			i += encode(e.Val(), b[i:])
		}
		return i
	default:
		panic("Unexpected node type")
	}
}

// EncodedLen returns the encoded length of the node in bytes.
func EncodedLen(n lisp.Val) int {
	if n == nil {
		return 0
	}
	switch x := n.(type) {
	case lisp.Lit:
		n := litLen(x)
		return 1 + 1 + varIntLen(uint64(n)) + n
	case lisp.Expr:
		size := 1
		for _, e := range x {
			size += EncodedLen(e.Val())
		}
		return varIntLen(uint64(size)) + size
	default:
		panic("Unexpected Val type")
	}
}

func litLen(v lisp.Lit) int {
	switch v.(type) {
	case lisp.IdLit, lisp.IntLit, lisp.FloatLit:
		return len(v.String())
	case lisp.StringLit:
		return len(v.String()) - 2
	default:
		panic("unexpected token")
	}
}

func varIntLen(x uint64) int {
	switch {
	case x < 1<<7:
		return 1
	case x < 1<<14:
		return 2
	case x < 1<<21:
		return 3
	case x < 1<<28:
		return 4
	case x < 1<<35:
		return 5
	case x < 1<<42:
		return 6
	case x < 1<<49:
		return 7
	case x < 1<<56:
		return 8
	case x < 1<<63:
		return 9
	}
	return 10
}
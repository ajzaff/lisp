package blisp

import (
	"encoding/binary"
	"io"

	"github.com/ajzaff/lisp"
)

const magic = "\x41blisp\n"

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
	if v == nil {
		return 0
	}
	switch v := v.(type) {
	case lisp.Lit:
		var i int
		switch v.Token {
		case lisp.Id:
			b[0] = byte(lisp.Id)
		case lisp.Int:
			b[0] = byte(lisp.Int)
		}
		i = copy(b[2:], []byte(v.String()))
		return 2 + i
	case *lisp.Cons:
		b[0] = byte(lisp.LParen)
		size := 0
		for e := v; e != nil; e = e.Cons {
			size += EncodedLen(e.Val)
		}
		i := binary.PutUvarint(b[1:], uint64(size))
		for e := v; e != nil; e = e.Cons {
			i += encode(e.Val, b[i:])
		}
		return i
	default:
		panic("Unexpected Val type")
	}
}

// EncodedLen returns the encoded length of the node in bytes.
func EncodedLen(n lisp.Val) int {
	if n == nil {
		return 0
	}
	switch x := n.(type) {
	case lisp.Lit:
		switch x.Token {
		case lisp.Id:
			return 1 + varIntLen(uint64(len(x.String()))) + len(x.String())
		case lisp.Int:
			return 1 + varIntLen(uint64(len(x.String()))) + len(x.String())
		default:
			panic("unknown Lit token")
		}
	case *lisp.Cons:
		size := 0
		for e := x; e != nil; e = e.Cons {
			size += EncodedLen(e.Val)
		}
		return 1 + varIntLen(uint64(size)) + size
	default:
		panic("Unexpected Val type")
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
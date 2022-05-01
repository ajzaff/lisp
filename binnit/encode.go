package binnit

import (
	"encoding/binary"
	"io"

	"github.com/ajzaff/innit"
)

const magic = "innit\n"

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

func (e *Encoder) Encode(n innit.Node) error {
	b := make([]byte, EncodedLen(n))
	encode(n, b)
	_, err := e.Writer.Write(b)
	return err
}

func encode(n innit.Node, b []byte) int {
	switch x := n.(type) {
	case *innit.Lit:
		b[0] = lit
		b[1] = byte(x.Tok)
		if x.Tok == innit.String {
			i := copy(b[2:], x.Value[1:len(x.Value)-1])
			return 2 + i
		}
		i := copy(b[2:], x.Value)
		return 2 + i
	case *innit.Expr:
		b[0] = expr
		return 1 + encode(x.X, b[1:])
	case innit.NodeList:
		size := 0
		for _, e := range x {
			size += EncodedLen(e)
		}
		i := binary.PutUvarint(b, uint64(size))
		for _, e := range x {
			i += encode(e, b[i:])
		}
		return i
	default:
		panic("Unexpected node type")
	}
}

// EncodedLen returns the encoded length of the node in bytes.
func EncodedLen(n innit.Node) int {
	if n == nil {
		return 0
	}
	switch x := n.(type) {
	case *innit.Lit:
		n := valLen(x.Tok, x.Value)
		return 1 + tokLen(x.Tok) + varIntLen(uint64(n)) + n
	case *innit.Expr:
		return 1 + EncodedLen(x.X)
	case innit.NodeList:
		size := 0
		for _, e := range x {
			size += EncodedLen(e)
		}
		return varIntLen(uint64(size)) + size
	default:
		panic("Unexpected node type")
	}
}

func tokLen(innit.Token) int { return 1 }

func valLen(t innit.Token, val string) int {
	switch t {
	case innit.Id, innit.Int, innit.Float:
		return len(val)
	case innit.String:
		return len(val) - 2
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

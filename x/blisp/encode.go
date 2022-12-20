package blisp

import (
	"bufio"
	"io"

	"github.com/ajzaff/lisp"
)

type Encoder struct {
	w *bufio.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{bufio.NewWriter(w)}
}

func (e *Encoder) EncodeMagic() int {
	e.w.Write([]byte(Magic))
	return len(Magic)
}

func (e *Encoder) Encode(v lisp.Val) {
	defer e.w.Flush()
	e.encode(v)
}

func (e *Encoder) encode(root lisp.Val) {
	if root == nil {
		return
	}
	switch root := root.(type) {
	case lisp.Lit:
		e.w.WriteByte(byte(root.Token))
		e.w.WriteString(root.Text)
		return
	case *lisp.Cons:
		e.EncodeCons(root)
	default:
		panic("Unexpected Val type")
	}
}

func (e *Encoder) EncodeCons(root *lisp.Cons) {
	e.encodeCons(root, true)
}

func (e *Encoder) encodeCons(root *lisp.Cons, first bool) {
	if root == nil {
		e.w.WriteByte(byte(lisp.RParen))
		return
	}
	if first {
		e.w.WriteByte(byte(lisp.LParen))
	}
	e.encode(root.Val)
	e.encodeCons(root.Cons, false)
}

// EncodedLen returns the encoded length of the node in bytes.
func EncodedLen(v lisp.Val) int {
	if v == nil {
		return 0
	}
	switch v := v.(type) {
	case lisp.Lit:
		n := len(v.Text)
		return 1 + varIntLen(uint64(n)) + n // "{Token}N{text}"
	case *lisp.Cons:
		return EncodedConsLen(v) // "(N{val}...N{val})"
	default:
		panic("Unexpected Val type")
	}
}

// EncodedConsLen returns the encoded length of the Cons in bytes.
func EncodedConsLen(root *lisp.Cons) int {
	var size int
	encodedConsLen(root, true, &size)
	return size
}

func encodedConsLen(root *lisp.Cons, first bool, size *int) {
	if root == nil {
		*size++ // ")"
		return
	}
	if first {
		*size++ // "("
	}
	*size += EncodedLen(root.Val)          // "N{Val}"
	encodedConsLen(root.Cons, false, size) // "..."
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

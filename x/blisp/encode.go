package blisp

import (
	"bufio"
	"io"

	"github.com/ajzaff/lisp"
)

type Encoder struct {
	w *bufio.Writer

	// Delimiter is needed for Ids.
	// Nats are self-delimiting.
	delim bool
}

func (e *Encoder) Reset(w io.Writer) {
	e.reset()
	if w, ok := w.(*bufio.Writer); ok {
		e.w = w
		return
	}
	if e.w == nil {
		e.w = new(bufio.Writer)
	}
	e.w.Reset(w)
}

func (e *Encoder) reset() {
	e.delim = false
}

func (e *Encoder) EncodeMagic() {
	e.w.Write([]byte(Magic))
}

func (e *Encoder) Encode(v lisp.Val) {
	e.reset()
	e.encode(v)
	e.w.Flush()
}

func (e *Encoder) encode(root lisp.Val) {
	if root == nil {
		return
	}
	switch root := root.(type) {
	case lisp.Lit:
		if e.delim {
			e.w.WriteByte(' ')
		}
		e.w.WriteString(string(root))
		e.delim = true // Set delim.
	case lisp.Group:
		e.EncodeGroup(root)
		e.delim = false // Clear delim.
	default:
		panic("Unexpected Val type")
	}
}

func (e *Encoder) EncodeGroup(root lisp.Group) {
	e.w.WriteByte(byte(lisp.LParen))
	for _, x := range root {
		e.encode(x)
	}
	e.w.WriteByte(byte(lisp.RParen))
}

type encodeLen struct {
	n     int
	delim bool
}

// Len returns the encoded length of the Val in bytes.
func Len(v lisp.Val) int {
	var e encodeLen
	e.Len(v)
	return e.n
}

func (e *encodeLen) Len(v lisp.Val) {
	if v == nil {
		return
	}
	switch v := v.(type) {
	case lisp.Lit:
		if e.delim {
			e.n++ // {delim}
		}
		e.n += len(v) // {text}
		e.delim = true
	case lisp.Group:
		e.GroupLen(v) // ({Val}...{Val})
		e.delim = false
	default:
		panic("Unexpected Val type")
	}
}

// GroupLen returns the encoded length of the Group in bytes.
func GroupLen(root lisp.Group) int {
	var e encodeLen
	e.GroupLen(root)
	return e.n
}

func (e *encodeLen) GroupLen(root lisp.Group) {
	e.n++ // "("
	for _, x := range root {
		e.Len(x) // {val}
	}
	e.n++ // ")"
}

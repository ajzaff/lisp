package blisp

import (
	"bufio"
	"encoding/binary"
	"io"
	"strconv"

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
		if root.Token == lisp.Nat {
			var b [10]byte
			i, _ := strconv.ParseUint(root.Text, 10, 64)
			n := binary.PutUvarint(b[:], i)
			e.w.WriteByte(byte(lisp.Nat))
			e.w.Write(b[:n])
			e.delim = false // Clear delim.
			return
		}
		if e.delim {
			e.w.WriteByte(' ')
		}
		e.w.WriteString(root.Text)
		e.delim = true // Set delim.
	case *lisp.Cons:
		e.EncodeCons(root)
		e.delim = false // Clear delim.
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

type encodeLen struct {
	n     int
	delim bool
}

// Len returns the encoded length of the node in bytes.
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
		if v.Token == lisp.Nat {
			i, _ := strconv.ParseUint(v.Text, 10, 64)
			e.n += 1 + varUintLen(i) // Nat{i}
			e.delim = false
			return
		}
		if e.delim {
			e.n++ // {delim}
		}
		e.n += len(v.Text) // {text}
		e.delim = true
	case *lisp.Cons:
		e.ConsLen(v, true) // ({Val}...{Val})
		e.delim = false
	default:
		panic("Unexpected Val type")
	}
}

// ConsLen returns the encoded length of the Cons in bytes.
func ConsLen(root *lisp.Cons) int {
	var e encodeLen
	e.ConsLen(root, true)
	return e.n
}

func (e *encodeLen) ConsLen(root *lisp.Cons, first bool) {
	if root == nil {
		e.n++ // ")"
		return
	}
	if first {
		e.n++ // "("
	}
	e.Len(root.Val)             // "{Val}"
	e.ConsLen(root.Cons, false) // "..."
}

func varUintLen(x uint64) int {
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
	default:
		return 10
	}
}

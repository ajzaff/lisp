package lispjson

import (
	"bytes"
	"io"
	"iter"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/scan"
)

type Decoder struct {
	sc scan.Scanner
}

func NewDecoder(r io.Reader) *Decoder {
	var d Decoder
	d.decodeSrc(r)
	return &d
}

func (d *Decoder) decodeSrc(r io.Reader) {
	// TODO: Implement JSON <-> Lisp readers
	//       That don't buffer the full reader.
	// Transform JSON source to Lisp in-place.
	var buf bytes.Buffer
	io.Copy(&buf, r)
	src := buf.Bytes()
	for i, b := range src {
		switch {
		case b == '[':
			src[i] = '('
		case b == ']':
			src[i] = ')'
		case b == ',':
			src[i] = ' '
		case b == '"':
			// Unquote literals by replacing '"".
			// FIXME: Creative, but this may cause issues.
			src[i] = ' '
		}
	}
	// Tokenize and parse normally.
	var sc scan.Scanner
	sc.Reset(bytes.NewReader(src))
	d.sc = sc
}

func (d *Decoder) Values() iter.Seq[lisp.Val] {
	return func(yield func(lisp.Val) bool) {
		for n := range d.sc.Nodes() {
			if !yield(n.Val) {
				break
			}
		}
	}
}

func (d *Decoder) Err() error { return d.sc.Err() }

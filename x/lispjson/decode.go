package lispjson

import (
	"bytes"
	"io"

	"github.com/ajzaff/lisp"
)

type Decoder struct {
	sc *lisp.NodeScanner
}

func NewDecoder(r io.Reader) *Decoder {
	var d Decoder
	d.decodeSrc(r)
	return &d
}

func (d *Decoder) decodeSrc(r io.Reader) {
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
	var s lisp.TokenScanner
	s.Reset(bytes.NewReader(src))
	var sc lisp.NodeScanner
	sc.Reset(&s)
	d.sc = &sc
}

func (d *Decoder) Decode() (lisp.Val, error) {
	d.sc.Scan()
	if err := d.sc.Err(); err != nil {
		return nil, err
	}
	return d.sc.Node().Val, nil
}

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
	var buf bytes.Buffer
	// FIXME: Do proper JSON token scanning.
	io.Copy(&buf, r)
	src := buf.Bytes()
	var (
		str    bool
		escape bool
	)
	for i, b := range src {
		switch {
		case !str && b == '[':
			src[i] = '('
		case !str && b == ']':
			src[i] = ')'
		case !str && b == ',':
			src[i] = ' '
		case b == '"':
			if str && escape {
				escape = false
				continue
			}
			str = !str
		case str && b == '\\':
			escape = !escape
			continue
		}
		if escape {
			escape = false
		}
	}
	var s lisp.TokenScanner
	s.Reset(bytes.NewReader(src))
	d.sc = lisp.NewNodeScanner(&s)
}

func (d *Decoder) Decode() (lisp.Val, error) {
	d.sc.Scan()
	if err := d.sc.Err(); err != nil {
		return nil, err
	}
	return d.sc.Node().Val, nil
}

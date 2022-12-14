package lisp

import (
	"bytes"
	"errors"
	"io"
	"unicode/utf8"
)

// DecodeStr decodes a String literal at the start of the input b.
//
// When the returned error is nil LitNode will contain a String token.
// Otherwise it returns an error when the input does not start with '"'
// And io.ErrUnexpectedEOF if the String terminates before the end of b.
func DecodeStr(b []byte) (LitNode, error) {
	n, err := decodeStr(b)
	if err != nil {
		return LitNode{}, &TokenError{
			Cause: err,
		}
	}
	return LitNode{
		Lit: Lit{
			Token: String,
			Text:  string(b[:n]),
		},
		EndPos: Pos(n),
	}, nil
}

var errStr = errors.New("expected '\"' to start String")

func decodeStr(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, io.ErrUnexpectedEOF
	}
	if b[0] != '"' {
		return 0, errStr
	}
	for n++; ; {
		i := bytes.IndexAny(b[n:], `\"`)
		if i < 0 {
			return len(b), io.ErrUnexpectedEOF
		}
		n += i
		t := b[n]
		n++
		switch t {
		case '"':
			return n, nil
		case '\\':
			if len(b) <= n {
				return len(b), io.ErrUnexpectedEOF
			}
			// Skip one escaped rune.
			_, size := utf8.DecodeRune(b[n:])
			n += size
		}
	}
}

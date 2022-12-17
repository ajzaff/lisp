package lisp

import (
	"unicode/utf8"
)

type IdDecoder struct{}

func (*IdDecoder) Deocde(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for {
		r, size := utf8.DecodeRune(data[advance:])
		if r == utf8.RuneError {
			return advance, nil, errRune
		}
		advance += size
		if !IsId(r) {
			if advance == size {
				return advance, nil, errRune
			}
			token = data[:advance]
			break
		}
		if len(data) <= advance {
			if atEOF {
				token = data[:advance]
			}
			break
		}
	}
	return advance, token, nil
}

package decoder

import (
	"errors"
	"unicode"
	"unicode/utf8"

	"github.com/ajzaff/lisp/scan"
	"golang.org/x/text/unicode/rangetable"
)

var errRune = errors.New("unexpected rune")

var idTab = rangetable.Merge(unicode.L, unicode.Digit)

type IdDecoder struct{}

func (*IdDecoder) Decode(data []byte, atEOF bool) (advance int, token []byte, err error) {
	r, size := utf8.DecodeRune(data)
	if r == utf8.RuneError {
		return 0, nil, errRune
	}
	if !unicode.IsLetter(r) {
		return 0, nil, nil
	}
	i := 0
	start := scan.Pos(i)
	i += size
	for {
		r, size := utf8.DecodeRune(data[i:])
		if r == utf8.RuneError {
			return i, nil, errRune
		}
		i += size
		if !unicode.Is(idTab, r) {
			break
		}
	}
	end := scan.Pos(i)
	return i, data[start:end], nil
}

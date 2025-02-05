package lisp

import (
	"unicode"

	"golang.org/x/text/unicode/rangetable"
)

var idTab *unicode.RangeTable

func init() {
	idTab = rangetable.Merge(
		unicode.Letter,
		rangetable.New('0', '1', '2', '3', '4', '5', '6', '7', '8', '9'),
	)
}

func IsLit(r rune) bool { return unicode.Is(idTab, r) }

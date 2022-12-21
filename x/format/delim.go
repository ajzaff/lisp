package format

import (
	"unicode"

	"github.com/ajzaff/lisp"
)

// delimClass represents symbols the same delimitable class must be spaced out when printing.
type delimClass int

const (
	delimNone   delimClass = iota // No delimiter needed when printing. Expr.
	delimNeeded                   // Delimiter needed. Id, Nat.
)

// delim returns the DelimClass for a lisp Value.
func delim(v lisp.Val) delimClass {
	switch x := v.(type) {
	case lisp.Lit:
		switch x.Token {
		case lisp.Id:
			return delimNeeded
		}
	}
	return delimNone
}

// delimByte returns the DelimClass for a rune.
func delimRune(r rune) delimClass {
	switch {
	case r == '(', r == ')':
		return delimNone
	case '0' <= r && r <= '9', unicode.IsLetter(r):
		return delimNeeded
	default:
		return delimNone
	}
}

// delimByte returns the DelimClass for a byte.
func delimByte(b byte) delimClass {
	return delimRune(rune(b))
}

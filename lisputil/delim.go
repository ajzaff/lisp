package lisputil

import (
	"unicode"

	"github.com/ajzaff/lisp"
)

// DelimClass represents symbols the same delimitable class must be spaced out when printing.
type DelimClass int

const (
	DelimNone DelimClass = iota // No delimiter needed when printing. Expr.
	DelimId                     // Delimiter needed. Id.
	DelimInt                    // Delimiter needed. Int.
	DelimExpr                   // Delimiter not needed. Expr.
)

// Delim returns the DelimClass for a lisp Value.
func Delim(v lisp.Val) DelimClass {
	switch x := v.(type) {
	case lisp.Lit:
		switch x.Token {
		case lisp.Id:
			return DelimId
		case lisp.Int:
			return DelimInt
		}
	case *lisp.Cons:
		return DelimExpr
	}
	return DelimNone
}

// DelimByte returns the DelimClass for a rune.
func DelimRune(r rune) DelimClass {
	switch {
	case unicode.IsLetter(r):
		return DelimId
	case '0' <= r && r <= '9':
		return DelimInt
	case r == '(', r == ')':
		return DelimExpr
	default:
		return DelimNone
	}
}

// DelimByte returns the DelimClass for a byte.
func DelimByte(b byte) DelimClass {
	return DelimRune(rune(b))
}

package lisp

import (
	"fmt"
	"unicode"
)

// Token is an enumeration which specifies a kind of AST token.
//
//go:generate stringer -type Token
type Token int

const (
	Invalid Token = iota

	Id     // x y z
	Number // 12345 -123.45 1.1e2

	LParen // (
	RParen // )

	Pre  // Pre content, see pre.go.
	Code // Code chunk is untokenized content, see code.go.
)

// TokenError implements an error at a specified line and column.
type TokenError struct {
	Cause error
	Pos   Pos
	Src   []byte
}

func (t *TokenError) LineCol() (line, col int) {
	// src := t.src
	line, col = 1, 1
	// for pos := int(t.Pos); ; {
	// 	n := bytes.IndexByte(src, '\n')
	// 	if n < 0 {
	// 		n = len(t.src)
	// 	} else {
	// 		n++
	// 	}
	// 	if pos < n {
	// 		return line, pos
	// 	}
	// 	line++
	// 	pos -= n
	// 	src = src[n:]
	// }
	return
}

func (t *TokenError) Error() string {
	line, col := t.LineCol()
	return fmt.Sprintf("%v: at line %d: col %d", t.Cause, line, col)
}

func (t *TokenError) Unwrap() error {
	return t.Cause
}

func IsId(r rune) bool {
	return unicode.IsLetter(r)
}

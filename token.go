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
	Id            // abc z3
	Int           // 12345
	LParen        // (
	RParen        // )
)

// TokenError implements an error at a specified line and column.
type TokenError struct {
	Cause error
	Pos   Pos
	Src   []byte
}

// LineCol returns a human readable line and column number based on a Pos and Src.
//
// Do not use these to index into Src.
func (t *TokenError) LineCol() (line, col Pos) {
	n := t.Pos
	line, col = 1, 1
	if n == 0 {
		return
	}
	for _, r := range string(t.Src) {
		if r == '\n' {
			line++
			col = 1
			continue
		}
		if unicode.IsPrint(r) {
			col++
		}
		if n--; n == 0 {
			return
		}
	}

	// Outside Src.
	return 0, 0
}

func (t *TokenError) Error() string {
	line, col := t.LineCol()
	return fmt.Sprintf("%v: at line %d: col %d", t.Cause, line, col)
}

func (t *TokenError) Unwrap() error {
	return t.Cause
}

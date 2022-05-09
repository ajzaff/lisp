package lisp

import (
	"fmt"
	"unicode"
)

// Token is an enumeration which specifies a kind of AST token.
//go:generate stringer -type Token
type Token int

const (
	Invalid Token = iota

	Id     // x y z + - / ++
	Int    // 12345
	Float  // 123.45
	String // "abc"

	LParen // (
	RParen // )
)

// TokenError implements an error at a specified line and column.
type TokenError struct {
	Cause     error
	Line, Col int
	Pos       Pos
}

func (t *TokenError) Error() string {
	return fmt.Sprintf("%v: at line %d: col %d", t.Cause, t.Line, t.Col)
}

func isExprOrStr(r rune) bool {
	return r == '(' || r == ')' || r == '"'
}

func IsSymbol(r rune) bool {
	return isExprOrStr(r) &&
		(unicode.IsPunct(r) || unicode.IsSymbol(r))
}

func IsIdent(r rune) bool {
	return isExprOrStr(r) &&
		(unicode.IsLetter(r) || unicode.IsMark(r))
}

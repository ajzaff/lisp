package scan

import (
	"fmt"
	"unicode"
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

// Error returns the error and cause.
func (t *TokenError) Error() string {
	line, col := t.LineCol()
	if t.Cause == nil {
		return fmt.Sprintf("TokenError: line %d, col %d", line, col)
	}
	return fmt.Sprintf("TokenError: line %d, col %d: %v", line, col, t.Cause)
}

// Unwrap implements error unwrapping, returning Cause.
func (t *TokenError) Unwrap() error {
	return t.Cause
}

package lisp

import (
	"strconv"

	"github.com/ajzaff/lisp"
)

// Id constructs an Id Node from the text.
func Id(text string) lisp.Lit { return lisp.Lit(text) }

// Nat constructs a Nat Node from the unsigned integer.
func Nat(i uint64) lisp.Lit {
	x := strconv.FormatUint(i, 10)
	return Id(x)
}

package lisp

import (
	"strconv"

	"github.com/ajzaff/lisp"
)

// Id constructs an Id Node from the text.
func Id(text string) lisp.Lit {
	return lisp.Lit{Token: lisp.Id, Text: text}
}

// Nat constructs a Nat Node from the unsigned integer.
func Nat(i uint64) lisp.Lit {
	return lisp.Lit{Token: lisp.Id, Text: strconv.FormatUint(i, 10)}
}

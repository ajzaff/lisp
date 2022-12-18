package lisp

import (
	"regexp"
)

// Quickly scan for invalid Lisp using a regular expression.

var badPattern = regexp.MustCompile(`[^\s]`)

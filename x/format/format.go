package format

import (
	"github.com/ajzaff/lisp/lisp"
)

// Source formats src Lisp code in-place.
//
// The returned slice is a formatted slice of the input.
func Source(src []byte) []byte {
	var space int
	var i int
	var delimClass lisp.DelimClass
	for j := 0; j < len(src); j++ {
		class := lisp.DelimByte(src[j])
		src[i] = src[j]
		switch src[j] {
		case '\t', '\n', ' ':
			space++
			if space < 2 && class != delimClass && class != lisp.DelimNone {
				i++
			}
		default:
			space = 0
			i++
		}
		if class != lisp.DelimNone {
			delimClass = class
		}
	}
	return src[:i]
}

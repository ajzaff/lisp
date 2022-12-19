package format

import (
	"github.com/ajzaff/lisp/lisputil"
)

// Source formats src Lisp code in-place.
//
// The returned slice is a formatted slice of the input.
func Source(src []byte) []byte {
	var space int
	var i int
	var delimClass lisputil.DelimClass
	for j := 0; j < len(src); j++ {
		class := lisputil.DelimByte(src[j])
		src[i] = src[j]
		switch src[j] {
		case '\t', '\n', ' ':
			space++
			if space < 2 && class != delimClass && class != lisputil.DelimNone {
				i++
			}
		default:
			space = 0
			i++
		}
		if class != lisputil.DelimNone {
			delimClass = class
		}
	}
	return src[:i]
}

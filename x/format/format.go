package format

// Source formats src Lisp code.
//
// Source will preserve one space between Id and Nat tokens but will not add them if not present.
// The returned slice is a formatted slice of the input.
func Source(src []byte) []byte {
	var space int
	var i int
	var delimClass delimClass
	for j := 0; j < len(src); j++ {
		class := delimByte(src[j])
		src[i] = src[j]
		switch src[j] {
		case '\t', '\r', '\n', ' ':
			space++
			if space < 2 && class != delimClass && class != delimNone {
				i++
			}
		default:
			delimClass = class
			space = 0
			i++
		}
	}
	return src[:i]
}

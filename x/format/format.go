package format

// Source formats src Lisp code.
//
// Source will preserve one space between Id and Nat tokens but will not add them if not present.
// The returned slice is a formatted slice of the input.
func Source(src []byte) []byte {
	var i int
	var delim bool
	for j := 0; j < len(src); j++ {
		src[i] = src[j]
		switch src[j] {
		case ' ', '\t', '\r', '\n':
			if delim {
				i++
				delim = false
			}
		case '(', ')':
			delim = false
			i++
		default:
			delim = true
			i++
		}
	}
	return src[:i]
}

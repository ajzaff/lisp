// Code generated by "stringer -type Token"; DO NOT EDIT.

package innit

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Invalid-0]
	_ = x[Id-1]
	_ = x[Int-2]
	_ = x[Float-3]
	_ = x[String-4]
	_ = x[LParen-5]
	_ = x[RParen-6]
}

const _Token_name = "InvalidIdIntFloatStringLParenRParen"

var _Token_index = [...]uint8{0, 7, 9, 12, 17, 23, 29, 35}

func (i Token) String() string {
	if i < 0 || i >= Token(len(_Token_index)-1) {
		return "Token(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Token_name[_Token_index[i]:_Token_index[i+1]]
}
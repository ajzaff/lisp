package lisp

const Space = " \t\r\n"

const Group = "()"

const Digit = "0123456789"

func IsSpace(b byte) bool {
	switch b {
	case ' ', '\t', '\r', '\n':
		return true
	default:
		return false
	}
}

func IsGroup(b byte) bool {
	switch b {
	case '(', ')':
		return true
	default:
		return false
	}
}

func IsNat(b byte) bool { return '0' <= b && b <= '9' }

func IsOther(b byte) bool { return !IsSpace(b) && !IsGroup(b) && !IsNat(b) }

func IsTokenBound(b byte) bool { return IsSpace(b) || IsGroup(b) }

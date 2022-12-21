package lisp

import (
	"strconv"
	"strings"
	"unicode"
)

// String returns Lisp string representation of the Lit.
//
// String falls back to GoString if the Lit is not valid to avoid conflating it with valid representations.
func (x Lit) String() string {
	var sb strings.Builder
	if !x.appendString(&sb) {
		// The Lit appears to be invalid.
		// Fall back to GoString instead.
		sb.Reset()
		x.appendGoString(&sb)
	}
	return sb.String()
}

func (x Lit) appendString(sb *strings.Builder) (valid bool) {
	if len(x.Text) == 0 {
		// Text is empty.
		// We shouldn't print this directly.
		return false
	}
	switch x.Token {
	case Id:
		for _, r := range x.Text {
			if !unicode.IsLetter(r) {
				// Lit is not valid.
				// We shouldn't print this directly.
				return false
			}
		}
	case Nat:
		if x.Text[0] == '0' && len(x.Text) > 1 {
			// Nat is not valid.
			// We shouldn't print this directly.
			return false
		}
		for _, b := range []byte(x.Text[1:]) {
			if b < '0' || '9' < b {
				// Nat is not valid.
				// We shouldn't print this directly.
				return false
			}
		}
	default:
		// Token is not obviously valid.
		// We shouldn't print this directly.
		return false
	}
	sb.WriteString(x.Text)
	return true
}

// String returns the Lisp string representation of this Cons.
//
// String falls back to GoString if the Cons is not valid to avoid conflating it with valid representations.
func (x *Cons) String() string {
	var sb strings.Builder
	if valid := x.appendString(&sb, true, false); !valid {
		// The Cons appears to be invalid.
		// Fall back to GoString instead.
		sb.Reset()
		x.appendGoString(&sb)
	}
	return sb.String()
}

func (x *Cons) appendString(sb *strings.Builder, first, delim bool) (valid bool) {
	if first {
		sb.WriteByte('(')
	}
	if x == nil {
		sb.WriteByte(')')
		return true
	}
	if x.Val == nil {
		if x.Cons != nil {
			// Cons is missing the Val element.
			// We shouldn't print this directly.
			return false
		}
	} else if cons, ok := x.Val.(*Cons); ok {
		cons.appendString(sb, true, false)
		delim = false
	} else if x != nil {
		if delim {
			sb.WriteByte(' ')
		}
		if !x.Val.(Lit).appendString(sb) {
			// The Lit in this Cons is invalid.
			// We shouldn't print this directly.
			return false
		}
		delim = true
	}
	return x.Cons.appendString(sb, false, delim)
}

// GoString returns the formatted GoString for this Lit.
func (x Lit) GoString() string {
	var sb strings.Builder
	x.appendGoString(&sb)
	return sb.String()
}

func (x Lit) appendGoString(sb *strings.Builder) {
	sb.WriteString("lisp.Lit{Token:")
	sb.WriteString(strconv.Itoa(int(x.Token)))
	sb.WriteString(", Text:")
	sb.WriteString(strconv.Quote(x.Text))
	sb.WriteByte('}')
}

func (x *Cons) GoString() string {
	var sb strings.Builder
	x.appendGoString(&sb)
	return sb.String()
}

func (x *Cons) appendGoString(sb *strings.Builder) {
	if x == nil {
		sb.WriteString("(*lisp.Cons)(nil)")
		return
	}
	sb.WriteString("&lisp.Cons{Val:")
	if x.Val == nil {
		sb.WriteString("(lisp.Val)(nil), Cons:")
	} else if cons, ok := x.Val.(*Cons); ok {
		cons.appendGoString(sb)
		sb.WriteString(", Cons:")
	} else {
		x.Val.(Lit).appendGoString(sb)
		sb.WriteString(", Cons:")
	}
	x.Cons.appendGoString(sb)
	sb.WriteByte('}')
}

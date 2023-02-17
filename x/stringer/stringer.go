package stringer

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/ajzaff/lisp"
)

// Val returns Lisp string representation of the Val.
//
// Val falls back to GoString when not valid to avoid conflating it with valid representations.
func Val(x lisp.Val) string {
	switch x := x.(type) {
	case lisp.Lit:
		return Lit(x)
	case *lisp.Cons:
		return Cons(x)
	case nil:
		return "<nil>"
	default:
		// Value Unknown type.
		return fmt.Sprintf("<Vunk>(%#v)", x)
	}
}

// Lit returns Lisp string representation of the Lit.
//
// Lit falls back to GoString if the Lit is not valid to avoid conflating it with valid representations.
func Lit(x lisp.Lit) string {
	var sb strings.Builder
	if !appendLit(x, &sb) {
		// The Lit appears to be invalid.
		// Fall back to GoString instead.
		sb.Reset()
		appendGoLit(x, &sb)
	}
	return sb.String()
}

func appendLit(x lisp.Lit, sb *strings.Builder) (valid bool) {
	if len(x.Text) == 0 {
		// Text is empty.
		// We shouldn't print this directly.
		return false
	}
	switch x.Token {
	case lisp.Id:
		for _, r := range x.Text {
			if !unicode.IsLetter(r) {
				// Lit is not valid.
				// We shouldn't print this directly.
				return false
			}
		}
	case lisp.Nat:
		if x.Text[0] == '0' && len(x.Text) > 1 {
			// Nat is not valid.
			// We shouldn't print this directly.
			return false
		}
		for _, b := range []byte(x.Text) {
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

// Cons returns the Lisp string representation of this Cons.
//
// Cons falls back to GoString if the Cons is not valid to avoid conflating it with valid representations.
func Cons(x *lisp.Cons) string {
	var sb strings.Builder
	if valid := appendCons(x, &sb, true, false); !valid {
		// The Cons appears to be invalid.
		// Fall back to GoString instead.
		sb.Reset()
		appendGoCons(x, &sb)
	}
	return sb.String()
}

func appendCons(x *lisp.Cons, sb *strings.Builder, first, delim bool) (valid bool) {
	if first {
		sb.WriteByte('(')
	}
	if x == nil {
		sb.WriteByte(')')
		return true
	}
	if x.Val == nil {
		if !first || x.Cons != nil {
			// Cons is missing the Val element.
			// We shouldn't print this directly.
			return false
		}
	} else if cons, ok := x.Val.(*lisp.Cons); ok {
		appendCons(cons, sb, true, false)
		delim = false
	} else if x != nil {
		if delim {
			sb.WriteByte(' ')
		}
		if !appendLit(x.Val.(lisp.Lit), sb) {
			// The Lit in this Cons is invalid.
			// We shouldn't print this directly.
			return false
		}
		delim = true
	}
	return appendCons(x.Cons, sb, false, delim)
}

func appendGoLit(x lisp.Lit, sb *strings.Builder) {
	sb.WriteString("lisp.Lit{Token:")
	sb.WriteString(strconv.Itoa(int(x.Token)))
	sb.WriteString(", Text:")
	sb.WriteString(strconv.Quote(x.Text))
	sb.WriteByte('}')
}

func appendGoCons(x *lisp.Cons, sb *strings.Builder) {
	if x == nil {
		sb.WriteString("(*lisp.Cons)(nil)")
		return
	}
	sb.WriteString("&lisp.Cons{Val:")
	if x.Val == nil {
		sb.WriteString("(lisp.Val)(nil), Cons:")
	} else if cons, ok := x.Val.(*lisp.Cons); ok {
		appendGoCons(cons, sb)
		sb.WriteString(", Cons:")
	} else {
		appendGoLit(x.Val.(lisp.Lit), sb)
		sb.WriteString(", Cons:")
	}
	appendGoCons(x.Cons, sb)
	sb.WriteByte('}')
}

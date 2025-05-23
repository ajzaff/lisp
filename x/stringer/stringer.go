package stringer

import (
	"fmt"
	"strings"

	"github.com/ajzaff/lisp"
	xlisp "github.com/ajzaff/lisp/x/lisp"
)

// Val returns Lisp string representation of the Val.
//
// Val falls back to GoString when not valid to avoid conflating it with valid representations.
func Val(x lisp.Val) string {
	switch x := x.(type) {
	case lisp.Lit:
		return Lit(x)
	case lisp.Group:
		return Group(x)
	case nil:
		return "<nil>"
	default:
		// Value Unknown type.
		return fmt.Sprintf("<Vunk>(%#v)", x)
	}
}

func appendVal(x lisp.Val, sb *strings.Builder, delim bool) (valid bool) {
	switch x := x.(type) {
	case lisp.Lit:
		return appendLit(x, sb, delim)
	case lisp.Group:
		return appendGroup(x, sb)
	case nil:
		sb.WriteString("<nil>")
		return true
	default:
		// Value Unknown type.
		fmt.Fprintf(sb, "<Vunk>(%#v)", x)
		return true
	}
}

// Lit returns Lisp string representation of the Lit.
//
// Lit falls back to GoString if the Lit is not valid to avoid conflating it with valid representations.
func Lit(x lisp.Lit) string {
	var sb strings.Builder
	if !appendLit(x, &sb, true) {
		// The Lit appears to be invalid.
		// Fall back to GoString instead.
		sb.Reset()
		appendGoLit(x, &sb)
	}
	return sb.String()
}

func appendLit(x lisp.Lit, sb *strings.Builder, delim bool) (valid bool) {
	if len(x) == 0 {
		// Text is empty.
		// We shouldn't print this directly.
		return false
	}
	if !delim {
		sb.WriteByte(' ')
	}
	for _, r := range x {
		if !xlisp.IsLit(r) {
			// Lit is not valid.
			// We shouldn't print this directly.
			return false
		}
	}
	sb.WriteString(string(x))
	return true
}

// Group returns the Lisp string representation of this Group.
//
// Group falls back to GoString if theGroups is not valid to avoid conflating it with valid representations.
func Group(x lisp.Group) string {
	var sb strings.Builder
	if valid := appendGroup(x, &sb); !valid {
		// The Group appears to be invalid.
		// Fall back to GoString instead.
		sb.Reset()
		appendGoGroup(x, &sb)
	}
	return sb.String()
}

func appendGroup(x lisp.Group, sb *strings.Builder) (valid bool) {
	sb.WriteByte('(')
	delim := true
	for _, e := range x {
		if !appendVal(e, sb, delim) {
			return false
		}
		switch e.(type) {
		case lisp.Lit:
			delim = false
		case lisp.Group:
			delim = true
		}
	}
	sb.WriteByte(')')
	return true
}

func appendGoVal(x lisp.Val, sb *strings.Builder) {
	switch x := x.(type) {
	case lisp.Lit:
		appendGoLit(x, sb)
	case lisp.Group:
		appendGoGroup(x, sb)
	case nil:
		sb.WriteString("<nil>")
	default:
		// Value Unknown type.
		fmt.Fprintf(sb, "<Vunk>(%#v)", x)
	}
}

func appendGoLit(x lisp.Lit, sb *strings.Builder) { fmt.Fprintf(sb, "lisp.Lit(%q)", x) }

func appendGoGroup(x lisp.Group, sb *strings.Builder) {
	if x == nil {
		sb.WriteString("(lisp.Group)(nil)")
		return
	}
	sb.WriteString("lisp.Group{")
	for _, e := range x {
		appendGoVal(e, sb)
	}
	sb.WriteByte('}')
}

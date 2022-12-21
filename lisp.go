// Package lisp implements minimal LISP-like expressions.
package lisp

// Token is an enumeration which specifies a kind of AST token.
type Token int

const (
	Invalid Token = iota
	Id            // abc. Strings of valid runes having the unicode Letter property.
	Nat           // 123. Natural number (plus 0), unsigned.
	LParen        // (
	RParen        // )
)

// Pos represents a byte position in a source file.
type Pos int

// NoPos is the canonical value for no position defined.
const NoPos Pos = -1

// Val is an interface for Lisp Values.
//
// Only allowed types are Lit and *Cons.
type Val interface {
	val()
}

// Cons is a construct used to build linked lists.
//
// It maintains pointers to a Val and the next Cons in the list.
type Cons struct {
	Val
	*Cons
}

// Lit is a basic literal type.
//
// Allowed Token types are Id, Number.
type Lit struct {
	Token Token
	Text  string
}

func (Lit) val()             {}
func (x Lit) String() string { return x.Text }

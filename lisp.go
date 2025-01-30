// Package lisp implements minimal LISP expressions.
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

// Val is an interface for Lisp Values.
//
// Only allowed types are Lit, Group, and Seq.
// Examples:
//
//	abc
//	123
//	abc 123
//	( abc 123 () )
type Val interface {
	val()
}

// Lit is a basic literal type.
//
// Allowed Token types are Id, Number.
// Examples:
//
//	abc
//	123
type Lit struct {
	Token Token
	Text  string
}

func (Lit) val() {}

// Group is a construct which groups a sequence of Vals together
// between parens.
//
// Example:
//
//	( abc 123 () )
type Group []Val

func (Group) val() {}

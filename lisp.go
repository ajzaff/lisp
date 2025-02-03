// Package lisp implements minimal LISP expressions.
package lisp

// Token is an enumeration which specifies a kind of AST token.
type Token int

const (
	Invalid Token = iota
	// Id tokens comprise unicode Letters, arabic numerals, plus underscore.
	//
	// Examples:
	//
	//	abc
	//	123
	//	abc_123
	Id
	LParen // (
	RParen // )
)

// Val is an interface for Lisp Values.
//
// Only allowed types are Lit, Group, and Seq.
// Examples:
//
//	abc
//	123
//	( abc 123 () )
type Val interface {
	val()
}

// Lit is a text identifier.
//
// Examples:
//
//	abc
//	123
//	abc_123
type Lit string

func (Lit) val() {}

// Group is a construct which groups a sequence of Vals together
// between parens.
//
// Example:
//
//	( abc 123 () )
type Group []Val

func (Group) val() {}

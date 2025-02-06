// Package lisp implements minimal LISP expressions.
package lisp

// Token is an enumeration which specifies a kind of symbol in Lisp.
type Token int

const (
	Invalid Token = iota
	// Id tokens comprise strings of consecutive unicode letters and digits.
	//
	// Examples:
	//
	//	abc
	//	123
	//	abc123
	Id
	LParen // (
	RParen // )
)

// Val is a closed interface for Lisp Values.
//
// The only allowed types are Lit or Group.
//
// Examples:
//
//	abc
//	123
//	( abc 123 () )
type Val interface {
	val()
}

// Lit is a text identifier comprising strings of consecutive unicode letters and digits.
//
// Examples:
//
//	abc
//	123
//	abc123
type Lit string

func (Lit) val() {}

// Group is a construct which encloses a sequence of Lisp values between parens.
//
// Example:
//
//	( )
//	( abc )
//	( abc 123 () )
type Group []Val

func (Group) val() {}

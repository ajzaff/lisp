// Package innit implements a simple LISP-like AST useful as a research
// language or wherever a bare-minimum language is required.
//
// It supports basic identifiers, numbers, strings, and expressions.
package innit

// Pos defines a position in the slice of code runes.
type Pos int

// NoPos is the flag value when no position is defined.
const NoPos = -1

// Node is an interface for AST nodes which have a start and end position.
type Node interface {
	Pos() Pos
	End() Pos
}

// Lit is a basic literal node.
//
// Lit can hold any token value. See token.go for more details.
type Lit struct {
	Tok      Token
	ValuePos Pos
	Value    string
}

// Expr defines a NodeList enclosed by parens.
type Expr struct {
	LParen Pos
	X      NodeList
	RParen Pos
}

func (x *Lit) Pos() Pos  { return x.ValuePos }
func (x *Lit) End() Pos  { return x.ValuePos + Pos(len(x.Value)) }
func (x *Expr) Pos() Pos { return x.LParen }
func (x *Expr) End() Pos { return x.RParen + 1 }

// NodeList defines a slice of nodes.
//
// The start and end positions correspond to the start and end positions
// of the first and last nodes respectively, otherwise, NoPos is used.
type NodeList []Node

func (x NodeList) Pos() Pos {
	if len(x) > 0 {
		return x[0].Pos()
	}
	return NoPos
}

func (x NodeList) End() Pos {
	if n := len(x); n > 0 {
		return x[n-1].End()
	}
	return NoPos
}

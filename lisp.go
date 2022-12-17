// Package lisp implements a minimal LISP-like expressions useful as a research
// language or wherever a bare-minimum language is required.
//
// The code here is split up into xxxNode types which implement Node and
// Value types which implement Val. Node types represent AST with source
// position information while Value types only contain context-free values.
//
// It supports basic identifiers, numbers, and compositional expressions.
package lisp

// Pos defines a position in the slice of code runes.
type Pos int

// NoPos is the flag value when no position is defined.
const NoPos = -1

// Val is an interface for Lisp Values.
//
// Only allowed types are Lit and Expr.
type Val interface {
	val()
}

// Expr is a slice of compound Nodes.
type Expr []Node

func (Expr) val() {}

// Lit is a basic literal type.
//
// Allowed Token types are Id, Number.
type Lit struct {
	Token Token
	Text  string
}

func (Lit) val()             {}
func (x Lit) String() string { return x.Text }

// Node represents an AST node in context with source file indices.
//
// See the Val interface for allowed AST types.
type Node struct {
	Pos Pos
	Val Val
	End Pos
}

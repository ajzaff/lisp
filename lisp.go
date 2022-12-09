// Package lisp implements a minimal LISP-like expressions useful as a research
// language or wherever a bare-minimum language is required.
//
// The code here is split up into xxxNode types which implement Node and
// Value types which implement Val. Node types represent AST with source
// position information while Value types only contain context-free values.
//
// It supports basic identifiers, numbers, strings, and expressions.
package lisp

// Pos defines a position in the slice of code runes.
type Pos int

// NoPos is the flag value when no position is defined.
const NoPos = -1

// Node is an interface for AST nodes which have a start and end source position.
//
// Only allowed types are *LitNode and *ExprNode.
//
// See Val for a context-free Value type.
type Node interface {
	Pos() Pos
	Val() Val
	End() Pos
}

// Val is an interface for Lisp Values.
//
// Only allowed types are IdLit, IntLit, FloatLit, StringLit and Expr.
type Val interface {
	val()
}

// Expr is a slice of compound Nodes.
type Expr []Node

func (Expr) val() {}

// LitNode is a basic literal node.
type LitNode struct {
	LitPos Pos
	Lit    Lit
	EndPos Pos
}

// ExprNode is a Expr enclosed by parens.
type ExprNode struct {
	LParen Pos
	Expr   Expr
	RParen Pos
}

func (x *LitNode) Pos() Pos  { return x.LitPos }
func (x *LitNode) Val() Val  { return x.Lit }
func (x *LitNode) End() Pos  { return x.EndPos }
func (x *ExprNode) Pos() Pos { return x.LParen }
func (x *ExprNode) Val() Val { return x.Expr }
func (x *ExprNode) End() Pos { return x.RParen + 1 }

// Lit is a basic literal type.
//
// Allowed type could be Id, Number, or String.
type Lit struct {
	Token
	Text string
}

func (Lit) val()             {}
func (x Lit) String() string { return string(x.Text) }

package lisp

import (
	"errors"
)

// Errors used to flag conditions in the Visitor.
var (
	errStop = errors.New("stop")
	errSkip = errors.New("skip")
)

// Visitor implements a Val visitor.
type Visitor struct {
	beforeValFn func(Val)
	afterValFn  func(Val)
	litFn       func(Lit)

	beforeExprFn func(Expr)
	afterExprFn  func(Expr)

	err error
}

// SetBeforeNodeVisitor sets the visitor called on every Node.
func (v *Visitor) SetBeforeValVisitor(fn func(Val)) { v.beforeValFn = fn }

// SetAfterNodeVisitor sets the visitor called on every Node.
func (v *Visitor) SetAfterValVisitor(fn func(Val)) { v.afterValFn = fn }

// SetLitVisitor sets the visitor called on every *Lit.
func (v *Visitor) SetLitVisitor(fn func(Lit)) { v.litFn = fn }

// SetNodeListVisitFunc sets the visitor called on every *Expr.
func (v *Visitor) SetBeforeExprVisitor(fn func(Expr)) { v.beforeExprFn = fn }

// SetNodeListVisitFunc sets the visitor called on every *Expr.
func (v *Visitor) SetAfterExprVisitor(fn func(Expr)) { v.afterExprFn = fn }

// Stop the visitor and return as soon as possible.
func (v *Visitor) Stop() {
	v.err = errStop
}

// Skip visiting the node recursively for compound nodes.
func (v *Visitor) Skip() {
	v.err = errSkip
}

// Visit the node recursively while calling visitor functions.
//
// Visit continues until all nodes are visited or Stop is called.
func (v *Visitor) Visit(root Val) {
	if root == nil {
		return
	}
	defer v.clearSkipErr()
	if v.hasErr() {
		return
	}
	if !callFn(v, v.beforeValFn, root) {
		return
	}
	defer callFn(v, v.afterValFn, root)
	switch x := root.(type) {
	case Lit:
		if !callFn(v, v.litFn, x) {
			return
		}
	case Expr:
		if !callFn(v, v.beforeExprFn, x) {
			return
		}
		defer callFn(v, v.afterExprFn, x)
		for _, e := range x {
			v.Visit(e.Val)
		}
	}
}

func (v *Visitor) hasErr() bool {
	return v.err != nil
}

func (v *Visitor) clearSkipErr() bool {
	if v.err == errSkip {
		v.err = nil
		return true
	}
	return false
}

func callFn[T Val](v *Visitor, fn func(T), e T) (ok bool) {
	if fn == nil {
		return true
	}
	fn(e)
	return !v.hasErr()
}

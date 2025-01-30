// Package visit provides Visitor which implements an efficient node visitor algorithm.
package visit

import (
	"errors"

	"github.com/ajzaff/lisp"
)

// Errors used to flag conditions in the Visitor.
var (
	errStop = errors.New("stop")
	errSkip = errors.New("skip")
)

// Visitor implements a Val visitor.
type Visitor struct {
	valFn         func(lisp.Val)
	litFn         func(lisp.Lit)
	beforeGroupFn func(lisp.Group)
	afterGroupFn  func(lisp.Group)

	err error
}

// SetValVisitor sets the visitor called on every Val.
func (v *Visitor) SetValVisitor(fn func(lisp.Val)) { v.valFn = fn }

// SetLitVisitor sets the visitor called on every Lit.
func (v *Visitor) SetLitVisitor(fn func(lisp.Lit)) { v.litFn = fn }

// SetBeforeGroupVisitor sets the visitor called on every Group before recursing on its elements.
func (v *Visitor) SetBeforeGroupVisitor(fn func(lisp.Group)) { v.beforeGroupFn = fn }

// SetAfterGroupVisitor sets the visitor called on every Group before recursing on its elements.
func (v *Visitor) SetAfterGroupVisitor(fn func(lisp.Group)) { v.afterGroupFn = fn }

// Stop stops the visitor and returns as soon as possible.
func (v *Visitor) Stop() {
	v.err = errStop
}

// Skip skips descending the next Val of a Group.
func (v *Visitor) Skip() {
	v.err = errSkip
}

// Visit the Val recursively while calling visitor functions.
//
// Visit continues in-order, descending Group links, or until Stop is called.
// Calling Skip will cause the next Val to not be descended.
func (v *Visitor) Visit(root lisp.Val) {
	if root == nil {
		return
	}
	defer v.clearSkipErr()
	if v.hasErr() {
		return
	}
	if !callFn(v, v.valFn, root) {
		return
	}
	switch x := root.(type) {
	case lisp.Lit:
		callFn(v, v.litFn, x)
	case lisp.Group:
		v.visitGroup(x)
	}
}

// VisitGroup visits the Group recursively.
func (v *Visitor) VisitGroup(root lisp.Group) {
	if !callFn[lisp.Val](v, v.valFn, root) {
		return
	}
	v.visitGroup(root)
}

func (v *Visitor) visitGroup(root lisp.Group) {
	if !callFn(v, v.beforeGroupFn, root) {
		return
	}
	for _, e := range root {
		if v.Visit(e); v.hasErr() {
			break
		}
	}
	callFn(v, v.afterGroupFn, root)
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

func callFn[T lisp.Val](v *Visitor, fn func(T), e T) (ok bool) {
	if fn == nil {
		return true
	}
	fn(e)
	return !v.hasErr()
}

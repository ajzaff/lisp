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
	valFn func(lisp.Val)
	litFn func(lisp.Lit)

	beforeConsFn func(*lisp.Cons)
	consFn       func(*lisp.Cons)
	afterConsFn  func(*lisp.Cons)

	err error
}

// SetValVisitor sets the visitor called on every Val.
func (v *Visitor) SetValVisitor(fn func(lisp.Val)) { v.valFn = fn }

// SetLitVisitor sets the visitor called on every Lit.
func (v *Visitor) SetLitVisitor(fn func(lisp.Lit)) { v.litFn = fn }

// SetBeforeConsVisitor sets the visitor called on the first Cons.
func (v *Visitor) SetBeforeConsVisitor(fn func(*lisp.Cons)) { v.beforeConsFn = fn }

// SetConsVisitor sets the visitor called on every cons before descending the Val.
func (v *Visitor) SetConsVisitor(fn func(*lisp.Cons)) { v.consFn = fn }

// SetAfterConsVisitor sets the visitor called on the last Cons.
func (v *Visitor) SetAfterConsVisitor(fn func(*lisp.Cons)) { v.afterConsFn = fn }

// Stop stops the visitor and returns as soon as possible.
func (v *Visitor) Stop() {
	v.err = errStop
}

// Skip skips descending the next Val of a Cons.
func (v *Visitor) Skip() {
	v.err = errSkip
}

// Visit the Val recursively while calling visitor functions.
//
// Visit continues in-order, descending Cons links, or until Stop is called.
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
		if !callFn(v, v.litFn, x) {
			return
		}
	case *lisp.Cons:
		v.visitConsNoVal(x)
	}
}

// VisitCons visits the Cons recursively.
func (v *Visitor) VisitCons(root *lisp.Cons) {
	if !callFn[lisp.Val](v, v.valFn, root) {
		return
	}
	v.visitConsNoVal(root)
}

// visitConsNoVal precondition: valFn(root) called.
func (v *Visitor) visitConsNoVal(root *lisp.Cons) {
	if !callFn(v, v.beforeConsFn, root) {
		return
	}
	v.visitCons(root)
}

// visitCons precondition: root != nil.
func (v *Visitor) visitCons(root *lisp.Cons) {
	if !callFn(v, v.consFn, root) {
		return
	}
	if root != nil {
		v.Visit(root.Val)
	}
	if root == nil || root.Cons == nil {
		callFn(v, v.afterConsFn, root)
		return
	}
	if !callFn[lisp.Val](v, v.valFn, root.Cons) {
		return
	}
	v.visitCons(root.Cons)
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

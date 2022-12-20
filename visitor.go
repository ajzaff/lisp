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
	valFn func(Val)
	litFn func(Lit)

	beforeConsFn func(*Cons)
	consFn       func(*Cons)
	afterConsFn  func(*Cons)

	err error
}

// SetValVisitor sets the visitor called on every Val.
func (v *Visitor) SetValVisitor(fn func(Val)) { v.valFn = fn }

// SetLitVisitor sets the visitor called on every Lit.
func (v *Visitor) SetLitVisitor(fn func(Lit)) { v.litFn = fn }

// SetBeforeConsVisitor sets the visitor called on the first Cons.
func (v *Visitor) SetBeforeConsVisitor(fn func(*Cons)) { v.beforeConsFn = fn }

// SetConsVisitor sets the visitor called on every cons before descending the Val.
func (v *Visitor) SetConsVisitor(fn func(*Cons)) { v.consFn = fn }

// SetAfterConsVisitor sets the visitor called on the last Cons.
func (v *Visitor) SetAfterConsVisitor(fn func(*Cons)) { v.afterConsFn = fn }

// Stop the visitor and return as soon as possible.
func (v *Visitor) Stop() {
	v.err = errStop
}

// Skip visiting the node recursively for compound nodes.
func (v *Visitor) Skip() {
	v.err = errSkip
}

// Visit the Val recursively while calling visitor functions.
//
// Visit continues in-order, descending Cons links, or until Stop is called.
// Calling Skip will cause the next Val to not be descended.
func (v *Visitor) Visit(root Val) {
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
	case Lit:
		if !callFn(v, v.litFn, x) {
			return
		}
	case *Cons:
		v.VisitCons(x)
	}
}

// VisitCons visits the Cons recursively.
func (v *Visitor) VisitCons(root *Cons) {
	if !callFn(v, v.beforeConsFn, root) {
		return
	}
	v.visitCons(root)
}

func (v *Visitor) visitCons(root *Cons) {
	if root == nil {
		return
	}
	if !callFn(v, v.consFn, root) {
		return
	}
	v.Visit(root.Val)
	if root.Cons == nil {
		callFn(v, v.afterConsFn, root)
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

func callFn[T Val](v *Visitor, fn func(T), e T) (ok bool) {
	if fn == nil {
		return true
	}
	fn(e)
	return !v.hasErr()
}

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

// SetBeforeConsVisitor sets the visitor called on every Cons.
func (v *Visitor) SetBeforeConsVisitor(fn func(*Cons)) { v.beforeConsFn = fn }

func (v *Visitor) SetConsVisitor(fn func(*Cons)) { v.consFn = fn }

// SetAfterConsVisitor sets the visitor called on every Cons after descending the Val.
func (v *Visitor) SetAfterConsVisitor(fn func(*Cons)) { v.afterConsFn = fn }

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
	if !callFn(v, v.valFn, root) {
		return
	}
	switch x := root.(type) {
	case Lit:
		if !callFn(v, v.litFn, x) {
			return
		}
	case *Cons:
		if !callFn(v, v.beforeConsFn, x) {
			return
		}
		e := x
		for ; e != nil; e = e.Cons {
			if !callFn(v, v.consFn, e) {
				return
			}
			v.Visit(e.Val)
		}
		if !callFn(v, v.afterConsFn, e) {
			return
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

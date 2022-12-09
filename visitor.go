package lisp

import (
	"errors"
	"fmt"
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

var (
	errStop = errors.New("stop")
	errSkip = errors.New("skip")
)

// Stop the visitor and return as soon as possible.
func (v *Visitor) Stop() {
	v.err = errStop
}

// Skip visiting the node recursively for compound nodes.
func (v *Visitor) Skip() {
	v.err = errSkip
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

func (v *Visitor) callBeforeValFn(e Val) bool {
	if v.beforeValFn != nil {
		v.beforeValFn(e)
		return true
	}
	return false
}

func (v *Visitor) callAfterValFn(e Val) bool {
	if v.afterValFn != nil {
		v.afterValFn(e)
		return true
	}
	return false
}

func (v *Visitor) callLitFn(e Lit) bool {
	if v.litFn != nil {
		v.litFn(e)
		return true
	}
	return false
}

func (v *Visitor) callBeforeExprFn(e Expr) bool {
	if v.beforeExprFn != nil {
		v.beforeExprFn(e)
		return true
	}
	return false
}

func (v *Visitor) callAfterExprFn(e Expr) bool {
	if v.afterExprFn != nil {
		v.afterExprFn(e)
		return true
	}
	return false
}

// Visit the node recursively while calling visitor functions.
//
// Visit continues until all nodes are visited or Stop is called.
func (v *Visitor) Visit(x Val) {
	if x == nil {
		return
	}
	if v.hasErr() {
		v.clearSkipErr()
		return
	}
	if v.callBeforeValFn(x) && v.hasErr() {
		v.clearSkipErr()
		return
	}
	defer func() {
		if v.callAfterValFn(x) && v.hasErr() {
			v.clearSkipErr()
		}
	}()
	switch x := x.(type) {
	case Lit:
		if v.callLitFn(x) && v.hasErr() {
			v.clearSkipErr()
			return
		}
	case Expr:
		if v.callBeforeExprFn(x) && v.hasErr() {
			v.clearSkipErr()
			return
		}
		defer func() {
			if v.callAfterExprFn(x) && v.hasErr() {
				v.clearSkipErr()
			}
		}()
		for _, e := range x {
			v.Visit(e.Val())
			if v.hasErr() {
				if v.clearSkipErr() {
					continue
				}
				return
			}
		}
	default: // unknown
		v.err = fmt.Errorf("unknown Val")
		return
	}
}

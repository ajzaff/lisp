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

	// Lit token type handlers.
	enableLitFn bool
	idFn        func(IdLit)
	numberFn    func(NumberLit)
	strFn       func(StringLit)

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

// SetIdVisitor sets the visitor called on every Id Lit.
func (v *Visitor) SetIdVisitor(fn func(IdLit)) { v.idFn = fn; v.enableLitFn = true }

// SetIntVisitor sets the visitor called on every Number Lit.
func (v *Visitor) SetNumberVisitor(fn func(NumberLit)) { v.numberFn = fn; v.enableLitFn = true }

// SetIntVisitor sets the visitor called on every String Lit.
func (v *Visitor) SetStringVisitor(fn func(StringLit)) { v.strFn = fn; v.enableLitFn = true }

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

func (v *Visitor) callIdFn(id IdLit) bool {
	if v.idFn != nil {
		v.idFn(id)
		return true
	}
	return false
}

func (v *Visitor) callNumberFn(e NumberLit) bool {
	if v.numberFn != nil {
		v.numberFn(e)
		return true
	}
	return false
}

func (v *Visitor) callStringFn(e StringLit) bool {
	if v.strFn != nil {
		v.strFn(e)
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
		if !v.enableLitFn {
			return
		}
		switch x := x.(type) {
		case IdLit:
			if v.callIdFn(x) && v.hasErr() {
				v.clearSkipErr()
			}
		case NumberLit:
			if v.callNumberFn(x) && v.hasErr() {
				v.clearSkipErr()
			}
		case StringLit:
			if v.callStringFn(x) && v.hasErr() {
				v.clearSkipErr()
			}
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

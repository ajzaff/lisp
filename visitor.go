package innit

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
)

// Visitor implements a Node visitor.
type Visitor struct {
	unknownFn        func(Node)
	beforeNodeFn     func(Node)
	afterNodeFn      func(Node)
	beforeNodeListFn func(NodeList)
	afterNodeListFn  func(NodeList)
	litFn            func(*Lit)

	// Lit token type handlers.
	enableLitFn bool
	idFn        func(string)
	intFn       func(*big.Int)
	floatFn     func(*big.Float)
	strFn       func(string)

	beforeExprFn func(*Expr)
	afterExprFn  func(*Expr)

	err error
}

// SetUnknownTypeVisitor sets the visitor called on unknown-typed Nodes.
func (v *Visitor) SetUnknownTypeVisitor(fn func(Node)) { v.unknownFn = fn }

// SetBeforeNodeVisitor sets the visitor called on every Node.
func (v *Visitor) SetBeforeNodeVisitor(fn func(Node)) { v.beforeNodeFn = fn }

// SetAfterNodeVisitor sets the visitor called on every Node.
func (v *Visitor) SetAfterNodeVisitor(fn func(Node)) { v.afterNodeFn = fn }

// SetNodeListVisitor sets the visitor called on every *NodeList.
func (v *Visitor) SetBeforeNodeListVisitor(fn func(NodeList)) { v.beforeNodeListFn = fn }

// SetNodeListVisitor sets the visitor called on every *NodeList.
func (v *Visitor) SetAfterNodeListVisitor(fn func(NodeList)) { v.afterNodeListFn = fn }

// SetLitVisitor sets the visitor called on every *Lit.
func (v *Visitor) SetLitVisitor(fn func(*Lit)) { v.litFn = fn }

// SetNodeListVisitFunc sets the visitor called on every *Expr.
func (v *Visitor) SetBeforeExprVisitor(fn func(*Expr)) { v.beforeExprFn = fn }

// SetNodeListVisitFunc sets the visitor called on every *Expr.
func (v *Visitor) SetAfterExprVisitor(fn func(*Expr)) { v.afterExprFn = fn }

// SetIdVisitor sets the visitor called on every Id Lit.
func (v *Visitor) SetIdVisitor(fn func(string)) { v.idFn = fn; v.enableLitFn = true }

// SetIntVisitor sets the visitor called on every Int Lit.
func (v *Visitor) SetIntVisitor(fn func(*big.Int)) { v.intFn = fn; v.enableLitFn = true }

// SetIntVisitor sets the visitor called on every Float Lit.
func (v *Visitor) SetFloatVisitor(fn func(*big.Float)) { v.floatFn = fn; v.enableLitFn = true }

// SetIntVisitor sets the visitor called on every String Lit.
func (v *Visitor) SetStringVisitor(fn func(string)) { v.strFn = fn; v.enableLitFn = true }

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

func (v *Visitor) callUnknownNodeFn(e Node) bool {
	if v.unknownFn != nil {
		v.unknownFn(e)
		return true
	}
	return false
}

func (v *Visitor) callBeforeNodeFn(e Node) bool {
	if v.beforeNodeFn != nil {
		v.beforeNodeFn(e)
		return true
	}
	return false
}

func (v *Visitor) callAfterNodeFn(e Node) bool {
	if v.afterNodeFn != nil {
		v.afterNodeFn(e)
		return true
	}
	return false
}

func (v *Visitor) callBeforeNodeListFn(e NodeList) bool {
	if v.beforeNodeListFn != nil {
		v.beforeNodeListFn(e)
		return true
	}
	return false
}

func (v *Visitor) callAfterNodeListFn(e NodeList) bool {
	if v.afterNodeListFn != nil {
		v.afterNodeListFn(e)
		return true
	}
	return false
}

func (v *Visitor) callLitFn(e *Lit) bool {
	if v.litFn != nil {
		v.litFn(e)
		return true
	}
	return false
}

func (v *Visitor) callBeforeExprFn(e *Expr) bool {
	if v.beforeExprFn != nil {
		v.beforeExprFn(e)
		return true
	}
	return false
}

func (v *Visitor) callAfterExprFn(e *Expr) bool {
	if v.afterExprFn != nil {
		v.afterExprFn(e)
		return true
	}
	return false
}

func (v *Visitor) callIdFn(id string) bool {
	if v.idFn != nil {
		v.idFn(id)
		return true
	}
	return false
}

func (v *Visitor) callIntFn(e string) bool {
	if v.intFn != nil {
		var i big.Int
		if _, ok := i.SetString(e, 10); !ok {
			v.err = fmt.Errorf("big.Int")
			return true
		}
		v.intFn(&i)
		return true
	}
	return false
}

func (v *Visitor) callFloatFn(e string) bool {
	if v.floatFn != nil {
		f, _, err := big.ParseFloat(e, 10, 64, big.ToZero)
		if err != nil {
			v.err = err
			return true
		}
		v.floatFn(f)
		return true
	}
	return false
}

func (v *Visitor) callStringFn(e string) bool {
	if v.strFn != nil {
		s, err := strconv.Unquote(e)
		if err != nil {
			v.err = err
			return true
		}
		v.strFn(s)
		return true
	}
	return false
}

// Visit the node recursively while calling visitor functions.
//
// Visit continues until all nodes are visited or Stop is called.
func (v *Visitor) Visit(node Node) {
	if v.hasErr() {
		v.clearSkipErr()
		return
	}
	if v.callBeforeNodeFn(node) && v.hasErr() {
		v.clearSkipErr()
		return
	}
	switch n := node.(type) {
	case *Lit:
		if v.callLitFn(n) && v.hasErr() {
			v.clearSkipErr()
			return
		}
		if !v.enableLitFn {
			return
		}
		switch n.Tok {
		case Id:
			if v.callIdFn(n.Value) && v.hasErr() {
				v.clearSkipErr()
			}
		case Int:
			if v.callIntFn(n.Value) && v.hasErr() {
				v.clearSkipErr()
			}
		case Float:
			if v.callFloatFn(n.Value) && v.hasErr() {
				v.clearSkipErr()
			}
		case String:
			if v.callStringFn(n.Value) && v.hasErr() {
				v.clearSkipErr()
			}
		}
	case *Expr:
		if v.callBeforeExprFn(n) && v.hasErr() {
			v.clearSkipErr()
			return
		}
		v.Visit(n.X)
		if v.callAfterExprFn(n) && v.hasErr() {
			v.clearSkipErr()
			return
		}
	case NodeList:
		if v.callBeforeNodeListFn(n) && v.hasErr() {
			v.clearSkipErr()
			return
		}
		for _, e := range n {
			v.Visit(e)
			if v.hasErr() {
				if v.clearSkipErr() {
					continue
				}
				return
			}
		}
		if v.callAfterNodeListFn(n) && v.hasErr() {
			v.clearSkipErr()
			return
		}
	default: // unknown
		if v.callUnknownNodeFn(n) && v.hasErr() {
			v.clearSkipErr()
		}
	}
	if v.callAfterNodeFn(node) && v.hasErr() {
		v.clearSkipErr()
		return
	}
}

package innit

import "errors"

// Visitor implements a Node visitor.
type Visitor struct {
	unknownFn    func(Node)
	beforeNodeFn func(Node)
	afterNodeFn  func(Node)
	nodeListFn   func(NodeList)
	litFn        func(*Lit)

	beforeExprFn func(*Expr)
	afterExprFn  func(*Expr)

	err error
}

// SetUnknownTypeVisitor sets the visitor called on unknown-typed Nodes.
func (v *Visitor) SetUnknownTypeVisitor(fn func(Node)) {
	v.unknownFn = fn
}

// SetBeforeNodeVisitor sets the visitor called on every Node.
func (v *Visitor) SetBeforeNodeVisitor(fn func(Node)) {
	v.beforeNodeFn = fn
}

// SetAfterNodeVisitor sets the visitor called on every Node.
func (v *Visitor) SetAfterNodeVisitor(fn func(Node)) {
	v.afterNodeFn = fn
}

// SetNodeListVisitor sets the visitor called on every *NodeList.
func (v *Visitor) SetNodeListVisitor(fn func(NodeList)) {
	v.nodeListFn = fn
}

// SetLitVisitor sets the visitor called on every *Lit.
func (v *Visitor) SetLitVisitor(fn func(*Lit)) {
	v.litFn = fn
}

// SetNodeListVisitFunc sets the visitor called on every *Expr.
func (v *Visitor) SetBeforeExprVisitor(fn func(*Expr)) {
	v.beforeExprFn = fn
}

// SetNodeListVisitFunc sets the visitor called on every *Expr.
func (v *Visitor) SetAfterExprVisitor(fn func(*Expr)) {
	v.afterExprFn = fn
}

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

func (v *Visitor) callNodeListFn(e NodeList) bool {
	if v.nodeListFn != nil {
		v.nodeListFn(e)
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
		if v.callNodeListFn(n) && v.hasErr() {
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

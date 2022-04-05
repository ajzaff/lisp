package innit

import "errors"

// Visitor implements a Node visitor.
type Visitor struct {
	unknownFn  func(Node)
	nodeFn     func(Node)
	nodeListFn func(NodeList)
	basicLitFn func(*BasicLit)
	exprFn     func(*Expr)

	err error
}

// SetNodeVisitFunc sets the visitor called on unknown-typed Nodes.
func (v *Visitor) SetUnknownNodeVisitFunc(fn func(Node)) {
	v.unknownFn = fn
}

// SetNodeVisitFunc sets the visitor called on every Node.
func (v *Visitor) SetNodeVisitFunc(fn func(Node)) {
	v.nodeFn = fn
}

// SetNodeListVisitFunc sets the visitor called on every *NodeList.
func (v *Visitor) SetNodeListVisitFunc(fn func(NodeList)) {
	v.nodeListFn = fn
}

// SetNodeListVisitFunc sets the visitor called on every *BasicLit.
func (v *Visitor) SetBasicLitVisitFunc(fn func(*BasicLit)) {
	v.basicLitFn = fn
}

// SetNodeListVisitFunc sets the visitor called on every *Expr.
func (v *Visitor) SetExprVisitFunc(fn func(*Expr)) {
	v.exprFn = fn
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

func (v *Visitor) callNodeFn(e Node) bool {
	if v.nodeFn != nil {
		v.nodeFn(e)
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

func (v *Visitor) callBasicLitFn(e *BasicLit) bool {
	if v.basicLitFn != nil {
		v.basicLitFn(e)
		return true
	}
	return false
}

func (v *Visitor) callExprFn(e *Expr) bool {
	if v.exprFn != nil {
		v.exprFn(e)
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
	if v.callNodeFn(node) && v.hasErr() {
		v.clearSkipErr()
		return
	}
	switch n := node.(type) {
	case *BasicLit:
		if v.callBasicLitFn(n) && v.hasErr() {
			v.clearSkipErr()
		}
	case *Expr:
		if v.callExprFn(n) && v.hasErr() {
			v.clearSkipErr()
			return
		}
		v.Visit(n.X)
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
}

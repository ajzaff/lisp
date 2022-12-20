package lisputil

import (
	"strconv"

	"github.com/ajzaff/lisp"
)

type nodeOptions struct {
	pos lisp.Pos
	end lisp.Pos
}

// NodeIndices used for setting Node indices in the Node function.
type NodeIndices func(*nodeOptions)

// Indices returns an option which sets Node Pos and End.
func Indices(pos, end lisp.Pos) NodeIndices {
	return func(opts *nodeOptions) {
		opts.pos = pos
		opts.end = end
	}
}

// Pos returns an option which sets the Node Pos.
func Pos(pos lisp.Pos) NodeIndices {
	return func(opts *nodeOptions) {
		opts.pos = pos
	}
}

// End returns an option which sets the Node End.
func End(end lisp.Pos) NodeIndices {
	return func(opts *nodeOptions) {
		opts.end = end
	}
}

// Id constructs an Id Node from the text and indices.
func Id(text string, indices ...NodeIndices) lisp.Node {
	var os nodeOptions
	for _, fn := range indices {
		fn(&os)
	}
	return lisp.Node{Pos: os.pos, Val: lisp.Lit{Token: lisp.Id, Text: text}, End: os.end}
}

// Int constructs an Int Node from the text and indices.
func Int(i uint64, indices ...NodeIndices) lisp.Node {
	var os nodeOptions
	for _, fn := range indices {
		fn(&os)
	}
	return lisp.Node{Pos: os.pos, Val: lisp.Lit{Token: lisp.Id, Text: strconv.FormatUint(i, 10)}, End: os.end}
}

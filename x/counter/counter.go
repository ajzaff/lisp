// Package counter implements efficient node counting algorithms.
package counter

import (
	"bufio"
	"io"

	"github.com/ajzaff/lisp"
)

// Counter counts the number of complete nodes per instance of the delimiter.
//
// Nodes spanning between delimiters will count as 0.
type NodeCounter struct {
	r     bufio.Reader
	delim byte

	depth int        // Current parser depth, preserved between Count calls.
	tok   lisp.Token // Current token, if any, preserved between Count calls.
}

// Reset the reader and delimiter.
func (c *NodeCounter) Reset(r io.Reader, delim byte) {
	c.r.Reset(r)
	c.delim = delim
	c.depth = 0
	c.tok = lisp.Invalid
}

func isSpace(b byte) bool {
	switch b {
	case ' ', '\t', '\r', '\n':
		return true
	default:
		return false
	}
}

// Count returns the number of complete nodes between the start and the next instance of delim.
func (c *NodeCounter) Count() (int, error) {
	line, err := c.r.ReadSlice(c.delim)
	complete := err == nil
	if !complete {
		if err == io.EOF || err == bufio.ErrBufferFull {
			err = nil
		}
		if err != nil {
			// Unhandled error.
			return 0, err
		}
	}
	// Skip spaces.
	off := 0
	for ; off < len(line) && isSpace(line[off]); off++ {
	}
	if off == len(line) {
		return 0, nil
	}

	// Count nodes.
	count := 0
	for _, b := range line {
		switch b {
		case '(':
			c.depth++
		case ')':
			if c.depth--; c.depth == 0 {
				count++
			}
		case ' ', '\t', '\r', '\n':

		}
		switch {
		case '0' <= b && b <= '9':
		}
	}
	if complete {
		c.depth = 0
		c.tok = lisp.Invalid
	}
	return count, nil
}

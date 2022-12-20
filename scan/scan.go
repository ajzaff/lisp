package scan

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"

	"github.com/ajzaff/lisp"
)

var errRune = errors.New("unexpected rune")

// Scanner scans the Lisp source for lisp.Tokens.
type TokenScanner struct {
	sc        *bufio.Scanner
	prev, off lisp.Pos   // absolute offsets
	pos       lisp.Pos   // relative lisp.Token lisp.Position
	tok       lisp.Token // last lisp.Token scanned
}

func (s *TokenScanner) Reset(r io.Reader) {
	*s = TokenScanner{}
	s.sc = bufio.NewScanner(r)
	s.sc.Split(s.scanTokens)
}

func (s *TokenScanner) Buffer(buf []byte, max int) {
	s.sc.Buffer(buf, max)
}

func (s *TokenScanner) Scan() bool {
	return s.sc.Scan()
}

func (s *TokenScanner) Err() error {
	return s.sc.Err()
}

func (s *TokenScanner) Text() string {
	return s.sc.Text()
}

func (s *TokenScanner) Token() (pos lisp.Pos, tok lisp.Token, text string) {
	return s.prev + s.pos, s.tok, s.Text()
}

func (s *TokenScanner) scanTokens(src []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(src) == 0 {
		return 0, nil, nil
	}
	defer func() {
		// Maintain absolute lisp.Positions into the scanner.
		s.prev = s.off
		s.off += lisp.Pos(advance)
	}()
	// Skip leading spaces.
	for size := 0; advance < len(src); advance += size {
		var r rune
		r, size = utf8.DecodeRune(src[advance:])
		if !unicode.IsSpace(r) {
			break
		}
	}
	if len(src) <= advance {
		// Request more data, if any.
		return len(src), nil, nil
	}
	// Decode length-1 lisp.Tokens.
	switch src[advance] {
	case '(': // LParen
		s.pos = lisp.Pos(advance)
		s.tok = lisp.LParen
		return advance + 1, src[advance : advance+1], nil
	case ')': // RParen
		s.pos = lisp.Pos(advance)
		s.tok = lisp.RParen
		return advance + 1, src[advance : advance+1], nil
	case '0': // Int(0)
		s.pos = lisp.Pos(advance)
		s.tok = lisp.Int
		return advance + 1, src[advance : advance+1], nil
	}
	// Decode longer lisp.Tokens.
	r, size := utf8.DecodeRune(src[advance:])
	switch {
	case '1' <= r && r <= '9': // Int
		s.pos = lisp.Pos(advance)
		// Int parsing may proceed byte-at-a-time since [0-9] <= RuneSelf.
		for advance++; advance < len(src); advance++ {
			b := src[advance]
			if b < '0' || '9' < b {
				break
			}
		}
		s.tok = lisp.Int
		return advance, src[s.pos:advance], nil
	case unicode.IsLetter(r): // Id
		s.pos = lisp.Pos(advance)
		advance += size
		for size := 0; advance < len(src); advance += size {
			var r rune
			r, size = utf8.DecodeRune(src[advance:])
			if !unicode.IsLetter(r) {
				break
			}
		}
		s.tok = lisp.Id
		return advance, src[s.pos:advance], nil
	}
	// Rune error.
	return advance, nil, &TokenError{
		Pos:   lisp.Pos(advance),
		Cause: fmt.Errorf("%w: %#q", errRune, r),
		Src:   src,
	}
}

type TokenScannerInterface interface {
	Reset(io.Reader)
	Scan() bool
	Token() (lisp.Pos, lisp.Token, string)
	Err() error
}

type consStackEntry struct {
	*lisp.Node            // Node for the cons.
	Root       *lisp.Cons // Root cons.
	Last       *lisp.Cons // Link to the last node, for fast insertion.
}

type NodeScanner struct {
	sc    TokenScannerInterface
	stack []*consStackEntry // FIXME: allow user-supplied buffer.
	err   error
	node  lisp.Node
}

func (s *NodeScanner) Reset(sc TokenScannerInterface) {
	*s = NodeScanner{}
	s.sc = sc
}

func (s *NodeScanner) Scan() bool {
	var n lisp.Node

	// Scan until a full Node is constructed.
	for s.sc.Scan() {
		var err error
		n, err = s.scan(s.sc.Token())
		if err != nil {
			s.err = err
			return false
		}
		// If the Node has valid indices, we're done.
		if n.Pos < n.End {
			break
		}
	}
	if err := s.sc.Err(); err != nil {
		s.err = err
		return false
	}
	s.node = n
	// Return true when valid Node scanned.
	// The final Scan call will be empty.
	return n.Pos < n.End
}

func (s *NodeScanner) scan(pos lisp.Pos, tok lisp.Token, text string) (lisp.Node, error) {
	switch tok {
	case lisp.Id, lisp.Int: // Id, Int
		n := lisp.Node{
			Pos: pos,
			Val: lisp.Lit{Token: tok, Text: text},
			End: pos + lisp.Pos(len(text)),
		}
		if i := len(s.stack); i > 0 {
			e := s.stack[i-1]
			if e.Last.Val != nil {
				e.Last.Cons = &lisp.Cons{}
				e.Last = e.Last.Cons
			}
			e.Last.Node = n
			// Need more scanning to finish this Cons.
			return lisp.Node{}, nil
		}
		return n, nil
	case lisp.LParen: // BEGIN Cons
		e := &consStackEntry{}
		cons := &lisp.Cons{}
		e.Node = &lisp.Node{Pos: pos, Val: cons}
		e.Root, e.Last = cons, cons
		s.stack = append(s.stack, e)
		// Need more scanning to finish this Cons.
		return lisp.Node{}, nil
	case lisp.RParen: // END Cons
		if len(s.stack) == 0 {
			// FIXME: Use  *lisp.NodeError.
			return lisp.Node{}, fmt.Errorf("unexpected ')'")
		}
		i := len(s.stack)
		x := s.stack[i-1]
		x.End = pos + 1
		s.stack = s.stack[:i-1]
		if i := len(s.stack); i > 0 {
			e := s.stack[i-1]
			if e.Last.Val != nil {
				e.Last.Cons = &lisp.Cons{}
				e.Last = e.Last.Cons
			}
			// Append the cons to the prev cons.
			e.Last.Node = *x.Node
			// Need more scanning to finish this Cons.
			return lisp.Node{}, nil
		}
		return *x.Node, nil
	default: // Unknown lisp.Token
		// FIXME: Use  *lisp.NodeError.
		return lisp.Node{}, fmt.Errorf("unexpected lisp.Token")
	}
}

// Node returns the last Node scanned.
//
// Valid scanned nodes always have indices set, i.e. lisp.Pos < End.
func (s *NodeScanner) Node() lisp.Node {
	return s.node
}

// Err returns the NodeScanner error.
func (s *NodeScanner) Err() error {
	if s.err != nil {
		return s.err
	}
	return s.sc.Err()
}

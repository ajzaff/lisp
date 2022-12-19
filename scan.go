package lisp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/rangetable"
)

var errRune = errors.New("unexpected rune")

var decTab = rangetable.New('0', '1', '2', '3', '4', '5', '6', '7', '8', '9')

var idTab = rangetable.Merge(
	unicode.L,
	decTab,
)

// Scanner scans the Lisp source for Tokens.
type TokenScanner struct {
	sc        *bufio.Scanner
	prev, off Pos   // absolute offsets
	pos       Pos   // relative token position
	tok       Token // last token scanned
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

func (s *TokenScanner) Token() (pos Pos, tok Token, text string) {
	return s.prev + s.pos, s.tok, s.Text()
}

func (s *TokenScanner) scanTokens(src []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(src) == 0 {
		return 0, nil, nil
	}
	defer func() {
		// Maintain absolute positions into the scanner.
		s.prev = s.off
		s.off += Pos(advance)
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
	// Decode length-1 tokens.
	switch src[advance] {
	case '(': // LParen
		s.pos = Pos(advance)
		s.tok = LParen
		return advance + 1, src[advance : advance+1], nil
	case ')': // RParen
		s.pos = Pos(advance)
		s.tok = RParen
		return advance + 1, src[advance : advance+1], nil
	case '0': // Int(0)
		s.pos = Pos(advance)
		s.tok = Int
		return advance + 1, src[advance : advance+1], nil
	}
	// Decode longer tokens.
	r, size := utf8.DecodeRune(src[advance:])
	switch {
	case '1' <= r && r <= '9': // Int
		s.pos = Pos(advance)
		// Int parsing may proceed byte-at-a-time since [0-9] <= RuneSelf.
		for advance++; advance < len(src); advance++ {
			b := src[advance]
			if b < '0' || '9' < b {
				break
			}
		}
		s.tok = Int
		return advance, src[s.pos:advance], nil
	case unicode.IsLetter(r): // Id
		s.pos = Pos(advance)
		advance += size
		for size := 0; advance < len(src); advance += size {
			var r rune
			r, size = utf8.DecodeRune(src[advance:])
			if !unicode.Is(idTab, r) {
				break
			}
		}
		s.tok = Id
		return advance, src[s.pos:advance], nil
	}
	// Rune error.
	return advance, nil, &TokenError{
		Pos:   Pos(advance),
		Cause: fmt.Errorf("%w: %#q", errRune, r),
		Src:   src,
	}
}

type TokenScannerInterface interface {
	Reset(io.Reader)
	Scan() bool
	Token() (Pos, Token, string)
	Err() error
}

type NodeScanner struct {
	sc    TokenScannerInterface
	stack []Node // FIXME: allow user-supplied buffer and Expr allocating config.
	err   error
	node  Node
}

func (s *NodeScanner) Reset(sc TokenScannerInterface) {
	*s = NodeScanner{}
	s.sc = sc
}

func (s *NodeScanner) Scan() bool {
	var n Node

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

func (s *NodeScanner) scan(pos Pos, tok Token, text string) (Node, error) {
	switch tok {
	case Id, Int: // Id, Int
		n := Node{
			Pos: pos,
			Val: Lit{Token: tok, Text: text},
			End: pos + Pos(len(text)),
		}
		if i := len(s.stack); i > 0 {
			s.stack[i-1].Val = append(s.stack[i-1].Val.(Expr), n)
			// Need more scanning to finish this Expr.
			return Node{}, nil
		}
		return n, nil
	case LParen: // BEGIN Expr
		// FIXME: Allow for tuning the Expr constructor (e.g. custom capacity).
		s.stack = append(s.stack, Node{Pos: pos, Val: Expr{}})
		// Need more scanning to finish this Expr.
		return Node{}, nil
	case RParen: // END Expr
		if len(s.stack) == 0 {
			// FIXME: Use *NodeError.
			return Node{}, fmt.Errorf("unexpected ')'")
		}
		i := len(s.stack)
		n := s.stack[i-1]
		n.End = pos + 1
		s.stack = s.stack[:i-1]
		if i := len(s.stack); i > 0 {
			s.stack[i-1].Val = append(s.stack[i-1].Val.(Expr), n)
			// Need more scanning to finish this Expr.
			return Node{}, nil
		}
		return n, nil
	default: // Unknown Token
		// FIXME: Use *NodeError.
		return Node{}, fmt.Errorf("unexpected token")
	}
}

// Node returns the last Node scanned.
//
// Valid scanned nodes always have indices set, i.e. Pos < End.
func (s *NodeScanner) Node() Node {
	return s.node
}

// Err returns the NodeScanner error.
func (s *NodeScanner) Err() error {
	if s.err != nil {
		return s.err
	}
	return s.sc.Err()
}

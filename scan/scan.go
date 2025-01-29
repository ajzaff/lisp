// Package scan provides a tokenizer and scanner implementation for Lisp.
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

// Pos represents a byte position in a source file.
type Pos int

// NoPos is the canonical value for no position defined.
const NoPos Pos = -1

// Scanner scans the Lisp source for lisp.Tokens.
type TokenScanner struct {
	sc        *bufio.Scanner
	prev, off Pos        // absolute offsets
	pos       Pos        // relative lisp.Token Position
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
	s.clearToken()
	return s.sc.Scan()
}

func (s *TokenScanner) Err() error {
	return s.sc.Err()
}

func (s *TokenScanner) Text() string {
	return s.sc.Text()
}

func (s *TokenScanner) clearToken() { s.setToken(NoPos, lisp.Invalid) }

func (s *TokenScanner) setToken(pos Pos, tok lisp.Token) {
	s.pos = pos
	s.tok = tok
}

func (s *TokenScanner) Token() (pos Pos, tok lisp.Token, text string) {
	return s.prev + s.pos, s.tok, s.Text()
}

func isDigit(b byte) bool { return '0' <= b && b <= '9' }

func isSpace(b byte) bool {
	switch b {
	case ' ', '\t', '\r', '\n':
		return true
	default:
		return false
	}
}

func isPunct(b byte) bool {
	switch b {
	case '(', ')':
		return true
	default:
		return false
	}
}

func isTokenBound(b byte) bool { return isSpace(b) || isPunct(b) }

func (s *TokenScanner) skipSpaces(src []byte) (advance int) {
	for ; advance < len(src) && isSpace(src[advance]); advance++ {
	}
	return advance
}

func (s *TokenScanner) scanTokenBound(src []byte, advance int, atEOF bool) (token bool, err error) {
	if advance == len(src) {
		// If we reached the end, but EOF is not set, we may have token data yet to process.
		// Request more data.
		if !atEOF {
			return false, nil
		}
	} else if !isTokenBound(src[advance]) {
		// If something caused us to end scanning early, make sure its a token boundary.
		// If not, that's a scan error.
		return false, &TokenError{
			Cause: fmt.Errorf("unexpected byte at token boundary: %#q", src[advance]),
			Pos:   Pos(advance),
			Src:   src,
		}
	}
	// We found either EOF or a valid token bound.
	// Either way the token is valid to emit.
	return true, nil
}

func (s *TokenScanner) scanNat(src []byte, atEOF bool) (advance int, token bool, err error) {
	for advance++; advance < len(src) && 0 <= src[advance] && src[advance] <= '9'; advance++ {
	}
	if valid, err := s.scanTokenBound(src, advance, atEOF); err != nil {
		return advance, false, err
	} else if !valid {
		return advance, false, nil
	}
	return advance, true, nil
}

func (s *TokenScanner) scanId(src []byte, atEOF bool) (advance int, token bool, err error) {
	for size := 0; advance < len(src); advance += size {
		var r rune
		r, size = utf8.DecodeRune(src[advance:])
		if !unicode.IsLetter(r) {
			break
		}
	}
	if valid, err := s.scanTokenBound(src, advance, atEOF); err != nil {
		return advance, false, err
	} else if !valid {
		return advance, false, nil
	}
	return advance, true, nil
}

func (s *TokenScanner) scanTokens(src []byte, atEOF bool) (advance int, token []byte, err error) {
	defer func() {
		// Maintain absolute Positions into the scanner.
		if token != nil {
			s.prev = s.off
		}
		s.off += Pos(advance)
	}()

	// Check in-progress token if any and fast-forward to new data.
	tok := s.tok
	switch s.tok {
	case lisp.Invalid: // No in-progress token; start a new token.
		// Skip leading spaces.
		if advance += s.skipSpaces(src); len(src) <= advance {
			return advance, nil, nil
		}
		// Decode cons operators directly, since they are only one byte.
		switch src[advance] {
		case '(': // LParen
			s.setToken(Pos(advance), lisp.LParen)
			return advance + 1, src[advance : advance+1], nil
		case ')': // RParen
			s.setToken(Pos(advance), lisp.RParen)
			return advance + 1, src[advance : advance+1], nil
		}
		switch {
		case isDigit(src[advance]): // Start a new Nat.
			tok = lisp.Nat
		default: // Start a new Id.
			tok = lisp.Id
		}

		s.setToken(Pos(advance), tok) // Mark new token start and position.
	default: // In progress token is set, fast-forward to new data.
		advance = int(s.pos)
	}

	var (
		n        int
		complete bool
	)

	switch tok {
	case lisp.Nat: // Nat
		n, complete, err = s.scanNat(src[advance:], atEOF)
	default: // lisp.Id
		n, complete, err = s.scanId(src[advance:], atEOF)
	}

	// Check error and request more data if the token is incomplete.
	advance += n
	if err != nil {
		return advance, nil, err
	}
	if !complete {
		// Request more data.
		return 0, nil, err
	}
	// Construct the complete token.
	token = src[s.pos : s.pos+Pos(advance)]
	return advance, token, err
}

type TokenScannerInterface interface {
	Reset(io.Reader)
	Scan() bool
	Token() (Pos, lisp.Token, string)
	Err() error
}

type consStackEntry struct {
	Pos, End Pos

	Root *lisp.Cons // Root cons.
	Last *lisp.Cons // Link to the last Cons, for fast insertion.
}

type NodeScanner struct {
	sc    TokenScannerInterface
	stack []*consStackEntry // FIXME: allow user-supplied buffer.
	err   error

	pos Pos
	end Pos
	val lisp.Val
}

func (s *NodeScanner) Reset(sc TokenScannerInterface) {
	*s = NodeScanner{}
	s.sc = sc
}

func (s *NodeScanner) Scan() bool {
	var end Pos
	var v lisp.Val

	// Scan until a full Val is constructed.
	for first := true; s.sc.Scan(); first = false {
		var err error
		var tok lisp.Token
		var text string
		pos, tok, text := s.sc.Token()
		if first {
			s.pos = pos
		}
		end, v, err = s.scan(pos, tok, text)
		if err != nil {
			s.err = err
			return false
		}
		// A Val is emitted, maybe return it.
		if v != nil {
			break
		}
	}
	if err := s.sc.Err(); err != nil {
		s.err = err
		return false
	}
	s.end = end
	s.val = v
	// Return true when valid Val scanned.
	// The final Scan call will be empty.
	return s.val != nil
}

func (s *NodeScanner) scan(pos Pos, tok lisp.Token, text string) (end Pos, v lisp.Val, err error) {
	switch tok {
	case lisp.Id, lisp.Nat: // Id, Nat
		end = pos + Pos(len(text))
		v = lisp.Lit{Token: tok, Text: text}
		if i := len(s.stack); i > 0 {
			e := s.stack[i-1]
			if e.Last.Val != nil {
				e.Last.Cons = &lisp.Cons{}
				e.Last = e.Last.Cons
			}
			e.Last.Val = v
			// Need more scanning to finish this Cons.
			return 0, nil, nil
		}
		return end, v, nil
	case lisp.LParen: // BEGIN Cons
		e := &consStackEntry{}
		cons := &lisp.Cons{}
		e.Root, e.Last = cons, cons
		s.stack = append(s.stack, e)
		// Need more scanning to finish this Cons.
		return 0, nil, nil
	case lisp.RParen: // END Cons
		if len(s.stack) == 0 {
			// FIXME: Use  *lisp.NodeError.
			return 0, nil, fmt.Errorf("unexpected ')'")
		}
		i := len(s.stack)
		x := s.stack[i-1]
		x.End = pos + 1
		s.stack = s.stack[:i-1]
		// Append to previous Cons.
		if i := len(s.stack); i > 0 {
			e := s.stack[i-1]
			if e.Last.Val != nil {
				e.Last.Cons = &lisp.Cons{}
				e.Last = e.Last.Cons
			}
			// Append the cons to the prev cons.
			e.Last.Val = x.Root
			// Need more scanning to finish this Cons.
			return 0, nil, nil
		}
		return x.End, x.Root, nil
	default: // Unknown lisp.Token
		// FIXME: Use  *lisp.NodeError.
		return 0, nil, fmt.Errorf("unexpected lisp.Token")
	}
}

// Node returns the last indices and Val scanned.
//
// When Scan returned true v will always be non-nil, and pos < end.
func (s *NodeScanner) Node() (pos, end Pos, v lisp.Val) {
	return s.pos, s.end, s.val
}

// Err returns the NodeScanner error.
func (s *NodeScanner) Err() error {
	if s.err != nil {
		return s.err
	}
	return s.sc.Err()
}

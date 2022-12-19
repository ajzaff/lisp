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
		s.pos += Pos(advance)
		s.tok = RParen
		return advance + 1, src[advance : advance+1], nil
	case '0': // Int(0)
		s.pos += Pos(advance)
		s.tok = Int
		return advance + 1, src[advance : advance+1], nil
	}
	// Decode longer tokens.
	r, size := utf8.DecodeRune(src[advance:])
	switch {
	case '1' <= r && r <= '9': // Int
		pos := Pos(advance)
		s.pos = pos
		for advance < len(src) {
			r, size := utf8.DecodeRune(src[advance:])
			if r < '0' || '9' < r {
				break
			}
			advance += size
		}
		s.tok = Int
		return advance, src[pos:advance], nil
	case unicode.IsLetter(r): // Id
		pos := Pos(advance)
		s.pos = pos
		advance += size
		for advance < len(src) {
			r, size := utf8.DecodeRune(src[advance:])
			if !unicode.Is(idTab, r) {
				break
			}
			advance += size
		}
		s.tok = Id
		return advance, src[pos:advance], nil
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
	sc   TokenScannerInterface
	err  error
	node Node
}

func (s *NodeScanner) Reset(sc TokenScannerInterface) {
	*s = NodeScanner{}
	s.sc = sc
}

func (s *NodeScanner) Scan() bool {
	if !s.sc.Scan() {
		return false
	}
	node, err := s.scan(s.sc.Token())
	if err != nil {
		s.err = err
		return false
	}
	s.node = node
	return true
}

func (s *NodeScanner) scan(pos Pos, tok Token, text string) (Node, error) {
	switch tok {
	case Id, Int:
		return Node{
			Pos: pos,
			Val: Lit{Token: tok, Text: text},
			End: pos + Pos(len(text)),
		}, nil
	case LParen, RParen:
		return s.scanExpr(pos)
	default:
		return Node{}, fmt.Errorf("unexpected token")
	}
}

func (s *NodeScanner) scanExpr(lParen Pos) (Node, error) {
	expr := Expr{}
	for {
		if !s.sc.Scan() {
			return Node{}, io.ErrUnexpectedEOF
		}
		pos, tok, text := s.sc.Token()
		switch tok {
		case RParen:
			return Node{
				Pos: lParen,
				Val: expr,
				End: pos + 1,
			}, nil
		}
		e, err := s.scan(pos, tok, text)
		if err != nil {
			return Node{}, err
		}
		expr = append(expr, e)
	}
}

func (s *NodeScanner) Node() Node {
	return s.node
}

func (s *NodeScanner) Err() error {
	if s.err != nil {
		return s.err
	}
	return s.sc.Err()
}

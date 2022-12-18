package lisp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

var (
	errRune = errors.New("unexpected rune")
)

// Scanner scans the Lisp source for Tokens.
type TokenScanner struct {
	sc    *bufio.Scanner
	start Pos   // absolute token start position
	end   Pos   // absolute token end position
	tok   Token // last token scanned
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
	return s.start, s.tok, s.Text()
}

func (s *TokenScanner) scanTokens(src []byte, atEOF bool) (advance int, token []byte, err error) {
	i := 0

	var r rune
	var size int

	// Skip space.
	for i < len(src) {
		r, size = utf8.DecodeRune(src[i:])
		if r == utf8.RuneError {
			if size == 0 {
				return 0, nil, nil
			}
			return 0, nil, errRune
		}
		i += size
		if !unicode.IsSpace(r) {
			break
		}
	}

	switch r {
	case '(':
		s.start = Pos(i)
		s.end = Pos(i + size)
		s.tok = LParen
		return size, src[s.start:s.end], nil
	case ')':
		s.start = Pos(i)
		s.end = Pos(i + size)
		s.tok = RParen
		return size, src[s.start:s.end], nil
	case '0':
		s.start = Pos(i)
		s.end = Pos(i + size)
		s.tok = Int
		return size, src[s.start:s.end], nil
	}
	switch {
	case unicode.IsDigit(r):
		s.start = Pos(i)
		i += size
		for i < len(src) {
			r, size := utf8.DecodeRune(src[i:])
			if r == utf8.RuneError {
				if size == 0 {
					return 0, nil, nil
				}
				return 0, nil, errRune
			}
			i += size
			if !unicode.IsDigit(r) {
				break
			}
		}
		s.end = Pos(i)
		s.tok = Int
		return i, src[s.start:s.end], nil
	case unicode.IsLetter(r):
		s.start = Pos(i)
		i += size
		for i < len(src) {
			r, size := utf8.DecodeRune(src[i:])
			if r == utf8.RuneError {
				if size == 0 {
					return 0, nil, nil
				}
				return 0, nil, errRune
			}
			i += size
			if !unicode.Is(idTab, r) {
				break
			}
		}
		s.end = Pos(i)
		s.tok = Id
		return i, src[s.start:s.end], nil
	}

	// Rune error!
	return i, src[i : i+size], errRune
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

func NewNodeScanner(sc TokenScannerInterface) *NodeScanner {
	var s NodeScanner
	s.sc = sc
	return &s
}

func (s *NodeScanner) Init(src io.Reader) {
	s.sc.Reset(src)
	s.err = nil
	s.node = Node{}
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
				End: pos,
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

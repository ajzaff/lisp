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
	var tok Token

	// Skip spaces.
	var i Pos
	for i < Pos(len(src)) {
		r, size := utf8.DecodeRune(src[i:])
		if !unicode.IsSpace(r) {
			break
		}
		i += Pos(size)
	}
	start := i
	if len(src[i:]) == 0 { // No token.
		return 0, nil, io.EOF
	}
	s.start = s.end + i // Update token start.

	// Get the first rune.
	r, size := utf8.DecodeRune(src[i:])
	i += Pos(size)
	switch {
	case r == '(': // LParen
		tok = LParen
	case r == ')': // RParen
		tok = RParen
	case r == '-', '0' <= r && r <= '9': // Int, Float
		tok = Number
		var dec bool
	num_loop:
		for i < Pos(len(src)) {
			r, size := utf8.DecodeRune(src[i:])
			switch {
			case '0' <= r && r <= '9':
			case r == '.':
				if dec {
					err = errRune
					break num_loop
				}
				dec = true
			default:
				break num_loop
			}
			i += Pos(size)
		}
	case IsId(r): // Id
		tok = Id
		var runeFunc func(rune) bool
		switch {
		case IsId(r):
			runeFunc = IsId
		case size == 0:
			err = io.EOF
		case r == utf8.RuneError:
			err = errRune
		default:
			panic("unreachable")
		}
		for i < Pos(len(src)) {
			r, size := utf8.DecodeRune(src[i:])
			if !runeFunc(r) {
				break
			}
			i += Pos(size)
		}
	default: // Unknown
		err = errRune
	}
	if i == Pos(len(src)) && err == nil {
		err = bufio.ErrFinalToken
	}
	advance, token = int(i), src[start:i]
	s.end += i // Update abs end position.
	s.tok = tok
	return advance, token, err
}

type NodeScanner struct {
	sc   *TokenScanner
	err  error
	node Node
}

func NewNodeScanner(sc *TokenScanner) *NodeScanner {
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
	case Id, Number:
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
	var expr Expr
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

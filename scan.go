package lisp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

// Scanner scans the Lisp source for Tokens.
type TokenScanner struct {
	sc  *bufio.Scanner
	tsc tokenScanner
}

func NewTokenScanner(src io.Reader) *TokenScanner {
	var sc TokenScanner
	sc.Reset(src)
	return &sc
}

func (sc *TokenScanner) Reset(src io.Reader) {
	sc.sc = bufio.NewScanner(src)
	sc.sc.Split(sc.tsc.scanTokens)
	sc.tsc = tokenScanner{}
}

func (sc *TokenScanner) Buffer(buf []byte, max int) {
	sc.sc.Buffer(buf, max)
}

func (sc *TokenScanner) Scan() bool {
	return sc.sc.Scan()
}

func (sc *TokenScanner) Err() error {
	return sc.sc.Err()
}

func (sc *TokenScanner) Text() string {
	return sc.sc.Text()
}

func (sc *TokenScanner) Token() (pos Pos, tok Token, text string) {
	return sc.tsc.start, sc.tsc.tok, sc.sc.Text()
}

type tokenScanner struct {
	start Pos   // absolute token start position
	end   Pos   // absolute token end position
	tok   Token // last token scanned
}

var (
	errRune = errors.New("unexpected rune")
)

func (t *tokenScanner) scanTokens(src []byte, atEOF bool) (advance int, token []byte, err error) {
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
	t.start = t.end + i // Update token start.

	// Get the first rune.
	r, size := utf8.DecodeRune(src[i:])
	i += Pos(size)
	switch r {
	case '(': // LParen
		tok = LParen
	case ')': // RParen
		tok = RParen
	case '"': // String
		tok = String
		n, scanErr := decodeStr(src[start:])
		if scanErr != nil {
			err = scanErr
			break
		}
		i += Pos(n)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // Int | Float
		tok = Number
		var dec bool
	num_loop:
		for i < Pos(len(src)) {
			r, size := utf8.DecodeRune(src[i:])
			switch r {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			case '.':
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
	default: // Id
		tok = Id
		var runeFunc func(rune) bool
		switch {
		case IsSymbol(r):
			runeFunc = IsSymbol
		case IsLetter(r):
			runeFunc = IsLetter
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
	}
	if i == Pos(len(src)) && err == nil {
		err = bufio.ErrFinalToken
	}
	advance, token = int(i), src[start:i]
	t.end += i // Update abs end position.
	t.tok = tok
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
	s.node = nil
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
	case Id, Number, String:
		return &LitNode{
			LitPos: pos,
			Lit:    Lit{Token: tok, Text: text},
			EndPos: pos + Pos(len(text)),
		}, nil
	case LParen, RParen:
		return s.scanExpr(pos)
	default:
		return nil, fmt.Errorf("unexpected token")
	}
}

func (s *NodeScanner) scanExpr(lParen Pos) (Node, error) {
	var expr Expr
	for {
		if !s.sc.Scan() {
			return nil, io.ErrUnexpectedEOF
		}
		pos, tok, text := s.sc.Token()
		switch tok {
		case RParen:
			return &ExprNode{
				LParen: lParen,
				Expr:   expr,
				RParen: pos,
			}, nil
		}
		e, err := s.scan(pos, tok, text)
		if err != nil {
			return nil, err
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

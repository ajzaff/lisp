// Package scan provides a tokenizer and scanner implementation for Lisp.
package scan

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"iter"
	"unicode"

	"github.com/ajzaff/lisp"
)

var errRune = errors.New("unexpected rune")

// Pos represents a byte position in a source file.
type Pos int

// NoPos is the canonical value for no position defined.
const NoPos Pos = -1

// Token emitted from TokenScanner.
type Token struct {
	Pos  Pos
	Tok  lisp.Token
	Text string
}

// Scanner scans the Lisp source for Lisp tokens and values.
type Scanner struct {
	r            bufio.Reader
	tb           bytes.Buffer
	pos          Pos
	err          error
	lastRuneSize int
	tokenFunc    func(Token) bool
	nodeFunc     func(Node) bool
}

func (s *Scanner) peekByte() (byte, error) {
	bs, err := s.r.Peek(1)
	if err != nil {
		return 0, err
	}
	return bs[0], nil
}

func (s *Scanner) readByte() (byte, error) {
	b, err := s.r.ReadByte()
	if err != nil {
		return 0, err
	}
	s.pos++
	return b, nil
}

func (s *Scanner) unreadByte() error {
	err := s.r.UnreadByte()
	if err != nil {
		return err
	}
	s.pos--
	return nil
}

func (s *Scanner) readRune() (rune, int, error) {
	r, size, err := s.r.ReadRune()
	if err != nil {
		return r, size, err
	}
	s.lastRuneSize = size
	s.pos += Pos(size)
	return r, size, nil
}

func (s *Scanner) unreadRune() error {
	err := s.r.UnreadRune()
	if err != nil {
		return err
	}
	s.pos -= Pos(s.lastRuneSize)
	s.lastRuneSize = -1
	return nil
}

func (s *Scanner) discardByte() error {
	_, err := s.r.Discard(1)
	if err == nil {
		s.pos++
	}
	return err
}

func (s *Scanner) setErr(err error) {
	if s.err == nil && err != io.EOF {
		s.err = err
	}
}

func (s *Scanner) hasErr() bool { return s.err != nil }

func (s *Scanner) Err() error { return s.err }

func (s *Scanner) resetToken() { s.tb.Reset() }

func (s *Scanner) Reset(r io.Reader) {
	s.r.Reset(r)
	s.lastRuneSize = -1
	s.pos = 0
	s.err = nil
	s.tokenFunc = nil
}

func (s *Scanner) yieldToken(tok Token) bool {
	if s.tokenFunc == nil {
		return true
	}
	return s.tokenFunc(tok)
}

func (s *Scanner) peekSpace0() bool {
	b, err := s.peekByte()
	if err != nil {
		s.setErr(err)
		return false
	}
	switch b {
	case ' ', '\t', '\r', '\n':
		return true
	default:
		return false
	}
}

func (s *Scanner) scanSpace0() bool {
	if !s.peekSpace0() {
		s.setErr(fmt.Errorf("expected space"))
		return false
	}
	s.discardByte()
	return true
}

func (s *Scanner) scanSpace1() {
	for s.peekSpace0() {
		s.discardByte()
	}
}

func (s *Scanner) scanSpace2() bool {
	if !s.scanSpace0() {
		return false
	}
	s.scanSpace1()
	return true
}

func (s *Scanner) scanLit0() bool {
	r, _, err := s.readRune()
	if err != nil {
		s.setErr(err)
		return false
	}
	if !unicode.IsLetter(r) {
		s.unreadRune()
		s.setErr(fmt.Errorf("expected letter"))
		return false
	}
	s.tb.WriteRune(r)
	return true
}

func (s *Scanner) peekDigit0() bool {
	b, err := s.peekByte()
	if err != nil {
		s.setErr(err)
		return false
	}
	switch {
	case '0' <= b && b <= '9':
		s.tb.WriteByte(b)
		return true
	default:
		return false
	}
}

func (s *Scanner) peekLit1() bool { return s.peekDigit0() || !s.peekGroup0() }

func (s *Scanner) scanLit1() bool {
	switch {
	case s.peekDigit0():
		s.discardByte()
		return true
	case s.peekGroup0():
		return false
	}
	return s.scanLit0()
}

func (s *Scanner) scanLit2() bool {
	tok := Token{
		Pos: s.pos,
		Tok: lisp.Id,
	}
	s.resetToken()
	if !s.scanLit1() {
		return false
	}
	for s.peekLit1() && s.scanLit1() {
	}
	tok.Text = s.tb.String()
	s.yieldToken(tok)
	return true
}

func (s *Scanner) scanLit3() {
	if !s.scanLit2() {
		return
	}
	for s.peekSpace0() {
		s.scanSpace2()
		if !s.peekLit1() {
			break
		}
		s.scanLit2()
	}
}

func (s *Scanner) peekGroup0() bool {
	b, err := s.peekByte()
	if err != nil {
		s.setErr(err)
		return false
	}
	return b == '('
}

func (s *Scanner) scanGroup0() bool {
	if !s.peekGroup0() {
		return false
	}
	s.yieldToken(Token{Pos: s.pos, Tok: lisp.LParen, Text: "("})
	s.discardByte()
	s.scanSpace1()
	s.scanExpr2()
	s.scanSpace1()
	if !s.peekGroup0() {
		return false
	}
	s.yieldToken(Token{Pos: s.pos, Tok: lisp.RParen, Text: ")"})
	s.discardByte()
	return true
}

func (s *Scanner) scanGroup1() bool {
	if !s.scanGroup0() {
		return false
	}
	for {
		s.scanSpace1()
		if !s.peekGroup0() {
			break
		}
	}
	return true
}

func (s *Scanner) peekExpr0() bool { return s.peekGroup0() || s.peekLit1() }

func (s *Scanner) scanExpr0() bool { return s.peekGroup0() && s.scanGroup0() || s.scanLit2() }

func (s *Scanner) scanExpr1() bool {
	if !s.scanExpr0() {
		return false
	}
	for s.peekExpr0() && s.scanExpr0() {
	}
	return true
}

func (s *Scanner) scanExpr2() bool {
	if s.peekExpr0() {
		return s.scanExpr1()
	}
	return true
}

func (s *Scanner) scanExpr3() {
	s.scanSpace1()
	if !s.scanExpr2() {
		return
	}
	s.scanSpace1()
}

// Tokens returns a iteration over tokens without respect for correct syntax.
func (s *Scanner) Tokens() iter.Seq[Token] {
	return func(yield func(Token) bool) {
		s.tokenFunc = yield
		defer func() { s.tokenFunc = nil }()
		for {
			s.scanSpace1()
			switch {
			case s.peekGroup0():
				s.yieldToken(Token{Pos: s.pos, Tok: lisp.LParen, Text: "("})
				s.discardByte()
			default:
				if b, err := s.peekByte(); err != nil {
					s.setErr(err)
					return
				} else if b == ')' {
					s.yieldToken(Token{Pos: s.pos, Tok: lisp.RParen, Text: ")"})
				} else if !s.scanLit2() {
					return
				}
			}
		}
	}
}

type Node struct {
	Pos Pos
	Val lisp.Val
	End Pos
}

func (s *Scanner) Nodes() iter.Seq[Node] {
	return func(yield func(Node) bool) {
		s.nodeFunc = yield
		s.scanExpr3()
		s.nodeFunc = nil
	}
}

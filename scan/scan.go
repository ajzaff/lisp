// Package scan provides a tokenizer and scanner implementation for Lisp.
package scan

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"iter"
	"unicode"

	"github.com/ajzaff/lisp"
)

// Pos represents a byte position in a source file.
type Pos int

// NoPos is the canonical value for no position defined.
const NoPos Pos = -1

// Scanner scans the Lisp source for Lisp tokens and values.
type Scanner struct {
	r            bufio.Reader
	tb           bytes.Buffer
	pos          Pos
	err          error
	lastRuneSize int
}

func (s *Scanner) peekByteErr() (byte, error) {
	bs, err := s.r.Peek(1)
	if err != nil {
		return 0, err
	}
	return bs[0], nil
}

func (s *Scanner) peekByte() byte {
	bs, err := s.r.Peek(1)
	if err != nil {
		return 0
	}
	return bs[0]
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
}

func (s *Scanner) peekSpace0(b byte) bool {
	switch b {
	case ' ', '\t', '\r', '\n':
		return true
	default:
		return false
	}
}

func (s *Scanner) skipSpace0() bool {
	b, err := s.readByte()
	if err != nil {
		s.setErr(err)
		return false
	}
	if !s.peekSpace0(b) {
		s.setErr(fmt.Errorf("expected space"))
		return false
	}
	return true
}

func (s *Scanner) skipSpace1() {
	for s.peekSpace0(s.peekByte()) {
		s.discardByte()
	}
}

func (s *Scanner) scanSpace2() bool {
	if !s.skipSpace0() {
		return false
	}
	s.skipSpace1()
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
		s.setErr(fmt.Errorf("expected LIT, got %q", r))
		return false
	}
	s.tb.WriteRune(r)
	return true
}

func (s *Scanner) peekDigit0(b byte) bool { return '0' <= b && b <= '9' }

func (s *Scanner) peekLit1(b byte) bool {
	return s.peekDigit0(b) || !s.peekGroup0(b) && !s.peekSpace0(b)
}

func (s *Scanner) scanLit1() bool {
	switch b := s.peekByte(); {
	case s.peekDigit0(b):
		s.tb.WriteByte(b)
		s.discardByte()
		return true
	case s.peekGroup0(b):
		return false
	}
	return s.scanLit0()
}

func (s *Scanner) scanLit2() bool {
	s.resetToken()
	if !s.scanLit1() {
		return false
	}
	for s.peekLit1(s.peekByte()) && s.scanLit1() {
	}
	return true
}

func (s *Scanner) scanLit3() {
	if !s.scanLit2() {
		return
	}
	for s.peekSpace0(s.peekByte()) && s.scanSpace2() {
		if !s.peekLit1(s.peekByte()) || !s.scanLit2() {
			break
		}
	}
}

func (s *Scanner) peekGroup0(b byte) bool { return b == '(' }

func (s *Scanner) peekGroupEnd(b byte) bool { return b == ')' }

func (s *Scanner) scanGroup0() bool {
	if !s.peekGroup0(s.peekByte()) {
		return false
	}
	s.discardByte()
	s.skipSpace1()
	s.scanExpr2()
	s.skipSpace1()
	if !s.peekGroup0(s.peekByte()) {
		return false
	}
	s.discardByte()
	return true
}

func (s *Scanner) scanGroup1() bool {
	if !s.scanGroup0() {
		return false
	}
	for {
		s.skipSpace1()
		if !s.scanGroup0() {
			break
		}
	}
	return true
}

func (s *Scanner) peekExpr0(b byte) bool { return s.peekGroup0(b) || s.peekLit1(b) }

func (s *Scanner) scanExpr0() bool {
	return s.peekGroup0(s.peekByte()) && s.scanGroup0() || s.scanLit2()
}

func (s *Scanner) scanExpr1() bool {
	if !s.scanExpr0() {
		return false
	}
	for s.peekExpr0(s.peekByte()) && s.scanExpr0() {
	}
	return true
}

func (s *Scanner) scanExpr2() bool {
	if s.peekExpr0(s.peekByte()) {
		return s.scanExpr1()
	}
	return true
}

func (s *Scanner) scanExpr3() {
	s.skipSpace1()
	if !s.scanExpr2() {
		return
	}
	s.skipSpace1()
}

// Token emitted from TokenScanner.
type Token struct {
	Pos  Pos
	Tok  lisp.Token
	Text string
}

// Tokens returns a iteration over tokens without respect for correct syntax.
func (s *Scanner) Tokens() iter.Seq[Token] {
	return func(yield func(Token) bool) {
		for {
			s.skipSpace1()
			switch b, err := s.peekByteErr(); {
			case err != nil:
				s.setErr(err)
				return
			case s.peekGroup0(b):
				if !yield(Token{Pos: s.pos, Tok: lisp.LParen, Text: "("}) {
					return
				}
				s.discardByte()
			case s.peekGroupEnd(b):
				if !yield(Token{Pos: s.pos, Tok: lisp.RParen, Text: ")"}) {
					return
				}
				s.discardByte()
			default:
				pos := s.pos
				if !s.scanLit2() || !yield(Token{Pos: pos, Tok: lisp.Id, Text: s.tb.String()}) {
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
		nodeStack := []Node{}
		for {
			s.skipSpace1()
			switch b, err := s.peekByteErr(); {
			case err != nil:
				s.setErr(err)
				return
			case s.peekGroup0(b):
				nodeStack = append(nodeStack, Node{
					Pos: s.pos,
					Val: lisp.Group{},
					End: NoPos,
				})
				s.discardByte()
			case s.peekGroupEnd(b):
				if len(nodeStack) == 0 {
					s.setErr(fmt.Errorf("unexpected )"))
					return
				}
				s.discardByte()
				n := len(nodeStack) - 1
				e := nodeStack[n]
				nodeStack = nodeStack[:n]
				if len(nodeStack) == 0 {
					e.End = s.pos
					if !yield(e) {
						return
					}
					continue
				}
				prev := nodeStack[len(nodeStack)-1]
				prev.Val = append(prev.Val.(lisp.Group), e.Val)
			default:
				pos := s.pos
				if !s.scanLit2() {
					return
				}
				text := s.tb.String()
				if len(nodeStack) > 0 {
					g := nodeStack[len(nodeStack)-1]
					g.Val = append(g.Val.(lisp.Group), lisp.Lit(s.tb.String()))
				} else if !yield(Node{
					Pos: pos,
					Val: lisp.Lit(text),
					End: pos + Pos(len(text)),
				}) {
					return
				}
			}
		}
	}
}

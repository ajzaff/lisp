// Package scan provides a tokenizer and scanner implementation for Lisp.
package scan

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

// Pos represents a byte position in a source file.
type Pos int

// NoPos is the canonical value for no position defined.
const NoPos Pos = -1

// Scanner scans the Lisp source for Lisp tokens and values.
type Scanner struct {
	r            bufio.Reader
	pos          Pos
	err          error
	lastRuneSize int
}

func (s *Scanner) Pos() Pos { return s.pos }

func (s *Scanner) PeekByteErr() (byte, error) {
	bs, err := s.r.Peek(1)
	if err != nil {
		return 0, err
	}
	return bs[0], nil
}

func (s *Scanner) PeekByte() byte {
	bs, err := s.r.Peek(1)
	if err != nil {
		return 0
	}
	return bs[0]
}

func (s *Scanner) ReadByte() (byte, error) {
	b, err := s.r.ReadByte()
	if err != nil {
		return 0, err
	}
	s.pos++
	return b, nil
}

func (s *Scanner) UnreadByte() error {
	err := s.r.UnreadByte()
	if err != nil {
		return err
	}
	s.pos--
	return nil
}

func (s *Scanner) ReadRune() (rune, int, error) {
	r, size, err := s.r.ReadRune()
	if err != nil {
		return r, size, err
	}
	s.lastRuneSize = size
	s.pos += Pos(size)
	return r, size, nil
}

func (s *Scanner) UnreadRune() error {
	err := s.r.UnreadRune()
	if err != nil {
		return err
	}
	s.pos -= Pos(s.lastRuneSize)
	s.lastRuneSize = -1
	return nil
}

func (s *Scanner) DiscardByte() error {
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

func (s *Scanner) Err() error { return s.err }

func (s *Scanner) Reset(r io.Reader) {
	s.r.Reset(r)
	s.lastRuneSize = -1
	s.pos = 0
	s.err = nil
}

func (s *Scanner) PeekSpace0(b byte) bool {
	switch b {
	case ' ', '\t', '\r', '\n':
		return true
	default:
		return false
	}
}

func (s *Scanner) SkipSpace0() bool {
	b, err := s.ReadByte()
	if err != nil {
		s.setErr(err)
		return false
	}
	if !s.PeekSpace0(b) {
		s.setErr(fmt.Errorf("expected space"))
		return false
	}
	return true
}

func (s *Scanner) SkipSpace1() (skip bool) {
	for s.PeekSpace0(s.PeekByte()) {
		s.DiscardByte()
		skip = true
	}
	return skip
}

func (s *Scanner) ScanSpace2() bool {
	if !s.SkipSpace0() {
		return false
	}
	s.SkipSpace1()
	return true
}

func (s *Scanner) ScanLit0() bool {
	r, _, err := s.ReadRune()
	if err != nil {
		s.setErr(err)
		return false
	}
	if !unicode.IsLetter(r) {
		s.UnreadRune()
		s.setErr(fmt.Errorf("expected LIT, got %q", r))
		return false
	}
	return true
}

func (s *Scanner) PeekDigit0(b byte) bool { return '0' <= b && b <= '9' }

func (s *Scanner) PeekLit1(b byte) bool {
	return s.PeekDigit0(b) || !s.PeekGroup0(b) && !s.PeekGroupEnd(b) && !s.PeekSpace0(b)
}

func (s *Scanner) ScanLit1() bool {
	switch b := s.PeekByte(); {
	case s.PeekDigit0(b):
		s.DiscardByte()
		return true
	case s.PeekGroup0(b), s.PeekGroupEnd(b), s.PeekSpace0(b):
		return false
	}
	return s.ScanLit0()
}

func (s *Scanner) ScanLit2() bool {
	if !s.ScanLit1() {
		return false
	}
	for s.PeekLit1(s.PeekByte()) && s.ScanLit1() {
	}
	return true
}

func (s *Scanner) ScanLit3() {
	if !s.ScanLit2() {
		return
	}
	for s.PeekSpace0(s.PeekByte()) && s.ScanSpace2() {
		if !s.PeekLit1(s.PeekByte()) || !s.ScanLit2() {
			break
		}
	}
}

func (s *Scanner) PeekGroup0(b byte) bool { return b == '(' }

func (s *Scanner) PeekGroupEnd(b byte) bool { return b == ')' }

func (s *Scanner) ScanGroup0() bool {
	if !s.PeekGroup0(s.PeekByte()) {
		return false
	}
	s.DiscardByte()
	s.SkipSpace1()
	s.ScanExpr2()
	s.SkipSpace1()
	if !s.PeekGroup0(s.PeekByte()) {
		return false
	}
	s.DiscardByte()
	return true
}

func (s *Scanner) ScanGroup1() bool {
	if !s.ScanGroup0() {
		return false
	}
	for {
		s.SkipSpace1()
		if !s.ScanGroup0() {
			break
		}
	}
	return true
}

func (s *Scanner) PeekExpr0(b byte) bool { return s.PeekGroup0(b) || s.PeekLit1(b) }

func (s *Scanner) ScanExpr0() bool {
	return s.PeekGroup0(s.PeekByte()) && s.ScanGroup0() || s.ScanLit2()
}

func (s *Scanner) ScanExpr1() bool {
	if !s.ScanExpr0() {
		return false
	}
	for s.PeekExpr0(s.PeekByte()) && s.ScanExpr0() {
	}
	return true
}

func (s *Scanner) ScanExpr2() bool {
	if s.PeekExpr0(s.PeekByte()) {
		return s.ScanExpr1()
	}
	return true
}

func (s *Scanner) ScanExpr3() {
	s.SkipSpace1()
	if !s.ScanExpr2() {
		return
	}
	s.SkipSpace1()
}

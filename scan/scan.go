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
	pos          Pos
	err          error
	lastRuneSize int
}

func (s *Scanner) peekByteErr() (byte, error) {
	bs, err := s.r.Peek(1)
	if err != nil {
		s.setErr(err)
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

func (s *Scanner) Err() error { return s.err }

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

func (s *Scanner) skipSpace1() {
	for s.peekSpace0(s.peekByte()) {
		s.discardByte()
	}
}

func (s *Scanner) peekDigit0(b byte) bool { return '0' <= b && b <= '9' }

func (s *Scanner) peekLit1(b byte) bool {
	return s.peekDigit0(b) || !s.peekGroup0(b) && !s.peekGroupEnd(b) && !s.peekSpace0(b)
}

func (s *Scanner) writeLit0(buf *bytes.Buffer) bool {
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
	buf.WriteRune(r)
	return true
}

func (s *Scanner) writeLit1(buf *bytes.Buffer) bool {
	switch b := s.peekByte(); {
	case s.peekDigit0(b):
		buf.WriteByte(b)
		s.discardByte()
		return true
	case s.peekGroup0(b), s.peekGroupEnd(b), s.peekSpace0(b):
		return false
	default:
		return s.writeLit0(buf)
	}
}

func (s *Scanner) writeLit2(buf *bytes.Buffer) bool {
	if !s.writeLit1(buf) {
		return false
	}
	for s.peekLit1(s.peekByte()) && s.writeLit1(buf) {
	}
	return true
}

func (s *Scanner) peekGroup0(b byte) bool { return b == '(' }

func (s *Scanner) peekGroupEnd(b byte) bool { return b == ')' }

// Token emitted from TokenScanner.
type Token struct {
	Pos  Pos
	Tok  lisp.Token
	Text string
}

// Tokens returns a iteration over tokens without respect for correct syntax.
func (s *Scanner) Tokens() iter.Seq[Token] {
	return func(yield func(Token) bool) {
		var buf bytes.Buffer
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
				buf.Reset()
				if !s.writeLit2(&buf) || !yield(Token{Pos: pos, Tok: lisp.Id, Text: buf.String()}) {
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
		nodeStack := []*Node{}
		var buf bytes.Buffer
		for {
			s.skipSpace1()
			switch b, err := s.peekByteErr(); {
			case err != nil:
				s.setErr(err)
				return
			case s.peekGroup0(b):
				nodeStack = append(nodeStack, &Node{
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
					if !yield(*e) {
						return
					}
					continue
				}
				prev := nodeStack[len(nodeStack)-1]
				prev.Val = append(prev.Val.(lisp.Group), e.Val)
			default:
				pos := s.pos
				buf.Reset()
				if !s.writeLit2(&buf) {
					return
				}
				text := buf.String()
				if len(nodeStack) > 0 {
					g := nodeStack[len(nodeStack)-1]
					g.Val = append(g.Val.(lisp.Group), lisp.Lit(buf.String()))
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

func (s *Scanner) Values() iter.Seq[lisp.Val] {
	return func(yield func(lisp.Val) bool) {
		groupStack := []*lisp.Group{}
		var buf bytes.Buffer
		for {
			s.skipSpace1()
			switch b, err := s.peekByteErr(); {
			case err != nil:
				s.setErr(err)
				return
			case s.peekGroup0(b):
				groupStack = append(groupStack, &lisp.Group{})
				s.discardByte()
			case s.peekGroupEnd(b):
				if len(groupStack) == 0 {
					s.setErr(fmt.Errorf("unexpected )"))
					return
				}
				s.discardByte()
				n := len(groupStack) - 1
				e := groupStack[n]
				groupStack = groupStack[:n]
				if len(groupStack) == 0 {
					if !yield(*e) {
						return
					}
					continue
				}
				prev := groupStack[len(groupStack)-1]
				*prev = append(*prev, *e)
			default:
				buf.Reset()
				if !s.writeLit2(&buf) {
					return
				}
				text := buf.String()
				if len(groupStack) > 0 {
					g := groupStack[len(groupStack)-1]
					*g = append(*g, lisp.Lit(buf.String()))
				} else if !yield(lisp.Lit(text)) {
					return
				}
			}
		}
	}
}

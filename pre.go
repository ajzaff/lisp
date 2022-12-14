package lisp

import (
	"bufio"
	"io"
)

type Pretokenizer struct {
	buf *bufio.Scanner
	sc  pretokenScanner
}

func NewPretokenizer(r io.Reader) *Pretokenizer {
	var t Pretokenizer
	t.Reset(r)
	return &t
}

func (t *Pretokenizer) Reset(r io.Reader) {
	t.buf = bufio.NewScanner(r)
	t.sc = pretokenScanner{}
	t.buf.Split(t.sc.scanRawToken)
}

type pretokenScanner struct {
	advance int
}

func (s *pretokenScanner) Reset() {
	s.advance = 0
}

func (s *pretokenScanner) scanRawToken(data []byte, atEOF bool) (advance int, token []byte, err error) {
	s.Reset()
	return 0, nil, nil
}

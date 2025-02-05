// Package fuzz implements random expression generation.
package fuzz

import (
	"bufio"
	"io"
	"math/rand/v2"
)

// Source writes a random program
type Source struct {
	w  bufio.Writer
	r  *rand.Rand
	sd int // expr stack depth
}

func (s *Source) Reset(w io.Writer) {
	s.sd = 0
	s.w.Reset(w)
	if s.r == nil {
		s.Seed(rand.NewPCG(1337, 0xbeef))
	}
}

func (s *Source) Seed(src rand.Source) {
	s.r = rand.New(src)
}

func (s *Source) Next() {
	switch s.r.IntN(4) {
	case 0: // 0
	case 1: // 123
	case 2: // abc
	default: //
	}
}

// Close the current expression.
func (s *Source) Close() {
	for ; s.sd > 0; s.sd-- {
		s.w.WriteByte(')')
	}
	s.w.WriteByte('\n')
}

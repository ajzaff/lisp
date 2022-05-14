// Package fdmap implements a generic frequency dictionary.
package fdmap

import (
	"io"
	"strconv"

	"github.com/ajzaff/lisp"
)

type Key interface {
	lisp.IdLit | lisp.IntLit | lisp.StringLit
}

type FreqMap[T Key] struct {
	sc   *lisp.TokenScanner
	err  error
	data map[T]int
}

func NewFreqMap[T Key]() *FreqMap[T] {
	return &FreqMap[T]{sc: lisp.NewTokenScanner(nil), data: make(map[T]int)}
}

func (m *FreqMap[T]) Init(r io.Reader) {
	m.sc.Init(r)
}

func (m *FreqMap[T]) Scan() bool {
	res := m.sc.Scan()
	if _, t, text := m.sc.Token(); t != lisp.Invalid {
		var lit lisp.Lit
		switch t {
		case lisp.Id:
			lit = lisp.IdLit(text)
		case lisp.Int:
			x, _ := strconv.ParseInt(text, 10, 64)
			lit = lisp.IntLit(x)
		case lisp.String:
			lit = lisp.StringLit(text)
		}
		if t, ok := lit.(T); ok {
			m.data[t]++
		}
	}
	return res
}

func (m *FreqMap[T]) Err() error {
	if err := m.sc.Err(); err != nil {
		return err
	}
	return m.err
}

func (m *FreqMap[T]) Put(key T, v int) {
	m.data[key] += v
}

func (m *FreqMap[T]) Get(key T) int {
	return m.data[key]
}

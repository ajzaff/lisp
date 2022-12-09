// Package fdmap implements a generic frequency dictionary.
package fdmap

import (
	"io"

	"github.com/ajzaff/lisp"
)

type FreqMap struct {
	sc   *lisp.TokenScanner
	err  error
	data map[lisp.Lit]int
}

func NewFreqMap() *FreqMap {
	return &FreqMap{sc: lisp.NewTokenScanner(nil), data: make(map[lisp.Lit]int)}
}

func (m *FreqMap) Init(r io.Reader) {
	m.sc.Init(r)
}

func (m *FreqMap) Scan() bool {
	res := m.sc.Scan()
	if _, t, text := m.sc.Token(); t != lisp.Invalid {
		var lit lisp.Lit
		switch t {
		case lisp.Id:
			lit = lisp.Lit{Token: lisp.Id, Text: text}
		case lisp.Number:
			lit = lisp.Lit{Token: lisp.Number, Text: text}
		case lisp.String:
			lit = lisp.Lit{Token: lisp.String, Text: text}
		}
		m.data[lit]++
	}
	return res
}

func (m *FreqMap) Err() error {
	if err := m.sc.Err(); err != nil {
		return err
	}
	return m.err
}

func (m *FreqMap) Put(key lisp.Lit, v int) {
	m.data[key] += v
}

func (m *FreqMap) Get(key lisp.Lit) int {
	return m.data[key]
}

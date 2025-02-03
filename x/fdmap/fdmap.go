// Package fdmap implements a generic frequency dictionary.
package fdmap

import (
	"io"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/scan"
)

type FreqMap struct {
	sc   *scan.TokenScanner
	err  error
	data map[lisp.Lit]int
}

func NewFreqMap() *FreqMap {
	var sc scan.TokenScanner
	return &FreqMap{sc: &sc, data: make(map[lisp.Lit]int)}
}

func (m *FreqMap) Init(r io.Reader) {
	m.sc.Reset(r)
}

func (m *FreqMap) Scan() bool {
	res := m.sc.Scan()
	if _, t, text := m.sc.Token(); t != lisp.Invalid {
		lit := lisp.Lit(text)
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

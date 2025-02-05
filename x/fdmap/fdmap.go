// Package fdmap implements a generic frequency dictionary.
package fdmap

import (
	"io"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/scan"
)

type FreqMap map[string]int

func (m FreqMap) Clear() {
	for k := range m {
		delete(m, k)
	}
}

func (m FreqMap) Count(r io.Reader) {
	var sc scan.Scanner
	sc.Reset(r)
	for t := range sc.Tokens() {
		if t.Tok == lisp.Id {
			m[t.Text]++
		}
	}
}

func (m FreqMap) Put(key string, v int) {
	m[key] += v
}

func (m FreqMap) Get(key string) int {
	return m[key]
}

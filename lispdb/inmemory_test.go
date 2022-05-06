package lispdb

import (
	"testing"

	"github.com/ajzaff/lisp"
)

func TestInMemory(t *testing.T) {
	m := NewInMemory()

	n, _ := lisp.Parser{}.Parse("(x y (z (1 2 3)) a)")

	Store(m, n[0].Val(), 1)

	t.Logf("%v\n", m)
}

package innitdb

import (
	"testing"

	"github.com/ajzaff/innit"
)

func TestInMemory(t *testing.T) {
	m := NewInMemory()

	n, _ := innit.Parse("(x y (z (1 2 3)) a)")

	Store(m, n[0].Val(), 1)

	t.Logf("%v\n", m)
}

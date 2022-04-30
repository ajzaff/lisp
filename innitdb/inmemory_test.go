package innitdb

import (
	"testing"

	"github.com/ajzaff/innit"
)

func TestInMemory(t *testing.T) {
	m := NewInMemory()

	n, _ := innit.Parse("(x y (z (1 2 3)) a)")

	Store(m, n, 1)

	t.Logf("%v\n", m)
}

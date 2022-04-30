package innitdb

import (
	"fmt"
	"testing"

	"github.com/ajzaff/innit"
)

func TestInMemory(t *testing.T) {
	m := NewInMemory()

	n, _ := innit.Parse("(x y (z (1 2 3)) a)")

	m.Store(n, 1)

	fmt.Printf("%v\n", m)
}

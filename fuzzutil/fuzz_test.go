package fuzzutil

import (
	"math/rand"
	"os"
	"testing"

	"github.com/ajzaff/lisp"
)

func TestFuzzExample(t *testing.T) {
	seeds := rand.New(rand.NewSource(1337))
	for i := 0; i < 10; i++ {
		g := NewGenerator(rand.New(rand.NewSource(seeds.Int63())))
		lisp.StdPrinter(os.Stdout).Print(g.Next().Val)
	}
}

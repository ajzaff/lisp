package fuzzutil

import (
	"math/rand"
	"os"
	"testing"

	"github.com/ajzaff/lisp/x/print"
)

func TestFuzzExample(t *testing.T) {
	seeds := rand.New(rand.NewSource(1337))
	for i := 0; i < 20; i++ {
		g := NewGenerator(rand.New(rand.NewSource(seeds.Int63())))
		g.ConsWeight = 2
		g.ConsMaxDepth = 3
		print.StdPrinter(os.Stdout).Print(g.Next())
	}
}

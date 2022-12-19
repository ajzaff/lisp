package main

import (
	"flag"
	"math/rand"
	"os"
	"time"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/fuzzutil"
)

var (
	seed         = flag.Int64("seed", 1337, "Seed the random number generator with the given seed")
	timeSeed     = flag.Bool("seedtime", false, "Seed the random number generator with the current time (overrides seed)")
	exprMaxDepth = flag.Int("expr_max_depth", 2, "Maximum depth of Expr")
	idWeight     = flag.Int("id_weight", 1, "Weight for emitting Id")
	intWeight    = flag.Int("int_weight", 1, "Weight for emitting Int")
	exprWeight   = flag.Int("expr_weight", 10, "Weight for emitting Expr")
)

func main() {
	flag.Parse()

	seed := *seed
	if *timeSeed {
		seed = time.Now().UnixNano()
	}

	r := rand.New(rand.NewSource(seed))
	g := fuzzutil.NewGenerator(r)
	g.ConsMaxDepth = *exprMaxDepth
	g.IdWeight = *idWeight
	g.IntWeight = *intWeight
	g.ConsWeight = *exprWeight

	lisp.StdPrinter(os.Stdout).Print(g.Next())
}

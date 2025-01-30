package main

import (
	"flag"
	"math/rand"
	"os"
	"time"

	"github.com/ajzaff/lisp/x/fuzzutil"
	"github.com/ajzaff/lisp/x/print"
)

var (
	seed         = flag.Int64("seed", 1337, "Seed the random number generator with the given seed")
	seedTime     = flag.Bool("seed_time", false, "Seed the random number generator with the current time (overrides seed)")
	idWeight     = flag.Int("id_weight", 1, "Weight for emitting Id")
	natWeight    = flag.Int("nat_weight", 1, "Weight for emitting Nat")
	consWeight   = flag.Int("cons_weight", 10, "Weight for emitting Group")
	consMaxDepth = flag.Int("cons_max_depth", 2, "Maximum depth of Group")
)

func main() {
	flag.Parse()

	seed := *seed
	if *seedTime {
		seed = time.Now().UnixNano()
	}

	r := rand.New(rand.NewSource(seed))
	g := fuzzutil.NewGenerator(r)
	g.GroupMaxDepth = *consMaxDepth
	g.IdWeight = *idWeight
	g.IntWeight = *natWeight
	g.GroupWeight = *consWeight

	print.StdPrinter(os.Stdout).Print(g.Next())
}

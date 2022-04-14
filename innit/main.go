package main

import (
	"flag"
	"fmt"
	"hash/maphash"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/ajzaff/innit"
	"github.com/ajzaff/innit/hash"
	"github.com/ajzaff/innit/innitdb"
)

var (
	order = flag.String("order", "", `Print order for AST print mode (Optional "reverse". Default uses in-order)`)
	print = flag.String("print", "", `Print mode (Optional "tok", "ast", "db". Default uses StdPrinter)`)
	file  = flag.String("file", "", "File to read innit code from.")
)

func main() {
	flag.Parse()

	if *file == "" {
		doRepl()
		return
	}

	src, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatal(err)
	}

	n, err := innit.Parse(string(src))
	if err != nil {
		log.Fatal(err)
	}

	switch *print {
	case "": // std
		innit.StdPrinter(os.Stdout).Print(n)
	case "tok":
		tokens, err := innit.Tokenize(string(src))
		if err != nil {
			panic(err)
		}
		for i := 0; i < len(tokens); i += 2 {
			token := string(src[tokens[i]:tokens[i+1]])
			println(token)
		}
	case "ast":
		var v innit.Visitor
		exprVisitor := func(e *innit.Expr) {
			var sb strings.Builder
			innit.StdPrinter(&sb).Print(e)
			fmt.Print("EXPR\t", sb.String())
		}
		switch *order {
		case "": // in-order
			v.SetBeforeExprVisitor(exprVisitor)
		case "reverse":
			v.SetAfterExprVisitor(exprVisitor)
		default:
			log.Fatalf("unexpected -order mode: %v", *order)
		}
		v.SetLitVisitor(func(e *innit.Lit) {
			fmt.Println("LIT\t", e.Tok.String(), "\t", e.Value)
		})
		v.Visit(n)
	case "db":
		db := innitdb.NewInMemory()
		innitdb.Store(db, n, 1)
		var h maphash.Hash
		h.SetSeed(db.Seed())
		hash.Node(&h, n)
		rootId := h.Sum64()
		var visited []uint64
		frontier := []uint64{rootId}
		for len(frontier) > 0 {
			id := frontier[0]
			visited = append(visited, id)
			frontier = frontier[1:]
			db.EachRef(id, func(childId uint64) {
				frontier = append(frontier, childId)
			})
		}
		printed := make(map[uint64]bool)
		for _, id := range visited {
			if printed[id] {
				continue
			}
			printed[id] = true
			fmt.Printf("%d\t", id)
			fmt.Println()
			// n := innitdb.Load(db, id)
			// innit.StdPrinter(os.Stdout).Print(n)
		}
	default:
		log.Fatalf("unexpected -print mode: %v", *print)
	}
}

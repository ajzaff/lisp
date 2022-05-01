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
	"github.com/ajzaff/innit/binnit"
	"github.com/ajzaff/innit/hash"
	"github.com/ajzaff/innit/innitdb"
)

var (
	order = flag.String("order", "", `Print order for AST print mode (Optional "reverse". Default uses in-order)`)
	print = flag.String("print", "", `Print mode (Optional "tok", "ast", "db", "bin". Default uses StdPrinter)`)
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

	ns, err := innit.Parse(string(src))
	if err != nil {
		log.Fatal(err)
	}

	switch *print {
	case "": // std
		for _, n := range ns {
			innit.StdPrinter(os.Stdout).Print(n.Val())
		}
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
		exprVisitor := func(e innit.Expr) {
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
		v.SetLitVisitor(func(e innit.Lit) {
			fmt.Println("LIT\t", e.String())
		})
		for _, n := range ns {
			v.Visit(n.Val())
		}
	case "db":
		db := innitdb.NewInMemory()
		for _, n := range ns {
			innitdb.Store(db, n.Val(), 1)
		}
		refs := make(map[innitdb.ID]struct {
			innit.Val
			Fc float64
		})
		var h maphash.Hash
		h.SetSeed(db.Seed())
		for _, n := range ns {
			h.Reset()
			hash.Val(&h, n.Val())
			id := h.Sum64()
			fc := innitdb.Load(db, n.Val())
			refs[id] = struct {
				innit.Val
				Fc float64
			}{n.Val(), fc}
		}
		for id := range refs {
			db.EachRef(id, func(id innitdb.ID) bool {
				n, fc := db.Load(id)
				refs[id] = struct {
					innit.Val
					Fc float64
				}{n, fc}
				return true
			})
		}
		visited := make(map[innitdb.ID]bool)
		for id, e := range refs {
			if visited[id] {
				continue
			}
			visited[id] = true
			fmt.Printf("%d\t", id)
			fmt.Printf("%f\t", e.Fc)
			innit.StdPrinter(os.Stdout).Print(e.Val)
		}
	case "bin":
		for _, n := range ns {
			binnit.NewEncoder(os.Stdout).Encode(n.Val())
		}
	default:
		log.Fatalf("unexpected -print mode: %v", *print)
	}
}

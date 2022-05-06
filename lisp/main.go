package main

import (
	"flag"
	"fmt"
	"hash/maphash"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/blisp"
	"github.com/ajzaff/lisp/hash"
	"github.com/ajzaff/lisp/lispdb"
)

var (
	order = flag.String("order", "", `Print order for AST print mode (Optional "reverse". Default uses in-order)`)
	print = flag.String("print", "", `Print mode (Optional "tok", "ast", "db", "bin", "none". Default uses StdPrinter)`)
	file  = flag.String("file", "", "File to read lisp code from.")
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

	ns, err := lisp.Parse(string(src))
	if err != nil {
		log.Fatal(err)
	}

	switch *print {
	case "": // std
		for _, n := range ns {
			lisp.StdPrinter(os.Stdout).Print(n.Val())
		}
	case "tok":
		tokens, err := lisp.Tokenize(string(src))
		if err != nil {
			panic(err)
		}
		for i := 0; i < len(tokens); i += 2 {
			token := string(src[tokens[i]:tokens[i+1]])
			println(token)
		}
	case "ast":
		var v lisp.Visitor
		exprVisitor := func(e lisp.Expr) {
			var sb strings.Builder
			lisp.StdPrinter(&sb).Print(e)
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
		v.SetLitVisitor(func(e lisp.Lit) {
			fmt.Println("LIT\t", e.String())
		})
		for _, n := range ns {
			v.Visit(n.Val())
		}
	case "db":
		db := lispdb.NewInMemory()
		for _, n := range ns {
			lispdb.Store(db, n.Val(), 1)
		}
		refs := make(map[lispdb.ID]struct {
			lisp.Val
			Fc float64
		})
		var h maphash.Hash
		h.SetSeed(db.Seed())
		for _, n := range ns {
			h.Reset()
			hash.Val(&h, n.Val())
			id := h.Sum64()
			fc := lispdb.Load(db, n.Val())
			refs[id] = struct {
				lisp.Val
				Fc float64
			}{n.Val(), fc}
		}
		for id := range refs {
			db.EachRef(id, func(id lispdb.ID) bool {
				n, fc := db.Load(id)
				refs[id] = struct {
					lisp.Val
					Fc float64
				}{n, fc}
				return true
			})
		}
		visited := make(map[lispdb.ID]bool)
		for id, e := range refs {
			if visited[id] {
				continue
			}
			visited[id] = true
			fmt.Printf("%d\t", id)
			fmt.Printf("%f\t", e.Fc)
			lisp.StdPrinter(os.Stdout).Print(e.Val)
		}
	case "bin":
		for _, n := range ns {
			blisp.NewEncoder(os.Stdout).Encode(n.Val())
		}
	case "none":
	default:
		log.Fatalf("unexpected -print mode: %v", *print)
	}
}

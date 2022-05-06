//go:build print
// +build print

package main

import (
	"flag"
	"fmt"
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
		log.Fatal("-file is required")
		return
	}

	src, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatal(err)
	}

	ns, err := lisp.Parser{}.Parse(string(src))
	if err != nil {
		log.Fatal(err)
	}

	switch *print {
	case "": // std
		for _, n := range ns {
			lisp.StdPrinter(os.Stdout).Print(n.Val())
		}
	case "tok":
		tokens, err := lisp.Tokenizer{}.Tokenize(string(src))
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
			Fc          float64
			Refs        []uint64
			InverseRefs []uint64
		})
		var h hash.MapHash
		h.SetSeed(db.Seed())
		for _, n := range ns {
			h.Reset()
			h.WriteVal(n.Val())
			root := h.Sum64()
			lispdb.EachTransRef(db, root, func(i lispdb.ID) bool {
				v, w := lispdb.QueryOneID(db, i)
				var idRefs []lispdb.ID
				db.EachRef(i, func(i lispdb.ID) bool {
					idRefs = append(idRefs, i)
					return true
				})
				var idInverseRefs []lispdb.ID
				db.EachInverseRef(i, func(i lispdb.ID) bool {
					idInverseRefs = append(idInverseRefs, i)
					return true
				})
				refs[i] = struct {
					lisp.Val
					Fc          float64
					Refs        []uint64
					InverseRefs []uint64
				}{v, w, idRefs, idInverseRefs}
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
			fmt.Print("Refs: ")
			for _, i := range e.Refs {
				fmt.Printf("%d,", i)
			}
			fmt.Println()
			fmt.Print("InverseRefs: ")
			for _, i := range e.InverseRefs {
				fmt.Printf("%d,", i)
			}
			fmt.Println()
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

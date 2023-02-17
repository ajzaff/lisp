// Binary
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/print"
	"github.com/ajzaff/lisp/scan"
	"github.com/ajzaff/lisp/visit"
	"github.com/ajzaff/lisp/x/blisp"
	"github.com/ajzaff/lisp/x/hash"
	"github.com/ajzaff/lisp/x/lispdb"
	"github.com/ajzaff/lisp/x/lispjson"
	"github.com/ajzaff/lisp/x/stringer"
	"golang.org/x/text/unicode/rangetable"
)

var (
	order = flag.String("order", "", `Print order for AST print mode (Optional "reverse". Default uses in-order)`)
	mode  = flag.String("mode", "", `Print mode (Optional "tok", "ast", "db", "bin", "json", "idtab", "none". Default uses StdPrinter)`)
	file  = flag.String("file", "", "File to read lisp code from.")
)

var tokStr = []string{"?", "Id", "Nat", "(", ")"}

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

	var vs []lisp.Val
	var s scan.TokenScanner
	s.Reset(bytes.NewReader(src))
	var sc scan.NodeScanner
	sc.Reset(&s)
	for sc.Scan() {
		_, _, v := sc.Node()
		vs = append(vs, v)
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	switch *mode {
	case "": // std
		for _, v := range vs {
			print.StdPrinter(os.Stdout).Print(v)
		}
	case "tok":
		var sc scan.TokenScanner
		sc.Reset(bytes.NewReader(src))
		for sc.Scan() {
			pos, tok, text := sc.Token()
			println(strconv.Itoa(int(pos)), "\t", tokStr[tok], "\t", text)
		}
	case "ast":
		var v visit.Visitor
		consVisitor := func(e *lisp.Cons) {
			var sb strings.Builder
			print.StdPrinter(&sb).Print(e)
			fmt.Print("EXPR\t", sb.String())
		}
		switch *order {
		case "": // in-order
			v.SetBeforeConsVisitor(consVisitor)
		case "reverse":
			v.SetAfterConsVisitor(consVisitor)
		default:
			log.Fatalf("unexpected -order mode: %v", *order)
		}
		v.SetLitVisitor(func(e lisp.Lit) {
			fmt.Println("LIT\t", stringer.Lit(e))
		})
		for _, x := range vs {
			v.Visit(x)
		}
	case "db":
		db := lispdb.NewInMemory()
		lispdb.Store(db, vs, 1)
		refs := make(map[lispdb.ID]struct {
			lisp.Val
			Fc          float64
			Refs        []uint64
			InverseRefs []uint64
		})
		var h hash.MapHash
		h.SetSeed(db.Seed())
		for _, v := range vs {
			h.Reset()
			h.WriteVal(v)
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
			print.StdPrinter(os.Stdout).Print(e.Val)
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
		var e blisp.Encoder
		e.Reset(os.Stdout)
		e.EncodeMagic()
		for _, v := range vs {
			e.Encode(v)
		}
	case "json":
		for _, v := range vs {
			lispjson.NewEncoder(os.Stdout).Encode(v)
			println()
		}
	case "idtab":
		t := rangetable.Merge(unicode.Letter)
		rangetable.Visit(t, func(r rune) {
			fmt.Print(string(r))
		})
	case "none":
	default:
		log.Fatalf("unexpected -print mode: %v", *mode)
	}
}

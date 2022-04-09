package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/ajzaff/innit"
)

var (
	order = flag.String("order", "", `Print order for AST print mode (Optional "reverse". Default uses in-order)`)
	print = flag.String("print", "", `Print mode (Optional "tok", "ast". Default uses StdPrinter)`)
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
	default:
		log.Fatalf("unexpected -print mode: %v", *print)
	}
}

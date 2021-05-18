package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ajzaff/innit"
)

var (
	tokenize = flag.Bool("tokenize", false, "Print tokens and exit")
	file     = flag.String("file", "", "File to read innit code from.")
)

func main() {
	flag.Parse()

	src, err := ioutil.ReadFile(*file)
	if err != nil {
		panic(err)
	}

	toks, err := innit.Tokenize(src)
	if err != nil {
		panic(err)
	}

	if *tokenize {
		for i := 0; i < len(toks); i += 2 {
			println(string(src[toks[i]:toks[i+1]]))
		}
		os.Exit(0)
	}

	n, err := innit.Parse(src)
	if err != nil {
		panic(err)
	}

	printNode(n, "", "  ")
}

func printNode(n innit.Node, prefix, indent string) {
	switch n := n.(type) {
	case *innit.BasicLit:
		fmt.Printf("%s%#v\n", prefix, *n)
	case *innit.Expr:
		fmt.Printf("%sinnit.Expr(\n", prefix)
		printNode(n.X, prefix+indent, indent)
		fmt.Printf("%s)\n", prefix)
	case innit.NodeList:
		for _, x := range n {
			printNode(x, prefix, indent)
		}
	}
}

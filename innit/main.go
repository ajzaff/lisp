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
	compact  = flag.Bool("compact", true, "Use compact printing.")
	file     = flag.String("file", "", "File to read innit code from.")
)

func main() {
	flag.Parse()

	if *file == "" {
		doRepl()
		return
	}

	src, err := ioutil.ReadFile(*file)
	if err != nil {
		panic(err)
	}

	toks, err := innit.Tokenize(string(src))
	if err != nil {
		panic(err)
	}

	if *tokenize {
		for i := 0; i < len(toks); i += 2 {
			println(string(src[toks[i]:toks[i+1]]))
		}
		os.Exit(0)
	}

	n, err := innit.Parse(string(src))
	if err != nil {
		panic(err)
	}

	if *compact {
		innit.CompactPrinter(os.Stdout).Print(n)
		fmt.Println()
		os.Exit(0)
	}

	innit.StdPrinter(os.Stdout).Print(n)
}

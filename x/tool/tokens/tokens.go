package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ajzaff/lisp/scan"
)

func main() {
	filename := flag.String("filename", "", "Filename to read Lisp")
	flag.Parse()

	var input string
	switch {
	case *filename != "":
		data, err := os.ReadFile(*filename)
		if err != nil {
			log.Fatal(err)
		}
		input = string(data)
	default:
		input = flag.Arg(0)
	}

	var sc scan.Scanner
	sc.Reset(strings.NewReader(input))
	for token := range sc.Tokens() {
		fmt.Printf("%d %-4d %-40s\n", token.Tok, token.Pos, token.Text)
	}
}

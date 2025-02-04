package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ajzaff/lisp/scan"
	"github.com/ajzaff/lisp/x/print"
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
	var p print.Printer
	p.Reset(os.Stdout)
	for node := range sc.Nodes() {
		fmt.Printf("%-4d %-4d ", node.Pos, node.End)
		p.Print(node.Val)
		fmt.Println()
	}
}

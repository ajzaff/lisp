// Binary utf8string encodes string data in canonical format.
package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()

	for _, s := range flag.Args() {
		fmt.Print("(")
		fmt.Printf("u")
		for _, b := range []byte(s) {
			fmt.Print(" ")
			fmt.Print(b)
		}
		fmt.Println(")")
	}
}

// Binary ebnfgen generates the ebnf grammar for the language.
//
// The grammar is complex because the definition of "unicode letter" is quite involved.
package main

import (
	"fmt"
	"unicode"
)

func main() {
	fmt.Println(`ws = { " " | "\t" | "\n" | "\r" }.`)
	fmt.Println(`nat = "0" | "1" … "9" { "0" … "9" }.`)
	outputIdProds()
	fmt.Println(`expr = nat | id | "(" { ws expr ws } ")".`)
	fmt.Println("lisp = { ws expr ws }.")
}

func outputIdProds() {
	prod := 0

	for _, r16 := range unicode.Letter.R16 {
		if r16.Stride == 1 {
			fmt.Printf("id%d = %q … %q.\n", prod, r16.Lo, r16.Hi)
			prod++
		} else {
			fmt.Printf("id%d = %q", prod, r16.Lo)
			for r := r16.Lo + r16.Stride; r <= r16.Hi; r += r16.Stride {
				fmt.Printf(" | %q", r)
			}
			fmt.Println(".")
			prod++
		}
	}

	for _, r32 := range unicode.Letter.R32 {
		if r32.Stride == 1 {
			fmt.Printf("id%d = %q … %q.\n", prod, r32.Lo, r32.Hi)
			prod++
		} else {
			fmt.Printf("id%d = %q", prod, r32.Lo)
			for r := r32.Lo + r32.Stride; r <= r32.Hi; r += r32.Stride {
				fmt.Printf(" | %q", r)
			}
			fmt.Println(".")
			prod++
		}
	}

	fmt.Print("id = id0")
	for i := 1; i < prod; i++ {
		fmt.Printf(" | id%d", i)
	}
	fmt.Println(".")
}

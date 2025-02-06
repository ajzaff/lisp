// Binary ebnfgen generates the ebnf grammar for the language.
//
// The grammar is complex because the definition of "unicode letter" is quite involved.
package main

import (
	"fmt"
	"unicode"
)

func main() {
	fmt.Println(`// Whitespace.
s0 = " " | "\t" | "\r" | "\n".
s1 = {s0}.
s2 = s0 s1.
d0 = "0" … "9".`)
	fmt.Println()

	outputIdProds()

	fmt.Println(`l1 = d0 | l0.
l2 = l1 { l1 }.
l3 = l2 { s2 l2 }.

// Groups.
g0 = "(" s1 e2 s1 ")".
g1 = g0 {s1 g0}. 

// Expressions.
e0  = g0 | l2.
e1 = e0 {s1 e0}.
e2 = "" | e1.
e3 = s1 e2 s1.`)
}

func outputIdProds() {
	fmt.Println("// Literals.")

	prod := 0
	for _, r16 := range unicode.Letter.R16 {
		if r16.Stride == 1 {
			fmt.Printf("u%d = %q … %q.\n", prod, r16.Lo, r16.Hi)
			prod++
		} else {
			fmt.Printf("u%d = %q", prod, r16.Lo)
			for r := r16.Lo + r16.Stride; r <= r16.Hi; r += r16.Stride {
				fmt.Printf(" | %q", r)
			}
			fmt.Println(".")
			prod++
		}
	}

	for _, r32 := range unicode.Letter.R32 {
		if r32.Stride == 1 {
			fmt.Printf("u%d = %q … %q.\n", prod, r32.Lo, r32.Hi)
			prod++
		} else {
			fmt.Printf("u%d = %q", prod, r32.Lo)
			for r := r32.Lo + r32.Stride; r <= r32.Hi; r += r32.Stride {
				fmt.Printf(" | %q", r)
			}
			fmt.Println(".")
			prod++
		}
	}

	fmt.Print("l0 = u0")
	for i := 1; i < prod; i++ {
		fmt.Printf(" | u%d", i)
	}
	fmt.Println(".")
}

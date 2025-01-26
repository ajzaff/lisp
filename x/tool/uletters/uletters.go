// Binary uletters prints a table of all valid unicode letter characters.
package main

import (
	"fmt"
	"unicode"

	"golang.org/x/text/unicode/rangetable"
)

func main() {
	pos := 0
	const cols = 40
	rangetable.Visit(unicode.Letter, func(r rune) {
		fmt.Print(string(r))
		pos++
		if pos > 0 && pos%cols == 0 {
			fmt.Println()
		}
	})
	if pos%cols != 0 {
		fmt.Println()
	}
}

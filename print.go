package innit

import (
	"fmt"
	"io"
)

type Printer struct {
	W io.Writer

	Prefix, Indent, NewLine string
}

func StdPrinter(w io.Writer) *Printer {
	return &Printer{w, "", "  ", "\n"}
}

func CompactPrinter(w io.Writer) *Printer {
	return &Printer{w, "", "", ""}
}

func (p *Printer) Print(n []Node) {
	for i, x := range n {
		printRec(x, p.W, p.Prefix, p.Indent)
		endl := p.NewLine
		if endl == "" && i < len(n) {
			if _, ok := x.(*BasicLit); ok {
				endl = " "
			}
		}
		fmt.Fprint(p.W, endl)
	}
}

func printRec(n Node, w io.Writer, prefix, indent string) {
	switch n := n.(type) {
	case *BasicLit:
		fmt.Fprint(w, n.Value)
	case *Expr:
		fmt.Fprintf(w, "%s(", prefix)
		for i, x := range n.X {
			printRec(x, w, prefix, indent)
			if i+1 < len(n.X) {
				fmt.Fprint(w, " ")
			}
		}
		fmt.Printf(")")
	}
}

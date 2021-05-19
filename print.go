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

func (p *Printer) Print(n Node) {
	list, ok := n.(NodeList)
	if !ok {
		printRec(n, p.W, p.Prefix, p.Indent)
		return
	}
	for i, x := range list {
		printRec(x, p.W, p.Prefix, p.Indent)
		endl := p.NewLine
		if endl == "" && i < len(list) {
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
		printRec(n.X, w, prefix, indent)
		fmt.Printf(")")
	case NodeList:
		for i, x := range n {
			printRec(x, w, prefix, indent)
			if i+1 < len(n) {
				fmt.Fprint(w, " ")
			}
		}
	default:
		panic(fmt.Errorf("innit.Print: internal error: unexpected node type: %T", n))
	}
}

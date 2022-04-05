package innit

import (
	"fmt"
	"io"
)

// Printer implements direct printing of AST nodes.
type Printer struct {
	io.Writer

	Prefix, Indent, NewLine string
}

// StdPrinter returns a printer which uses spaces and new lines.
func StdPrinter(w io.Writer) *Printer {
	return &Printer{
		Writer:  w,
		Indent:  "  ",
		NewLine: "\n",
	}
}

// CompactPrinter returns a printer using the least characters possible.
func CompactPrinter(w io.Writer) *Printer {
	return &Printer{Writer: w}
}

// Print the Node n.
func (p *Printer) Print(n Node) {
	var stack []int

	var v Visitor
	v.SetExprVisitor(func(e *Expr) {
		stack = append(stack, len(e.X))
		fmt.Fprint(p.Writer, p.Prefix, "(")
	})
	v.SetLitVisitor(func(e *Lit) {
		switch n := len(stack); n {
		case 0:
			fmt.Fprint(p.Writer, p.Prefix, e.Value)
		default:
			fmt.Fprint(p.Writer, e.Value)
			if stack[n-1]--; stack[n-1] <= 0 {
				fmt.Fprint(p.Writer, ")")
				stack = stack[:n-1]
				if len(stack) == 0 {
					fmt.Fprint(p.Writer, p.NewLine)
				}
			}
		}
		fmt.Fprint(p.Writer, " ")
	})

	v.Visit(n)

	for i := 0; i < len(stack); i++ {
		fmt.Fprint(p.Writer, ")")
	}
}

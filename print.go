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
	var (
		firstExpr = true
		exprDepth int
	)
	var v Visitor
	v.SetBeforeExprVisitor(func(e *Expr) {
		if exprDepth > 0 {
			fmt.Fprint(p.Writer, " ")
		}
		fmt.Fprint(p.Writer, p.Prefix, "(")
		firstExpr = true
		exprDepth++
	})
	v.SetAfterExprVisitor(func(e *Expr) {
		exprDepth--
		fmt.Fprint(p.Writer, ")")
		if exprDepth == 0 {
			fmt.Fprint(p.Writer, p.NewLine)
		}
	})
	v.SetLitVisitor(func(e *Lit) {
		if exprDepth == 0 {
			fmt.Fprint(p.Writer, p.Prefix, e.Value, p.NewLine)
			return
		}
		if !firstExpr {
			fmt.Fprint(p.Writer, " ")
		}
		firstExpr = false
		fmt.Fprint(p.Writer, e.Value)
	})

	v.Visit(n)
}

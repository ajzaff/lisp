package innit

import (
	"fmt"
	"io"
)

// Printer implements direct printing of AST nodes.
type Printer struct {
	io.Writer

	Nil string

	Prefix, Indent, NewLine string
}

// StdPrinter returns a printer which uses spaces and new lines.
func StdPrinter(w io.Writer) *Printer {
	return &Printer{
		Writer:  w,
		Nil:     "<nil>",
		Indent:  "  ",
		NewLine: "\n",
	}
}

// Print the Node n.
func (p *Printer) Print(n Node) {
	if n == nil {
		p.Write([]byte(p.Nil))
		return
	}
	var (
		exprDepth  int
		firstWrite = true
		firstLit   = true
		newLine    = fmt.Sprint(p.NewLine, p.Prefix)
	)
	var v Visitor
	v.SetBeforeExprVisitor(func(e *Expr) {
		if !firstWrite {
			fmt.Fprint(p.Writer, p.Prefix)
			firstWrite = true
		}
		fmt.Fprint(p.Writer, "(")
		firstLit = false
		exprDepth++
	})
	v.SetAfterExprVisitor(func(e *Expr) {
		exprDepth--
		fmt.Fprint(p.Writer, ")")
		if exprDepth == 0 {
			fmt.Fprint(p.Writer, newLine)
		}
		firstLit = false
	})
	v.SetLitVisitor(func(e *Lit) {
		if !firstWrite {
			fmt.Fprint(p.Writer, p.Prefix)
			firstWrite = true
		}
		if exprDepth == 0 {
			fmt.Fprint(p.Writer, e.Value, newLine)
			return
		}
		if firstLit {
			fmt.Fprint(p.Writer, " ")
		}
		firstLit = true
		fmt.Fprint(p.Writer, e.Value)
	})

	v.Visit(n)
}

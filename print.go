package lisp

import (
	"fmt"
	"io"
	"unicode/utf8"
)

// Printer implements direct printing of AST nodes.
type Printer struct {
	io.Writer

	Nil string

	Prefix, NewLine string
}

// StdPrinter returns a printer which uses spaces and new lines.
func StdPrinter(w io.Writer) *Printer {
	return &Printer{
		Writer:  w,
		Nil:     "<nil>",
		NewLine: "\n",
	}
}

// Print the Node n.
func (p *Printer) Print(n Val) {
	if n == nil {
		p.Write([]byte(p.Nil))
		p.Write([]byte(p.NewLine))
		return
	}
	var (
		exprDepth  int
		firstWrite = true
		lastLetter = true
		newLine    = fmt.Sprint(p.NewLine, p.Prefix)
	)
	var v Visitor
	v.SetBeforeExprVisitor(func(e Expr) {
		if !firstWrite {
			fmt.Fprint(p.Writer, p.Prefix)
			firstWrite = true
		}
		fmt.Fprint(p.Writer, "(")
		lastLetter = false
		exprDepth++
	})
	v.SetAfterExprVisitor(func(e Expr) {
		exprDepth--
		fmt.Fprint(p.Writer, ")")
		if exprDepth == 0 {
			fmt.Fprint(p.Writer, newLine)
		}
		lastLetter = false
	})
	v.SetLitVisitor(func(e Lit) {
		if !firstWrite {
			fmt.Fprint(p.Writer, p.Prefix)
			firstWrite = true
		}
		if exprDepth == 0 {
			fmt.Fprint(p.Writer, e.String(), newLine)
			return
		}
		r, _ := utf8.DecodeRuneInString(e.String())
		currLetter := IsLetter(r)
		if lastLetter && currLetter {
			fmt.Fprint(p.Writer, " ")
		}
		lastLetter = currLetter
		fmt.Fprint(p.Writer, e.String())
	})
	v.Visit(n)
}

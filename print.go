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
func (p *Printer) Print(n Val) {
	if n == nil {
		p.Write([]byte(p.Nil))
		p.Write([]byte(p.NewLine))
		return
	}
	var (
		exprDepth  int
		firstWrite = true
		lastIdent  = true
		newLine    = fmt.Sprint(p.NewLine, p.Prefix)
	)
	var v Visitor
	v.SetBeforeExprVisitor(func(e Expr) {
		if !firstWrite {
			fmt.Fprint(p.Writer, p.Prefix)
			firstWrite = true
		}
		fmt.Fprint(p.Writer, "(")
		lastIdent = false
		exprDepth++
	})
	v.SetAfterExprVisitor(func(e Expr) {
		exprDepth--
		fmt.Fprint(p.Writer, ")")
		if exprDepth == 0 {
			fmt.Fprint(p.Writer, newLine)
		}
		lastIdent = false
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
		currIdent := IsIdent(r)
		if lastIdent && currIdent {
			fmt.Fprint(p.Writer, " ")
		}
		lastIdent = currIdent
		fmt.Fprint(p.Writer, e.String())
	})
	v.Visit(n)
}

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
		exprDepth       int
		firstWrite      = true
		lastDelimitable delimitable
		newLine         = fmt.Sprint(p.NewLine, p.Prefix)
	)
	var v Visitor
	v.SetBeforeExprVisitor(func(e Expr) {
		if !firstWrite {
			fmt.Fprint(p.Writer, p.Prefix)
			firstWrite = true
		}
		fmt.Fprint(p.Writer, "(")
		lastDelimitable = delimitableNone
		exprDepth++
	})
	v.SetAfterExprVisitor(func(e Expr) {
		exprDepth--
		fmt.Fprint(p.Writer, ")")
		if exprDepth == 0 {
			fmt.Fprint(p.Writer, newLine)
		}
		lastDelimitable = delimitableNone
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
		delim := delimitableLitType(e)
		if lastDelimitable != delimitableNone && lastDelimitable == delim {
			fmt.Fprint(p.Writer, " ")
		}
		lastDelimitable = delim
		fmt.Fprint(p.Writer, e.String())
	})
	v.Visit(n)
}

// Lits in the same delimitable class must be spaced out.
type delimitable int

const (
	delimitableNone   delimitable = iota
	delimitableClass1             // Number
	delimitableClass2             // Symbol
)

func delimitableLitType(e Lit) delimitable {
	switch e.Token {
	case Id:
		r, _ := utf8.DecodeRuneInString(e.String())
		if IsLetter(r) {
			return delimitableClass1
		}
		return delimitableClass2
	case Number:
		return delimitableClass1
	default:
		return delimitableNone
	}
}

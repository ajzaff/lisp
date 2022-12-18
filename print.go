package lisp

import (
	"fmt"
	"io"
	"sync"
)

// Printer implements direct printing of AST nodes.
type Printer struct {
	w io.Writer

	v    Visitor // init once
	once sync.Once

	PrinterOptions
}

// PrinterOptions supplied to the Printer.
type PrinterOptions struct {
	Nil             string
	Prefix, NewLine string
}

func makeStdPrinterOptions() PrinterOptions {
	return PrinterOptions{
		Nil:     "(nil)",
		NewLine: "\n",
	}
}

func (p *Printer) initVisitor() {
	var (
		exprDepth       int
		firstWrite      = true
		lastDelimitable delimitable
	)
	p.v.SetBeforeExprVisitor(func(e Expr) {
		if !firstWrite {
			fmt.Fprint(p.w, p.Prefix)
			firstWrite = true
		}
		fmt.Fprint(p.w, "(")
		lastDelimitable = delimitableNone
		exprDepth++
	})
	p.v.SetAfterExprVisitor(func(e Expr) {
		exprDepth--
		fmt.Fprint(p.w, ")")
		if exprDepth == 0 {
			fmt.Fprint(p.w, p.NewLine, p.Prefix)
		}
		lastDelimitable = delimitableNone
	})
	p.v.SetLitVisitor(func(e Lit) {
		if !firstWrite {
			fmt.Fprint(p.w, p.Prefix)
			firstWrite = true
		}
		if exprDepth == 0 {
			fmt.Fprint(p.w, e.String(), p.NewLine, p.Prefix)
			return
		}
		delim := delimitableLitType(e)
		if lastDelimitable != delimitableNone && lastDelimitable == delim {
			fmt.Fprint(p.w, " ")
		}
		lastDelimitable = delim
		fmt.Fprint(p.w, e.String())
	})
}

// Reset resets the Printer to use the given writer.
func (p *Printer) Reset(w io.Writer) {
	p.w = w
}

// StdPrinter returns a printer which uses spaces and new lines.
func StdPrinter(w io.Writer) *Printer {
	var p Printer
	p.PrinterOptions = makeStdPrinterOptions()
	p.Reset(w)
	return &p
}

// Print the Node n.
func (p *Printer) Print(n Val) {
	if n == nil {
		p.w.Write([]byte(p.Nil))
		p.w.Write([]byte(p.NewLine))
		return
	}
	p.once.Do(p.initVisitor)
	p.v.Visit(n)
}

// Lits in the same delimitable class must be spaced out.
type delimitable int

const (
	delimitableNone   delimitable = iota
	delimitableClass1             // Id, Number
)

func delimitableLitType(e Lit) delimitable {
	switch e.Token {
	case Id, Int:
		return delimitableClass1
	default:
		return delimitableNone
	}
}

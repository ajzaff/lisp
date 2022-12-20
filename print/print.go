package print

import (
	"bufio"
	"io"
	"sync"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/visit"
)

// Printer implements direct printing of AST nodes.
type Printer struct {
	w *bufio.Writer

	v    visit.Visitor // init once
	once sync.Once

	PrinterOptions
}

// PrinterOptions supplied to the Printer.
type PrinterOptions struct {
	Nil     string
	Prefix  string // Prefix added before every line.
	NewLine string // NewLine added after every top-level expression.
}

func makeStdPrinterOptions() PrinterOptions {
	return PrinterOptions{
		Nil:     "()",
		NewLine: "\n",
	}
}

func (p *Printer) initVisitor() {
	var (
		consDepth       int
		lastDelimitable delimitable
	)
	p.v.SetBeforeConsVisitor(func(e *lisp.Cons) {
		p.w.WriteByte('(')
		lastDelimitable = delimitableNone
		consDepth++
	})
	p.v.SetAfterConsVisitor(func(e *lisp.Cons) {
		consDepth--
		p.w.WriteByte(')')
		if consDepth == 0 {
			p.w.WriteString(p.NewLine)
			p.w.WriteString(p.Prefix)
		}
		lastDelimitable = delimitableNone
	})
	p.v.SetLitVisitor(func(e lisp.Lit) {
		delim := delimitableLitType(e)
		if lastDelimitable != delimitableNone && lastDelimitable == delim {
			p.w.WriteByte(' ')
		}
		lastDelimitable = delim
		p.w.WriteString(e.Text)
		if consDepth == 0 {
			p.w.WriteString(p.NewLine)
			p.w.WriteString(p.Prefix)
		}
	})
}

// Reset resets the Printer to use the given writer.
func (p *Printer) Reset(w io.Writer) {
	p.w = bufio.NewWriter(w)
}

// StdPrinter returns a printer which uses spaces and new lines.
func StdPrinter(w io.Writer) *Printer {
	var p Printer
	p.PrinterOptions = makeStdPrinterOptions()
	p.Reset(w)
	return &p
}

// Print the Node n.
func (p *Printer) Print(n lisp.Val) {
	defer p.w.Flush()
	if n == nil {
		p.w.Write([]byte(p.Nil))
		p.w.Write([]byte(p.NewLine))
		return
	}
	p.once.Do(p.initVisitor)
	p.w.WriteString(p.Prefix)
	p.v.Visit(n)
}

// Lits in the same delimitable class must be spaced out.
type delimitable int

const (
	delimitableNone   delimitable = iota
	delimitableClass1             // Id, Number
)

func delimitableLitType(e lisp.Lit) delimitable {
	switch e.Token {
	case lisp.Id, lisp.Int:
		return delimitableClass1
	default:
		return delimitableNone
	}
}

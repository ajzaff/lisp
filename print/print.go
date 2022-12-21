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
	NewLine bool   // Whether a new line is added after every top-level expression.
}

func makeStdPrinterOptions() PrinterOptions {
	return PrinterOptions{
		Nil:     "()",
		NewLine: true,
	}
}

func (p *Printer) initVisitor() {
	var (
		consDepth int
		delim     bool
	)
	p.v.SetLitVisitor(func(x lisp.Lit) {
		if delim {
			p.w.WriteByte(' ')
		}
		p.w.WriteString(x.Text)
		switch consDepth == 0 && p.NewLine {
		case true:
			p.w.WriteByte('\n')
			p.w.WriteString(p.Prefix)
			delim = false
		default:
			delim = true
		}
	})
	p.v.SetBeforeConsVisitor(func(*lisp.Cons) { p.w.WriteByte('('); delim = false; consDepth++ })
	p.v.SetAfterConsVisitor(func(*lisp.Cons) {
		p.w.WriteByte(')')
		delim = false
		if consDepth--; consDepth == 0 && p.NewLine {
			p.w.WriteByte('\n')
			p.w.WriteString(p.Prefix)
		}
	})
}

// Reset resets the Printer to use the given writer.
func (p *Printer) Reset(w io.Writer) {
	if w, ok := w.(*bufio.Writer); ok {
		p.w = w
		return
	}
	if p.w == nil {
		p.w = new(bufio.Writer)
	}
	p.w.Reset(w)
}

// StdPrinter returns a printer which uses spaces and new lines.
func StdPrinter(w io.Writer) *Printer {
	var p Printer
	p.PrinterOptions = makeStdPrinterOptions()
	p.Reset(w)
	return &p
}

// Print the Val v.
func (p *Printer) Print(v lisp.Val) {
	defer p.w.Flush()
	if v == nil {
		p.w.WriteString(p.Nil)
		if p.NewLine {
			p.w.WriteByte('\n')
		}
		return
	}
	p.once.Do(p.initVisitor)
	p.w.WriteString(p.Prefix)
	p.v.Visit(v)
}

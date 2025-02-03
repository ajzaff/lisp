package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/rangetable"
)

var (
	rules     map[string][]rule
	terminals map[string]string
)

func main() {
	var (
		filename = flag.String("filename", "", "Input filename")
		start    = flag.String("start", "(e3)", "Start symbol")
	)
	flag.Parse()

	var src string

	switch *filename {
	case "":
		src = flag.Arg(0)
	default:
		data, err := os.ReadFile(*filename)
		if err != nil {
			log.Fatal(err)
		}
		src = string(data)
	}

	stack := make([]entry, 0, 512)
	stack = append(stack, entry{deriv: "", nonterminals: *start})

	for len(stack) > 0 {
		n := len(stack) - 1
		e := stack[n]
		stack = stack[:n]

		// fmt.Printf("%q %q\n", e.deriv, e.nonterminals)

		if strings.HasPrefix(src, e.deriv) {
			if e.nonterminals == "" && e.deriv == src {
				e.PrintDerivation()
				break
			}
		} else {
			// Give up on this derivation.
			continue
		}

		elem := e.Elem()
		rs := rules[elem]

		if rs == nil {
			// No nonterminal found.
			// Try letter or terminal.
			switch elem {
			case "(a unicode)":
				// Special handling for unicode alternate.
				// Avoids large branches.
				r, _ := utf8.DecodeRuneInString(src[len(e.deriv):])
				next, ok := e.ApplyLetter(r)
				if ok {
					stack = append(stack, next)
				}
				continue
			default:
				// Try looking up terminals.
				t, ok := terminals[elem]
				if !ok {
					// Terminal not found?
					continue
				}
				// All terminals derived here.
				if !strings.HasPrefix(src[len(e.deriv):], t) {
					// Skip this impossible derivation.
					continue
				}
				next, _ := e.ApplyTerminal(elem, t)
				stack = append(stack, next)
			}
		} else {
			for _, r := range rs {
				next, ok := e.Apply(r)
				if ok {
					stack = append(stack, next)
				}
			}
		}
	}
}

type entry struct {
	prev         *entry
	deriv        string
	nonterminals string
}

func (e entry) Elem() string {
	i := strings.IndexByte(e.nonterminals, ')')
	return e.nonterminals[:i+1]
}

func (e entry) ApplyLetter(r rune) (entry, bool) {
	if !unicode.IsLetter(r) {
		return entry{}, false
	}
	next := entry{}
	next.prev = &e
	next.deriv = e.deriv
	const u = "(a unicode)"
	const ul = len(u)
	ut := fmt.Sprintf("(%c)", r)
	next.nonterminals = fmt.Sprint(ut, e.nonterminals[ul:])
	return next, true
}

func (e entry) ApplyTerminal(lhs, rhs string) (entry, bool) {
	next := entry{}
	next.prev = &e
	next.deriv = fmt.Sprint(e.deriv, rhs)
	next.nonterminals = e.nonterminals[len(lhs):]
	return next, true
}

func (e entry) Apply(r rule) (entry, bool) {
	next := entry{}
	if !strings.HasPrefix(e.nonterminals, r.lhs) {
		return entry{}, false
	}
	next.prev = &e
	next.deriv = e.deriv
	next.nonterminals = fmt.Sprint(r.rhs, e.nonterminals[len(r.lhs):])
	return next, true
}

func (e entry) Print(i int) {
	fmt.Printf("%-4d%-20q%40q\n", i, e.deriv, e.nonterminals)
}

func (e entry) PrintDerivation() {
	entries := make([]entry, 0, 40)
	for e := &e; e != nil; e = e.prev {
		entries = append(entries, *e)
	}
	fmt.Printf("%-4s%-20s%40s\n", "NUM", "TXT", "RULES")
	for i := len(entries) - 1; i >= 0; i-- {
		entries[i].Print(len(entries) - i - 1)
	}
}

func init() {
	// Add unicode terminals.
	terminals = make(map[string]string, 1000)
	rangetable.Visit(unicode.Letter, func(r rune) {
		terminals[fmt.Sprintf("(%c)", r)] = string(r)
	})
	for _, t := range terminalRules {
		terminals[t.lhs] = t.rhs
	}
	// Append rules.
	rules = make(map[string][]rule, 200)
	for _, r := range nonterminalRules {
		rules[r.lhs] = append(rules[r.lhs], r)
	}
	for _, r := range repeatRules {
		rules[r.lhs] = append(rules[r.lhs], r)
	}
	for _, r := range alternateRules {
		rules[r.lhs] = append(rules[r.lhs], r)
	}
}

type rule struct {
	lhs string
	rhs string
}

// {lhs, rhs}
var nonterminalRules = []rule{
	{"(e0)", "(a g0 l2)"},
	{"(e1)", "(e0)(r s1 e0)"},
	{"(e2)", "(a ε e1)"},
	{"(e3)", "(s1)(e2)(s1)"},
	{"(d0)", "(a 0 1 2 3 4 5 6 7 8 9)"},
	{"(g0)", "(lb)(s1)(e2)(s1)(rb)"},
	{"(g1)", "(g0)(r s1 g0)"},
	{"(l0)", "(a unicode)"},
	{"(l1)", "(a d0 l0)"},
	{"(l2)", "(l1)(r l1)"},
	{"(l3)", "(l2)(r s2 l2)"},
	{"(s0)", "(a sp tb cr nl)"},
	{"(s1)", "(r s0)"},
	{"(s2)", "(s0)(s1)"},
}

// {lhs, rhs}
var repeatRules = []rule{
	{"(r l1)", "(l1)(r l1)"},
	{"(r l1)", "(ε)"},
	{"(r s0)", "(s0)(r s0)"},
	{"(r s0)", "(ε)"},
	{"(r s1 e0)", "(s1)(e0)(r s1 e0)"},
	{"(r s1 e0)", "(ε)"},
	{"(r s1 g0)", "(s1)(g0)(r s1 g0)"},
	{"(r s1 g0)", "(ε)"},
	{"(r s2 l2)", "(s2)(l2)(r s2 l2)"},
	{"(r s2 l2)", "(ε)"},
}

// {lhs, rhs}
var alternateRules = []rule{
	// Expr.
	{"(a ε e1)", "(ε)"},
	{"(a ε e1)", "(e1)"},
	{"(a g0 l2)", "(g0)"},
	{"(a g0 l2)", "(l2)"},
	// Id.
	{"(a d0 l0)", "(d0)"},
	{"(a d0 l0)", "(l0)"},
	// Numbers.
	{"(a 0 1 2 3 4 5 6 7 8 9)", "(0)"},
	{"(a 0 1 2 3 4 5 6 7 8 9)", "(1)"},
	{"(a 0 1 2 3 4 5 6 7 8 9)", "(2)"},
	{"(a 0 1 2 3 4 5 6 7 8 9)", "(3)"},
	{"(a 0 1 2 3 4 5 6 7 8 9)", "(4)"},
	{"(a 0 1 2 3 4 5 6 7 8 9)", "(5)"},
	{"(a 0 1 2 3 4 5 6 7 8 9)", "(6)"},
	{"(a 0 1 2 3 4 5 6 7 8 9)", "(7)"},
	{"(a 0 1 2 3 4 5 6 7 8 9)", "(8)"},
	{"(a 0 1 2 3 4 5 6 7 8 9)", "(9)"},
	// Whitespace.
	{"(a sp tb cr nl)", "(sp)"},
	{"(a sp tb cr nl)", "(tb)"},
	{"(a sp tb cr nl)", "(cr)"},
	{"(a sp tb cr nl)", "(nl)"},
	// Letter (a unicode) added above.
}

type terminal struct {
	lhs string
	rhs string
}

var terminalRules = []terminal{
	{"(ε)", ""},
	{"(lb)", "("},
	{"(rb)", ")"},
	{"(0)", "0"},
	{"(1)", "1"},
	{"(2)", "2"},
	{"(3)", "3"},
	{"(4)", "4"},
	{"(5)", "5"},
	{"(6)", "6"},
	{"(7)", "7"},
	{"(8)", "8"},
	{"(9)", "9"},
	{"(sp)", " "},
	{"(tb)", "\t"},
	{"(cr)", "\r"},
	{"(nl)", "\n"},
	// unicode letters added above.
}

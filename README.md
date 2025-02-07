# lisp

Lisp is a research language.

It features:

* Minimal language constructs and Unicode support.
* Simple API: Tokenizer, Parser, Printer, and Visitor.
* Bring your own libraries.

## Syntax

```
// Whitespace.
s0 = " " | "\t" | "\r" | "\n".
s1 = { s0 }.
s2 = s0 s1.

// Literals.
d0 = "0" … "9".
l0    = ... // Definition of "unicode Letter" omitted.
l1 = d0 | l0.
l2 = l1 { l1 }.
l3 = l2 { s2 l2 }.

// Groups.
g0 = "(" e3 ")".
g1 = g0 { s1 g0 }.

// Expressions.
e0  = g0 | l2.
e1 = e0 { s1 e0 }.
e2 = "" | e1.
e3 = s1 e2 s1.
```

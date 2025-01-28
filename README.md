# lisp

Lisp is a research language.

It features:

* Minimal language constructs and Unicode support.
* Simple API: Tokenizer, Parser, Printer, and Visitor.
* Bring your own libraries.

## Syntax

```
id   = unicode_letter { unicode_letter }.
nat  = "0" … "9" { "0" … "9" }.
cons = "(" { expr } ")"
expr = id | nat | cons.
```
# lisp

Lisp is a research language.

It features:

* Minimal language constructs and Unicode support.
* Simple API: Tokenizer, Parser, Printer, and Visitor.
* Bring your own libraries.

## Syntax

```
id          => \p{Letter}[\p{Letter}\d]*
number      => 0 | [1-9]\d*
cons        => '(' (id | number | cons)* ')'
```
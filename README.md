# lisp

Lisp is a research language.

It features:

* Minimal language constructs and Unicode support.
* Simple API: Tokenizer, Parser, Printer, and Visitor.
* Datastore Optimized for Expression Multimapping, Counting and Querying.

## Syntax

```
id          => \p{Letter}+
number      => 0 | [1-9]\d*
expr        => '(' (id | number | expr)* ')'
```
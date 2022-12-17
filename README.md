# lisp

Lisp is a research language.

It features:

* Minimal language constructs and Unicode support.
* Simple API: Tokenizer, Parser, Printer, and Visitor.
* Datastore Optimized for Expression Multimapping, Counting and Querying.

## Syntax

```
id          => [[:letter:]]+
number      => -? [[:digit:]]+ ('.' [[:digit:]]*)? ([eE] [[:digit]]+)?
expr        => '(' (id | number | expr)+ ')'
```
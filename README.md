# lisp

Lisp is a research language.

It features:

* Minimal language constructs and Unicode support.
* Simple API: Tokenizer, Parser, Printer, and Visitor.
* Datastore Optimized for Expression Multimapping, Counting and Querying.

## Syntax

```
val         => expr | id | number
expr        => '(' val ')'
id          => [[:punct:]]+ | [^[[:punct:]]()"]+
number      => -? [[:digit:]]+ ('.' [[:digit:]]*)? ([eE] [[:digit]]+)?
str         => '"' ([^"] | escape)* '"'
escape      => '\' ('\' | 't' | 'n' | 'x' hex_digit{2})
hex_digit   => [0-9A-Za-z]
space       => ' ' | '\n' | '\t' | '\r'
```
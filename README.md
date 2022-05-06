# lisp

Lisp is a research language.

It features:

* Minimal language constructs and Unicode support.
* Simple API: Tokenizer, Parser, Printer, and Visitor.
* Datastore Optimized for Expression Multimapping, Counting and Querying.

## Syntax

```
val         => expr | id | int | float
expr        => '(' val ')'
id          => id_norm | id_punct
id_norm     => [^[:punct:]()"]+
id_punct    => [[:punct:]]+
int         => dec_digit+
float       => int '.' int? | '.' int
str         => '"' ([^"] | escape)* '"'
escape      => '\' ('\' | 't' | 'n' | 'x' hex_digit{2})
dec_digit   => [0-9]
hex_digit   => [0-9A-Za-z]
space       => ' ' | '\n' | '\t' | '\r'
```
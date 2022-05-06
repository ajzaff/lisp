# innit

Innit is a bare-minimum Lisp-like research language.

It features:

* Minimal language constructs.
* UTF-8 support.
* Built-in Tokenizer, Parser, Printer, and Visitor.
* Some support for extensibility.

## Syntax

```
space       => ' ' | '\n' | '\t' | '\r';
val         => expr | id | int | float;
expr        => '(' val ')';
id          => [[:punct:]]+ | [^[[:punct:]]\s()"]
int         => dec_digit | dec_digit int;
float       => '.' int | int '.' int | int '.';
dec_digit   => [0-9];
str         => '"' str_body '"';
str_body    => '' | str_char str_body;
str_char    => [^"];
escape      => '\' escape_char;
escape_char => '\' | 't' | 'n' | 'x' hex hex;
hex_digit   => [0-9a-z];
```
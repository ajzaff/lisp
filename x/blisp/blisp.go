// Package blisp implements binary Lisp encoding.
// The blisp encoding currently increases message size,
// so it should not be generally used over the wire.
package blisp

const Magic = "\x41blisp\n"

const MagicLen = len(Magic)

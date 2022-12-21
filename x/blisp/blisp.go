// Package blisp implements binary Lisp encoding.
//
// The blisp uses varint encoding for Nats and a representative form for Cons and Ids.
// This can improve the compactness of Nat as well as minimizing use of delimiters.
package blisp

const Magic = "blisp1\n"

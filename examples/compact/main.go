package main

import "math/big"

type expr struct {
	X nodeList
}

type nodeList []interface{}

type ID string

type Constraint interface {
	ID | big.Int | big.Float | string
}

type id struct {
	Value string
}

type int struct {
	big.Int
}

type float struct {
	big.Float
}

type str struct {
	Value string
}

func main() {
}

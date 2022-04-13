package innitdb

import "github.com/ajzaff/innit"

type InnitDB interface {
	Load(uint64) innit.Node
	EachRef(uint64, func(uint64))
	InverseRef(uint64) uint64
}

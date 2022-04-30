package innitdb

import (
	"hash/maphash"

	"github.com/ajzaff/innit"
)

type ID = uint64

type InnitDB interface {
	Seed() maphash.Seed
	Store(innit.Node) ID
	Load(ID) innit.Node
	EachRef(ID, func(ID))
	EachInverseRef(ID) ID
}

func Store(db InnitDB, n innit.Node) ID {
	return db.Store(n)
}

func Load(db InnitDB, id ID) innit.Node {
	return db.Load(id)
}

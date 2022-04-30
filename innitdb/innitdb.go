package innitdb

import (
	"hash/maphash"

	"github.com/ajzaff/innit"
	"github.com/ajzaff/innit/hash"
)

type ID = uint64

type InnitDB interface {
	Seed() maphash.Seed
	Store(innit.Node, float64) ID
	Load(ID) (innit.Node, float64)
	EachRef(ID, func(ID) bool)
	EachInverseRef(ID, func(ID) bool)
}

func Store(db InnitDB, n innit.Node, w float64) ID {
	return db.Store(n, w)
}

func Load(db InnitDB, n innit.Node) (innit.Node, float64) {
	var h maphash.Hash
	h.SetSeed(db.Seed())
	id := hash.Hash(&h, n)
	return db.Load(id)
}

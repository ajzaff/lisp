package lispdb

import (
	"hash/maphash"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/hash"
)

type ID = uint64

type InnitDB interface {
	Seed() maphash.Seed
}

type LoadInterface interface {
	InnitDB
	Load(ID) (lisp.Lit, float64)
}

func Load(db LoadInterface, v lisp.Val) float64 {
	var h maphash.Hash
	h.SetSeed(db.Seed())
	hash.Val(&h, v)
	_, w := db.Load(h.Sum64())
	return w
}

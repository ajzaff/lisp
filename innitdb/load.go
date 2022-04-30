package innitdb

import (
	"hash/maphash"

	"github.com/ajzaff/innit"
	"github.com/ajzaff/innit/hash"
)

type LoadInterface interface {
	Seed() maphash.Seed
	Load(id uint64) (fc int, refs, irefs []uint64)
}

func Load(d LoadInterface, n innit.Node) (fc int, refs, irefs []uint64) {
	var h maphash.Hash
	h.SetSeed(d.Seed())
	hash.Node(&h, n)
	id := h.Sum64()
	return d.Load(id)
}

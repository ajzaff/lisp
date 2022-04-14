package innitdb

import (
	"hash/maphash"
	"strings"

	"github.com/ajzaff/innit"
)

type QueryInterface interface {
	Seed() maphash.Seed
	Load(id uint64) (fc int, refs, irefs []float64)
}

type QueryResult struct {
	Names []string
	Rows  [][]innit.Node
}

func Query(db QueryInterface, q string) (*QueryResult, error) {
	query, err := innit.Parse(q)
	if err != nil {
		return nil, err
	}
	var res QueryResult
	res.Names = queryNames(query)
	return &res, nil
}

func queryNames(q innit.Node) (names []string) {
	var v innit.Visitor
	v.SetLitVisitor(func(e *innit.Lit) {
		if e.Tok == innit.Id && strings.HasPrefix(e.Value, "?") {
			names = append(names, e.Value)
		}
	})
	v.Visit(q)
	return
}

package hash

import (
	"hash/maphash"
	"strings"
	"testing"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/scan"
)

func TestDistictHashes(t *testing.T) {
	for _, tc := range []struct {
		name         string
		input1       string
		input2       string
		wantDistinct bool
	}{{
		name:   "id is independent of src pos",
		input1: "a",
		input2: "   a",
	}, {
		name:   "int is independent of src pos",
		input1: "12",
		input2: "   12",
	}, {
		name:         "id and group are distinct",
		input1:       "a",
		input2:       "(a)",
		wantDistinct: true,
	}, {
		name:         "spaces inserted between id lits",
		input1:       "a b c",
		input2:       "abc",
		wantDistinct: true,
	}, {
		name:         "spaces inserted between id lits in group",
		input1:       "(a b c)",
		input2:       "(abc)",
		wantDistinct: true,
	}} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			v1 := mustParse(t, tc.input1)
			v2 := mustParse(t, tc.input2)
			seed := maphash.MakeSeed()
			var h1, h2 uint64
			{
				var h MapHash
				h.SetSeed(seed)
				h.WriteVal(v1)
				h1 = h.Sum64()
			}
			{
				var h MapHash
				h.SetSeed(seed)
				h.WriteVal(v2)
				h2 = h.Sum64()
			}

			if tc.wantDistinct && h1 == h2 {
				t.Errorf("Hash(%q) == Hash(%q) but wanted distinct hashes", tc.input1, tc.input2)
			}

			if !tc.wantDistinct && h1 != h2 {
				t.Errorf("Hash(%q) != Hash(%q) but wanted equivalent hashes", tc.input1, tc.input2)
			}
		})
	}
}

func mustParse(t *testing.T, input string) (val lisp.Val) {
	t.Helper()
	var sc scan.Scanner
	sc.Reset(strings.NewReader(input))
	for n := range sc.Nodes() {
		if val != nil {
			t.Fatalf("mustParse(%q): consumed more than one node", input)
		}
		val = n.Val
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("mustParse(%q): failed: %v", input, err)
	}
	return nil
}

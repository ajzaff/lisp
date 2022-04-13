package hash

import (
	"hash/maphash"
	"testing"

	"github.com/ajzaff/innit"
)

func TestDistictHash(t *testing.T) {
	for _, tc := range []struct {
		name         string
		input1       string
		input2       string
		wantDistinct bool
	}{{
		name:         "id and expr are distinct",
		input1:       "a",
		input2:       "(a)",
		wantDistinct: true,
	}, {
		name:         "spaces inserted between id lits",
		input1:       "a b c",
		input2:       "abc",
		wantDistinct: true,
	}, {
		name:         "spaces inserted between id lits in expr",
		input1:       "(a b c)",
		input2:       "(abc)",
		wantDistinct: true,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			n1, err := innit.Parse(tc.input1)
			if err != nil {
				t.Fatalf("Parse(%q): fails: %v", tc.input1, err)
			}
			n2, err := innit.Parse(tc.input2)
			if err != nil {
				t.Fatalf("Parse(%q): fails: %v", tc.input2, err)
			}

			var h maphash.Hash
			h1 := Hash(&h, n1)
			h.Reset()
			h2 := Hash(&h, n2)

			if tc.wantDistinct && h1 == h2 {
				t.Errorf("Hash(%q) == Hash(%q) but wanted distinct hashes", tc.input1, tc.input2)
			}

			if !tc.wantDistinct && h1 != h2 {
				t.Errorf("Hash(%q) != Hash(%q) but wanted equivalent hashes", tc.input1, tc.input2)
			}
		})
	}
}

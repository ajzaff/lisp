package hash

import (
	"strings"
	"testing"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/scan"
)

func TestDistictHash(t *testing.T) {
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
		name:         "id and cons are distinct",
		input1:       "a",
		input2:       "(a)",
		wantDistinct: true,
	}, {
		name:         "spaces inserted between id lits",
		input1:       "a b c",
		input2:       "abc",
		wantDistinct: true,
	}, {
		name:         "spaces inserted between id lits in cons",
		input1:       "(a b c)",
		input2:       "(abc)",
		wantDistinct: true,
	}} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var v1 lisp.Val
			var s scan.TokenScanner
			s.Reset(strings.NewReader(tc.input1))
			var sc scan.NodeScanner
			sc.Reset(&s)
			for sc.Scan() {
				_, _, v1 = sc.Node()
				break
			}
			if err := sc.Err(); err != nil {
				t.Fatalf("Parse(%q): fails: %v", tc.input1, err)
			}
			var v2 lisp.Val
			s.Reset(strings.NewReader(tc.input2))
			sc.Reset(&s)
			for sc.Scan() {
				_, _, v2 = sc.Node()
				break
			}
			if err := sc.Err(); err != nil {
				t.Fatalf("Parse(%q): fails: %v", tc.input2, err)
			}

			var h MapHash
			h.WriteVal(v1)
			h1 := h.Sum64()
			h.Reset()
			h.WriteVal(v2)
			h2 := h.Sum64()

			if tc.wantDistinct && h1 == h2 {
				t.Errorf("Hash(%q) == Hash(%q) but wanted distinct hashes", tc.input1, tc.input2)
			}

			if !tc.wantDistinct && h1 != h2 {
				t.Errorf("Hash(%q) != Hash(%q) but wanted equivalent hashes", tc.input1, tc.input2)
			}
		})
	}
}

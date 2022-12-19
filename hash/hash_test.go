package hash

import (
	"strings"
	"testing"

	"github.com/ajzaff/lisp"
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
		name:   "float is independent of src pos",
		input1: ".44",
		input2: "   .44",
	}, {
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var n1 lisp.Node
			var s lisp.TokenScanner
			s.Reset(strings.NewReader(tc.input1))
			var sc lisp.NodeScanner
			sc.Reset(&s)
			for sc.Scan() {
				n1 = sc.Node()
				break
			}
			if err := sc.Err(); err != nil {
				t.Fatalf("Parse(%q): fails: %v", tc.input1, err)
			}
			var n2 lisp.Node
			s.Reset(strings.NewReader(tc.input2))
			sc.Reset(&s)
			for sc.Scan() {
				n2 = sc.Node()
				break
			}
			if err := sc.Err(); err != nil {
				t.Fatalf("Parse(%q): fails: %v", tc.input2, err)
			}

			var h MapHash
			h.WriteVal(n1.Val)
			h1 := h.Sum64()
			h.Reset()
			h.WriteVal(n2.Val)
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

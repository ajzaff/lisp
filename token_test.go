package innit

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTokenizeBasic(t *testing.T) {
	for _, tc := range []struct {
		name    string
		src     string
		want    []Pos
		wantErr bool
	}{{
		name: "empty",
	}, {
		name: "whitespace",
		src:  "  \t\n",
	}, {
		name: "id",
		src:  "foo",
		want: []Pos{0, 3},
	}, {
		name: "id 3",
		src:  "  \t\n x",
		want: []Pos{5, 6},
	}, {
		name: "id 4",
		src:  "a b c",
		want: []Pos{0, 1, 2, 3, 4, 5},
	}, {
		name: "id op",
		src:  "a...",
		want: []Pos{0, 1, 1, 4},
	}, {
		name: "id string",
		src:  `a"abc"`,
		want: []Pos{0, 1, 1, 6},
	}, {
		name: "id op id",
		src:  "foo-bar",
		want: []Pos{0, 3, 3, 4, 4, 7},
	}, {
		name: "id op int",
		src:  `a.0`,
		want: []Pos{0, 1, 1, 2, 2, 3},
	}, {
		name: "op op",
		src:  `.+`,
		want: []Pos{0, 2},
	}, {
		name: "int",
		src:  "0",
		want: []Pos{0, 1},
	}, {
		name: "int 2",
		src:  "0 1 2",
		want: []Pos{0, 1, 2, 3, 4, 5},
	}, {
		name: "float",
		src:  "1.0",
		want: []Pos{0, 3},
	}, {
		name: "float 2",
		src:  "1.",
		want: []Pos{0, 2},
	}, {
		name: "float 3",
		src:  ".1",
		want: []Pos{0, 2},
	}, {
		name: "float 4",
		src:  "1. 2. 3.",
		want: []Pos{0, 2, 3, 5, 6, 8},
	}, {
		name: "string",
		src:  `"a"`,
		want: []Pos{0, 3},
	}, {
		name: "string 2",
		src:  `"a b c"`,
		want: []Pos{0, 7},
	}, {
		name: "string 3",
		src:  `"a" "b" "c"`,
		want: []Pos{0, 3, 4, 7, 8, 11},
	}, {
		name: "string (multiline)",
		src: `"
"`,
		want: []Pos{0, 3},
	}, {
		name: "expr",
		src:  "(abc)",
		want: []Pos{0, 1, 1, 4, 4, 5},
	}, {
		name: "expr 2",
		src:  "(add 1 2)",
		want: []Pos{0, 1, 1, 4, 5, 6, 7, 8, 8, 9},
	}, {
		name: "expr 3",
		src:  "(add (sub 3 2) 2)",
		want: []Pos{0, 1, 1, 4, 5, 6, 6, 9, 10, 11, 12, 13, 13, 14, 15, 16, 16, 17},
	}, {
		name: "expr 4",
		src:  "((a))",
		want: []Pos{0, 1, 1, 2, 2, 3, 3, 4, 4, 5},
	}, {
		name: "expr 5",
		src:  "(a)(b) (c)",
		want: []Pos{0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 7, 8, 8, 9, 9, 10},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.want)%2 != 0 {
				t.Fatalf("Tokenize() wants invalid result (cannot have odd length): %v", tc.want)
			}
			got, gotErr := Tokenize(tc.src)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Tokenize() got diff (-want, +got):\n%s", diff)
			}
			if (gotErr == nil) == tc.wantErr {
				t.Errorf("Tokenize() got err: %q, want err? %v", gotErr, tc.wantErr)
			}
		})
	}
}

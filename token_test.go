package lisp

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTokenizeLit(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   string
		want    []Pos
		wantErr bool
	}{{
		name:  "id",
		input: "foo",
		want:  []Pos{0, 3},
	}, {
		name:  "symbol",
		input: `.+`,
		want:  []Pos{0, 2},
	}, {
		name:  "id symbol id",
		input: "foo-bar",
		want:  []Pos{0, 3, 3, 4, 4, 7},
	}, {
		name:  "space id",
		input: "  \t\n x",
		want:  []Pos{5, 6},
	}, {
		name:  "id id id",
		input: "a b c",
		want:  []Pos{0, 1, 2, 3, 4, 5},
	}, {
		name:  "symbol 2",
		input: ".",
		want:  []Pos{0, 1},
	}, {
		name:  "id float",
		input: `a.0`,
		want:  []Pos{0, 1, 1, 3},
	}, {
		name:  "id symbol",
		input: "a...",
		want:  []Pos{0, 1, 1, 4},
	}, {
		name:  "symbol id",
		input: "...a",
		want:  []Pos{0, 3, 3, 4},
	}, {
		name:  "id string",
		input: `a"abc"`,
		want:  []Pos{0, 1, 1, 6},
	}, {
		name:  "id id id 2",
		input: "ab cd ef",
		want:  []Pos{0, 2, 3, 5, 6, 8},
	}, {
		name:  "int",
		input: "0",
		want:  []Pos{0, 1},
	}, {
		name:  "int 2",
		input: "0 1 2",
		want:  []Pos{0, 1, 2, 3, 4, 5},
	}, {
		name:  "float",
		input: "1.0",
		want:  []Pos{0, 3},
	}, {
		name:  "float 2",
		input: "1.",
		want:  []Pos{0, 2},
	}, {
		name:  "float 3",
		input: ".1",
		want:  []Pos{0, 2},
	}, {
		name:  "float 4",
		input: "1. 2. 3.",
		want:  []Pos{0, 2, 3, 5, 6, 8},
	}, {
		name:  "float 4_2",
		input: ".1 .2 .3",
		want:  []Pos{0, 2, 3, 5, 6, 8},
	}, {
		name:  "string",
		input: `"a"`,
		want:  []Pos{0, 3},
	}, {
		name:  "string 2",
		input: `"a b c"`,
		want:  []Pos{0, 7},
	}, {
		name:  "string 3",
		input: `"a" "b" "c"`,
		want:  []Pos{0, 3, 4, 7, 8, 11},
	}, {
		name: "string (multiline)",
		input: `"
"`,
		want: []Pos{0, 3},
	}, {
		name:  "string (double escape)",
		input: `"\\"`,
		want:  []Pos{0, 4},
	}, {
		name:  "byte lit",
		input: `"abc\x00\x11\xff"`,
		want:  []Pos{0, 17},
	}, {
		name:    "byte lit",
		input:   `"abc\x00\x1"`,
		want:    []Pos{0},
		wantErr: true,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.wantErr && len(tc.want)%2 != 0 {
				t.Fatalf("Tokenize(%q) wants invalid result (cannot have odd length when wantErr=true): %v", tc.name, tc.want)
			}
			got, gotErr := Tokenize(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Tokenize(%q) got diff (-want, +got):\n%s", tc.name, diff)
			}
			if (gotErr == nil) == tc.wantErr {
				t.Errorf("Tokenize(%q) got err: %q, want err? %v", tc.name, gotErr, tc.wantErr)
			}
		})
	}
}

func TestTokenizeExpr(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   string
		want    []Pos
		wantErr bool
	}{{
		name: "empty",
	}, {
		name:  "whitespace",
		input: "  \t\n",
	}, {
		name:  "expr",
		input: "(abc)",
		want:  []Pos{0, 1, 1, 4, 4, 5},
	}, {
		name:  "expr symbol",
		input: "(.)",
		want:  []Pos{0, 1, 1, 2, 2, 3},
	}, {
		name:  "expr 2",
		input: "(add 1 2)",
		want:  []Pos{0, 1, 1, 4, 5, 6, 7, 8, 8, 9},
	}, {
		name:  "expr 3",
		input: "(add (sub 3 2) 2)",
		want:  []Pos{0, 1, 1, 4, 5, 6, 6, 9, 10, 11, 12, 13, 13, 14, 15, 16, 16, 17},
	}, {
		name:  "expr 4",
		input: "((a))",
		want:  []Pos{0, 1, 1, 2, 2, 3, 3, 4, 4, 5},
	}, {
		name:  "expr 5",
		input: "(a)(b) (c)",
		want:  []Pos{0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 7, 8, 8, 9, 9, 10},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.wantErr && len(tc.want)%2 != 0 {
				t.Fatalf("Tokenize(%q) wants invalid result (cannot have odd length when wantErr=true): %v", tc.name, tc.want)
			}
			got, gotErr := Tokenize(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Tokenize(%q) got diff (-want, +got):\n%s", tc.input, diff)
			}
			if (gotErr == nil) == tc.wantErr {
				t.Errorf("Tokenize(%q) got err: %q, want err? %v", tc.input, gotErr, tc.wantErr)
			}
		})
	}
}

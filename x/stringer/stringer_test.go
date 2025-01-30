package stringer

import (
	"testing"

	"github.com/ajzaff/lisp"
	"github.com/google/go-cmp/cmp"
)

func TestVal(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input lisp.Val
		want  string
	}{{
		name: "Go nil",
		want: "<nil>",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			got := Val(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Val(%q): got diff:\n%s", tc.name, diff)
			}
		})
	}
}

func TestLit(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input lisp.Lit
		want  string
	}{{
		name:  "Invalid Lit uses GoString",
		input: lisp.Lit{Text: "a"},
		want:  `lisp.Lit{Token:0, Text:"a"}`,
	}, {
		name:  "unknown Invalid Lit uses GoString",
		input: lisp.Lit{Token: 999, Text: "a"},
		want:  `lisp.Lit{Token:999, Text:"a"}`,
	}, {
		name:  "empty Id uses GoString",
		input: lisp.Lit{Token: lisp.Id},
		want:  `lisp.Lit{Token:1, Text:""}`,
	}, {
		name:  "Nat as invalid Id uses GoString",
		input: lisp.Lit{Token: lisp.Id, Text: "0"},
		want:  `lisp.Lit{Token:1, Text:"0"}`,
	}, {
		name:  "Invalid Id uses GoString",
		input: lisp.Lit{Token: lisp.Id, Text: "\x00"},
		want:  `lisp.Lit{Token:1, Text:"\x00"}`,
	}, {
		name:  "Invalid unicode Id uses GoString",
		input: lisp.Lit{Token: lisp.Id, Text: "⍟"},
		want:  `lisp.Lit{Token:1, Text:"⍟"}`,
	}, {
		name:  "Invalid Id uses GoString",
		input: lisp.Lit{Token: lisp.Id, Text: "abc123"},
		want:  `lisp.Lit{Token:1, Text:"abc123"}`,
	}, {
		name:  "Id a",
		input: lisp.Lit{Token: lisp.Id, Text: "a"},
		want:  "a",
	}, {
		name:  "Id abc",
		input: lisp.Lit{Token: lisp.Id, Text: "abc"},
		want:  "abc",
	}, {
		name:  "empty Nat uses GoString",
		input: lisp.Lit{Token: lisp.Nat},
		want:  `lisp.Lit{Token:2, Text:""}`,
	}, {
		name:  "Id as invalid Nat uses GoString",
		input: lisp.Lit{Token: lisp.Nat, Text: "a"},
		want:  `lisp.Lit{Token:2, Text:"a"}`,
	}, {
		name:  "Invalid Nat uses GoString",
		input: lisp.Lit{Token: lisp.Nat, Text: "01"},
		want:  `lisp.Lit{Token:2, Text:"01"}`,
	}, {
		name:  "Invalid Nat uses GoString",
		input: lisp.Lit{Token: lisp.Nat, Text: "1a"},
		want:  `lisp.Lit{Token:2, Text:"1a"}`,
	}, {
		name:  "Nat 0",
		input: lisp.Lit{Token: lisp.Nat, Text: "0"},
		want:  "0",
	}, {
		name:  "Nat 42",
		input: lisp.Lit{Token: lisp.Nat, Text: "42"},
		want:  "42",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			got := Lit(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Lit(%q): got diff:\n%s", tc.name, diff)
			}
		})
	}
}

func TestGroup(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input lisp.Group
		want  string
	}{{
		name: "Go nil",
		want: "()",
	}, {
		name:  "nil Group",
		input: (lisp.Group)(nil),
		want:  "()",
	}, {
		name:  "empty Group",
		input: lisp.Group{},
		want:  "()",
	}, {
		name:  "group with invalid Lit uses GoString",
		input: lisp.Group{lisp.Lit{}},
		want:  `lisp.Group{lisp.Lit{Token:0, Text:""}}`,
	}, {
		name:  "nested group",
		input: lisp.Group{(lisp.Group)(nil), lisp.Group{}},
		want:  "(()())",
	}, {
		name: "valid Group with some Lits",
		input: lisp.Group{
			lisp.Lit{Token: lisp.Id, Text: "a"},
			lisp.Lit{Token: lisp.Id, Text: "b"},
			lisp.Lit{Token: lisp.Id, Text: "c"},
		},
		want: "(a b c)",
	}, {
		name: "valid nested Group with some mixed Lits",
		input: lisp.Group{
			lisp.Lit{Token: lisp.Id, Text: "a"},
			lisp.Lit{Token: lisp.Id, Text: "b"},
			lisp.Group{
				lisp.Lit{Token: lisp.Nat, Text: "1"},
				lisp.Lit{Token: lisp.Nat, Text: "2"},
				lisp.Lit{Token: lisp.Nat, Text: "3"},
			},
			lisp.Lit{Token: lisp.Id, Text: "c"},
		},
		want: "(a b(1 2 3)c)",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			got := Group(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Group(%q): got diff:\n%s", tc.name, diff)
			}
		})
	}
}

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

func TestCons(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input *lisp.Cons
		want  string
	}{{
		name: "Go nil",
		want: "()",
	}, {
		name:  "nil Cons",
		input: (*lisp.Cons)(nil),
		want:  "()",
	}, {
		name:  "empty Cons",
		input: &lisp.Cons{},
		want:  "()",
	}, {
		name:  "invalid Cons struct uses GoString",
		input: &lisp.Cons{Cons: &lisp.Cons{}},
		want:  "&lisp.Cons{Val:(lisp.Val)(nil), Cons:&lisp.Cons{Val:(lisp.Val)(nil), Cons:(*lisp.Cons)(nil)}}",
	}, {
		name:  "nested invalid Cons struct uses GoString",
		input: &lisp.Cons{Cons: &lisp.Cons{Val: &lisp.Cons{}}},
		want:  "&lisp.Cons{Val:(lisp.Val)(nil), Cons:&lisp.Cons{Val:&lisp.Cons{Val:(lisp.Val)(nil), Cons:(*lisp.Cons)(nil)}, Cons:(*lisp.Cons)(nil)}}",
	}, {
		name:  "empty value in linked cons uses GoString",
		input: &lisp.Cons{Val: lisp.Lit{Token: lisp.Id, Text: "a"}, Cons: &lisp.Cons{}},
		want:  `&lisp.Cons{Val:lisp.Lit{Token:1, Text:"a"}, Cons:&lisp.Cons{Val:(lisp.Val)(nil), Cons:(*lisp.Cons)(nil)}}`,
	}, {
		name:  "cons with invalid Lit uses GoString",
		input: &lisp.Cons{Val: lisp.Lit{}},
		want:  `&lisp.Cons{Val:lisp.Lit{Token:0, Text:""}, Cons:(*lisp.Cons)(nil)}`,
	}, {
		name:  "nested cons",
		input: &lisp.Cons{Val: (*lisp.Cons)(nil), Cons: &lisp.Cons{Val: &lisp.Cons{}}},
		want:  "(()())",
	}, {
		name: "valid Cons with some Lits",
		input: &lisp.Cons{
			Val: lisp.Lit{Token: lisp.Id, Text: "a"},
			Cons: &lisp.Cons{
				Val: lisp.Lit{Token: lisp.Id, Text: "b"},
				Cons: &lisp.Cons{
					Val: lisp.Lit{Token: lisp.Id, Text: "c"},
				},
			},
		},
		want: "(a b c)",
	}, {
		name: "valid nested Cons with some mixed Lits",
		input: &lisp.Cons{
			Val: lisp.Lit{Token: lisp.Id, Text: "a"},
			Cons: &lisp.Cons{
				Val: lisp.Lit{Token: lisp.Id, Text: "b"},
				Cons: &lisp.Cons{
					Val: &lisp.Cons{
						Val: lisp.Lit{Token: lisp.Nat, Text: "1"},
						Cons: &lisp.Cons{
							Val: lisp.Lit{Token: lisp.Nat, Text: "2"},
							Cons: &lisp.Cons{
								Val: lisp.Lit{Token: lisp.Nat, Text: "3"},
							},
						},
					},
					Cons: &lisp.Cons{
						Val: lisp.Lit{Token: lisp.Id, Text: "c"},
					},
				},
			},
		},
		want: "(a b(1 2 3)c)",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			got := Cons(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Cons(%q): got diff:\n%s", tc.name, diff)
			}
		})
	}
}

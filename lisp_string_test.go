package lisp

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestValString(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input Val
		want  string
	}{{
		name: "Go nil",
		want: "<nil>",
	}, {
		name:  "Invalid Lit uses GoString",
		input: Lit{Text: "a"},
		want:  `lisp.Lit{Token:0, Text:"a"}`,
	}, {
		name:  "unknown Invalid Lit uses GoString",
		input: Lit{Token: 999, Text: "a"},
		want:  `lisp.Lit{Token:999, Text:"a"}`,
	}, {
		name:  "empty Id uses GoString",
		input: Lit{Token: Id},
		want:  `lisp.Lit{Token:1, Text:""}`,
	}, {
		name:  "empty Nat uses GoString",
		input: Lit{Token: Nat},
		want:  `lisp.Lit{Token:2, Text:""}`,
	}, {
		name:  "nil Cons",
		input: (*Cons)(nil),
		want:  "()",
	}, {
		name:  "empty Cons",
		input: &Cons{},
		want:  "()",
	}, {
		name:  "invalid Cons struct uses GoString",
		input: &Cons{Cons: &Cons{}},
		want:  "&lisp.Cons{Val:(lisp.Val)(nil), Cons:&lisp.Cons{Val:(lisp.Val)(nil), Cons:(*lisp.Cons)(nil)}}",
	}, {
		name:  "nested invalid Cons struct uses GoString",
		input: &Cons{Cons: &Cons{Val: &Cons{}}},
		want:  "&lisp.Cons{Val:(lisp.Val)(nil), Cons:&lisp.Cons{Val:&lisp.Cons{Val:(lisp.Val)(nil), Cons:(*lisp.Cons)(nil)}, Cons:(*lisp.Cons)(nil)}}",
	}, {
		name:  "cons with invalid Lit uses GoString",
		input: &Cons{Val: Lit{}},
		want:  `&lisp.Cons{Val:lisp.Lit{Token:0, Text:""}, Cons:(*lisp.Cons)(nil)}`,
	}, {
		name: "valid Cons with some Lits",
		input: &Cons{
			Val: Lit{Token: Id, Text: "a"},
			Cons: &Cons{
				Val: Lit{Token: Id, Text: "b"},
				Cons: &Cons{
					Val: Lit{Token: Id, Text: "c"},
				},
			},
		},
		want: "(a b c)",
	}, {
		name: "valid nested Cons with some mixed Lits",
		input: &Cons{
			Val: Lit{Token: Id, Text: "a"},
			Cons: &Cons{
				Val: Lit{Token: Id, Text: "b"},
				Cons: &Cons{
					Val: &Cons{
						Val: Lit{Token: Nat, Text: "1"},
						Cons: &Cons{
							Val: Lit{Token: Nat, Text: "2"},
							Cons: &Cons{
								Val: Lit{Token: Nat, Text: "3"},
							},
						},
					},
					Cons: &Cons{
						Val: Lit{Token: Id, Text: "c"},
					},
				},
			},
		},
		want: "(a b(1 2 3)c)",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			got := fmt.Sprint(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("String(%q): got diff:\n%s", tc.name, diff)
			}
		})
	}
}

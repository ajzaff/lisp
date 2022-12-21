package print

import (
	"strings"
	"testing"

	"github.com/ajzaff/lisp"
	"github.com/google/go-cmp/cmp"
)

func TestStdPrint(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input lisp.Val
		want  string
	}{{
		name: "empty",
		want: "()\n",
	}, {
		name:  "empty cons",
		input: &lisp.Cons{},
		want:  "()\n",
	}, {
		name:  "lit",
		input: lisp.Lit{Token: lisp.Id, Text: "hello"},
		want:  "hello\n",
	}, {
		name:  "cons",
		input: &lisp.Cons{Val: lisp.Lit{Token: lisp.Id, Text: "x"}},
		want:  "(x)\n",
	}, {
		name: "cons3",
		input: &lisp.Cons{
			Val: lisp.Lit{Token: lisp.Id, Text: "x"},
			Cons: &lisp.Cons{
				Val: lisp.Lit{Token: lisp.Id, Text: "y"},
				Cons: &lisp.Cons{
					Val: lisp.Lit{Token: lisp.Id, Text: "z"},
				},
			},
		},
		want: "(x y z)\n",
	}, {
		name: "nested cons",
		input: &lisp.Cons{
			Val: lisp.Lit{Token: lisp.Id, Text: "x"},
			Cons: &lisp.Cons{
				Val: &lisp.Cons{Val: lisp.Lit{Token: lisp.Id, Text: "y"}},
				Cons: &lisp.Cons{
					Val: lisp.Lit{Token: lisp.Id, Text: "z"},
				},
			},
		},
		want: "(x(y)z)\n",
	}, {
		name: "numbers andlisp.Ids are delimitable",
		input: &lisp.Cons{
			Val: lisp.Lit{Token: lisp.Id, Text: "add"},
			Cons: &lisp.Cons{
				Val: lisp.Lit{Token: lisp.Nat, Text: "1"},
				Cons: &lisp.Cons{
					Val: lisp.Lit{Token: lisp.Nat, Text: "2"},
				},
			},
		},
		want: "(add 1 2)\n",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			var sb strings.Builder
			StdPrinter(&sb).Print(tc.input)
			got := sb.String()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Print(%q): got diff (-want, +got):\n%v", tc.name, diff)
			}
		})
	}
}

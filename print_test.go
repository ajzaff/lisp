package lisp

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type unknownVal struct{}

func (unknownVal) val() {}

func TestStdPrint(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input Val
		want  string
	}{{
		name: "empty",
		want: "(nil)\n",
	}, {
		name:  "unknown node",
		input: unknownVal{},
	}, {
		name:  "empty cons",
		input: &Cons{},
		want:  "()\n",
	}, {
		name:  "lit",
		input: Lit{Token: Id, Text: "hello"},
		want:  "hello\n",
	}, {
		name:  "cons",
		input: &Cons{Node: Node{Val: Lit{Token: Id, Text: "x"}}},
		want:  "(x)\n",
	}, {
		name: "cons3",
		input: &Cons{
			Node: Node{Val: Lit{Token: Id, Text: "x"}},
			Cons: &Cons{
				Node: Node{Val: Lit{Token: Id, Text: "y"}},
				Cons: &Cons{
					Node: Node{Val: Lit{Token: Id, Text: "z"}},
				},
			},
		},
		want: "(x y z)\n",
	}, {
		name: "nested cons",
		input: &Cons{
			Node: Node{Val: Lit{Token: Id, Text: "x"}},
			Cons: &Cons{
				Node: Node{Val: &Cons{Node: Node{Val: Lit{Token: Id, Text: "y"}}}},
				Cons: &Cons{
					Node: Node{Val: Lit{Token: Id, Text: "z"}},
				},
			},
		},
		want: "(x(y)z)\n",
	}, {
		name: "numbers and ids are delimitable",
		input: &Cons{
			Node: Node{Val: Lit{Token: Id, Text: "add"}},
			Cons: &Cons{
				Node: Node{Val: Lit{Token: Int, Text: "1"}},
				Cons: &Cons{
					Node: Node{Val: Lit{Token: Int, Text: "2"}},
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

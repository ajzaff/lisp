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
		want: "<nil>\n",
	}, {
		name:  "unknown node",
		input: unknownVal{},
	}, {
		name:  "empty expr",
		input: Expr{},
		want:  "()\n",
	}, {
		name:  "lit",
		input: Lit{Token: Id, Text: "hello"},
		want:  "hello\n",
	}, {
		name:  "expr",
		input: Expr{&LitNode{Lit: Lit{Token: Id, Text: "x"}}},
		want:  "(x)\n",
	}, {
		name: "expr3",
		input: Expr{
			&LitNode{Lit: Lit{Token: Id, Text: "x"}},
			&LitNode{Lit: Lit{Token: Id, Text: "y"}},
			&LitNode{Lit: Lit{Token: Id, Text: "z"}},
		},
		want: "(x y z)\n",
	}, {
		name: "nested expr",
		input: Expr{
			&LitNode{Lit: Lit{Token: Id, Text: "x"}},
			&ExprNode{Expr: Expr{&LitNode{Lit: Lit{Token: Id, Text: "y"}}}},
			&LitNode{Lit: Lit{Token: Id, Text: "z"}},
		},
		want: "(x(y)z)\n",
	}, {
		name: "squashed ids and symbols in expr",
		input: Expr{
			&LitNode{Lit: Lit{Token: Id, Text: "?"}},
			&LitNode{Lit: Lit{Token: Id, Text: "x"}},
			&LitNode{Lit: Lit{Token: Id, Text: "/"}},
			&LitNode{Lit: Lit{Token: Id, Text: "y"}},
		},
		want: "(?x/y)\n",
	}, {
		name: "numbers and ids are delimitable",
		input: Expr{
			&LitNode{Lit: Lit{Token: Id, Text: "add"}},
			&LitNode{Lit: Lit{Token: Number, Text: "1"}},
			&LitNode{Lit: Lit{Token: Number, Text: "2"}},
		},
		want: "(add 1 2)\n",
	}, {
		name: "numbers and symbols are not delimitable",
		input: Expr{
			&LitNode{Lit: Lit{Token: Number, Text: "1"}},
			&LitNode{Lit: Lit{Token: Id, Text: "+"}},
			&LitNode{Lit: Lit{Token: Number, Text: "2"}},
		},
		want: "(1+2)\n",
	}, {
		name: "ids and symbols are not delimitable",
		input: Expr{
			&LitNode{Lit: Lit{Token: Id, Text: "a"}},
			&LitNode{Lit: Lit{Token: Id, Text: "+"}},
			&LitNode{Lit: Lit{Token: Id, Text: "b"}},
		},
		want: "(a+b)\n",
	}, {
		name: "repeated distinct symbols are delimitable",
		input: Expr{
			&LitNode{Lit: Lit{Token: Id, Text: "+"}},
			&LitNode{Lit: Lit{Token: Id, Text: "-"}},
			&LitNode{Lit: Lit{Token: Id, Text: "/"}},
		},
		want: "(+ - /)\n",
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

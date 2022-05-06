package innit

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
		name:  "lit",
		input: IdLit("hello"),
		want:  "hello\n",
	}, {
		name:  "expr",
		input: Expr{&LitNode{Lit: IdLit("x")}},
		want:  "(x)\n",
	}, {
		name: "expr3",
		input: Expr{
			&LitNode{Lit: IdLit("x")},
			&LitNode{Lit: IdLit("y")},
			&LitNode{Lit: IdLit("z")},
		},
		want: "(x y z)\n",
	}, {
		name: "nested expr",
		input: Expr{
			&LitNode{Lit: IdLit("x")},
			&ExprNode{Expr: Expr{&LitNode{Lit: IdLit("y")}}},
			&LitNode{Lit: IdLit("z")},
		},
		want: "(x(y)z)\n",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			var sb strings.Builder
			StdPrinter(&sb).Print(tc.input)
			got := sb.String()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Print(%q): got diff:\n%v", tc.name, diff)
			}
		})
	}
}

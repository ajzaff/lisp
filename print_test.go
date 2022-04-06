package innit

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStdPrint(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input Node
		want  string
	}{{
		name: "empty",
	}, {
		name:  "lit",
		input: &Lit{Tok: Id, Value: "hello"},
		want:  "hello\n",
	}, {
		name:  "expr",
		input: &Expr{X: NodeList{&Lit{Tok: Id, Value: "x"}}},
		want:  "(x)\n",
	}, {
		name: "expr",
		input: &Expr{X: NodeList{
			&Lit{Tok: Id, Value: "x"},
			&Lit{Tok: Id, Value: "y"},
			&Lit{Tok: Id, Value: "z"},
		}},
		want: "(x y z)\n",
	}, {
		name: "nested expr",
		input: &Expr{X: NodeList{
			&Lit{Tok: Id, Value: "x"},
			&Expr{X: NodeList{&Lit{Tok: Id, Value: "y"}}},
			&Lit{Tok: Id, Value: "z"},
		}},
		want: "(x (y) z)\n",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var sb strings.Builder
			StdPrinter(&sb).Print(tc.input)
			got := sb.String()
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("Print(%q) = got diff:\n%v", tc.name, diff)
			}
		})
	}
}

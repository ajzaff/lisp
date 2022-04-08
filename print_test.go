package innit

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type unknownNode struct{}

func (unknownNode) Pos() Pos { return NoPos }
func (unknownNode) End() Pos { return NoPos }

func TestStdPrint(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input Node
		want  string
	}{{
		name: "empty",
	}, {
		name:  "unknown node",
		input: unknownNode{},
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
	}, {
		name: "node list",
		input: NodeList{
			&Expr{X: NodeList{&Lit{Tok: Id, Value: "x"}}},
			&Expr{X: NodeList{&Lit{Tok: Id, Value: "y"}}},
			&Expr{X: NodeList{&Lit{Tok: Id, Value: "z"}}},
		},
		want: "(x)\n(y)\n(z)\n",
	}, {
		name: "weird nested list",
		input: NodeList{
			NodeList{
				&Lit{Tok: String, Value: `"str1"`},
				&Lit{Tok: String, Value: `"str2"`},
			},
			&Lit{Tok: String, Value: `"str3"`},
		},
		want: "\"str1\"\n\"str2\"\n\"str3\"\n",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var sb strings.Builder
			StdPrinter(&sb).Print(tc.input)
			got := sb.String()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Print(%q): got diff:\n%v", tc.name, diff)
			}
		})
	}
}

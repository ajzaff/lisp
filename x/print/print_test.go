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
		name:  "nil group",
		input: (lisp.Group)(nil),
		want:  "()\n",
	}, {
		name:  "empty group",
		input: lisp.Group{},
		want:  "()\n",
	}, {
		name:  "lit",
		input: lisp.Lit("hello"),
		want:  "hello\n",
	}, {
		name:  "group",
		input: lisp.Group{lisp.Lit("x")},
		want:  "(x)\n",
	}, {
		name: "cons3",
		input: lisp.Group{
			lisp.Lit("x"),
			lisp.Lit("y"),
			lisp.Lit("z"),
		},
		want: "(x y z)\n",
	}, {
		name: "nested group",
		input: lisp.Group{
			lisp.Lit("x"),
			lisp.Group{lisp.Lit("y")},
			lisp.Lit("z"),
		},
		want: "(x(y)z)\n",
	}, {
		name: "numbers andlisp.Ids are delimitable",
		input: lisp.Group{
			lisp.Lit("add"),
			lisp.Lit("1"),
			lisp.Lit("2"),
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

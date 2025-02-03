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
		name:  "Empty Lit uses GoString",
		input: lisp.Lit(""),
		want:  `lisp.Lit("")`,
	}, {
		name:  "Invalid Id uses GoString",
		input: lisp.Lit("\x00"),
		want:  `lisp.Lit("\x00")`,
	}, {
		name:  "Invalid unicode Id uses GoString",
		input: lisp.Lit("⍟"),
		want:  `lisp.Lit("⍟")`,
	}, {
		name:  "Id a",
		input: lisp.Lit("a"),
		want:  "a",
	}, {
		name:  "Id abc",
		input: lisp.Lit("abc"),
		want:  "abc",
	}, {
		name:  "Nat 0",
		input: lisp.Lit("0"),
		want:  "0",
	}, {
		name:  "Nat 42",
		input: lisp.Lit("42"),
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
		input: lisp.Group{lisp.Lit("")},
		want:  `lisp.Group{lisp.Lit("")}`,
	}, {
		name:  "nested group",
		input: lisp.Group{(lisp.Group)(nil), lisp.Group{}},
		want:  "(()())",
	}, {
		name: "valid Group with some Lits",
		input: lisp.Group{
			lisp.Lit("a"),
			lisp.Lit("b"),
			lisp.Lit("c"),
		},
		want: "(a b c)",
	}, {
		name: "valid nested Group with some mixed Lits",
		input: lisp.Group{
			lisp.Lit("a"),
			lisp.Lit("b"),
			lisp.Group{
				lisp.Lit("1"),
				lisp.Lit("2"),
				lisp.Lit("3"),
			},
			lisp.Lit("c"),
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

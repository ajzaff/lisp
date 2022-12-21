package builder

import (
	"strings"
	"testing"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/scan"
	"github.com/google/go-cmp/cmp"
)

type builderTestCase struct {
	name      string
	input     func() *Builder
	want      string
	wantPanic bool
}

var builderTestCases = []builderTestCase{{
	name:  "nil Builder returns empty Cons",
	input: func() *Builder { return nil },
	want:  "()",
}, {
	name:  "new empty Builder returns empty Cons",
	input: func() *Builder { return new(Builder) },
	want:  "()",
}, {
	name: "BeginFrame on nil Builder causes panic",
	input: func() *Builder {
		var b *Builder
		b.BeginFrame()
		return b
	},
	wantPanic: true,
}, {
	name: "BeginFrame on new empty Builder returns empty Cons",
	input: func() *Builder {
		var b Builder
		b.BeginFrame()
		return &b
	},
	want: "()",
}, {
	name: "new empty Builder with Id",
	input: func() *Builder {
		var b Builder
		b.AppendId("a")
		return &b
	},
	want: "(a)",
}, {
	name: "BeginFrame on new empty Builder with Id",
	input: func() *Builder {
		var b Builder
		b.BeginFrame()
		b.AppendId("a")
		return &b
	},
	want: "(a)",
}}

func TestBuilder(t *testing.T) {
	for _, tc := range builderTestCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					if !tc.wantPanic {
						t.Errorf("TestBuilder(%q): got an unexpected panic", tc.name)
					}
				}
			}()
			b := tc.input()
			got := b.Build()
			want := mustParse(t, tc.want)
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("Build(%q): got diff:\n%s", tc.name, diff)
			}
			gotSecond := b.Build()
			if diff := cmp.Diff(got, gotSecond); diff != "" {
				t.Errorf("Build(%q) was not idempotent, got diff:\n%v", tc.name, diff)
			}
		})
	}
}

func mustParse(t *testing.T, input string) lisp.Val {
	var s scan.NodeScanner
	var sc scan.TokenScanner
	sc.Reset(strings.NewReader(input))
	s.Reset(&sc)
	for s.Scan() {
		_, _, v := s.Node()
		return v
	}
	t.Fatalf("mustParse(%q): failed: %v", input, s.Err())
	return nil
}

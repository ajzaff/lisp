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
	want      func(t *testing.T) lisp.Val
	wantPanic bool
}

var builderTestCases = []builderTestCase{{
	name:  "nil Builder returns empty Group",
	input: func() *Builder { return nil },
	want:  mustParseFunc("()"),
}, {
	name:  "new empty Builder returns empty Group",
	input: func() *Builder { return new(Builder) },
	want:  mustParseFunc("()"),
}, {
	name: "BeginFrame on nil Builder causes panic",
	input: func() *Builder {
		var b *Builder
		b.BeginFrame()
		return b
	},
	wantPanic: true,
}, {
	name: "EndFrame on nil Builder causes panic",
	input: func() *Builder {
		var b *Builder
		b.EndFrame()
		return b
	},
	wantPanic: true,
}, {
	name: "Append* on nil Builder causes panic",
	input: func() *Builder {
		var b *Builder
		b.AppendId("a")
		return b
	},
	wantPanic: true,
}, {
	name: "BeginFrame on new empty Builder returns empty Group",
	input: func() *Builder {
		var b Builder
		b.BeginFrame()
		return &b
	},
	want: mustParseFunc("()"),
}, {
	name: "new empty Builder with Id",
	input: func() *Builder {
		var b Builder
		b.AppendId("a")
		return &b
	},
	want: mustParseFunc("(a)"),
}, {
	name: "BeginFrame on new empty Builder with Id",
	input: func() *Builder {
		var b Builder
		b.BeginFrame()
		b.AppendId("a")
		return &b
	},
	want: mustParseFunc("(a)"),
}, {
	name: "extra EndFrame has no effect",
	input: func() *Builder {
		var b Builder
		b.BeginFrame()
		b.AppendId("a")
		b.EndFrame()
		b.EndFrame()
		return &b
	},
	want: mustParseFunc("(a)"),
}, {
	name: "Nested Group",
	input: func() *Builder {
		var b Builder
		b.BeginFrame()
		b.BeginFrame()
		return &b
	},
	want: mustParseFunc("(())"),
}, {
	name: "multiple Nested Group",
	input: func() *Builder {
		var b Builder
		b.BeginFrame()
		for i := 0; i < 3; i++ {
			b.BeginFrame()
			b.EndFrame()
		}
		return &b
	},
	want: mustParseFunc("(()()())"),
}, {
	name: "multiple mixed nested types",
	input: func() *Builder {
		var b Builder
		b.AppendId("a")
		b.BeginFrame()
		b.AppendId("b")
		b.AppendId("c")
		b.AppendNat(1)
		b.EndFrame()
		b.BeginFrame()
		b.AppendNat(2)
		b.AppendNat(3)
		b.EndFrame()
		b.AppendId("g")
		b.EndFrame()
		return &b
	},
	want: mustParseFunc("(a(b c 1)(2 3)g)"),
}, {
	name: "raw text Lit",
	input: func() *Builder {
		var b Builder
		b.AppendText("a")
		return &b
	},
	want: func(t *testing.T) lisp.Val { return lisp.Group{lisp.Lit("a")} },
}}

func TestBuilder(t *testing.T) {
	for _, tc := range builderTestCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					if !tc.wantPanic {
						t.Errorf("TestBuilder(%q): got an unexpected panic:", tc.name)
						panic(err)
					}
				}
			}()
			b := tc.input()
			got := b.Build()
			want := tc.want(t)
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("Build(%q): got diff:\n%s", tc.name, diff)
			}
			gotSecond := b.Build()
			if diff := cmp.Diff(got, gotSecond); diff != "" {
				t.Errorf("Build(%q) was not idempotent, got diff:\n%v", tc.name, diff)
			}
			if tc.wantPanic {
				t.Errorf("TestBuilder(%q): wanted a panic but did not get one", tc.name)
			}
		})
	}
}

func mustParseFunc(input string) func(t *testing.T) lisp.Val {
	return func(t *testing.T) lisp.Val { return mustParse(t, input) }
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
	if err := s.Err(); err != nil {
		t.Fatalf("mustParse(%q): failed: %v", input, err)
	}
	return nil
}

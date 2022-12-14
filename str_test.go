package lisp

import (
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestStr(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    string
		wantText string
		wantErr  error
	}{{
		name:    "empty",
		input:   "",
		wantErr: io.ErrUnexpectedEOF,
	}, {
		name:    "not str",
		input:   "a",
		wantErr: errStr,
	}, {
		name:    "UnexpectedEOF",
		input:   `"`,
		wantErr: io.ErrUnexpectedEOF,
	}, {
		name:     "replacement rune",
		input:    "\"\uFFFD\"",
		wantText: "\"\uFFFD\"",
	}, {
		name:     "whitespace",
		input:    "\" \n\r\t\"",
		wantText: "\" \n\r\t\"",
	}, {
		name:     "empty str",
		input:    `""`,
		wantText: `""`,
	}, {
		name:     "simple str",
		input:    `"a"`,
		wantText: `"a"`,
	}, {
		name:     "escape",
		input:    `"\""`,
		wantText: `"\""`,
	}, {
		name:     "multiple str",
		input:    `"a" "b" "c"`,
		wantText: `"a"`,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			got, gotErr := DecodeStr([]byte(tc.input))
			want := LitNode{}
			if gotErr == nil {
				want = LitNode{
					Lit: Lit{
						Token: String,
						Text:  tc.wantText,
					},
					EndPos: Pos(len(tc.wantText)),
				}
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("DecodeStr(%q): got diff:\n%s", tc.name, diff)
			}
			if errDiff := cmp.Diff(tc.wantErr, gotErr, cmpopts.EquateErrors()); errDiff != "" {
				t.Errorf("DecodeStr(%q): got err diff:\n%s", tc.name, errDiff)
			}
		})
	}
}

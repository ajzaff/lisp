package blisp

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/scan"
	"github.com/google/go-cmp/cmp"
)

func mustParse(t *testing.T, src string) lisp.Val {
	t.Helper()
	var v lisp.Val
	var s scan.Scanner
	s.Reset(strings.NewReader(src))
	var sc scan.NodeScanner
	sc.Reset(&s)
	for sc.Scan() {
		_, _, v = sc.Node()
		break
	}
	if err := sc.Err(); err != nil {
		panic(fmt.Sprintf("mustParse: failed to parse: %v", src))
	}
	return v
}

func TestEncode(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		want  []byte
	}{{
		name: "empty",
	}, {
		name:  "Id",
		input: "a",
		want: []byte{
			'a',
		},
	}, {
		name:  "Nat",
		input: "1",
		want: []byte{
			'1',
		},
	}, {
		name:  "empty Group",
		input: "()",
		want: []byte{
			byte(lisp.LParen),
			byte(lisp.RParen),
		},
	}, {
		name:  "Ids use delimiters",
		input: "(a b c)",
		want: []byte{
			byte(lisp.LParen),
			'a',
			' ',
			'b',
			' ',
			'c',
			byte(lisp.RParen),
		},
	}, {
		name:  "group id",
		input: "(abc)",
		want: []byte{
			byte(lisp.LParen),
			'a',
			'b',
			'c',
			byte(lisp.RParen),
		},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			v := mustParse(t, tc.input)
			if gotLen, wantLen := Len(v), len(tc.want); gotLen != wantLen {
				t.Errorf("EncodedLen(%q): got %d, want %d", tc.name, gotLen, wantLen)
			}
			var buf bytes.Buffer
			var e Encoder
			e.Reset(&buf)
			e.Encode(v)
			got := buf.Bytes()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Encode(%q): got diff (-want, +got):\n%v", tc.name, diff)
			}
		})
	}
}

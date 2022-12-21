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
	var n lisp.Node
	var s scan.TokenScanner
	s.Reset(strings.NewReader(src))
	var sc scan.NodeScanner
	sc.Reset(&s)
	for sc.Scan() {
		n = sc.Node()
		break
	}
	if err := sc.Err(); err != nil {
		panic(fmt.Sprintf("mustParse: failed to parse: %v", src))
	}
	return n.Val
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
			byte(lisp.Nat),
			1,
		},
	}, {
		name:  "empty Cons",
		input: "()",
		want: []byte{
			byte(lisp.LParen),
			byte(lisp.RParen),
		},
	}, {
		name:  "Nats are self-delimiting",
		input: "(1 2 3)",
		want: []byte{
			byte(lisp.LParen),
			byte(lisp.Nat),
			1,
			byte(lisp.Nat),
			2,
			byte(lisp.Nat),
			3,
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
		name:  "mixed Nat and Id minimizes delimiters",
		input: "(1 a 2 b 3 c)",
		want: []byte{
			byte(lisp.LParen),
			byte(lisp.Nat),
			1,
			'a',
			byte(lisp.Nat),
			2,
			'b',
			byte(lisp.Nat),
			3,
			'c',
			byte(lisp.RParen),
		},
	}, {
		name:  "cons id",
		input: "(abc)",
		want: []byte{
			byte(lisp.LParen),
			'a',
			'b',
			'c',
			byte(lisp.RParen),
		},
	}, {
		name:  "nested nats",
		input: "(1 (2 (3 4000)) abc)",
		want: []byte{
			byte(lisp.LParen),
			byte(lisp.Nat),
			1,
			byte(lisp.LParen),
			byte(lisp.Nat),
			2,
			byte(lisp.LParen),
			byte(lisp.Nat),
			3,
			byte(lisp.Nat),
			0xa0,
			0x1f,
			byte(lisp.RParen),
			byte(lisp.RParen),
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

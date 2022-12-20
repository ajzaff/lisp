package blisp

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/print"
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

func TestEncodedLen(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input lisp.Val
		want  int
	}{{
		name: "empty",
	}, {
		name:  "id",
		input: mustParse(t, "main"),
		want:  6,
	}, {
		name:  "int",
		input: mustParse(t, "1"),
		want:  2,
	}, {
		name:  "cons",
		input: mustParse(t, "(a)"),
		want:  5,
	}, {
		name:  "nested cons",
		input: mustParse(t, "(a(b)c)"),
		want:  7,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			if got := EncodedLen(tc.input); got != tc.want {
				var buf bytes.Buffer
				print.StdPrinter(&buf).Print(tc.input)
				t.Errorf("EncodedLen(%v): got %d, want %d", buf.String(), got, tc.want)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		want  []byte
	}{{
		name: "empty",
	}, {
		name:  "complex nested cons",
		input: "(1 (2 (3 4)))",
		want:  []byte{0x05, 0x02, 0x00, 0x31, 0x05, 0x02, 0x00, 0x32, 0x05, 0x02, 0x00, 0x33, 0x02, 0x00},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			v := mustParse(t, tc.input)
			var buf bytes.Buffer
			e := NewEncoder(&buf)
			e.Encode(v)
			got := buf.Bytes()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Encode(%q): got diff (-want, +got):\n%v", tc.name, diff)
			}
		})
	}
}

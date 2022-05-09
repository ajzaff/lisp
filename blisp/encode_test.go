package blisp

import (
	"bytes"
	"testing"

	"github.com/ajzaff/lisp"
	"github.com/google/go-cmp/cmp"
)

func mustParse(t *testing.T, src string) lisp.Val {
	n, err := lisp.Parser{}.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	return n[0].Val()
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
		name:  "float",
		input: mustParse(t, "1.125"),
		want:  10,
	}, {
		name:  "string",
		input: mustParse(t, `"abc"`),
		want:  5,
	}, {
		name:  "expr",
		input: mustParse(t, "(a)"),
		want:  5,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			if got := EncodedLen(tc.input); got != tc.want {
				var buf bytes.Buffer
				lisp.StdPrinter(&buf).Print(tc.input)
				t.Errorf("EncodedLen(%v): got %d, want %d", buf.String(), got, tc.want)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   string
		want    []byte
		wantErr bool
	}{{
		name: "empty",
		want: []byte{0x05, 0x02, 0x00, 0x31, 0x05, 0x02, 0x00, 0x32, 0x05, 0x02, 0x00, 0x33, 0x02, 0x00},
	}, {
		name:  "complex nested expr",
		input: "(1 (2 (3 4)))",
		want:  []byte{0x05, 0x02, 0x00, 0x31, 0x05, 0x02, 0x00, 0x32, 0x05, 0x02, 0x00, 0x33, 0x02, 0x00},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			ns, _ := lisp.Parser{}.Parse("(1 (2 (3 4)))")
			if len(ns) > 1 {
				t.Fatalf("Parse(%q): test expects single Val", tc.name)
			}
			v := ns[0].Val()
			var buf bytes.Buffer
			e := NewEncoder(&buf)
			if gotErr := e.Encode(v); (gotErr != nil) != tc.wantErr {
				t.Fatalf("Encode(%q): got err = %v, want err = %v", tc.name, gotErr, tc.wantErr)
			}
			got := buf.Bytes()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Encode(%q): got diff (-want, +got):\n%v", tc.name, diff)
			}
		})
	}
}

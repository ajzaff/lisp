package binnit

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ajzaff/innit"
)

func mustParse(t *testing.T, src string) innit.Node {
	n, err := innit.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

func TestEncodedLen(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input innit.Node
		want  int
	}{{
		name: "empty",
	}, {
		name:  "id",
		input: mustParse(t, "main"),
		want:  7,
	}, {
		name:  "int",
		input: mustParse(t, "1"),
		want:  4,
	}, {
		name:  "float",
		input: mustParse(t, "1.125"),
		want:  8,
	}, {
		name:  "string",
		input: mustParse(t, `"abc"`),
		want:  6,
	}, {
		name:  "expr",
		input: mustParse(t, "(a)"),
		want:  6,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			if got := EncodedLen(tc.input); got != tc.want {
				var buf bytes.Buffer
				innit.StdPrinter(&buf).Print(tc.input)
				t.Errorf("EncodedLen(%v): got %d, want %d", buf.String(), got, tc.want)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	n, _ := innit.Parse("(1 (2 (3 4)))")
	var buf bytes.Buffer
	e := NewEncoder(&buf)
	e.Encode(n)
	fmt.Println(EncodedLen(n))
	fmt.Println(buf.Bytes())
}

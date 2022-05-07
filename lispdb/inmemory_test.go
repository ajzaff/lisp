package lispdb

import (
	"testing"

	"github.com/ajzaff/lisp"
)

func mustParse(t *testing.T, src string) []lisp.Val {
	nodes, err := lisp.Parser{}.Parse(src)
	if err != nil {
		t.Fatalf("mustParse(%q) failed: %v", src, err)
	}
	vals := make([]lisp.Val, 0, len(nodes))
	for _, n := range nodes {
		vals = append(vals, n.Val())
	}
	return vals
}

func TestGenericStoreMultipleInMemory(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   []lisp.Val
		wantErr bool
		wantLen int
	}{{
		name: "empty",
	}, {
		name:  "empty parse",
		input: mustParse(t, ""),
	}, {
		name:    "int",
		input:   mustParse(t, "1"),
		wantLen: 1,
	}, {
		name:    "int{3}",
		input:   mustParse(t, "1 2 3"),
		wantLen: 3,
	}, {
		name:    "1{3}",
		input:   mustParse(t, "1 1 1"),
		wantLen: 1,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInMemory()
			gotErr := Store(db, tc.input, 1)
			if (gotErr != nil) != tc.wantErr {
				t.Fatalf("Store(%q): got err = %v, want err = %v", tc.name, gotErr, tc.wantErr)
			}
			if gotLen := db.Len(); gotLen != tc.wantLen {
				t.Fatalf("db.Len(): got len = %v, want len = %v", gotLen, tc.wantLen)
			}
		})
	}
}

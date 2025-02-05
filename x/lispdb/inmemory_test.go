package lispdb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ajzaff/lisp"
	"github.com/ajzaff/lisp/scan"
)

func mustParseMultiple(t *testing.T, src string) []lisp.Val {
	t.Helper()
	var vs []lisp.Val
	var sc scan.Scanner
	sc.Reset(strings.NewReader(src))
	for n := range sc.Nodes() {
		vs = append(vs, n.Val)
	}
	if err := sc.Err(); err != nil {
		panic(fmt.Sprintf("mustParse: failed to parse: %v", src))
	}
	return vs
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
		input: mustParseMultiple(t, ""),
	}, {
		name:    "int",
		input:   mustParseMultiple(t, "1"),
		wantLen: 1,
	}, {
		name:    "int{3}",
		input:   mustParseMultiple(t, "1 2 3"),
		wantLen: 3,
	}, {
		name:    "1{3}",
		input:   mustParseMultiple(t, "1 1 1"),
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

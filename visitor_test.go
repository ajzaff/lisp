package lisp

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testVisits struct {
	Visitor string
	Val     Val
}

func TestVisitor(t *testing.T) {
	for _, tc := range []struct {
		name       string
		input      Val
		overrides  func(v *Visitor)
		wantVisits []testVisits
	}{{
		name: "empty",
	}, {
		name:  "int",
		input: Lit{Token: Int, Text: "1"},
		wantVisits: []testVisits{{
			Visitor: "Val",
			Val:     Lit{Token: Int, Text: "1"},
		}, {
			Visitor: "Lit",
			Val:     Lit{Token: Int, Text: "1"},
		}},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			var gotVisits []testVisits
			var v Visitor
			v.SetValVisitor(func(e Val) {
				gotVisits = append(gotVisits, testVisits{
					Visitor: "Val",
					Val:     e,
				})
			})
			v.SetLitVisitor(func(e Lit) {
				gotVisits = append(gotVisits, testVisits{
					Visitor: "Lit",
					Val:     e,
				})
			})
			v.SetBeforeConsVisitor(func(e *Cons) {
				gotVisits = append(gotVisits, testVisits{
					Visitor: "BeforeCons",
					Val:     e,
				})
			})
			v.SetConsVisitor(func(e *Cons) {
				gotVisits = append(gotVisits, testVisits{
					Visitor: "Cons",
					Val:     e,
				})
			})
			v.SetAfterConsVisitor(func(e *Cons) {
				gotVisits = append(gotVisits, testVisits{
					Visitor: "AfterCons",
					Val:     e,
				})
			})
			if overrides := tc.overrides; overrides != nil {
				tc.overrides(&v)
			}
			v.Visit(tc.input)
			if diff := cmp.Diff(tc.wantVisits, gotVisits); diff != "" {
				t.Errorf("Visit(%q): got visit diff: (-want, +got):\n%v", tc.name, diff)
			}
		})
	}
}

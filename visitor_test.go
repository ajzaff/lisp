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
		input: Lit{Token: Number, Text: "1"},
		wantVisits: []testVisits{{
			Visitor: "BeforeVal",
			Val:     Lit{Token: Number, Text: "1"},
		}, {
			Visitor: "Lit",
			Val:     Lit{Token: Number, Text: "1"},
		}, {
			Visitor: "AfterVal",
			Val:     Lit{Token: Number, Text: "1"},
		}},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			var gotVisits []testVisits
			var v Visitor
			v.SetBeforeValVisitor(func(e Val) {
				gotVisits = append(gotVisits, testVisits{
					Visitor: "BeforeVal",
					Val:     e,
				})
			})
			v.SetAfterValVisitor(func(e Val) {
				gotVisits = append(gotVisits, testVisits{
					Visitor: "AfterVal",
					Val:     e,
				})
			})
			v.SetLitVisitor(func(e Lit) {
				gotVisits = append(gotVisits, testVisits{
					Visitor: "Lit",
					Val:     e,
				})
			})
			v.SetBeforeExprVisitor(func(e Expr) {
				gotVisits = append(gotVisits, testVisits{
					Visitor: "BeforeExpr",
					Val:     e,
				})
			})
			v.SetAfterExprVisitor(func(e Expr) {
				gotVisits = append(gotVisits, testVisits{
					Visitor: "AfterExpr",
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

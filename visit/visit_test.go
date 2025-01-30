package visit

import (
	"testing"

	"github.com/ajzaff/lisp"
	"github.com/google/go-cmp/cmp"
)

type testVisits struct {
	Visitor string
	Val     lisp.Val
}

func TestVisitor(t *testing.T) {
	for _, tc := range []struct {
		name       string
		input      lisp.Val
		overrides  func(v *Visitor)
		wantVisits []testVisits
	}{{
		name: "nil node has no visits",
	}, {
		name:  "nat",
		input: lisp.Lit{Token: lisp.Nat, Text: "1"},
		wantVisits: []testVisits{{
			Visitor: "Val",
			Val:     lisp.Lit{Token: lisp.Nat, Text: "1"},
		}, {
			Visitor: "Lit",
			Val:     lisp.Lit{Token: lisp.Nat, Text: "1"},
		}},
	}, {
		name:  "nil group",
		input: (lisp.Group)(nil),
		wantVisits: []testVisits{{
			Visitor: "Val",
			Val:     (lisp.Group)(nil),
		}, {
			Visitor: "BeforeGroup",
			Val:     (lisp.Group)(nil),
		}, {
			Visitor: "AfterGroup",
			Val:     (lisp.Group)(nil),
		}},
	}, {
		name:  "empty group",
		input: lisp.Group{},
		wantVisits: []testVisits{{
			Visitor: "Val",
			Val:     lisp.Group{},
		}, {
			Visitor: "BeforeGroup",
			Val:     lisp.Group{},
		}, {
			Visitor: "AfterGroup",
			Val:     lisp.Group{},
		}},
	}, {
		name: "simple nested group",
		// (a b(c))
		input: lisp.Group{
			lisp.Lit{Token: lisp.Id, Text: "a"},
			lisp.Lit{Token: lisp.Id, Text: "b"},
			lisp.Group{
				lisp.Lit{Token: lisp.Id, Text: "c"},
			},
		},
		wantVisits: []testVisits{{
			Visitor: "Val",
			Val: lisp.Group{
				lisp.Lit{Token: lisp.Id, Text: "a"},
				lisp.Lit{Token: lisp.Id, Text: "b"},
				lisp.Group{
					lisp.Lit{Token: lisp.Id, Text: "c"},
				},
			},
		}, {
			Visitor: "BeforeGroup",
			Val: lisp.Group{
				lisp.Lit{Token: lisp.Id, Text: "a"},
				lisp.Lit{Token: lisp.Id, Text: "b"},
				lisp.Group{
					lisp.Lit{Token: lisp.Id, Text: "c"},
				},
			},
		}, {
			Visitor: "Val",
			Val:     lisp.Lit{Token: lisp.Id, Text: "a"},
		}, {
			Visitor: "Lit",
			Val:     lisp.Lit{Token: lisp.Id, Text: "a"},
		}, {
			Visitor: "Val",
			Val:     lisp.Lit{Token: lisp.Id, Text: "b"},
		}, {
			Visitor: "Lit",
			Val:     lisp.Lit{Token: lisp.Id, Text: "b"},
		}, {
			Visitor: "Val",
			Val: lisp.Group{
				lisp.Lit{Token: lisp.Id, Text: "c"},
			},
		}, {
			Visitor: "BeforeGroup",
			Val: lisp.Group{
				lisp.Lit{Token: lisp.Id, Text: "c"},
			},
		}, {
			Visitor: "Val",
			Val:     lisp.Lit{Token: lisp.Id, Text: "c"},
		}, {
			Visitor: "Lit",
			Val:     lisp.Lit{Token: lisp.Id, Text: "c"},
		}, {
			Visitor: "AfterGroup",
			Val: lisp.Group{
				lisp.Lit{Token: lisp.Id, Text: "c"},
			},
		}, {
			Visitor: "AfterGroup",
			Val: lisp.Group{
				lisp.Lit{Token: lisp.Id, Text: "a"},
				lisp.Lit{Token: lisp.Id, Text: "b"},
				lisp.Group{
					lisp.Lit{Token: lisp.Id, Text: "c"},
				},
			},
		}},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			var gotVisits []testVisits
			var v Visitor
			v.SetValVisitor(func(e lisp.Val) {
				gotVisits = append(gotVisits, testVisits{
					Visitor: "Val",
					Val:     e,
				})
			})
			v.SetLitVisitor(func(e lisp.Lit) {
				gotVisits = append(gotVisits, testVisits{
					Visitor: "Lit",
					Val:     e,
				})
			})
			v.SetBeforeGroupVisitor(func(e lisp.Group) {
				gotVisits = append(gotVisits, testVisits{
					Visitor: "BeforeGroup",
					Val:     e,
				})
			})
			v.SetAfterGroupVisitor(func(e lisp.Group) {
				gotVisits = append(gotVisits, testVisits{
					Visitor: "AfterGroup",
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

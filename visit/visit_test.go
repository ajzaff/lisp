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
		name: "nil has no visits",
	}, {
		name:  "int",
		input: lisp.Lit{Token: lisp.Nat, Text: "1"},
		wantVisits: []testVisits{{
			Visitor: "Val",
			Val:     lisp.Lit{Token: lisp.Nat, Text: "1"},
		}, {
			Visitor: "Lit",
			Val:     lisp.Lit{Token: lisp.Nat, Text: "1"},
		}},
	}, {
		name:  "nil cons",
		input: (*lisp.Cons)(nil),
		wantVisits: []testVisits{{
			Visitor: "Val",
			Val:     (*lisp.Cons)(nil),
		}, {
			Visitor: "BeforeCons",
			Val:     (*lisp.Cons)(nil),
		}, {
			Visitor: "Cons",
			Val:     (*lisp.Cons)(nil),
		}, {
			Visitor: "AfterCons",
			Val:     (*lisp.Cons)(nil),
		}},
	}, {
		name:  "nil and empty cons are the same",
		input: &lisp.Cons{},
		wantVisits: []testVisits{{
			Visitor: "Val",
			Val:     &lisp.Cons{},
		}, {
			Visitor: "BeforeCons",
			Val:     &lisp.Cons{},
		}, {
			Visitor: "Cons",
			Val:     &lisp.Cons{},
		}, {
			Visitor: "AfterCons",
			Val:     &lisp.Cons{},
		}},
	}, {
		name: "cons",
		// (a b(c))
		input: &lisp.Cons{
			Val: lisp.Lit{Token: lisp.Id, Text: "a"},
			Cons: &lisp.Cons{
				Val: lisp.Lit{Token: lisp.Id, Text: "b"},
				Cons: &lisp.Cons{
					Val: &lisp.Cons{
						Val: lisp.Lit{Token: lisp.Id, Text: "c"},
					},
				},
			},
		},
		wantVisits: []testVisits{{
			Visitor: "Val",
			Val: &lisp.Cons{
				Val: lisp.Lit{Token: lisp.Id, Text: "a"},
				Cons: &lisp.Cons{
					Val: lisp.Lit{Token: lisp.Id, Text: "b"},
					Cons: &lisp.Cons{
						Val: &lisp.Cons{
							Val: lisp.Lit{Token: lisp.Id, Text: "c"},
						},
					},
				},
			},
		}, {
			Visitor: "BeforeCons",
			Val: &lisp.Cons{
				Val: lisp.Lit{Token: lisp.Id, Text: "a"},
				Cons: &lisp.Cons{
					Val: lisp.Lit{Token: lisp.Id, Text: "b"},
					Cons: &lisp.Cons{
						Val: &lisp.Cons{
							Val: lisp.Lit{Token: lisp.Id, Text: "c"},
						},
					},
				},
			},
		}, {
			Visitor: "Cons",
			Val: &lisp.Cons{
				Val: lisp.Lit{Token: lisp.Id, Text: "a"},
				Cons: &lisp.Cons{
					Val: lisp.Lit{Token: lisp.Id, Text: "b"},
					Cons: &lisp.Cons{
						Val: &lisp.Cons{
							Val: lisp.Lit{Token: lisp.Id, Text: "c"},
						},
					},
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
			Val: &lisp.Cons{
				Val: lisp.Lit{Token: lisp.Id, Text: "b"},
				Cons: &lisp.Cons{
					Val: &lisp.Cons{
						Val: lisp.Lit{Token: lisp.Id, Text: "c"},
					},
				},
			},
		}, {
			Visitor: "Cons",
			Val: &lisp.Cons{
				Val: lisp.Lit{Token: lisp.Id, Text: "b"},
				Cons: &lisp.Cons{
					Val: &lisp.Cons{
						Val: lisp.Lit{Token: lisp.Id, Text: "c"},
					},
				},
			},
		}, {
			Visitor: "Val",
			Val:     lisp.Lit{Token: lisp.Id, Text: "b"},
		}, {
			Visitor: "Lit",
			Val:     lisp.Lit{Token: lisp.Id, Text: "b"},
		}, {
			Visitor: "Val",
			Val: &lisp.Cons{
				Val: &lisp.Cons{
					Val: lisp.Lit{Token: lisp.Id, Text: "c"},
				},
			},
		}, {
			Visitor: "Cons",
			Val: &lisp.Cons{
				Val: &lisp.Cons{
					Val: lisp.Lit{Token: lisp.Id, Text: "c"},
				},
			},
		}, {
			Visitor: "Val",
			Val: &lisp.Cons{
				Val: lisp.Lit{Token: lisp.Id, Text: "c"},
			},
		}, {
			Visitor: "BeforeCons",
			Val: &lisp.Cons{
				Val: lisp.Lit{Token: lisp.Id, Text: "c"},
			},
		}, {
			Visitor: "Cons",
			Val: &lisp.Cons{
				Val: lisp.Lit{Token: lisp.Id, Text: "c"},
			},
		}, {
			Visitor: "Val",
			Val:     lisp.Lit{Token: lisp.Id, Text: "c"},
		}, {
			Visitor: "Lit",
			Val:     lisp.Lit{Token: lisp.Id, Text: "c"},
		}, {
			Visitor: "AfterCons",
			Val: &lisp.Cons{
				Val: lisp.Lit{Token: lisp.Id, Text: "c"},
			},
		}, {
			Visitor: "AfterCons",
			Val: &lisp.Cons{
				Val: &lisp.Cons{
					Val: lisp.Lit{Token: lisp.Id, Text: "c"},
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
			v.SetBeforeConsVisitor(func(e *lisp.Cons) {
				gotVisits = append(gotVisits, testVisits{
					Visitor: "BeforeCons",
					Val:     e,
				})
			})
			v.SetConsVisitor(func(e *lisp.Cons) {
				gotVisits = append(gotVisits, testVisits{
					Visitor: "Cons",
					Val:     e,
				})
			})
			v.SetAfterConsVisitor(func(e *lisp.Cons) {
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

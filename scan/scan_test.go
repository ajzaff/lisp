package scan

import (
	"strings"
	"testing"

	"github.com/ajzaff/lisp"
	"github.com/google/go-cmp/cmp"
)

type scanTestCase struct {
	name         string
	input        string
	wantPos      []lisp.Pos
	wantTok      []lisp.Token
	wantText     []string
	wantNodePos  []lisp.Pos
	wantNode     []lisp.Val
	wantTokenErr bool
	wantNodeErr  bool
}

func (tc scanTestCase) scanTokenTest(t *testing.T) {
	if !tc.wantTokenErr && len(tc.wantPos)%2 != 0 {
		t.Fatalf("Tokenize(%q) wants invalid result (cannot have odd length when wantErr=true): %v", tc.name, tc.wantPos)
	}
	var (
		gotPos  []lisp.Pos
		gotTok  []lisp.Token
		gotText []string
	)
	var sc TokenScanner
	sc.Reset(strings.NewReader(tc.input))
	for sc.Scan() {
		pos, tok, text := sc.Token()
		gotPos = append(gotPos, pos, pos+lisp.Pos(len(text)))
		gotTok = append(gotTok, tok)
		gotText = append(gotText, text)
	}
	gotTokenErr := sc.Err()
	if diff := cmp.Diff(tc.wantPos, gotPos); diff != "" {
		t.Errorf("Token(%q) got pos diff (-want, +got):\n%s", tc.name, diff)
	}
	if diff := cmp.Diff(tc.wantTok, gotTok); diff != "" {
		t.Errorf("Token(%q) got Token diff (-want, +got):\n%s", tc.name, diff)
	}
	if diff := cmp.Diff(tc.wantText, gotText); diff != "" {
		t.Errorf("Token(%q) got text diff (-want, +got):\n%s", tc.name, diff)
	}
	if (gotTokenErr != nil) != tc.wantTokenErr {
		t.Errorf("Token(%q) got err: %v, want err? %v", tc.name, gotTokenErr, tc.wantTokenErr)
	}

	// TODO: Split this to another test case.
	sc.Reset(strings.NewReader(tc.input))
	var s NodeScanner
	s.Reset(&sc)
	var gotNodePos []lisp.Pos
	var gotVal []lisp.Val
	for s.Scan() {
		pos, end, v := s.Node()
		gotNodePos = append(gotNodePos, pos, end)
		gotVal = append(gotVal, v)
	}
	gotNodeErr := s.Err()
	if diff := cmp.Diff(tc.wantNodePos, gotNodePos); diff != "" {
		t.Errorf("Node(%q) got pos diff (-want, +got):\n%s", tc.name, diff)
	}
	if diff := cmp.Diff(tc.wantNode, gotVal); diff != "" {
		t.Errorf("Node(%q) got Val diff (-want, +got):\n%s", tc.name, diff)
	}
	if (gotNodeErr != nil) != tc.wantNodeErr {
		t.Errorf("Node(%q) got err: %v, want err? %v", tc.name, gotNodeErr, tc.wantNodeErr)
	}
}

func TestTokenizeLit(t *testing.T) {
	for _, tc := range []scanTestCase{{
		name:        "id",
		input:       "foo",
		wantPos:     []lisp.Pos{0, 3},
		wantTok:     []lisp.Token{lisp.Id},
		wantText:    []string{"foo"},
		wantNodePos: []lisp.Pos{0, 3},
		wantNode:    []lisp.Val{lisp.Lit{Token: lisp.Id, Text: "foo"}},
	}, {
		name:        "space id",
		input:       "  \t\n x",
		wantPos:     []lisp.Pos{5, 6},
		wantTok:     []lisp.Token{lisp.Id},
		wantText:    []string{"x"},
		wantNodePos: []lisp.Pos{5, 6},
		wantNode:    []lisp.Val{lisp.Lit{Token: lisp.Id, Text: "x"}},
	}, {
		name:        "id id id",
		input:       "a b c",
		wantPos:     []lisp.Pos{0, 1, 2, 3, 4, 5},
		wantTok:     []lisp.Token{lisp.Id, lisp.Id, lisp.Id},
		wantText:    []string{"a", "b", "c"},
		wantNodePos: []lisp.Pos{0, 1, 2, 3, 4, 5},
		wantNode: []lisp.Val{
			lisp.Lit{Token: lisp.Id, Text: "a"},
			lisp.Lit{Token: lisp.Id, Text: "b"},
			lisp.Lit{Token: lisp.Id, Text: "c"},
		},
	}, {
		name:        "id id id 2",
		input:       "ab cd ef",
		wantPos:     []lisp.Pos{0, 2, 3, 5, 6, 8},
		wantTok:     []lisp.Token{lisp.Id, lisp.Id, lisp.Id},
		wantText:    []string{"ab", "cd", "ef"},
		wantNodePos: []lisp.Pos{0, 2, 3, 5, 6, 8},
		wantNode: []lisp.Val{
			lisp.Lit{Token: lisp.Id, Text: "ab"},
			lisp.Lit{Token: lisp.Id, Text: "cd"},
			lisp.Lit{Token: lisp.Id, Text: "ef"},
		},
	}, {
		name:        "int",
		input:       "0",
		wantPos:     []lisp.Pos{0, 1},
		wantTok:     []lisp.Token{lisp.Nat},
		wantText:    []string{"0"},
		wantNodePos: []lisp.Pos{0, 1},
		wantNode:    []lisp.Val{lisp.Lit{Token: lisp.Nat, Text: "0"}},
	}, {
		name:        "int 2",
		input:       "0 1 2",
		wantPos:     []lisp.Pos{0, 1, 2, 3, 4, 5},
		wantTok:     []lisp.Token{lisp.Nat, lisp.Nat, lisp.Nat},
		wantText:    []string{"0", "1", "2"},
		wantNodePos: []lisp.Pos{0, 1, 2, 3, 4, 5},
		wantNode: []lisp.Val{
			lisp.Lit{Token: lisp.Nat, Text: "0"},
			lisp.Lit{Token: lisp.Nat, Text: "1"},
			lisp.Lit{Token: lisp.Nat, Text: "2"},
		},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			tc.scanTokenTest(t)
		})
	}
}

func TestTokenizeLitErrors(t *testing.T) {
	for _, tc := range []scanTestCase{} {
		t.Run(tc.name, func(t *testing.T) {
			tc.scanTokenTest(t)
		})
	}
}

func TestTokenizeCons(t *testing.T) {
	for _, tc := range []scanTestCase{{
		name: "empty",
	}, {
		name:  "whitespace",
		input: "  \t\r\n",
	}, {
		name:        "empty cons",
		input:       "()",
		wantPos:     []lisp.Pos{0, 1, 1, 2},
		wantTok:     []lisp.Token{lisp.LParen, lisp.RParen},
		wantText:    []string{"(", ")"},
		wantNodePos: []lisp.Pos{0, 2},
		wantNode:    []lisp.Val{&lisp.Cons{}}, // FIXME: Expect canonial empty Cons.
	}, {
		name:        "nested cons",
		input:       "(())",
		wantPos:     []lisp.Pos{0, 1, 1, 2, 2, 3, 3, 4},
		wantTok:     []lisp.Token{lisp.LParen, lisp.LParen, lisp.RParen, lisp.RParen},
		wantText:    []string{"(", "(", ")", ")"},
		wantNodePos: []lisp.Pos{0, 4},
		wantNode:    []lisp.Val{&lisp.Cons{Val: &lisp.Cons{}}},
	}, {
		name:        "cons",
		input:       "(abc)",
		wantPos:     []lisp.Pos{0, 1, 1, 4, 4, 5},
		wantTok:     []lisp.Token{lisp.LParen, lisp.Id, lisp.RParen},
		wantText:    []string{"(", "abc", ")"},
		wantNodePos: []lisp.Pos{0, 5},
		wantNode: []lisp.Val{&lisp.Cons{
			Val: lisp.Lit{Token: lisp.Id, Text: "abc"},
		}},
	}, {
		name:        "cons 2",
		input:       "(add 1 2)",
		wantPos:     []lisp.Pos{0, 1, 1, 4, 5, 6, 7, 8, 8, 9},
		wantTok:     []lisp.Token{lisp.LParen, lisp.Id, lisp.Nat, lisp.Nat, lisp.RParen},
		wantText:    []string{"(", "add", "1", "2", ")"},
		wantNodePos: []lisp.Pos{0, 9},
		wantNode: []lisp.Val{&lisp.Cons{
			Val: lisp.Lit{Token: lisp.Id, Text: "add"},
			Cons: &lisp.Cons{
				Val: lisp.Lit{Token: lisp.Nat, Text: "1"},
				Cons: &lisp.Cons{
					Val: lisp.Lit{Token: lisp.Nat, Text: "2"},
				},
			},
		}},
	}, {
		name:        "cons 3",
		input:       "(add (sub 3 2) 2)",
		wantPos:     []lisp.Pos{0, 1, 1, 4, 5, 6, 6, 9, 10, 11, 12, 13, 13, 14, 15, 16, 16, 17},
		wantTok:     []lisp.Token{lisp.LParen, lisp.Id, lisp.LParen, lisp.Id, lisp.Nat, lisp.Nat, lisp.RParen, lisp.Nat, lisp.RParen},
		wantText:    []string{"(", "add", "(", "sub", "3", "2", ")", "2", ")"},
		wantNodePos: []lisp.Pos{0, 17},
		wantNode: []lisp.Val{&lisp.Cons{
			Val: lisp.Lit{Token: lisp.Id, Text: "add"},
			Cons: &lisp.Cons{
				Val: &lisp.Cons{
					Val: lisp.Lit{Token: lisp.Id, Text: "sub"},
					Cons: &lisp.Cons{
						Val: lisp.Lit{Token: lisp.Nat, Text: "3"},
						Cons: &lisp.Cons{
							Val: lisp.Lit{Token: lisp.Nat, Text: "2"},
						},
					},
				},
				Cons: &lisp.Cons{
					Val: lisp.Lit{Token: lisp.Nat, Text: "2"},
				},
			},
		}},
	}, {
		name:        "cons 4",
		input:       "((a))",
		wantPos:     []lisp.Pos{0, 1, 1, 2, 2, 3, 3, 4, 4, 5},
		wantTok:     []lisp.Token{lisp.LParen, lisp.LParen, lisp.Id, lisp.RParen, lisp.RParen},
		wantText:    []string{"(", "(", "a", ")", ")"},
		wantNodePos: []lisp.Pos{0, 5},
		wantNode: []lisp.Val{&lisp.Cons{
			Val: &lisp.Cons{
				Val: lisp.Lit{Token: lisp.Id, Text: "a"},
			},
		}},
	}, {
		name:        "cons 5",
		input:       "(a)(b) (c)",
		wantPos:     []lisp.Pos{0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 7, 8, 8, 9, 9, 10},
		wantTok:     []lisp.Token{lisp.LParen, lisp.Id, lisp.RParen, lisp.LParen, lisp.Id, lisp.RParen, lisp.LParen, lisp.Id, lisp.RParen},
		wantText:    []string{"(", "a", ")", "(", "b", ")", "(", "c", ")"},
		wantNodePos: []lisp.Pos{0, 3, 3, 6, 7, 10},
		wantNode: []lisp.Val{
			&lisp.Cons{Val: lisp.Lit{Token: lisp.Id, Text: "a"}},
			&lisp.Cons{Val: lisp.Lit{Token: lisp.Id, Text: "b"}},
			&lisp.Cons{Val: lisp.Lit{Token: lisp.Id, Text: "c"}},
		},
	}, {
		name:        "cons 6",
		input:       "(div (q x) y)\n",
		wantPos:     []lisp.Pos{0, 1, 1, 4, 5, 6, 6, 7, 8, 9, 9, 10, 11, 12, 12, 13},
		wantTok:     []lisp.Token{lisp.LParen, lisp.Id, lisp.LParen, lisp.Id, lisp.Id, lisp.RParen, lisp.Id, lisp.RParen},
		wantText:    []string{"(", "div", "(", "q", "x", ")", "y", ")"},
		wantNodePos: []lisp.Pos{0, 13},
		wantNode: []lisp.Val{&lisp.Cons{
			Val: lisp.Lit{Token: lisp.Id, Text: "div"},
			Cons: &lisp.Cons{
				Val: &lisp.Cons{
					Val: lisp.Lit{Token: lisp.Id, Text: "q"},
					Cons: &lisp.Cons{
						Val: lisp.Lit{Token: lisp.Id, Text: "x"},
					},
				},
				Cons: &lisp.Cons{
					Val: lisp.Lit{Token: lisp.Id, Text: "y"},
				},
			},
		}},
	}, {
		name:        "cons 7",
		input:       "((a b))\n",
		wantPos:     []lisp.Pos{0, 1, 1, 2, 2, 3, 4, 5, 5, 6, 6, 7},
		wantTok:     []lisp.Token{lisp.LParen, lisp.LParen, lisp.Id, lisp.Id, lisp.RParen, lisp.RParen},
		wantText:    []string{"(", "(", "a", "b", ")", ")"},
		wantNodePos: []lisp.Pos{0, 7},
		wantNode: []lisp.Val{&lisp.Cons{
			Val: &lisp.Cons{
				Val: lisp.Lit{Token: lisp.Id, Text: "a"},
				Cons: &lisp.Cons{
					Val: lisp.Lit{Token: lisp.Id, Text: "b"},
				},
			},
		}},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			tc.scanTokenTest(t)
		})
	}
}

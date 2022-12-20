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
	wantNodes    []lisp.Node
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
		t.Errorf("Tokenize(%q) got diff (-want, +got):\n%s", tc.name, diff)
	}
	if diff := cmp.Diff(tc.wantTok, gotTok); diff != "" {
		t.Errorf("Tokenize(%q) got token diff (-want, +got):\n%s", tc.name, diff)
	}
	if diff := cmp.Diff(tc.wantText, gotText); diff != "" {
		t.Errorf("Tokenize(%q) got text diff (-want, +got):\n%s", tc.name, diff)
	}
	if (gotTokenErr != nil) != tc.wantTokenErr {
		t.Errorf("Tokenize(%q) got token err: %v, want err? %v", tc.name, gotTokenErr, tc.wantTokenErr)
	}

	// TODO: Split this to another test case.
	sc.Reset(strings.NewReader(tc.input))
	var s NodeScanner
	s.Reset(&sc)
	var gotNodes []lisp.Node
	for s.Scan() {
		gotNodes = append(gotNodes, s.Node())
	}
	gotNodeErr := s.Err()
	if diff := cmp.Diff(tc.wantNodes, gotNodes); diff != "" {
		t.Errorf("Parse(%q) got node diff (-want, +got):\n%s", tc.name, diff)
	}
	if (gotNodeErr != nil) != tc.wantNodeErr {
		t.Errorf("Parse(%q) got node err: %v, want err? %v", tc.name, gotNodeErr, tc.wantNodeErr)
	}
}

func TestTokenizeLit(t *testing.T) {
	for _, tc := range []scanTestCase{{
		name:      "id",
		input:     "foo",
		wantPos:   []lisp.Pos{0, 3},
		wantTok:   []lisp.Token{lisp.Id},
		wantText:  []string{"foo"},
		wantNodes: []lisp.Node{{Val: lisp.Lit{Token: lisp.Id, Text: "foo"}, End: 3}},
	}, {
		name:      "space id",
		input:     "  \t\n x",
		wantPos:   []lisp.Pos{5, 6},
		wantTok:   []lisp.Token{lisp.Id},
		wantText:  []string{"x"},
		wantNodes: []lisp.Node{{Pos: 5, Val: lisp.Lit{Token: lisp.Id, Text: "x"}, End: 6}},
	}, {
		name:     "id id id",
		input:    "a b c",
		wantPos:  []lisp.Pos{0, 1, 2, 3, 4, 5},
		wantTok:  []lisp.Token{lisp.Id, lisp.Id, lisp.Id},
		wantText: []string{"a", "b", "c"},
		wantNodes: []lisp.Node{
			{Val: lisp.Lit{Token: lisp.Id, Text: "a"}, End: 1},
			{Pos: 2, Val: lisp.Lit{Token: lisp.Id, Text: "b"}, End: 3},
			{Pos: 4, Val: lisp.Lit{Token: lisp.Id, Text: "c"}, End: 5},
		},
	}, {
		name:     "id id id 2",
		input:    "ab cd ef",
		wantPos:  []lisp.Pos{0, 2, 3, 5, 6, 8},
		wantTok:  []lisp.Token{lisp.Id, lisp.Id, lisp.Id},
		wantText: []string{"ab", "cd", "ef"},
		wantNodes: []lisp.Node{
			{Val: lisp.Lit{Token: lisp.Id, Text: "ab"}, End: 2},
			{Pos: 3, Val: lisp.Lit{Token: lisp.Id, Text: "cd"}, End: 5},
			{Pos: 6, Val: lisp.Lit{Token: lisp.Id, Text: "ef"}, End: 8},
		},
	}, {
		name:      "int",
		input:     "0",
		wantPos:   []lisp.Pos{0, 1},
		wantTok:   []lisp.Token{lisp.Nat},
		wantText:  []string{"0"},
		wantNodes: []lisp.Node{{Val: lisp.Lit{Token: lisp.Nat, Text: "0"}, End: 1}},
	}, {
		name:     "int 2",
		input:    "0 1 2",
		wantPos:  []lisp.Pos{0, 1, 2, 3, 4, 5},
		wantTok:  []lisp.Token{lisp.Nat, lisp.Nat, lisp.Nat},
		wantText: []string{"0", "1", "2"},
		wantNodes: []lisp.Node{
			{Val: lisp.Lit{Token: lisp.Nat, Text: "0"}, End: 1},
			{Pos: 2, Val: lisp.Lit{Token: lisp.Nat, Text: "1"}, End: 3},
			{Pos: 4, Val: lisp.Lit{Token: lisp.Nat, Text: "2"}, End: 5},
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
		name:      "empty cons",
		input:     "()",
		wantPos:   []lisp.Pos{0, 1, 1, 2},
		wantTok:   []lisp.Token{lisp.LParen, lisp.RParen},
		wantText:  []string{"(", ")"},
		wantNodes: []lisp.Node{{Val: &lisp.Cons{}, End: 2}},
	}, {
		name:      "nested cons",
		input:     "(())",
		wantPos:   []lisp.Pos{0, 1, 1, 2, 2, 3, 3, 4},
		wantTok:   []lisp.Token{lisp.LParen, lisp.LParen, lisp.RParen, lisp.RParen},
		wantText:  []string{"(", "(", ")", ")"},
		wantNodes: []lisp.Node{{Val: &lisp.Cons{Node: lisp.Node{Pos: 1, Val: &lisp.Cons{}, End: 3}}, End: 4}},
	}, {
		name:     "cons",
		input:    "(abc)",
		wantPos:  []lisp.Pos{0, 1, 1, 4, 4, 5},
		wantTok:  []lisp.Token{lisp.LParen, lisp.Id, lisp.RParen},
		wantText: []string{"(", "abc", ")"},
		wantNodes: []lisp.Node{{Val: &lisp.Cons{
			Node: lisp.Node{Pos: 1, Val: lisp.Lit{Token: lisp.Id, Text: "abc"}, End: 4},
		}, End: 5}},
	}, {
		name:     "cons 2",
		input:    "(add 1 2)",
		wantPos:  []lisp.Pos{0, 1, 1, 4, 5, 6, 7, 8, 8, 9},
		wantTok:  []lisp.Token{lisp.LParen, lisp.Id, lisp.Nat, lisp.Nat, lisp.RParen},
		wantText: []string{"(", "add", "1", "2", ")"},
		wantNodes: []lisp.Node{{Val: &lisp.Cons{
			Node: lisp.Node{Pos: 1, Val: lisp.Lit{Token: lisp.Id, Text: "add"}, End: 4},
			Cons: &lisp.Cons{
				Node: lisp.Node{Pos: 5, Val: lisp.Lit{Token: lisp.Nat, Text: "1"}, End: 6},
				Cons: &lisp.Cons{
					Node: lisp.Node{Pos: 7, Val: lisp.Lit{Token: lisp.Nat, Text: "2"}, End: 8},
				},
			},
		}, End: 9}},
	}, {
		name:     "cons 3",
		input:    "(add (sub 3 2) 2)",
		wantPos:  []lisp.Pos{0, 1, 1, 4, 5, 6, 6, 9, 10, 11, 12, 13, 13, 14, 15, 16, 16, 17},
		wantTok:  []lisp.Token{lisp.LParen, lisp.Id, lisp.LParen, lisp.Id, lisp.Nat, lisp.Nat, lisp.RParen, lisp.Nat, lisp.RParen},
		wantText: []string{"(", "add", "(", "sub", "3", "2", ")", "2", ")"},
		wantNodes: []lisp.Node{{Val: &lisp.Cons{
			Node: lisp.Node{Pos: 1, Val: lisp.Lit{Token: lisp.Id, Text: "add"}, End: 4},
			Cons: &lisp.Cons{
				Node: lisp.Node{Pos: 5, Val: &lisp.Cons{
					Node: lisp.Node{Pos: 6, Val: lisp.Lit{Token: lisp.Id, Text: "sub"}, End: 9},
					Cons: &lisp.Cons{
						Node: lisp.Node{Pos: 10, Val: lisp.Lit{Token: lisp.Nat, Text: "3"}, End: 11},
						Cons: &lisp.Cons{
							Node: lisp.Node{Pos: 12, Val: lisp.Lit{Token: lisp.Nat, Text: "2"}, End: 13},
						},
					},
				}, End: 14},
				Cons: &lisp.Cons{
					Node: lisp.Node{Pos: 15, Val: lisp.Lit{Token: lisp.Nat, Text: "2"}, End: 16},
				},
			}}, End: 17}},
	}, {
		name:     "cons 4",
		input:    "((a))",
		wantPos:  []lisp.Pos{0, 1, 1, 2, 2, 3, 3, 4, 4, 5},
		wantTok:  []lisp.Token{lisp.LParen, lisp.LParen, lisp.Id, lisp.RParen, lisp.RParen},
		wantText: []string{"(", "(", "a", ")", ")"},
		wantNodes: []lisp.Node{{Val: &lisp.Cons{
			Node: lisp.Node{Pos: 1, Val: &lisp.Cons{
				Node: lisp.Node{Pos: 2, Val: lisp.Lit{Token: lisp.Id, Text: "a"}, End: 3},
			}, End: 4},
		}, End: 5}},
	}, {
		name:     "cons 5",
		input:    "(a)(b) (c)",
		wantPos:  []lisp.Pos{0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 7, 8, 8, 9, 9, 10},
		wantTok:  []lisp.Token{lisp.LParen, lisp.Id, lisp.RParen, lisp.LParen, lisp.Id, lisp.RParen, lisp.LParen, lisp.Id, lisp.RParen},
		wantText: []string{"(", "a", ")", "(", "b", ")", "(", "c", ")"},
		wantNodes: []lisp.Node{
			{Val: &lisp.Cons{Node: lisp.Node{Pos: 1, Val: lisp.Lit{Token: lisp.Id, Text: "a"}, End: 2}}, End: 3},
			{Pos: 3, Val: &lisp.Cons{Node: lisp.Node{Pos: 4, Val: lisp.Lit{Token: lisp.Id, Text: "b"}, End: 5}}, End: 6},
			{Pos: 7, Val: &lisp.Cons{Node: lisp.Node{Pos: 8, Val: lisp.Lit{Token: lisp.Id, Text: "c"}, End: 9}}, End: 10},
		},
	}, {
		name:     "cons 6",
		input:    "(div (q x) y)\n",
		wantPos:  []lisp.Pos{0, 1, 1, 4, 5, 6, 6, 7, 8, 9, 9, 10, 11, 12, 12, 13},
		wantTok:  []lisp.Token{lisp.LParen, lisp.Id, lisp.LParen, lisp.Id, lisp.Id, lisp.RParen, lisp.Id, lisp.RParen},
		wantText: []string{"(", "div", "(", "q", "x", ")", "y", ")"},
		wantNodes: []lisp.Node{{Val: &lisp.Cons{
			Node: lisp.Node{Pos: 1, Val: lisp.Lit{Token: lisp.Id, Text: "div"}, End: 4},
			Cons: &lisp.Cons{
				Node: lisp.Node{Pos: 5, Val: &lisp.Cons{
					Node: lisp.Node{Pos: 6, Val: lisp.Lit{Token: lisp.Id, Text: "q"}, End: 7},
					Cons: &lisp.Cons{
						Node: lisp.Node{Pos: 8, Val: lisp.Lit{Token: lisp.Id, Text: "x"}, End: 9},
					},
				}, End: 10},
				Cons: &lisp.Cons{
					Node: lisp.Node{Pos: 11, Val: lisp.Lit{Token: lisp.Id, Text: "y"}, End: 12},
				},
			},
		}, End: 13}},
	}, {
		name:     "cons 7",
		input:    "((a b))\n",
		wantPos:  []lisp.Pos{0, 1, 1, 2, 2, 3, 4, 5, 5, 6, 6, 7},
		wantTok:  []lisp.Token{lisp.LParen, lisp.LParen, lisp.Id, lisp.Id, lisp.RParen, lisp.RParen},
		wantText: []string{"(", "(", "a", "b", ")", ")"},
		wantNodes: []lisp.Node{{Val: &lisp.Cons{
			Node: lisp.Node{Pos: 1, Val: &lisp.Cons{
				Node: lisp.Node{Pos: 2, Val: lisp.Lit{Token: lisp.Id, Text: "a"}, End: 3},
				Cons: &lisp.Cons{
					Node: lisp.Node{Pos: 4, Val: lisp.Lit{Token: lisp.Id, Text: "b"}, End: 5},
				},
			}, End: 6},
		}, End: 7}},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			tc.scanTokenTest(t)
		})
	}
}

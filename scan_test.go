package lisp

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type scanTestCase struct {
	name         string
	input        string
	wantPos      []Pos
	wantTok      []Token
	wantText     []string
	wantNodes    []Node
	wantTokenErr bool
	wantNodeErr  bool
}

func (tc scanTestCase) scanTokenTest(t *testing.T) {
	if !tc.wantTokenErr && len(tc.wantPos)%2 != 0 {
		t.Fatalf("Tokenize(%q) wants invalid result (cannot have odd length when wantErr=true): %v", tc.name, tc.wantPos)
	}
	var (
		gotPos  []Pos
		gotTok  []Token
		gotText []string
	)
	var sc TokenScanner
	sc.Reset(strings.NewReader(tc.input))
	for sc.Scan() {
		pos, tok, text := sc.Token()
		gotPos = append(gotPos, pos, pos+Pos(len(text)))
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
	var gotNodes []Node
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
		wantPos:   []Pos{0, 3},
		wantTok:   []Token{Id},
		wantText:  []string{"foo"},
		wantNodes: []Node{{Val: Lit{Token: Id, Text: "foo"}, End: 3}},
	}, {
		name:      "space id",
		input:     "  \t\n x",
		wantPos:   []Pos{5, 6},
		wantTok:   []Token{Id},
		wantText:  []string{"x"},
		wantNodes: []Node{{Pos: 5, Val: Lit{Token: Id, Text: "x"}, End: 6}},
	}, {
		name:     "id id id",
		input:    "a b c",
		wantPos:  []Pos{0, 1, 2, 3, 4, 5},
		wantTok:  []Token{Id, Id, Id},
		wantText: []string{"a", "b", "c"},
		wantNodes: []Node{
			{Val: Lit{Token: Id, Text: "a"}, End: 1},
			{Pos: 2, Val: Lit{Token: Id, Text: "b"}, End: 3},
			{Pos: 4, Val: Lit{Token: Id, Text: "c"}, End: 5},
		},
	}, {
		name:     "id id id 2",
		input:    "ab cd ef",
		wantPos:  []Pos{0, 2, 3, 5, 6, 8},
		wantTok:  []Token{Id, Id, Id},
		wantText: []string{"ab", "cd", "ef"},
		wantNodes: []Node{
			{Val: Lit{Token: Id, Text: "ab"}, End: 2},
			{Pos: 3, Val: Lit{Token: Id, Text: "cd"}, End: 5},
			{Pos: 6, Val: Lit{Token: Id, Text: "ef"}, End: 8},
		},
	}, {
		name:      "int",
		input:     "0",
		wantPos:   []Pos{0, 1},
		wantTok:   []Token{Int},
		wantText:  []string{"0"},
		wantNodes: []Node{{Val: Lit{Token: Int, Text: "0"}, End: 1}},
	}, {
		name:     "int 2",
		input:    "0 1 2",
		wantPos:  []Pos{0, 1, 2, 3, 4, 5},
		wantTok:  []Token{Int, Int, Int},
		wantText: []string{"0", "1", "2"},
		wantNodes: []Node{
			{Val: Lit{Token: Int, Text: "0"}, End: 1},
			{Pos: 2, Val: Lit{Token: Int, Text: "1"}, End: 3},
			{Pos: 4, Val: Lit{Token: Int, Text: "2"}, End: 5},
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
		wantPos:   []Pos{0, 1, 1, 2},
		wantTok:   []Token{LParen, RParen},
		wantText:  []string{"(", ")"},
		wantNodes: []Node{{Val: &Cons{}, End: 2}},
	}, {
		name:      "nested cons",
		input:     "(())",
		wantPos:   []Pos{0, 1, 1, 2, 2, 3, 3, 4},
		wantTok:   []Token{LParen, LParen, RParen, RParen},
		wantText:  []string{"(", "(", ")", ")"},
		wantNodes: []Node{{Val: &Cons{Node: Node{Pos: 1, Val: &Cons{}, End: 3}}, End: 4}},
	}, {
		name:     "cons",
		input:    "(abc)",
		wantPos:  []Pos{0, 1, 1, 4, 4, 5},
		wantTok:  []Token{LParen, Id, RParen},
		wantText: []string{"(", "abc", ")"},
		wantNodes: []Node{{Val: &Cons{
			Node: Node{Pos: 1, Val: Lit{Token: Id, Text: "abc"}, End: 4},
		}, End: 5}},
	}, {
		name:     "cons 2",
		input:    "(add 1 2)",
		wantPos:  []Pos{0, 1, 1, 4, 5, 6, 7, 8, 8, 9},
		wantTok:  []Token{LParen, Id, Int, Int, RParen},
		wantText: []string{"(", "add", "1", "2", ")"},
		wantNodes: []Node{{Val: &Cons{
			Node: Node{Pos: 1, Val: Lit{Token: Id, Text: "add"}, End: 4},
			Cons: &Cons{
				Node: Node{Pos: 5, Val: Lit{Token: Int, Text: "1"}, End: 6},
				Cons: &Cons{
					Node: Node{Pos: 7, Val: Lit{Token: Int, Text: "2"}, End: 8},
				},
			},
		}, End: 9}},
	}, {
		name:     "cons 3",
		input:    "(add (sub 3 2) 2)",
		wantPos:  []Pos{0, 1, 1, 4, 5, 6, 6, 9, 10, 11, 12, 13, 13, 14, 15, 16, 16, 17},
		wantTok:  []Token{LParen, Id, LParen, Id, Int, Int, RParen, Int, RParen},
		wantText: []string{"(", "add", "(", "sub", "3", "2", ")", "2", ")"},
		wantNodes: []Node{{Val: &Cons{
			Node: Node{Pos: 1, Val: Lit{Token: Id, Text: "add"}, End: 4},
			Cons: &Cons{
				Node: Node{Pos: 5, Val: &Cons{
					Node: Node{Pos: 6, Val: Lit{Token: Id, Text: "sub"}, End: 9},
					Cons: &Cons{
						Node: Node{Pos: 10, Val: Lit{Token: Int, Text: "3"}, End: 11},
						Cons: &Cons{
							Node: Node{Pos: 12, Val: Lit{Token: Int, Text: "2"}, End: 13},
						},
					},
				}, End: 14},
				Cons: &Cons{
					Node: Node{Pos: 15, Val: Lit{Token: Int, Text: "2"}, End: 16},
				},
			}}, End: 17}},
	}, {
		name:     "cons 4",
		input:    "((a))",
		wantPos:  []Pos{0, 1, 1, 2, 2, 3, 3, 4, 4, 5},
		wantTok:  []Token{LParen, LParen, Id, RParen, RParen},
		wantText: []string{"(", "(", "a", ")", ")"},
		wantNodes: []Node{{Val: &Cons{
			Node: Node{Pos: 1, Val: &Cons{
				Node: Node{Pos: 2, Val: Lit{Token: Id, Text: "a"}, End: 3},
			}, End: 4},
		}, End: 5}},
	}, {
		name:     "cons 5",
		input:    "(a)(b) (c)",
		wantPos:  []Pos{0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 7, 8, 8, 9, 9, 10},
		wantTok:  []Token{LParen, Id, RParen, LParen, Id, RParen, LParen, Id, RParen},
		wantText: []string{"(", "a", ")", "(", "b", ")", "(", "c", ")"},
		wantNodes: []Node{
			{Val: &Cons{Node: Node{Pos: 1, Val: Lit{Token: Id, Text: "a"}, End: 2}}, End: 3},
			{Pos: 3, Val: &Cons{Node: Node{Pos: 4, Val: Lit{Token: Id, Text: "b"}, End: 5}}, End: 6},
			{Pos: 7, Val: &Cons{Node: Node{Pos: 8, Val: Lit{Token: Id, Text: "c"}, End: 9}}, End: 10},
		},
	}, {
		name:     "cons 6",
		input:    "(div (q x) y)\n",
		wantPos:  []Pos{0, 1, 1, 4, 5, 6, 6, 7, 8, 9, 9, 10, 11, 12, 12, 13},
		wantTok:  []Token{LParen, Id, LParen, Id, Id, RParen, Id, RParen},
		wantText: []string{"(", "div", "(", "q", "x", ")", "y", ")"},
		wantNodes: []Node{{Val: &Cons{
			Node: Node{Pos: 1, Val: Lit{Token: Id, Text: "div"}, End: 4},
			Cons: &Cons{
				Node: Node{Pos: 5, Val: &Cons{
					Node: Node{Pos: 6, Val: Lit{Token: Id, Text: "q"}, End: 7},
					Cons: &Cons{
						Node: Node{Pos: 8, Val: Lit{Token: Id, Text: "x"}, End: 9},
					},
				}, End: 10},
				Cons: &Cons{
					Node: Node{Pos: 11, Val: Lit{Token: Id, Text: "y"}, End: 12},
				},
			},
		}, End: 13}},
	}, {
		name:     "cons 7",
		input:    "((a b))\n",
		wantPos:  []Pos{0, 1, 1, 2, 2, 3, 4, 5, 5, 6, 6, 7},
		wantTok:  []Token{LParen, LParen, Id, Id, RParen, RParen},
		wantText: []string{"(", "(", "a", "b", ")", ")"},
		wantNodes: []Node{{Val: &Cons{
			Node: Node{Pos: 1, Val: &Cons{
				Node: Node{Pos: 2, Val: Lit{Token: Id, Text: "a"}, End: 3},
				Cons: &Cons{
					Node: Node{Pos: 4, Val: Lit{Token: Id, Text: "b"}, End: 5},
				},
			}, End: 6},
		}, End: 7}},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			tc.scanTokenTest(t)
		})
	}
}

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
	wantPos      []Pos
	wantTok      []lisp.Token
	wantText     []string
	wantNodePos  []Pos
	wantNode     []lisp.Val
	wantTokenErr bool
	wantNodeErr  bool
}

func (tc scanTestCase) scanTokenTest(t *testing.T) {
	t.Helper()
	if !tc.wantTokenErr && len(tc.wantPos)%2 != 0 {
		t.Fatalf("Tokenize(%q) wants invalid result (cannot have odd length when wantErr=true): %v", tc.name, tc.wantPos)
	}
	var (
		gotPos  []Pos
		gotTok  []lisp.Token
		gotText []string
	)
	var sc Scanner
	sc.Reset(strings.NewReader(tc.input))
	for tok := range sc.Tokens() {
		pos, tok, text := tok.Pos, tok.Tok, tok.Text
		gotPos = append(gotPos, pos, pos+Pos(len(text)))
		gotTok = append(gotTok, tok)
		gotText = append(gotText, text)
	}
	gotTokenErr := sc.Err()
	if diff := cmp.Diff(tc.wantPos, gotPos); diff != "" {
		t.Errorf("TestScanTokens(%q) got pos diff (-want, +got):\n%s", tc.name, diff)
	}
	if diff := cmp.Diff(tc.wantTok, gotTok); diff != "" {
		t.Errorf("TestScanTokens(%q) got Token diff (-want, +got):\n%s", tc.name, diff)
	}
	if diff := cmp.Diff(tc.wantText, gotText); diff != "" {
		t.Errorf("TestScanTokens(%q) got text diff (-want, +got):\n%s", tc.name, diff)
	}
	if gotErr := gotTokenErr != nil; gotErr != tc.wantTokenErr {
		t.Errorf("TestScanTokens(%q) got token err: %v, want err? %v", tc.name, gotTokenErr, tc.wantTokenErr)
	}

	sc.Reset(strings.NewReader(tc.input))
	var gotNodePos []Pos
	var gotVal []lisp.Val
	for n := range sc.Nodes() {
		pos, v, end := n.Pos, n.Val, n.End
		gotNodePos = append(gotNodePos, pos, end)
		gotVal = append(gotVal, v)
	}
	gotNodeErr := sc.Err()
	if diff := cmp.Diff(tc.wantNodePos, gotNodePos); diff != "" {
		t.Errorf("TestScanNodes(%q) got pos diff (-want, +got):\n%s", tc.name, diff)
	}
	if diff := cmp.Diff(tc.wantNode, gotVal); diff != "" {
		t.Errorf("TestScanNodes(%q) got Val diff (-want, +got):\n%s", tc.name, diff)
	}
	if gotErr := gotNodeErr != nil; gotErr != tc.wantNodeErr {
		t.Errorf("TestScanNodes(%q) got node err: %v, want err? %v", tc.name, gotNodeErr, tc.wantNodeErr)
	}
}

func TestTokenizeLit(t *testing.T) {
	for _, tc := range []scanTestCase{{
		name:        "id",
		input:       "foo",
		wantPos:     []Pos{0, 3},
		wantTok:     []lisp.Token{lisp.Id},
		wantText:    []string{"foo"},
		wantNodePos: []Pos{0, 3},
		wantNode:    []lisp.Val{lisp.Lit("foo")},
	}, {
		name:        "space id",
		input:       "  \t\n x",
		wantPos:     []Pos{5, 6},
		wantTok:     []lisp.Token{lisp.Id},
		wantText:    []string{"x"},
		wantNodePos: []Pos{5, 6},
		wantNode:    []lisp.Val{lisp.Lit("x")},
	}, {
		name:        "id id id",
		input:       "a b c",
		wantPos:     []Pos{0, 1, 2, 3, 4, 5},
		wantTok:     []lisp.Token{lisp.Id, lisp.Id, lisp.Id},
		wantText:    []string{"a", "b", "c"},
		wantNodePos: []Pos{0, 1, 2, 3, 4, 5},
		wantNode: []lisp.Val{
			lisp.Lit("a"),
			lisp.Lit("b"),
			lisp.Lit("c"),
		},
	}, {
		name:        "id id id 2",
		input:       "ab cd ef",
		wantPos:     []Pos{0, 2, 3, 5, 6, 8},
		wantTok:     []lisp.Token{lisp.Id, lisp.Id, lisp.Id},
		wantText:    []string{"ab", "cd", "ef"},
		wantNodePos: []Pos{0, 2, 3, 5, 6, 8},
		wantNode: []lisp.Val{
			lisp.Lit("ab"),
			lisp.Lit("cd"),
			lisp.Lit("ef"),
		},
	}, {
		name:        "int",
		input:       "0",
		wantPos:     []Pos{0, 1},
		wantTok:     []lisp.Token{lisp.Id},
		wantText:    []string{"0"},
		wantNodePos: []Pos{0, 1},
		wantNode:    []lisp.Val{lisp.Lit("0")},
	}, {
		name:        "int 2",
		input:       "0 1 2",
		wantPos:     []Pos{0, 1, 2, 3, 4, 5},
		wantTok:     []lisp.Token{lisp.Id, lisp.Id, lisp.Id},
		wantText:    []string{"0", "1", "2"},
		wantNodePos: []Pos{0, 1, 2, 3, 4, 5},
		wantNode: []lisp.Val{
			lisp.Lit("0"),
			lisp.Lit("1"),
			lisp.Lit("2"),
		},
	}, {
		name:        "zero sequence",
		input:       "00000",
		wantPos:     []Pos{0, 5},
		wantTok:     []lisp.Token{lisp.Id},
		wantText:    []string{"00000"},
		wantNodePos: []Pos{0, 5},
		wantNode: []lisp.Val{
			lisp.Lit("00000"),
		},
	}, {
		name:        "token nat",
		input:       "token0",
		wantPos:     []Pos{0, 6},
		wantTok:     []lisp.Token{lisp.Id},
		wantText:    []string{"token0"},
		wantNodePos: []Pos{0, 6},
		wantNode: []lisp.Val{
			lisp.Lit("token0"),
		},
	}, {
		name:        "nat token",
		input:       "0token",
		wantPos:     []Pos{0, 6},
		wantTok:     []lisp.Token{lisp.Id},
		wantText:    []string{"0token"},
		wantNodePos: []Pos{0, 6},
		wantNode: []lisp.Val{
			lisp.Lit("0token"),
		},
	}, {
		name:        "nat in token",
		input:       "tok0tok",
		wantPos:     []Pos{0, 7},
		wantTok:     []lisp.Token{lisp.Id},
		wantText:    []string{"tok0tok"},
		wantNodePos: []Pos{0, 7},
		wantNode: []lisp.Val{
			lisp.Lit("tok0tok"),
		},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			tc.scanTokenTest(t)
		})
	}
}

func TestTokenizeGroup(t *testing.T) {
	for _, tc := range []scanTestCase{{
		name: "empty",
	}, {
		name:  "whitespace",
		input: "  \t\r\n",
	}, {
		name:        "empty group",
		input:       "()",
		wantPos:     []Pos{0, 1, 1, 2},
		wantTok:     []lisp.Token{lisp.LParen, lisp.RParen},
		wantText:    []string{"(", ")"},
		wantNodePos: []Pos{0, 2},
		wantNode:    []lisp.Val{lisp.Group{}},
	}, {
		name:        "nested group",
		input:       "(())",
		wantPos:     []Pos{0, 1, 1, 2, 2, 3, 3, 4},
		wantTok:     []lisp.Token{lisp.LParen, lisp.LParen, lisp.RParen, lisp.RParen},
		wantText:    []string{"(", "(", ")", ")"},
		wantNodePos: []Pos{0, 4},
		wantNode:    []lisp.Val{lisp.Group{lisp.Group{}}},
	}, {
		name:        "group",
		input:       "(abc)",
		wantPos:     []Pos{0, 1, 1, 4, 4, 5},
		wantTok:     []lisp.Token{lisp.LParen, lisp.Id, lisp.RParen},
		wantText:    []string{"(", "abc", ")"},
		wantNodePos: []Pos{0, 5},
		wantNode: []lisp.Val{lisp.Group{
			lisp.Lit("abc"),
		}},
	}, {
		name:        "group 2",
		input:       "(add 1 2)",
		wantPos:     []Pos{0, 1, 1, 4, 5, 6, 7, 8, 8, 9},
		wantTok:     []lisp.Token{lisp.LParen, lisp.Id, lisp.Id, lisp.Id, lisp.RParen},
		wantText:    []string{"(", "add", "1", "2", ")"},
		wantNodePos: []Pos{0, 9},
		wantNode: []lisp.Val{lisp.Group{
			lisp.Lit("add"),
			lisp.Lit("1"),
			lisp.Lit("2"),
		}},
	}, {
		name:        "group 3",
		input:       "(add (sub 3 2) 2)",
		wantPos:     []Pos{0, 1, 1, 4, 5, 6, 6, 9, 10, 11, 12, 13, 13, 14, 15, 16, 16, 17},
		wantTok:     []lisp.Token{lisp.LParen, lisp.Id, lisp.LParen, lisp.Id, lisp.Id, lisp.Id, lisp.RParen, lisp.Id, lisp.RParen},
		wantText:    []string{"(", "add", "(", "sub", "3", "2", ")", "2", ")"},
		wantNodePos: []Pos{0, 17},
		wantNode: []lisp.Val{lisp.Group{
			lisp.Lit("add"),
			lisp.Group{
				lisp.Lit("sub"),
				lisp.Lit("3"),
				lisp.Lit("2"),
			},
			lisp.Lit("2"),
		}},
	}, {
		name:        "group 4",
		input:       "((a))",
		wantPos:     []Pos{0, 1, 1, 2, 2, 3, 3, 4, 4, 5},
		wantTok:     []lisp.Token{lisp.LParen, lisp.LParen, lisp.Id, lisp.RParen, lisp.RParen},
		wantText:    []string{"(", "(", "a", ")", ")"},
		wantNodePos: []Pos{0, 5},
		wantNode: []lisp.Val{lisp.Group{
			lisp.Group{
				lisp.Lit("a"),
			},
		}},
	}, {
		name:        "group 5",
		input:       "(a)(b) (c)",
		wantPos:     []Pos{0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 7, 8, 8, 9, 9, 10},
		wantTok:     []lisp.Token{lisp.LParen, lisp.Id, lisp.RParen, lisp.LParen, lisp.Id, lisp.RParen, lisp.LParen, lisp.Id, lisp.RParen},
		wantText:    []string{"(", "a", ")", "(", "b", ")", "(", "c", ")"},
		wantNodePos: []Pos{0, 3, 3, 6, 7, 10},
		wantNode: []lisp.Val{
			lisp.Group{lisp.Lit("a")},
			lisp.Group{lisp.Lit("b")},
			lisp.Group{lisp.Lit("c")},
		},
	}, {
		name:        "group 6",
		input:       "(div (q x) y)\n",
		wantPos:     []Pos{0, 1, 1, 4, 5, 6, 6, 7, 8, 9, 9, 10, 11, 12, 12, 13},
		wantTok:     []lisp.Token{lisp.LParen, lisp.Id, lisp.LParen, lisp.Id, lisp.Id, lisp.RParen, lisp.Id, lisp.RParen},
		wantText:    []string{"(", "div", "(", "q", "x", ")", "y", ")"},
		wantNodePos: []Pos{0, 13},
		wantNode: []lisp.Val{lisp.Group{
			lisp.Lit("div"),
			lisp.Group{
				lisp.Lit("q"),
				lisp.Lit("x"),
			},
			lisp.Lit("y"),
		}},
	}, {
		name:        "group 7",
		input:       "((a b))\n",
		wantPos:     []Pos{0, 1, 1, 2, 2, 3, 4, 5, 5, 6, 6, 7},
		wantTok:     []lisp.Token{lisp.LParen, lisp.LParen, lisp.Id, lisp.Id, lisp.RParen, lisp.RParen},
		wantText:    []string{"(", "(", "a", "b", ")", ")"},
		wantNodePos: []Pos{0, 7},
		wantNode: []lisp.Val{lisp.Group{lisp.Group{
			lisp.Lit("a"),
			lisp.Lit("b"),
		}}},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			tc.scanTokenTest(t)
		})
	}
}

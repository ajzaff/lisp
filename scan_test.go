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
	sc := NewTokenScanner(strings.NewReader(tc.input))
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
	sc.Init(strings.NewReader(tc.input))
	s := NewNodeScanner(sc)
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
		wantNodes: []Node{&LitNode{LitPos: 0, Lit: IdLit("foo"), EndPos: 3}},
	}, {
		name:      "symbol",
		input:     `.+`,
		wantPos:   []Pos{0, 2},
		wantTok:   []Token{Id},
		wantText:  []string{".+"},
		wantNodes: []Node{&LitNode{LitPos: 0, Lit: IdLit(".+"), EndPos: 2}},
	}, {
		name:     "id symbol id",
		input:    "foo-bar",
		wantPos:  []Pos{0, 3, 3, 4, 4, 7},
		wantTok:  []Token{Id, Id, Id},
		wantText: []string{"foo", "-", "bar"},
		wantNodes: []Node{
			&LitNode{LitPos: 0, Lit: IdLit("foo"), EndPos: 3},
			&LitNode{LitPos: 3, Lit: IdLit("-"), EndPos: 4},
			&LitNode{LitPos: 4, Lit: IdLit("bar"), EndPos: 7},
		},
	}, {
		name:      "space id",
		input:     "  \t\n x",
		wantPos:   []Pos{5, 6},
		wantTok:   []Token{Id},
		wantText:  []string{"x"},
		wantNodes: []Node{&LitNode{LitPos: 5, Lit: IdLit("x"), EndPos: 6}},
	}, {
		name:     "id id id",
		input:    "a b c",
		wantPos:  []Pos{0, 1, 2, 3, 4, 5},
		wantTok:  []Token{Id, Id, Id},
		wantText: []string{"a", "b", "c"},
		wantNodes: []Node{
			&LitNode{LitPos: 0, Lit: IdLit("a"), EndPos: 1},
			&LitNode{LitPos: 2, Lit: IdLit("b"), EndPos: 3},
			&LitNode{LitPos: 4, Lit: IdLit("c"), EndPos: 5},
		},
	}, {
		name:      "symbol 2",
		input:     ".",
		wantPos:   []Pos{0, 1},
		wantTok:   []Token{Id},
		wantText:  []string{"."},
		wantNodes: []Node{&LitNode{LitPos: 0, Lit: IdLit("."), EndPos: 1}},
	}, {
		name:     "id id int",
		input:    `a.0`,
		wantPos:  []Pos{0, 1, 1, 2, 2, 3},
		wantTok:  []Token{Id, Id, Number},
		wantText: []string{"a", ".", "0"},
		wantNodes: []Node{
			&LitNode{LitPos: 0, Lit: IdLit("a"), EndPos: 1},
			&LitNode{LitPos: 1, Lit: IdLit("."), EndPos: 2},
			&LitNode{LitPos: 2, Lit: NumberLit("0"), EndPos: 3},
		},
	}, {
		name:     "id symbol",
		input:    "a...",
		wantPos:  []Pos{0, 1, 1, 4},
		wantTok:  []Token{Id, Id},
		wantText: []string{"a", "..."},
		wantNodes: []Node{
			&LitNode{LitPos: 0, Lit: IdLit("a"), EndPos: 1},
			&LitNode{LitPos: 1, Lit: IdLit("..."), EndPos: 4},
		},
	}, {
		name:     "symbol id",
		input:    "...a",
		wantPos:  []Pos{0, 3, 3, 4},
		wantTok:  []Token{Id, Id},
		wantText: []string{"...", "a"},
		wantNodes: []Node{
			&LitNode{LitPos: 0, Lit: IdLit("..."), EndPos: 3},
			&LitNode{LitPos: 3, Lit: IdLit("a"), EndPos: 4},
		},
	}, {
		name:     "id string",
		input:    `a"abc"`,
		wantPos:  []Pos{0, 1, 1, 6},
		wantTok:  []Token{Id, String},
		wantText: []string{"a", `"abc"`},
		wantNodes: []Node{
			&LitNode{LitPos: 0, Lit: IdLit("a"), EndPos: 1},
			&LitNode{LitPos: 1, Lit: StringLit(`"abc"`), EndPos: 6},
		},
	}, {
		name:     "id id id 2",
		input:    "ab cd ef",
		wantPos:  []Pos{0, 2, 3, 5, 6, 8},
		wantTok:  []Token{Id, Id, Id},
		wantText: []string{"ab", "cd", "ef"},
		wantNodes: []Node{
			&LitNode{LitPos: 0, Lit: IdLit("ab"), EndPos: 2},
			&LitNode{LitPos: 3, Lit: IdLit("cd"), EndPos: 5},
			&LitNode{LitPos: 6, Lit: IdLit("ef"), EndPos: 8},
		},
	}, {
		name:      "int",
		input:     "0",
		wantPos:   []Pos{0, 1},
		wantTok:   []Token{Number},
		wantText:  []string{"0"},
		wantNodes: []Node{&LitNode{LitPos: 0, Lit: NumberLit("0"), EndPos: 1}},
	}, {
		name:     "int 2",
		input:    "0 1 2",
		wantPos:  []Pos{0, 1, 2, 3, 4, 5},
		wantTok:  []Token{Number, Number, Number},
		wantText: []string{"0", "1", "2"},
		wantNodes: []Node{
			&LitNode{LitPos: 0, Lit: NumberLit("0"), EndPos: 1},
			&LitNode{LitPos: 2, Lit: NumberLit("1"), EndPos: 3},
			&LitNode{LitPos: 4, Lit: NumberLit("2"), EndPos: 5},
		},
	}, {
		name:      "float",
		input:     "1.0",
		wantPos:   []Pos{0, 3},
		wantTok:   []Token{Number},
		wantText:  []string{"1.0"},
		wantNodes: []Node{&LitNode{LitPos: 0, Lit: NumberLit("1.0"), EndPos: 3}},
	}, {
		name:      "float 2",
		input:     "1.",
		wantPos:   []Pos{0, 2},
		wantTok:   []Token{Number},
		wantText:  []string{"1."},
		wantNodes: []Node{&LitNode{LitPos: 0, Lit: NumberLit("1."), EndPos: 2}},
	}, {
		name:      "float 3",
		input:     "0.",
		wantPos:   []Pos{0, 2},
		wantTok:   []Token{Number},
		wantText:  []string{"0."},
		wantNodes: []Node{&LitNode{LitPos: 0, Lit: NumberLit("0."), EndPos: 2}},
	}, {
		name:     "float 4",
		input:    "1. 2. 3.",
		wantPos:  []Pos{0, 2, 3, 5, 6, 8},
		wantTok:  []Token{Number, Number, Number},
		wantText: []string{"1.", "2.", "3."},
		wantNodes: []Node{
			&LitNode{LitPos: 0, Lit: NumberLit("1."), EndPos: 2},
			&LitNode{LitPos: 3, Lit: NumberLit("2."), EndPos: 5},
			&LitNode{LitPos: 6, Lit: NumberLit("3."), EndPos: 8},
		},
	}, {
		name:     "float 4_2",
		input:    "0.1 0.2 0.3",
		wantPos:  []Pos{0, 3, 4, 7, 8, 11},
		wantTok:  []Token{Number, Number, Number},
		wantText: []string{"0.1", "0.2", "0.3"},
		wantNodes: []Node{
			&LitNode{LitPos: 0, Lit: NumberLit("0.1"), EndPos: 3},
			&LitNode{LitPos: 4, Lit: NumberLit("0.2"), EndPos: 7},
			&LitNode{LitPos: 8, Lit: NumberLit("0.3"), EndPos: 11},
		},
	}, {
		name:      "string",
		input:     `"a"`,
		wantPos:   []Pos{0, 3},
		wantTok:   []Token{String},
		wantText:  []string{`"a"`},
		wantNodes: []Node{&LitNode{LitPos: 0, Lit: StringLit(`"a"`), EndPos: 3}},
	}, {
		name:      "string 2",
		input:     `"a b c"`,
		wantPos:   []Pos{0, 7},
		wantTok:   []Token{String},
		wantText:  []string{`"a b c"`},
		wantNodes: []Node{&LitNode{LitPos: 0, Lit: StringLit(`"a b c"`), EndPos: 7}},
	}, {
		name:     "string 3",
		input:    `"a" "b" "c"`,
		wantPos:  []Pos{0, 3, 4, 7, 8, 11},
		wantTok:  []Token{String, String, String},
		wantText: []string{`"a"`, `"b"`, `"c"`},
		wantNodes: []Node{
			&LitNode{LitPos: 0, Lit: StringLit(`"a"`), EndPos: 3},
			&LitNode{LitPos: 4, Lit: StringLit(`"b"`), EndPos: 7},
			&LitNode{LitPos: 8, Lit: StringLit(`"c"`), EndPos: 11},
		},
	}, {
		name: "string (multiline)",
		input: `"
"`,
		wantPos: []Pos{0, 3},
		wantTok: []Token{String},
		wantText: []string{`"
"`},
		wantNodes: []Node{&LitNode{LitPos: 0, Lit: StringLit(`"
"`), EndPos: 3}},
	}, {
		name:      "string (double escape)",
		input:     `"\\"`,
		wantPos:   []Pos{0, 4},
		wantTok:   []Token{String},
		wantText:  []string{`"\\"`},
		wantNodes: []Node{&LitNode{LitPos: 0, Lit: StringLit(`"\\"`), EndPos: 4}},
	}, {
		name:      "byte lit",
		input:     `"abc\x00\x11\xff"`,
		wantPos:   []Pos{0, 17},
		wantTok:   []Token{String},
		wantText:  []string{`"abc\x00\x11\xff"`},
		wantNodes: []Node{&LitNode{LitPos: 0, Lit: StringLit(`"abc\x00\x11\xff"`), EndPos: 17}},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			tc.scanTokenTest(t)
		})
	}
}

func TestTokenizeLitErrors(t *testing.T) {
	for _, tc := range []scanTestCase{{
		name:         "byte lit EOF",
		input:        `"abc\x00\x1"`,
		wantTokenErr: true,
		wantNodeErr:  true,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			tc.scanTokenTest(t)
		})
	}
}

func TestTokenizeExpr(t *testing.T) {
	for _, tc := range []scanTestCase{{
		name: "empty",
	}, {
		name:  "whitespace",
		input: "  \t\r\n",
	}, {
		name:     "expr",
		input:    "(abc)",
		wantPos:  []Pos{0, 1, 1, 4, 4, 5},
		wantTok:  []Token{LParen, Id, RParen},
		wantText: []string{"(", "abc", ")"},
		wantNodes: []Node{&ExprNode{LParen: 0, Expr: Expr{
			&LitNode{LitPos: 1, Lit: IdLit("abc"), EndPos: 4},
		}, RParen: 4}},
	}, {
		name:     "expr symbol",
		input:    "(.)",
		wantPos:  []Pos{0, 1, 1, 2, 2, 3},
		wantTok:  []Token{LParen, Id, RParen},
		wantText: []string{"(", ".", ")"},
		wantNodes: []Node{&ExprNode{LParen: 0, Expr: Expr{
			&LitNode{LitPos: 1, Lit: IdLit("."), EndPos: 2},
		}, RParen: 2}},
	}, {
		name:     "expr 2",
		input:    "(add 1 2)",
		wantPos:  []Pos{0, 1, 1, 4, 5, 6, 7, 8, 8, 9},
		wantTok:  []Token{LParen, Id, Number, Number, RParen},
		wantText: []string{"(", "add", "1", "2", ")"},
		wantNodes: []Node{&ExprNode{LParen: 0, Expr: Expr{
			&LitNode{LitPos: 1, Lit: IdLit("add"), EndPos: 4},
			&LitNode{LitPos: 5, Lit: NumberLit("1"), EndPos: 6},
			&LitNode{LitPos: 7, Lit: NumberLit("2"), EndPos: 8},
		}, RParen: 8}},
	}, {
		name:     "expr 3",
		input:    "(add (sub 3 2) 2)",
		wantPos:  []Pos{0, 1, 1, 4, 5, 6, 6, 9, 10, 11, 12, 13, 13, 14, 15, 16, 16, 17},
		wantTok:  []Token{LParen, Id, LParen, Id, Number, Number, RParen, Number, RParen},
		wantText: []string{"(", "add", "(", "sub", "3", "2", ")", "2", ")"},
		wantNodes: []Node{&ExprNode{LParen: 0, Expr: Expr{
			&LitNode{LitPos: 1, Lit: IdLit("add"), EndPos: 4},
			&ExprNode{LParen: 5, Expr: Expr{
				&LitNode{LitPos: 6, Lit: IdLit("sub"), EndPos: 9},
				&LitNode{LitPos: 10, Lit: NumberLit("3"), EndPos: 11},
				&LitNode{LitPos: 12, Lit: NumberLit("2"), EndPos: 13},
			}, RParen: 13},
			&LitNode{LitPos: 15, Lit: NumberLit("2"), EndPos: 16},
		}, RParen: 16}},
	}, {
		name:     "expr 4",
		input:    "((a))",
		wantPos:  []Pos{0, 1, 1, 2, 2, 3, 3, 4, 4, 5},
		wantTok:  []Token{LParen, LParen, Id, RParen, RParen},
		wantText: []string{"(", "(", "a", ")", ")"},
		wantNodes: []Node{&ExprNode{LParen: 0, Expr: Expr{
			&ExprNode{LParen: 1, Expr: Expr{
				&LitNode{LitPos: 2, Lit: IdLit("a"), EndPos: 3},
			}, RParen: 3},
		}, RParen: 4}},
	}, {
		name:     "expr 5",
		input:    "(a)(b) (c)",
		wantPos:  []Pos{0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 7, 8, 8, 9, 9, 10},
		wantTok:  []Token{LParen, Id, RParen, LParen, Id, RParen, LParen, Id, RParen},
		wantText: []string{"(", "a", ")", "(", "b", ")", "(", "c", ")"},
		wantNodes: []Node{
			&ExprNode{LParen: 0, Expr: Expr{&LitNode{LitPos: 1, Lit: IdLit("a"), EndPos: 2}}, RParen: 2},
			&ExprNode{LParen: 3, Expr: Expr{&LitNode{LitPos: 4, Lit: IdLit("b"), EndPos: 5}}, RParen: 5},
			&ExprNode{LParen: 7, Expr: Expr{&LitNode{LitPos: 8, Lit: IdLit("c"), EndPos: 9}}, RParen: 9},
		},
	}, {
		name:     "expr 6",
		input:    "(?x/y)\n",
		wantPos:  []Pos{0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6},
		wantTok:  []Token{LParen, Id, Id, Id, Id, RParen},
		wantText: []string{"(", "?", "x", "/", "y", ")"},
		wantNodes: []Node{&ExprNode{LParen: 0, Expr: Expr{
			&LitNode{LitPos: 1, Lit: IdLit("?"), EndPos: 2},
			&LitNode{LitPos: 2, Lit: IdLit("x"), EndPos: 3},
			&LitNode{LitPos: 3, Lit: IdLit("/"), EndPos: 4},
			&LitNode{LitPos: 4, Lit: IdLit("y"), EndPos: 5},
		}, RParen: 5}},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			tc.scanTokenTest(t)
		})
	}
}

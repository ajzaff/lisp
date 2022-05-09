package lisp

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type scanTestCase struct {
	name     string
	input    string
	wantPos  []Pos
	wantTok  []Token
	wantText []string
	wantErr  bool
}

func scanTest(t *testing.T, tc scanTestCase) {
	if !tc.wantErr && len(tc.wantPos)%2 != 0 {
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
	gotErr := sc.Err()
	if diff := cmp.Diff(tc.wantPos, gotPos); diff != "" {
		t.Errorf("Tokenize(%q) got diff (-want, +got):\n%s", tc.name, diff)
	}
	if diff := cmp.Diff(tc.wantTok, gotTok); diff != "" {
		t.Errorf("Tokenize(%q) got token diff (-want, +got):\n%s", tc.name, diff)
	}
	if diff := cmp.Diff(tc.wantText, gotText); diff != "" {
		t.Errorf("Tokenize(%q) got text diff (-want, +got):\n%s", tc.name, diff)
	}
	if (gotErr != nil) != tc.wantErr {
		t.Errorf("Tokenize(%q) got err: %v, want err? %v", tc.name, gotErr, tc.wantErr)
	}
}

func TestTokenizeLit(t *testing.T) {
	for _, tc := range []scanTestCase{{
		name:     "id",
		input:    "foo",
		wantPos:  []Pos{0, 3},
		wantTok:  []Token{Id},
		wantText: []string{"foo"},
	}, {
		name:     "symbol",
		input:    `.+`,
		wantPos:  []Pos{0, 2},
		wantTok:  []Token{Id},
		wantText: []string{".+"},
	}, {
		name:     "id symbol id",
		input:    "foo-bar",
		wantPos:  []Pos{0, 3, 3, 4, 4, 7},
		wantTok:  []Token{Id, Id, Id},
		wantText: []string{"foo", "-", "bar"},
	}, {
		name:     "space id",
		input:    "  \t\n x",
		wantPos:  []Pos{5, 6},
		wantTok:  []Token{Id},
		wantText: []string{"x"},
	}, {
		name:     "id id id",
		input:    "a b c",
		wantPos:  []Pos{0, 1, 2, 3, 4, 5},
		wantTok:  []Token{Id, Id, Id},
		wantText: []string{"a", "b", "c"},
	}, {
		name:     "symbol 2",
		input:    ".",
		wantPos:  []Pos{0, 1},
		wantTok:  []Token{Id},
		wantText: []string{"."},
	}, {
		name:     "id id int",
		input:    `a.0`,
		wantPos:  []Pos{0, 1, 1, 2, 2, 3},
		wantTok:  []Token{Id, Id, Int},
		wantText: []string{"a", ".", "0"},
	}, {
		name:     "id symbol",
		input:    "a...",
		wantPos:  []Pos{0, 1, 1, 4},
		wantTok:  []Token{Id, Id},
		wantText: []string{"a", "..."},
	}, {
		name:     "symbol id",
		input:    "...a",
		wantPos:  []Pos{0, 3, 3, 4},
		wantTok:  []Token{Id, Id},
		wantText: []string{"...", "a"},
	}, {
		name:     "id string",
		input:    `a"abc"`,
		wantPos:  []Pos{0, 1, 1, 6},
		wantTok:  []Token{Id, String},
		wantText: []string{"a", `"abc"`},
	}, {
		name:     "id id id 2",
		input:    "ab cd ef",
		wantPos:  []Pos{0, 2, 3, 5, 6, 8},
		wantTok:  []Token{Id, Id, Id},
		wantText: []string{"ab", "cd", "ef"},
	}, {
		name:     "int",
		input:    "0",
		wantPos:  []Pos{0, 1},
		wantTok:  []Token{Int},
		wantText: []string{"0"},
	}, {
		name:     "int 2",
		input:    "0 1 2",
		wantPos:  []Pos{0, 1, 2, 3, 4, 5},
		wantTok:  []Token{Int, Int, Int},
		wantText: []string{"0", "1", "2"},
	}, {
		name:     "float",
		input:    "1.0",
		wantPos:  []Pos{0, 3},
		wantTok:  []Token{Float},
		wantText: []string{"1.0"},
	}, {
		name:     "float 2",
		input:    "1.",
		wantPos:  []Pos{0, 2},
		wantTok:  []Token{Float},
		wantText: []string{"1."},
	}, {
		name:     "float 3",
		input:    "0.",
		wantPos:  []Pos{0, 2},
		wantTok:  []Token{Float},
		wantText: []string{"0."},
	}, {
		name:     "float 4",
		input:    "1. 2. 3.",
		wantPos:  []Pos{0, 2, 3, 5, 6, 8},
		wantTok:  []Token{Float, Float, Float},
		wantText: []string{"1.", "2.", "3."},
	}, {
		name:     "float 4_2",
		input:    "0.1 0.2 0.3",
		wantPos:  []Pos{0, 3, 4, 7, 8, 11},
		wantTok:  []Token{Float, Float, Float},
		wantText: []string{"0.1", "0.2", "0.3"},
	}, {
		name:     "string",
		input:    `"a"`,
		wantPos:  []Pos{0, 3},
		wantTok:  []Token{String},
		wantText: []string{`"a"`},
	}, {
		name:     "string 2",
		input:    `"a b c"`,
		wantPos:  []Pos{0, 7},
		wantTok:  []Token{String},
		wantText: []string{`"a b c"`},
	}, {
		name:     "string 3",
		input:    `"a" "b" "c"`,
		wantPos:  []Pos{0, 3, 4, 7, 8, 11},
		wantTok:  []Token{String, String, String},
		wantText: []string{`"a"`, `"b"`, `"c"`},
	}, {
		name: "string (multiline)",
		input: `"
"`,
		wantPos: []Pos{0, 3},
		wantTok: []Token{String},
		wantText: []string{`"
"`},
	}, {
		name:     "string (double escape)",
		input:    `"\\"`,
		wantPos:  []Pos{0, 4},
		wantTok:  []Token{String},
		wantText: []string{`"\\"`},
	}, {
		name:     "byte lit",
		input:    `"abc\x00\x11\xff"`,
		wantPos:  []Pos{0, 17},
		wantTok:  []Token{String},
		wantText: []string{`"abc\x00\x11\xff"`},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			scanTest(t, tc)
		})
	}
}

func TestTokenizeLitErrors(t *testing.T) {
	for _, tc := range []scanTestCase{{
		name:    "byte lit EOF",
		input:   `"abc\x00\x1"`,
		wantErr: true,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			scanTest(t, tc)
		})
	}
}

func TestTokenizeExpr(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    string
		wantPos  []Pos
		wantTok  []Token
		wantText []string
		wantErr  bool
	}{{
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
	}, {
		name:     "expr symbol",
		input:    "(.)",
		wantPos:  []Pos{0, 1, 1, 2, 2, 3},
		wantTok:  []Token{LParen, Id, RParen},
		wantText: []string{"(", ".", ")"},
	}, {
		name:     "expr 2",
		input:    "(add 1 2)",
		wantPos:  []Pos{0, 1, 1, 4, 5, 6, 7, 8, 8, 9},
		wantTok:  []Token{LParen, Id, Int, Int, RParen},
		wantText: []string{"(", "add", "1", "2", ")"},
	}, {
		name:     "expr 3",
		input:    "(add (sub 3 2) 2)",
		wantPos:  []Pos{0, 1, 1, 4, 5, 6, 6, 9, 10, 11, 12, 13, 13, 14, 15, 16, 16, 17},
		wantTok:  []Token{LParen, Id, LParen, Id, Int, Int, RParen, Int, RParen},
		wantText: []string{"(", "add", "(", "sub", "3", "2", ")", "2", ")"},
	}, {
		name:     "expr 4",
		input:    "((a))",
		wantPos:  []Pos{0, 1, 1, 2, 2, 3, 3, 4, 4, 5},
		wantTok:  []Token{LParen, LParen, Id, RParen, RParen},
		wantText: []string{"(", "(", "a", ")", ")"},
	}, {
		name:     "expr 5",
		input:    "(a)(b) (c)",
		wantPos:  []Pos{0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 7, 8, 8, 9, 9, 10},
		wantTok:  []Token{LParen, Id, RParen, LParen, Id, RParen, LParen, Id, RParen},
		wantText: []string{"(", "a", ")", "(", "b", ")", "(", "c", ")"},
	}, {
		name:     "expr 6",
		input:    "(?x/y)\n",
		wantPos:  []Pos{0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6},
		wantTok:  []Token{LParen, Id, Id, Id, Id, RParen},
		wantText: []string{"(", "?", "x", "/", "y", ")"},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			scanTest(t, tc)
		})
	}
}

package innit

import "testing"

func checkTokens(got, want []Pos) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func TestTokenizeBasic(t *testing.T) {
	for _, tc := range []struct {
		name    string
		src     []byte
		want    []Pos
		wantErr bool
	}{{
		name: "id",
		src:  []byte("foo"),
		want: []Pos{0, 3},
	}, {
		name: "int",
		src:  []byte("0"),
		want: []Pos{0, 1},
	}, {
		name: "float",
		src:  []byte("1.0"),
		want: []Pos{0, 3},
	}, {
		name: "expr",
		src:  []byte("(abc)"),
		want: []Pos{0, 1, 1, 4, 4, 5},
	}, {
		name: "compound",
		src:  []byte("(add 1 2)"),
		want: []Pos{0, 1, 1, 4, 5, 6, 7, 8, 8, 9},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			got, gotErr := Tokenize(tc.src)
			if !checkTokens(got, tc.want) {
				t.Errorf("Tokenize() got %v, want %v", got, tc.want)
			}
			if (gotErr == nil) == tc.wantErr {
				t.Errorf("Tokenize() got err %v, want err? %v", gotErr, tc.wantErr)
			}
		})
	}
}

func TestTokenizeError(t *testing.T) {
	for _, tc := range []struct {
		name    string
		src     []byte
		want    []Pos
		wantErr bool
	}{{
		name:    "missing end",
		src:     []byte("("),
		wantErr: true,
	}, {
		name:    "unexpected end",
		src:     []byte(")"),
		wantErr: true,
	}, {
		name:    "bad id",
		src:     []byte("a-"),
		wantErr: true,
	}, {
		name:    "bad id 2",
		src:     []byte("-a"),
		wantErr: true,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			got, gotErr := Tokenize(tc.src)
			if !checkTokens(got, tc.want) {
				t.Errorf("Tokenize() got %v, want %v", got, tc.want)
			}
			if (gotErr == nil) == tc.wantErr {
				t.Errorf("Tokenize() got err %v, want err? %v", gotErr, tc.wantErr)
			}
		})
	}
}

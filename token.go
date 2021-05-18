package innit

import (
	"fmt"
	"unicode"
)

type Token int

const (
	Invalid Token = iota

	Id     // main
	Int    // 12345
	Float  // 123.45
	String // "abc"

	LParen // (
	RParen // )
)

type Pos int

const NoPos = -1

type TokenError struct {
	Line, Col int
	Pos       Pos
}

func (t *TokenError) Error() string {
	return fmt.Sprintf("tokenize error: at line %d: col %d", t.Line, t.Col)
}

type tokenState int

const (
	tokenStart tokenState = iota
	tokenArgs
	tokenIdent
	tokenInt
	tokenFloat
	tokenString
	tokenEscape
	tokenByte
	tokenUnicode
)

func Tokenize(src []byte) ([]Pos, error) {
	var pos []Pos
	var (
		line = 1
		col  = 1
	)
	state := tokenStart
	for i, r := range string(src) {
		switch {
		case unicode.IsSpace(r):
			switch state {
			case tokenStart, tokenArgs, tokenString:
			case tokenIdent, tokenInt, tokenFloat:
				pos = append(pos, Pos(i))
				state = tokenArgs
			default:
				return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
			}
			if r == '\n' {
				line++
				col = 1
			}
		case r == '"':
			switch state {
			case tokenStart, tokenArgs:
				pos = append(pos, Pos(i))
				state = tokenString
			case tokenString:
				pos = append(pos, Pos(i+1))
				state = tokenArgs
			case tokenEscape:
				state = tokenString
			default:
			}
		case r == '\\':
			switch state {
			case tokenString:
				state = tokenEscape
			case tokenEscape:
				state = tokenString
			default:
				return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
			}
		case r == '.':
			switch state {
			case tokenStart:
				pos = append(pos, Pos(i))
				state = tokenIdent
			case tokenArgs:
				pos = append(pos, Pos(i))
				state = tokenFloat
			case tokenInt:
				state = tokenFloat
			default:
				return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
			}
		case r == '(' || r == ')':
			switch state {
			case tokenStart, tokenArgs:
			case tokenIdent, tokenInt, tokenFloat:
				pos = append(pos, Pos(i))
			default:
				return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
			}
			pos = append(pos, Pos(i), Pos(i+1))
			state = tokenStart
		case unicode.IsNumber(r):
			switch state {
			case tokenStart, tokenArgs:
				state = tokenInt
				pos = append(pos, Pos(i))
			case tokenIdent, tokenInt, tokenFloat, tokenByte, tokenUnicode:
			default:
				return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
			}
		case unicode.IsPrint(r):
			switch state {
			case tokenStart, tokenArgs:
				state = tokenIdent
				pos = append(pos, Pos(i))
			case tokenIdent, tokenString:
			case tokenEscape:
				switch r {
				case 't', 'n':
					state = tokenString
				case 'u':
					state = tokenUnicode
				case 'x':
					state = tokenByte
				default:
					return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
				}
			case tokenByte, tokenUnicode:
				switch {
				case r >= 'A' && r <= 'F' || r >= 'a' && r <= 'f':
				default:
					return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
				}
			default:
				return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
			}
		}
		col++
	}
	switch state {
	case tokenIdent, tokenInt, tokenFloat:
		pos = append(pos, Pos(len(src)))
	}
	return pos, nil
}

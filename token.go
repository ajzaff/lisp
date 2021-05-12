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
	Char   // 'a'
	String // "abc"

	LParen // (
	RParen // )
)

type TokenError struct {
	Line, Col int
	Pos       Pos
}

func (t *TokenError) Error() string {
	return fmt.Sprintf("tokenize error: at line %d: col %d", t.Line, t.Col)
}

type state int

const (
	stateBegin state = iota
	stateIdent
	stateIdentNoDash
	stateInt
	stateFloat
)

func Tokenize(src []byte) ([]Pos, error) {
	var pos []Pos
	var (
		indent = 0
		line   = 1
		col    = 1
	)
	state := stateBegin
	for i, r := range string(src) {
		if unicode.IsSpace(r) {
			switch state {
			case stateBegin, stateIdent, stateInt, stateFloat:
				pos = append(pos, Pos(i))
				state = stateBegin
			default:
				return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
			}
			if r == '\n' {
				line++
				col = 1
			}
		} else if r == '-' || r == '_' {
			switch state {
			case stateIdent:
				state = stateIdentNoDash
			default:
				return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
			}
		} else if r == '.' {
			switch state {
			case stateInt:
				state = stateFloat
			default:
				return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
			}
		} else if r == '(' || r == ')' {
			switch state {
			case stateBegin:
			case stateIdent, stateInt, stateFloat:
				pos = append(pos, Pos(i))
			default:
				return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
			}
			pos = append(pos, Pos(i), Pos(i+1))
			state = stateBegin
			if r == '(' {
				indent++
			} else if r == ')' {
				if indent--; indent < 0 {
					return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
				}
			}
		} else if unicode.IsLetter(r) {
			switch state {
			case stateBegin:
				state = stateIdent
				pos = append(pos, Pos(i))
			case stateIdent:
			case stateIdentNoDash:
				state = stateIdent
			default:
				return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
			}
		} else if unicode.IsNumber(r) {
			switch state {
			case stateBegin:
				state = stateInt
				pos = append(pos, Pos(i))
			case stateIdentNoDash:
				state = stateIdent
			case stateIdent, stateInt, stateFloat:
			default:
				return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
			}
		}
		col++
	}
	switch state {
	case stateIdentNoDash:
		return nil, &TokenError{Line: line, Col: col, Pos: Pos(len(src))}
	case stateIdent, stateInt, stateFloat:
		pos = append(pos, Pos(len(src)))
	}
	if indent != 0 {
		return nil, &TokenError{Line: line, Col: col, Pos: Pos(len(src))}
	}
	return pos, nil
}

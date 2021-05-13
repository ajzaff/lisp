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

type TokenError struct {
	Line, Col int
	Pos       Pos
}

func (t *TokenError) Error() string {
	return fmt.Sprintf("tokenize error: at line %d: col %d", t.Line, t.Col)
}

type state int

const (
	stateStart state = iota
	stateArgs
	stateIdent
	stateInt
	stateFloat
	stateString
	stateEscape
	stateByte
	stateUnicode
)

func Tokenize(src []byte) ([]Pos, error) {
	var pos []Pos
	var (
		line = 1
		col  = 1
	)
	state := stateStart
	for i, r := range string(src) {
		if unicode.IsSpace(r) {
			switch state {
			case stateStart, stateArgs, stateString:
			case stateIdent, stateInt, stateFloat:
				pos = append(pos, Pos(i))
				state = stateArgs
			default:
				return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
			}
			if r == '\n' {
				line++
				col = 1
			}
		} else if r == '"' {
			switch state {
			case stateStart, stateArgs:
				pos = append(pos, Pos(i))
				state = stateString
			case stateString:
				pos = append(pos, Pos(i+1))
				state = stateArgs
			case stateEscape:
				state = stateString
			default:
			}
		} else if r == '\\' {
			switch state {
			case stateString:
				state = stateEscape
			case stateEscape:
				state = stateString
			default:
				return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
			}
		} else if r == '.' {
			switch state {
			case stateStart:
				pos = append(pos, Pos(i))
				state = stateIdent
			case stateArgs:
				pos = append(pos, Pos(i))
				state = stateFloat
			case stateInt:
				state = stateFloat
			default:
				return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
			}
		} else if r == '(' || r == ')' {
			switch state {
			case stateStart, stateArgs:
			case stateIdent, stateInt, stateFloat:
				pos = append(pos, Pos(i))
			default:
				return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
			}
			pos = append(pos, Pos(i), Pos(i+1))
			state = stateStart
		} else if unicode.IsNumber(r) {
			switch state {
			case stateStart, stateArgs:
				state = stateInt
				pos = append(pos, Pos(i))
			case stateIdent, stateInt, stateFloat, stateByte, stateUnicode:
			default:
				return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
			}
		} else if unicode.IsPrint(r) {
			switch state {
			case stateStart, stateArgs:
				state = stateIdent
				pos = append(pos, Pos(i))
			case stateIdent, stateString:
			case stateEscape:
				switch r {
				case 't', 'n':
					state = stateString
				case 'u':
					state = stateUnicode
				case 'x':
					state = stateByte
				default:
					return nil, &TokenError{Line: line, Col: col, Pos: Pos(i)}
				}
			case stateByte, stateUnicode:
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
	case stateIdent, stateInt, stateFloat:
		pos = append(pos, Pos(len(src)))
	}
	return pos, nil
}

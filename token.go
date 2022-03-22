package innit

import (
	"errors"
	"fmt"
	"unicode"
	"unicode/utf8"
)

// Token is an enumeration which specifies a kind of AST token.
type Token int

const (
	Invalid Token = iota

	Id     // main
	Op     // + - / ++
	Int    // 12345
	Float  // 123.45
	String // "abc"

	LParen // (
	RParen // )
)

// Tokenize takes source code and returns a set of token positions.
func Tokenize(src string) ([]Pos, error) {
	s := &tokenState{src: []byte(src), line: 1, col: 1}
	for fn := tokenStart(s); fn != nil; fn = fn() {
	}
	return s.pos, s.err
}

// TokenError implements an error at a specified line and column.
type TokenError struct {
	Line, Col int
	Pos       Pos
}

func (t *TokenError) Error() string {
	return fmt.Sprintf("at line %d: col %d", t.Line, t.Col)
}

func isOp(r rune) bool {
	return unicode.IsPunct(r) || unicode.IsSymbol(r)
}

type tokenState struct {
	src []byte

	pos []Pos
	err error

	line, col int
	last, p   Pos
}

func (s *tokenState) markLast() {
	s.pos = append(s.pos, s.last)
}

func (s *tokenState) mark() {
	s.pos = append(s.pos, s.p)
}

func (s *tokenState) decodeNext() (rune, int) {
	r, size := utf8.DecodeRune(s.src[s.p:])
	s.last = s.p
	s.p += Pos(size)
	if r != '\r' {
		s.col++
	}
	if r == '\n' {
		s.line++
		s.col = 1
	}
	return r, size
}

func (s *tokenState) next() (r rune, ok bool) {
	if r, size := s.decodeNext(); size != 0 {
		return r, true
	}
	return 0, false
}

func (s *tokenState) skipSpace() {
	for {
		r, _ := utf8.DecodeRune(s.src[s.p:])
		if !unicode.IsSpace(r) {
			break
		}
		s.decodeNext()
	}
}

var (
	errRune = errors.New("unexpected rune")
	errEOF  = errors.New("unexpected EOF")
)

func (s *tokenState) setErr(cause error) {
	s.err = fmt.Errorf("%v: %w", cause, &TokenError{s.line, s.col, s.p})
}

type tokenFunc func() tokenFunc

func tokenStart(s *tokenState) tokenFunc {
	s.skipSpace()
	return func() tokenFunc {
		switch r, ok := s.next(); {
		case r == '"':
			s.markLast()
			return tokenString(s)
		case r == '.':
			s.markLast()
			return tokenFloatOrOp(s)
		case r == '(', r == ')':
			s.markLast()
			s.mark()
			return tokenStart(s)
		case unicode.IsNumber(r):
			s.markLast()
			return tokenInt(s)
		case isOp(r): // op
			s.markLast()
			return tokenOp(s)
		case unicode.IsGraphic(r): // id
			s.markLast()
			return tokenId(s)
		case !ok: // EOF
			return nil
		default:
			s.setErr(fmt.Errorf("%v: %v", errRune, r))
			return nil
		}
	}
}

func tokenString(s *tokenState) tokenFunc {
	return func() tokenFunc {
		switch r, ok := s.next(); {
		case r == '"':
			s.mark()
			return tokenStart(s)
		case r == '\\':
			return tokenEscape(s)
		case !ok:
			s.setErr(errEOF)
			return nil
		default:
			return tokenString(s)
		}
	}
}

func tokenEscape(s *tokenState) tokenFunc {
	return func() tokenFunc {
		switch r, ok := s.next(); r {
		case 'n', 't':
			return tokenString(s)
		case 'x', 'u':
			return tokenCharLit(s)
		default:
			if !ok {
				s.setErr(errEOF)
				return nil
			}
			s.setErr(errRune)
			return nil
		}
	}
}

func tokenCharLit(s *tokenState) tokenFunc {
	return func() tokenFunc {
		switch r, ok := s.next(); {
		case r >= '0' && r < '9', r >= 'a' && r <= 'f', r >= 'A' && r <= 'F':
			return tokenCharLit(s)
		default:
			if !ok {
				s.setErr(errEOF)
				return nil
			}
			return tokenString(s)
		}
	}
}

func tokenInt(s *tokenState) tokenFunc {
	return func() tokenFunc {
		switch r, ok := s.next(); {
		case r >= '0' && r <= '9':
			return tokenInt(s)
		case r == '(', r == ')':
			s.markLast()
			s.markLast()
			s.mark()
			return tokenStart(s)
		case r == '.':
			return tokenFloat(s)
		case unicode.IsSpace(r):
			s.markLast()
			return tokenStart(s)
		case !ok:
			s.markLast()
			return nil
		default:
			s.setErr(fmt.Errorf("%v: %v", errRune, r))
			return nil
		}
	}
}

func tokenFloatOrOp(s *tokenState) tokenFunc {
	return func() tokenFunc {
		switch r, ok := s.next(); {
		case r >= '0' && r <= '9':
			return tokenFloat(s)
		case r == '(', r == ')':
			s.markLast()
			s.markLast()
			s.mark()
			return tokenStart(s)
		case unicode.IsSpace(r):
			s.markLast()
			return tokenStart(s)
		case isOp(r):
			return tokenOp(s)
		case !ok:
			s.markLast()
			return nil
		default:
			s.setErr(fmt.Errorf("%v: %v", errRune, r))
			return nil
		}
	}
}

func tokenFloat(s *tokenState) tokenFunc {
	return func() tokenFunc {
		switch r, ok := s.next(); {
		case r >= '0' && r <= '9':
			return tokenFloat(s)
		case r == '(', r == ')':
			s.markLast()
			s.markLast()
			s.mark()
			return tokenStart(s)
		case unicode.IsSpace(r):
			s.markLast()
			return tokenStart(s)
		case !ok:
			s.markLast()
			return nil
		default:
			s.setErr(fmt.Errorf("%v: %v", errRune, r))
			return nil
		}
	}
}

func tokenId(s *tokenState) tokenFunc {
	return func() tokenFunc {
		switch r, ok := s.next(); {
		case r == '"':
			s.markLast()
			s.markLast()
			return tokenString(s)
		case r == '(', r == ')':
			s.markLast()
			s.markLast()
			s.mark()
			return tokenStart(s)
		case r == '.':
			s.markLast()
			s.markLast()
			return tokenOp(s)
		case unicode.IsSpace(r):
			s.markLast()
			return tokenStart(s)
		case isOp(r):
			s.markLast()
			s.markLast()
			return tokenOp(s)
		case unicode.IsGraphic(r):
			return tokenId(s)
		case !ok:
			s.markLast()
			return nil
		default:
			s.setErr(fmt.Errorf("%v: %v", errRune, r))
			return nil
		}
	}
}

func tokenOp(s *tokenState) tokenFunc {
	return func() tokenFunc {
		switch r, ok := s.next(); {
		case r == '"':
			s.markLast()
			s.markLast()
			return tokenString(s)
		case r == '(', r == ')':
			s.markLast()
			s.markLast()
			s.mark()
			return tokenStart(s)
		case unicode.IsSpace(r):
			s.markLast()
			return tokenStart(s)
		case isOp(r):
			return tokenOp(s)
		case unicode.IsGraphic(r):
			s.markLast()
			s.markLast()
			return tokenId(s)
		case !ok:
			s.markLast()
			return nil
		default:
			s.setErr(fmt.Errorf("%v: %v", errRune, r))
			return nil
		}
	}
}

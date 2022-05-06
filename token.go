package lisp

import (
	"errors"
	"fmt"
	"unicode"
	"unicode/utf8"
)

// Token is an enumeration which specifies a kind of AST token.
//go:generate stringer -type Token
type Token int

const (
	Invalid Token = iota

	Id     // x y z + - / ++
	Int    // 12345
	Float  // 123.45
	String // "abc"

	LParen // (
	RParen // )
)

// Tokenizer implements the Lisp tokenizer.
type Tokenizer struct{}

// Tokenize takes source code and returns a set of token positions.
func (Tokenizer) Tokenize(src string) ([]Pos, error) {
	s := &tokenState{src: []byte(src), line: 1, col: 1}
	for fn := tokenStart(s); fn != nil; fn = fn() {
	}
	return s.pos, s.err
}

// TokenError implements an error at a specified line and column.
type TokenError struct {
	Cause     error
	Line, Col int
	Pos       Pos
}

func (t *TokenError) Error() string {
	return fmt.Sprintf("%v: at line %d: col %d", t.Cause, t.Line, t.Col)
}

type tokenState struct {
	src []byte

	line int
	col  int
	err  error

	pos []Pos
	p   Pos
}

// mark is called when starting a token.
func (s *tokenState) mark() {
	s.pos = append(s.pos, s.p)
}

// markEnd is called when ending a token.
func (s *tokenState) markEnd() {
	s.mark()
}

func (s *tokenState) decode() (r rune, size int) {
	return utf8.DecodeRune(s.src[s.p:])
}

func (s *tokenState) advance(r rune, size int) (rune, int) {
	s.p += Pos(size)
	if r == '\n' {
		s.line++
		s.col = 1
	}
	return r, size
}

func (s *tokenState) advanceSpaces() {
	for r, size := s.decode(); unicode.IsSpace(r); r, size = s.decode() {
		s.advance(r, size)
	}
}

func (s *tokenState) advanceDigits() {
	for r, size := s.decode(); unicode.IsDigit(r); r, size = s.decode() {
		s.advance(r, size)
	}
}

func isExprOrStr(r rune) bool {
	return r == '(' || r == ')' || r == '"'
}

// idIdSymbolic checks whether r is a legal symbolic Id rune.
func isIdSymbolic(r rune) bool {
	return r != utf8.RuneError &&
		!isExprOrStr(r) &&
		(unicode.IsPunct(r) || unicode.IsSymbol(r))
}

// isIdStart checks whether r can start an Id (not a number).
func isIdStart(r rune) bool {
	return isID(r) && !unicode.IsNumber(r)
}

// isID checks whether r is a legal Id rune.
func isID(r rune) bool {
	return r != utf8.RuneError &&
		!isIdSymbolic(r) &&
		!isExprOrStr(r) &&
		!unicode.IsSpace(r)
}

var (
	errRune = errors.New("unexpected rune")
	errEOF  = errors.New("unexpected EOF")
)

func (s *tokenState) setErr(cause error) {
	s.err = &TokenError{cause, s.line, s.col, s.p}
}

type tokenFunc func() tokenFunc

func tokenStart(s *tokenState) tokenFunc {
	s.advanceSpaces()
	return func() tokenFunc {
		r, size := s.decode()
		if size == 0 {
			return nil // EOF
		}
		if r == utf8.RuneError {
			s.setErr(fmt.Errorf("%v: %v", errRune, r))
			return nil
		}
		switch {
		case r == '(', r == ')': // Expr
			s.mark()
			s.advance(r, size)
			s.markEnd()
			return tokenStart(s)
		case r == '"': // Str
			s.mark()
			s.advance(r, size)
			return tokenString(s)
		case r == '.': // Id | Float
			s.mark()
			s.advance(r, size)
			return tokenIdPunctOrFloat(s)
		case unicode.IsNumber(r): // Int | Float
			s.mark()
			s.advance(r, size)
			return tokenNumber(s)
		case isIdSymbolic(r): // symbolic Id
			s.mark()
			s.advance(r, size)
			return tokenIdSymbolic(s)
		case isIdStart(r): // Id
			s.mark()
			s.advance(r, size)
			return tokenId(s)
		default:
			s.setErr(fmt.Errorf("%v: %v", errRune, r))
			return nil
		}
	}
}

func tokenString(s *tokenState) tokenFunc {
	return func() tokenFunc {
		switch r, size := s.advance(s.decode()); {
		case r == '"':
			s.markEnd()
			return tokenStart(s)
		case r == '\\':
			return tokenEscape(s)
		case size == 0:
			s.setErr(errEOF)
			return nil
		case r == utf8.RuneError:
			s.setErr(errRune)
			return nil
		default:
			return tokenString(s)
		}
	}
}

func tokenEscape(s *tokenState) tokenFunc {
	return func() tokenFunc {
		switch r, size := s.advance(s.decode()); r {
		case 'n', 't', '\\':
			return tokenString(s)
		case 'x':
			return tokenByteLit(s)
		case utf8.RuneError:
			if size == 0 {
				s.setErr(errEOF)
				return nil
			}
			s.setErr(errRune)
			return nil
		default:
			s.setErr(fmt.Errorf("unexpected escape: \\%v", string(r)))
			return nil
		}
	}
}

// Byte literals are supported from \x00 to \xff.
// tokenByteLit2 ensures the right length.
func tokenByteLit(s *tokenState) tokenFunc {
	return func() tokenFunc {
		switch r, size := s.decode(); {
		case unicode.In(r, unicode.N, unicode.ASCII_Hex_Digit):
			s.advance(r, size)
			return tokenByteLit2(s)
		case size == 0:
			s.setErr(errEOF)
			return nil
		default:
			s.setErr(fmt.Errorf("%v: %v", errRune, r))
			return nil
		}
	}
}

func tokenByteLit2(s *tokenState) tokenFunc {
	return func() tokenFunc {
		switch r, size := s.decode(); {
		case unicode.In(r, unicode.N, unicode.ASCII_Hex_Digit):
			s.advance(r, size)
			return tokenString(s)
		case size == 0:
			s.setErr(errEOF)
			return nil
		default:
			s.setErr(fmt.Errorf("%v: %v", errRune, r))
			return nil
		}
	}
}

func tokenNumber(s *tokenState) tokenFunc {
	s.advanceDigits()
	return func() tokenFunc {
		switch r, size := s.decode(); {
		case r == '.':
			s.advance(r, size)
			return tokenFloat(s)
		default:
			s.markEnd()
			return tokenStart(s)
		}
	}
}

// tokenIdOrFloat is invoked on "."
func tokenIdPunctOrFloat(s *tokenState) tokenFunc {
	return func() tokenFunc {
		switch r, size := s.decode(); {
		case unicode.IsNumber(r):
			s.advance(r, size)
			return tokenFloat(s)
		case isIdSymbolic(r):
			s.advance(r, size)
			return tokenIdSymbolic(s)
		default:
			s.markEnd()
			return tokenStart(s)
		}
	}
}

func tokenFloat(s *tokenState) tokenFunc {
	return func() tokenFunc {
		s.advanceDigits()
		s.markEnd()
		return tokenStart(s)
	}
}

func tokenIdSymbolic(s *tokenState) tokenFunc {
	return func() tokenFunc {
		switch r, size := s.decode(); {
		case isIdSymbolic(r):
			s.advance(r, size)
			return tokenIdSymbolic(s)
		default:
			s.markEnd()
			return tokenStart(s)
		}
	}
}

func tokenId(s *tokenState) tokenFunc {
	return func() tokenFunc {
		switch r, size := s.decode(); {
		case isID(r):
			s.advance(r, size)
			return tokenId(s)
		default:
			s.markEnd()
			return tokenStart(s)
		}
	}
}

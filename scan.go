package lisp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sync"
	"unicode"
	"unicode/utf8"
)

// Scanner scans the Lisp source for Tokens.
type TokenScanner struct {
	sc   *bufio.Scanner
	tsc  tokenScanner
	once sync.Once
}

func NewTokenScanner(src io.Reader) *TokenScanner {
	var sc TokenScanner
	sc.Init(src)
	return &sc
}

func (sc *TokenScanner) Init(src io.Reader) {
	sc.sc = bufio.NewScanner(src)
	sc.sc.Split(sc.tsc.scanTokens)
}

func (sc *TokenScanner) Buffer(buf []byte, max int) {
	sc.sc.Buffer(buf, max)
}

func (sc *TokenScanner) Scan() bool {
	return sc.sc.Scan()
}

func (sc *TokenScanner) Err() error {
	return sc.sc.Err()
}

func (sc *TokenScanner) Text() string {
	return sc.sc.Text()
}

func (sc *TokenScanner) Token() (pos Pos, tok Token, text string) {
	return sc.tsc.start, sc.tsc.tok, sc.sc.Text()
}

type tokenScanner struct {
	start Pos   // absolute token start position
	end   Pos   // absolute token end position
	tok   Token // last token scanned
}

var (
	errEOF  = errors.New("unexpected EOF")
	errRune = errors.New("unexpected rune")
)

func (t *tokenScanner) scanTokens(src []byte, atEOF bool) (advance int, token []byte, err error) {
	t.start = t.end

	var tok Token

	// Skip spaces.
	var i Pos
	for i < Pos(len(src)) {
		r, size := utf8.DecodeRune(src[i:])
		if !unicode.IsSpace(r) {
			break
		}
		i += Pos(size)
	}
	start := i
	t.start += i           // Update abs start position.
	if len(src[i:]) == 0 { // No token.
		return 0, nil, io.EOF
	}

	// Get the first rune.
	r, size := utf8.DecodeRune(src[i:])
	i += Pos(size)
	switch r {
	case '(': // LParen
		tok = LParen
	case ')': // RParen
		tok = RParen
	case '"': // String
		tok = String
	string_loop:
		for i < Pos(len(src)) {
			r, size := utf8.DecodeRune(src[i:])
			i += Pos(size)
			switch r {
			case '"':
				break string_loop
			case '\\':
				if Pos(len(src)) < i {
					err = errEOF
					break string_loop
				}
				r, size := utf8.DecodeRune(src[i:])
				i += Pos(size)
				switch r {
				case 'n', 't', '\\':
				case 'x':
					for j := 0; j < 2; j++ {
						r, size := utf8.DecodeRune(src[i+Pos(j):])
						i += Pos(size)
						switch {
						case r >= '0' && r <= '9' || r >= 'A' && r <= 'F' || r >= 'a' && r <= 'f':
						case size == 0:
							err = errEOF
							break string_loop
						case r == utf8.RuneError:
							err = errRune
							break string_loop
						}
					}
				default:
					switch {
					case size == 0:
						err = errEOF
						break string_loop
					case r == utf8.RuneError:
						err = errRune
						break string_loop
					default:
						err = fmt.Errorf("unexpected escape: \\%v", string(r))
						break string_loop
					}
				}
			default:
				switch {
				case size == 0:
					err = errEOF
					break string_loop
				case r == utf8.RuneError:
					err = errRune
					break string_loop
				}
			}
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // Int | Float
		tok = Int
		var dec bool
	num_loop:
		for i < Pos(len(src)) {
			r, size := utf8.DecodeRune(src[i:])
			switch r {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			case '.':
				if dec {
					err = errRune
					break num_loop
				}
				tok = Float
				dec = true
			default:
				break num_loop
			}
			i += Pos(size)
		}
	default: // Id
		tok = Id
		var runeFunc func(rune) bool
		switch {
		case IsSymbol(r):
			runeFunc = IsSymbol
		case IsLetter(r):
			runeFunc = IsLetter
		case size == 0:
			err = io.EOF
		case r == utf8.RuneError:
			err = errRune
		default:
			panic("unreachable")
		}
		for i < Pos(len(src)) {
			r, size := utf8.DecodeRune(src[i:])
			if !runeFunc(r) {
				break
			}
			i += Pos(size)
		}
	}
	if i == Pos(len(src)) && err == nil {
		err = bufio.ErrFinalToken
	}
	advance, token = int(i), src[start:i]
	t.end += i // Update abs end position.
	t.tok = tok
	return advance, token, err
}

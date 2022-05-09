package lisp

import (
	"bufio"
	"bytes"
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

func NewTokenScanner(r io.Reader) *TokenScanner {
	return &TokenScanner{
		sc: bufio.NewScanner(r),
	}
}

func (sc *TokenScanner) Init(src []byte) {
	if sc.sc == nil {
		sc.sc = bufio.NewScanner(bytes.NewReader(src))
		sc.sc.Split(sc.tsc.scanTokens)
	}
}

func (sc *TokenScanner) Buffer(buf []byte, max int) {
	sc.once.Do(func() { sc.Init(nil) })
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
	t.start += i // Update abs start position.

	// Get the first rune.
	r, size := utf8.DecodeRune(src[i:])
	i += Pos(size)
rune_switch:
	switch r {
	case '(':
		tok = LParen
	case ')':
		tok = RParen
	case '"':
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
						r, size := utf8.DecodeRune(src[i:])
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
			i += Pos(size)
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
		}
	default: // Id
		tok = Id
		var runeFunc func(rune) bool
		switch {
		case IsSymbol(r):
			runeFunc = IsSymbol
		case IsIdent(r):
			runeFunc = IsIdent
		case size == 0:
			err = errEOF
			break rune_switch
		case r == utf8.RuneError:
			err = errRune
			break rune_switch
		default:
			panic("unreachable")
		}
		for i < Pos(len(src)) {
			r, size := utf8.DecodeRune(src[i:])
			i += Pos(size)
			if !runeFunc(r) {
				break
			}
		}
	}
	advance, token = int(i), src[start:i]
	if atEOF {
		err = bufio.ErrFinalToken
	}
	t.end += i // Update abs end position.
	t.tok = tok
	return
}

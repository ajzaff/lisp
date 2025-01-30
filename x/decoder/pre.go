package decoder

import (
	"bufio"
	"io"
)

// PreDecoder extracts expressions with raw untokenized content.
//
// Raw untokenized content is extracted as Id with Token type Invalid.
type PreDecoder struct {
	buf   *bufio.Scanner
	depth int
}

func (t *PreDecoder) Reset(r io.Reader) {
	t.buf = bufio.NewScanner(r)
	t.buf.Split(t.scanRawToken)
}

func (s *PreDecoder) scanRawToken(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 {
		return 0, nil, nil
	}

	if s.depth == 0 {
		// Depth=0: Scan until the first Group.
		for _, r := range string(data) {
			if r == '(' {
				s.depth++
				break
			}
			advance++
		}
		// Need more data.
		if s.depth == 0 {
			return advance, nil, nil
		}
	}

	// Scan for the next Paren.
	if s.depth > 0 {
		// Depth>0: Scan until the RParen or the next Group start.
		for _, r := range string(data) {
			advance++
			switch r {
			case ')':
				s.depth--
				if s.depth == 0 {
					return advance, data[:advance], nil
				}
			case '(':
				s.depth++
			}
		}

		// Need more data.
		return advance, nil, nil
	}

	// No token extracted.
	return 0, nil, nil
}

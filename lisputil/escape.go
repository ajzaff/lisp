package lisputil

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/rangetable"
)

var tab = rangetable.Merge(
	unicode.Space,
	unicode.Letter,
	unicode.Digit,
	rangetable.New('(', ')'),
)

// Escape escapes the literal to incorporate unicode in place of unsupported code points.
func Escape(s string) string {
	var sb strings.Builder
	sb.Grow(len(s))
	for _, r := range s {
		if !unicode.Is(tab, r) {
			fmt.Fprintf(&sb, "(u %v)", r)
			continue
		}
		sb.WriteRune(r)
	}
	return sb.String()
}

var escapePattern = regexp.MustCompile(`\(\s*u\s+(\d+)\)`)

// Unescape unescapes the string by removing unicode.
func Unescape(s string) string {
	var sb strings.Builder
	sb.Grow(len(s))

	i := 0
	for i < len(s) {
		match := escapePattern.FindStringSubmatchIndex(s[i:])
		if len(match) == 0 {
			break
		}
		begin, end, d0, d1 := i+match[0], i+match[1], i+match[2], i+match[3]
		u, _ := strconv.ParseInt(s[d0:d1], 10, 32)

		sb.WriteString(s[i:begin])
		sb.WriteRune(rune(u))

		i = end
	}
	sb.WriteString(s[i:])

	return sb.String()
}

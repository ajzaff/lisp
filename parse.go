package innit

import (
	"fmt"
	"strings"
)

func Parse(src string) (Node, error) {
	tokens, err := Tokenize(src)
	if err != nil {
		return nil, err
	}
	return parseTokens(src, tokens)
}

func parseTokens(src string, tokens []Pos) (Node, error) {
	var out NodeList
	var stack []*Expr
	for i := 0; i < len(tokens); i += 2 {
		pos := tokens[i]
		end := tokens[i+1]
		tok := string(src[pos:end])
		switch {
		case tok == "(":
			stack = append(stack, &Expr{LParen: pos})
		case tok == ")":
			if len(stack) == 0 {
				return nil, fmt.Errorf("innit.Parse: internal error: unexpected ')' at pos %d", pos)
			}
			stack[len(stack)-1].RParen = pos
			if len(stack) == 1 {
				out = append(out, stack[0])
			} else {
				stack[len(stack)-2].X = append(stack[len(stack)-2].X, stack[len(stack)-1])
			}
			stack = stack[:len(stack)-1]
		case strings.HasPrefix(tok, `"`):
			lit := &BasicLit{
				Tok:      String,
				ValuePos: Pos(i),
				Value:    tok,
			}
			if len(stack) == 0 {
				out = append(out, lit)
			} else {
				stack[len(stack)-1].X = append(stack[len(stack)-1].X, lit)
			}
		case (tok[0] == '.' && len(tok) > 1) || tok[0] >= '0' && tok[0] <= '9':
			lit := &BasicLit{
				ValuePos: Pos(i),
				Value:    tok,
			}
			if strings.ContainsRune(tok, '.') {
				lit.Tok = Float
			} else {
				lit.Tok = Int
			}
			if len(stack) == 0 {
				out = append(out, lit)
			} else {
				stack[len(stack)-1].X = append(stack[len(stack)-1].X, lit)
			}
		default:
			id := &BasicLit{Tok: Id, ValuePos: Pos(i), Value: tok}
			if len(stack) == 0 {
				out = append(out, id)
			} else {
				stack[len(stack)-1].X = append(stack[len(stack)-1].X, id)
			}
		}
	}
	if len(stack) > 0 {
		err := fmt.Errorf("innit.Parse: unexpected EOF")
		if len(tokens) >= 2 {
			pos, end := tokens[len(tokens)-2], tokens[len(tokens)-1]
			err = fmt.Errorf("%v: at %q", err, string(src[pos:end]))
		}
		return nil, err
	}
	return out, nil
}

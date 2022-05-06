package lisp

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser struct{}

func (Parser) Parse(src string) ([]Node, error) {
	tokens, err := Tokenizer{}.Tokenize(src)
	if err != nil {
		return nil, err
	}
	return parseTokens(src, tokens)
}

func parseTokens(src string, tokens []Pos) ([]Node, error) {
	var out []Node
	var stack []*ExprNode
	for i := 0; i < len(tokens); i += 2 {
		pos, end := tokens[i], tokens[i+1]
		tok := string(src[pos:end])
		switch {
		case tok == "(":
			stack = append(stack, &ExprNode{LParen: pos})
		case tok == ")":
			if len(stack) == 0 {
				return nil, fmt.Errorf("lisp.Parse: internal error: unexpected ')' at pos %d", pos)
			}
			stack[len(stack)-1].RParen = pos
			if len(stack) == 1 {
				out = append(out, stack[0])
			} else {
				stack[len(stack)-2].Expr = append(stack[len(stack)-2].Expr, stack[len(stack)-1])
			}
			stack = stack[:len(stack)-1]
		case strings.HasPrefix(tok, `"`):
			lit := &LitNode{
				LitPos: Pos(i),
				Lit:    StringLit(tok),
			}
			if len(stack) == 0 {
				out = append(out, lit)
			} else {
				stack[len(stack)-1].Expr = append(stack[len(stack)-1].Expr, lit)
			}
		default:
			lit := &LitNode{LitPos: Pos(i)}

			if x, err := strconv.ParseInt(tok, 10, 64); err == nil {
				lit.Lit = IntLit(x)
			} else if x, err := strconv.ParseFloat(tok, 64); err == nil {
				lit.Lit = FloatLit(x)
			} else {
				lit.Lit = IdLit(tok)
			}
			if len(stack) == 0 {
				out = append(out, lit)
			} else {
				stack[len(stack)-1].Expr = append(stack[len(stack)-1].Expr, lit)
			}
		}
	}
	if len(stack) > 0 {
		err := fmt.Errorf("lisp.Parse: unexpected EOF")
		if len(tokens) >= 2 {
			pos, end := tokens[len(tokens)-2], tokens[len(tokens)-1]
			err = fmt.Errorf("%v: at %q", err, string(src[pos:end]))
		}
		return nil, err
	}
	return out, nil
}

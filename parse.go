package lisp

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser struct{}

func (Parser) Parse(src string) ([]Node, error) {
	return parseTokens(NewTokenScanner(strings.NewReader(src)))
}

func parseTokens(sc *TokenScanner) ([]Node, error) {
	var out []Node
	var stack []*ExprNode
	for sc.Scan() {
		pos, tok, text := sc.Token()
		switch tok {
		case LParen:
			stack = append(stack, &ExprNode{LParen: pos})
		case RParen:
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
		case Int:
			x, err := strconv.ParseInt(text, 10, 64)
			if err != nil {
				panic("parseInt")
			}
			lit := &LitNode{
				LitPos: pos,
				Lit:    IntLit(x),
			}
			if len(stack) == 0 {
				out = append(out, lit)
			} else {
				stack[len(stack)-1].Expr = append(stack[len(stack)-1].Expr, lit)
			}
		case Float:
			x, err := strconv.ParseFloat(text, 64)
			if err != nil {
				panic("parseFloat")
			}
			lit := &LitNode{
				LitPos: pos,
				Lit:    FloatLit(x),
			}
			if len(stack) == 0 {
				out = append(out, lit)
			} else {
				stack[len(stack)-1].Expr = append(stack[len(stack)-1].Expr, lit)
			}
		case String:
			lit := &LitNode{
				LitPos: pos,
				Lit:    StringLit(text),
			}
			if len(stack) == 0 {
				out = append(out, lit)
			} else {
				stack[len(stack)-1].Expr = append(stack[len(stack)-1].Expr, lit)
			}
		case Id:
			lit := &LitNode{
				LitPos: pos,
				Lit:    IdLit(text),
			}
			if len(stack) == 0 {
				out = append(out, lit)
			} else {
				stack[len(stack)-1].Expr = append(stack[len(stack)-1].Expr, lit)
			}
		default:
			panic("unreachable")
		}
	}
	if len(stack) > 0 {
		err := fmt.Errorf("lisp.Parse: unexpected EOF")
		return nil, err
	}
	return out, nil
}

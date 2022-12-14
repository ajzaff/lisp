package lisp

func DecodeNum(b []byte) (LitNode, error) {
	n, err := decodeNum(b)
	if err != nil {
		return LitNode{}, err
	}
	return LitNode{
		Lit: Lit{
			Token: Number,
			Text:  string(b[:n]),
		},
		EndPos: Pos(n),
	}, nil
}

func decodeNum(b []byte) (n int, err error) {
	return 0, nil
}

package lisp

func DecodeNum(b []byte) (Node, error) {
	n, err := decodeNum(b)
	if err != nil {
		return Node{}, err
	}
	return Node{
		Val: Lit{
			Token: Number,
			Text:  string(b[:n]),
		},
		End: Pos(n),
	}, nil
}

func decodeNum(b []byte) (n int, err error) {
	return 0, nil
}

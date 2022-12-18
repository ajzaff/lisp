package decoder

import "github.com/ajzaff/lisp"

func DecodeNum(b []byte) (lisp.Node, error) {
	n, err := decodeNum(b)
	if err != nil {
		return lisp.Node{}, err
	}
	return lisp.Node{
		Val: lisp.Lit{
			Token: lisp.Int,
			Text:  string(b[:n]),
		},
		End: lisp.Pos(n),
	}, nil
}

func decodeNum(b []byte) (n int, err error) {
	return 0, nil
}

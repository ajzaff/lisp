package decoder

import "github.com/ajzaff/lisp"

func DecodeNat(b []byte) (lisp.Val, error) {
	n, err := decodeNum(b)
	if err != nil {
		return nil, err
	}
	return lisp.Lit{
		Token: lisp.Nat,
		Text:  string(b[:n]),
	}, nil
}

// FIXME: Implement this!
func decodeNum(b []byte) (n int, err error) {
	return 0, nil
}

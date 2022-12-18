package lisputil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBigEscape(t *testing.T) {
	input := `Some years before his untimely death M. de Rougé read his translation of
	this chapter before the Académie des Sciences. It is much to be lamented
	that this has never been published. I have, in addition to the versions
	of other scholars, a copy of one by Mr. Goodwin, with whom I read this
	and other chapters nearly thirty years ago. But this kind of literature
	is not one of those in which his marvellous sagacity showed to
	advantage.
`
	got := Escape(input)
	t.Log(got)
}

func TestEscape(t *testing.T) {
	input := "!@#$%^&*()1234567890"
	got := Escape(input)
	want := "(u 33)(u 64)(u 35)(u 36)(u 37)(u 94)(u 38)(u 42)()1234567890"
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Escape() got diff:\n%s", diff)
	}
}

func TestUnescape(t *testing.T) {
	input := "(u 33)(u 64)(u 35)(u 36)(u 37)(u 94)(u 38)(u 42)()1234567890"
	got := Unescape(input)
	want := "!@#$%^&*()1234567890"
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unescape() got diff:\n%s", diff)
	}
}

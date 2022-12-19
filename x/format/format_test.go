package format

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSource(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		want  string
	}{{
		name:  "de-dupe spaces across delimiter class",
		input: "abc (abc)  ( abc )  abc    1234",
		want:  "abc(abc)(abc)abc 1234",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			got := Source([]byte(tc.input))
			if diff := cmp.Diff(tc.want, string(got)); diff != "" {
				t.Errorf("Source() got diff:\n%s", diff)
			}
		})
	}
}

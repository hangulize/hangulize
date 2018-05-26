package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type replacer1 struct{}
type replacer2 struct{}

func (*replacer1) Replacements(word string) []Replacement {
	return []Replacement{Replacement{0, 1, []string{"1"}}}
}

func (*replacer2) Replacements(word string) []Replacement {
	return []Replacement{Replacement{1, 2, []string{"2"}}}
}

func TestRewrite(t *testing.T) {
	rep1 := &replacer1{}
	rep2 := &replacer2{}
	assert.Equal(t, "12llo", Rewrite("hello", rep1, rep2)[0])
}

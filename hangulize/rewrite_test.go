package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func rewrite(word string, fromExpr string, toExpr string) string {
	spec := parseSpec(`
	vars:
		vowels = "a", "e", "i", "o", "u"
		abc    = "a", "b", "c"
		def    = "d", "e", "f"

	macros:
		"@" = "<vowels>"
	`)
	r := NewRule(
		newPattern(fromExpr, spec),
		newRPattern(toExpr, spec),
	)
	return Rewrite(word, r)[0]
}

func TestVarToVar(t *testing.T) {
	assert.Equal(t, "d", rewrite("a", "<abc>", "<def>"))
	assert.Equal(t, "def", rewrite("abc", "<abc>", "<def>"))
	assert.Equal(t, "fde", rewrite("cab", "<abc>", "<def>"))
	assert.Equal(t, "XfdeX", rewrite("XcabX", "<abc>", "<def>"))
}

func TestCaret(t *testing.T) {
	assert.Equal(t, "baa", rewrite("aaa", "^a", "b"))
}

// -----------------------------------------------------------------------------

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

package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func rewrite(word string, fromExpr string, toExpr string) string {
	r := Rule{from: fixturePattern(fromExpr), to: fixtureRPatterns(toExpr)}
	return r.Rewrite(word)
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

package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func rewrite(word string, fromExpr string, toExpr string) string {
	r := Rule{from: compile(fromExpr), to: []string{toExpr}}
	return r.Rewrite(word)
}

func TestVarToVar(t *testing.T) {
	assert.Equal(t, "d", rewrite("a", "<abc>", "<def>"))
}

func TestCaret(t *testing.T) {
	assert.Equal(t, "baa", rewrite("aaa", "^a", "b"))
}

package hangulize

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func newRule(fromExpr string, toExpr string) *Rule {
// 	spec := parseSpec(`
// 	vars:
// 		vowels = "a", "e", "i", "o", "u"
// 		abc    = "a", "b", "c"
// 		def    = "d", "e", "f"

// 	macros:
// 		"@" = "<vowels>"
// 	`)
// 	return &Rule{
// 		newPattern(fromExpr, spec),
// 		newRPattern(toExpr, spec),
// 	}
// }

// func rewrite(word string, fromExpr string, toExpr string) string {
// 	r := newRule(fromExpr, toExpr)
// 	return Rewrite(word, []*Rule{r})
// }

// func TestVarToVar(t *testing.T) {
// 	assert.Equal(t, "d", rewrite("a", "<abc>", "<def>"))
// 	assert.Equal(t, "def", rewrite("abc", "<abc>", "<def>"))
// 	assert.Equal(t, "fde", rewrite("cab", "<abc>", "<def>"))
// 	assert.Equal(t, "XfdeX", rewrite("XcabX", "<abc>", "<def>"))
// }

// func TestCaret(t *testing.T) {
// 	assert.Equal(t, "baa", rewrite("aaa", "^a", "b"))
// }

// func TestRewrite(t *testing.T) {
// 	r1 := newRule(`h`, `1`)
// 	r2 := newRule(`el|lo`, `2`)

// 	word := Rewrite("hello?", []*Rule{r1, r2})
// 	assert.Equal(t, "122?", word)
// }

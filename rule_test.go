package hangulize

import (
	"testing"

	"github.com/hangulize/hre"
	"github.com/stretchr/testify/assert"
)

func TestRuleString(t *testing.T) {
	p, _ := hre.NewPattern("foo", nil, nil)
	rp := hre.NewRPattern("bar", nil, nil)
	r := Rule{0, p, rp}
	assert.Equal(t, `"foo" -> "bar"`, r.String())
}

func TestRuleReplacements(t *testing.T) {
	p, _ := hre.NewPattern("foo", nil, nil)
	rp := hre.NewRPattern("bar", nil, nil)
	r := Rule{0, p, rp}

	repls := r.replacements("abcfoodef")

	assert.Len(t, repls, 1)
	assert.Equal(t, 3, repls[0].Start)
	assert.Equal(t, 6, repls[0].Stop)
	assert.Equal(t, "bar", repls[0].Word)
}

func TestRuleReplace(t *testing.T) {
	p, _ := hre.NewPattern("foo", nil, nil)
	rp := hre.NewRPattern("bar", nil, nil)
	r := Rule{0, p, rp}
	assert.Equal(t, "abcbardef", r.Replace("abcfoodef"))
}

func TestRuleUnmatchedVar(t *testing.T) {
	vars := map[string][]string{
		"foo": {"foo"},
		"bar": {"b", "a", "r"},
		"baz": {"b", "a", "z"},
	}
	p, _ := hre.NewPattern("<foo>", nil, vars)
	rp := hre.NewRPattern("<bar><baz>", nil, vars)
	r := Rule{0, p, rp}

	// Silently, keep the original.
	assert.Equal(t, "abcfoodef", r.Replace("abcfoodef"))
}

package hre

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexpLetters(t *testing.T) {
	assert.Equal(t, "abc", regexpLetters(`abc`))
	assert.Equal(t, "abc", regexpLetters(`(a|b|c)`))
	assert.Equal(t, "", regexpLetters(`\n`))
}

func TestSubstr(t *testing.T) {
	assert.Equal(t, "", substr("abc", -1, 1))
	assert.Equal(t, "", substr("abc", 1, -1))
	assert.Equal(t, "a", substr("abc", 0, 1))
	assert.Equal(t, "bc", substr("abc", 1, 3))
	assert.Equal(t, "bc", substr("abc", 1, 10))
	assert.Equal(t, "", substr("abc", 1, 0))
}

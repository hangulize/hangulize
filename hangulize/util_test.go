package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexpLetters(t *testing.T) {
	assert.Equal(t, "abc", regexpLetters(`abc`))
	assert.Equal(t, "abc", regexpLetters(`(a|b|c)`))
	assert.Equal(t, "", regexpLetters(`\n`))
}

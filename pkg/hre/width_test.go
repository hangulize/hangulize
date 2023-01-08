package hre

import (
	"regexp/syntax"
	"testing"

	"github.com/stretchr/testify/assert"
)

func regexpMaxWidth(expr string) int {
	re, err := syntax.Parse(expr, syntax.Perl)

	if err != nil {
		// Failed to parse Regexp.
		panic(err)
	}

	return RegexpMaxWidth(re)
}

func TestRegexpMaxWidth(t *testing.T) {
	assert.Equal(t, 1, regexpMaxWidth(`a`))
	assert.Equal(t, 4, regexpMaxWidth(`a...`))
	assert.Equal(t, 1, regexpMaxWidth(`a|b`))
	assert.Equal(t, 1, regexpMaxWidth(`((a)|b)`))
	assert.Equal(t, 1, regexpMaxWidth(`((a)|b)?`))
	assert.Equal(t, 3, regexpMaxWidth(`a|bc|def|g|hi|jkl`))
	assert.Equal(t, 1, regexpMaxWidth(`[abcde]`))
	assert.Equal(t, 1, regexpMaxWidth(`[\]abcde]`))
	assert.Equal(t, 1, regexpMaxWidth(`[(abcde)]`))
	assert.Equal(t, 126, regexpMaxWidth(`(...){1,42}`))
	assert.Equal(t, -1, regexpMaxWidth(`.*`))
	assert.Equal(t, -1, regexpMaxWidth(`.+`))
	assert.Equal(t, -1, regexpMaxWidth(`.*?`))
	assert.Equal(t, -1, regexpMaxWidth(`.+?`))
	assert.Equal(t, -1, regexpMaxWidth(`(.+|...)`))
	assert.Equal(t, -1, regexpMaxWidth(`.+...`))
}

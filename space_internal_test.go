package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasSpace(t *testing.T) {
	assert.True(t, hasSpace(" "))
	assert.True(t, hasSpace(" prefix"))
	assert.True(t, hasSpace("suffix "))
	assert.True(t, hasSpace("in fix"))
	assert.True(t, hasSpace(" wrap "))
	assert.True(t, hasSpace(" mix ed "))
	assert.True(t, hasSpace("\u3000")) // IDEOGRAPHIC SPACE
	assert.True(t, hasSpace("\u00A0")) // NO-BREAK SPACE
	assert.True(t, hasSpace("\u2000")) // EN QUAD
	assert.False(t, hasSpace(""))
	assert.False(t, hasSpace("Hello"))
}

func TestHasSpaceOnly(t *testing.T) {
	assert.True(t, hasSpaceOnly(" "))
	assert.True(t, hasSpaceOnly("   "))
	assert.True(t, hasSpaceOnly("\u3000\u00A0\u2000")) // IDEOGRAPHIC SPACE, NO-BREAK SPACE, EN QUAD
	assert.False(t, hasSpaceOnly(" !"))
	assert.False(t, hasSpaceOnly(""))
}

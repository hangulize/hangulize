package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeRoman(t *testing.T) {
	assert.Equal(t, "hello", NormalizeRoman("Hello"))
	assert.Equal(t, "cafe", NormalizeRoman("Café"))
	assert.Equal(t, "melee", NormalizeRoman("Mêlée"))
	assert.Equal(t, " cafe, hi! ", NormalizeRoman(" Café, Hi! "))
}

package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// The most basic case of Hangulize.
func TestItaGloria(t *testing.T) {
	gloria := Hangulize("ita", "gloria")
	assert.Equal(t, "글로리아", gloria)
}

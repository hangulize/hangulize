package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlash(t *testing.T) {
	// It was failed in the original Hangulize.
	// The result was "글로르이아" without the slash.
	assert.Equal(t, "글로르/이아", Hangulize("ita", "glor/ia"))
}

package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHangulizeTrace(t *testing.T) {
	spec, _ := LoadSpec("ita")
	h := NewHangulizer(spec)
	transcribed, traces := h.HangulizeTrace("Cappuccino")

	assert.Equal(t, "카푸치노", transcribed)
	assert.NotEqual(t, 0, len(traces))
}

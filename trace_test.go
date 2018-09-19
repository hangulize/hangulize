package hangulize

import (
	"bytes"
	"strings"
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

func TestTraceString(t *testing.T) {
	spec, _ := LoadSpec("ita")
	h := NewHangulizer(spec)
	_, traces := h.HangulizeTrace("Cappuccino")

	assert.Equal(t, `[Input] "Cappuccino"`, traces[0].String())
	assert.Equal(t, `[Normalize] "cappuccino" | (latin)`, traces[1].String())
}

func TestTracesRender(t *testing.T) {
	spec, _ := LoadSpec("ita")
	h := NewHangulizer(spec)
	_, traces := h.HangulizeTrace("Cappuccino")

	var b bytes.Buffer
	traces.Render(&b)
	rendered := b.String()

	assert.True(t, strings.HasPrefix(rendered, "[Input]"))
	assert.True(t, strings.Contains(rendered, "Cappuccino"))
	assert.True(t, strings.Contains(rendered, "카푸치노"))
}

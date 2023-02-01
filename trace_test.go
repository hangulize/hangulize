package hangulize_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/hangulize/hangulize"
	"github.com/stretchr/testify/assert"
)

func TestHangulizeTrace(t *testing.T) {
	spec, _ := hangulize.LoadSpec("ita")
	h := hangulize.New(spec)
	result, traces, err := h.HangulizeTrace("Cappuccino")

	assert.NoError(t, err)
	assert.Equal(t, "카푸치노", result)
	assert.NotEqual(t, 0, len(traces))
}

func TestTraceString(t *testing.T) {
	spec, _ := hangulize.LoadSpec("ita")
	h := hangulize.New(spec)
	_, traces, _ := h.HangulizeTrace("Cappuccino")

	assert.Equal(t, `[Input] "Cappuccino"`, traces[0].String())
	assert.Equal(t, `[Normalize] "cappuccino" | (Latn)`, traces[1].String())
}

func TestTracesRender(t *testing.T) {
	spec, _ := hangulize.LoadSpec("ita")
	h := hangulize.New(spec)
	_, traces, _ := h.HangulizeTrace("Cappuccino")

	var b bytes.Buffer
	traces.Render(&b)
	rendered := b.String()

	assert.True(t, strings.HasPrefix(rendered, "[Input]"))
	assert.True(t, strings.Contains(rendered, "Cappuccino"))
	assert.True(t, strings.Contains(rendered, "카푸치노"))
}

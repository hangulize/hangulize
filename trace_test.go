package hangulize_test

import (
	"testing"

	"github.com/hangulize/hangulize"
	"github.com/stretchr/testify/assert"
)

func TestHangulizeTrace(t *testing.T) {
	spec, _ := hangulize.LoadSpec("ita")
	h := hangulize.New(spec)

	traces := make([]hangulize.Trace, 0)
	h.Trace(func(t hangulize.Trace) {
		traces = append(traces, t)
	})

	_, _ = h.Hangulize("Cappuccino")
	assert.NotEmpty(t, traces)

	prevLength := len(traces)
	h.Trace(nil)

	_, _ = h.Hangulize("Cappuccino")
	assert.Equal(t, prevLength, len(traces))
}

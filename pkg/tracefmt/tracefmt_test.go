package tracefmt_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/hangulize/hangulize"
	"github.com/hangulize/hangulize/pkg/tracefmt"
	"github.com/stretchr/testify/assert"
)

func TestFprintTraces(t *testing.T) {
	spec, _ := hangulize.LoadSpec("ita")
	h := hangulize.New(spec)

	traces := make([]hangulize.Trace, 0)
	h.Trace(func(t hangulize.Trace) {
		traces = append(traces, t)
	})

	_, _ = h.Hangulize("Cappuccino")

	var b bytes.Buffer
	tracefmt.FprintTraces(&b, traces)
	rendered := b.String()

	assert.True(t, strings.HasPrefix(rendered, "[Input]"))
	assert.True(t, strings.Contains(rendered, "Cappuccino"))
	assert.True(t, strings.Contains(rendered, "카푸치노"))
}

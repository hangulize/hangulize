//go:build js

package main

import (
	"syscall/js"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHangulize(t *testing.T) {
	promise := jsHangulize.Invoke("ita", "Cappuccino")
	assert.True(t, promise.InstanceOf(promiseClass))

	result := await(promise)
	assert.Equal(t, "카푸치노", result.String())
}

func TestHangulizeTrace(t *testing.T) {
	traces := make([]js.Value, 0)

	await(jsHangulize.Invoke("ita", "Cappuccino", js.FuncOf(func(this js.Value, args []js.Value) any {
		trace := args[0]
		traces = append(traces, trace)
		return nil
	})))

	assert.NotEmpty(t, traces)
}

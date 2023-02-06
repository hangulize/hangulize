//go:build js

package main

import (
	"syscall/js"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUseTranslit(t *testing.T) {
	translit := js.FuncOf(func(this js.Value, args []js.Value) any {
		return promiseClass.New(js.FuncOf(func(this js.Value, args []js.Value) any {
			resolve := args[0]
			resolve.Invoke("TRANSLIT")
			return nil
		}))
	})

	jsUseTranslit.Invoke("furigana", translit)
	defer jsUnuseTranslit.Invoke("furigana")

	promise := jsHangulize.Invoke("jpn", "INPUT")
	result := await(promise)
	assert.Equal(t, "TRANSLIT", result.String())
}

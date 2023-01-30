//go:build js

package main

import (
	"syscall/js"

	"github.com/hangulize/hangulize/phonemize/furigana"
)

var jsPhonemize = js.FuncOf(func(this js.Value, args []js.Value) any {
	word := args[0].String()

	result, err := furigana.P.Phonemize(word)
	if err != nil {
		return js.Global().Get("Error").New(err.Error())
	}
	return result
})

var version string

func main() {
	js.Global().Set("phonemize", jsPhonemize)
	js.Global().Get("phonemize").Set("version", version)

	furigana.P.Load()

	<-make(chan struct{})
}

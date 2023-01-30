//go:build js

package main

import (
	"syscall/js"

	"github.com/hangulize/hangulize/phonemize/pinyin"
)

var jsPhonemize = js.FuncOf(func(this js.Value, args []js.Value) any {
	word := args[0].String()

	result, err := pinyin.P.Phonemize(word)
	if err != nil {
		return js.Global().Get("Error").New(err.Error())
	}
	return result
})

var version string

func main() {
	js.Global().Set("phonemize", jsPhonemize)
	js.Global().Get("phonemize").Set("version", version)

	pinyin.P.Load()

	<-make(chan struct{})
}

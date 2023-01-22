//go:build js

package main

import (
	"syscall/js"

	"github.com/hangulize/hangulize"
	"github.com/hangulize/hangulize/phonemize/furigana"
	"github.com/hangulize/hangulize/phonemize/pinyin"
)

var phonemizers = map[string]hangulize.Phonemizer{
	furigana.P.ID(): &furigana.P,
	pinyin.P.ID():   &pinyin.P,
}

var jsPhonemize = js.FuncOf(func(this js.Value, args []js.Value) any {
	id := args[0].String()
	word := args[1].String()

	result, err := phonemizers[id].Phonemize(word)
	if err != nil {
		return js.Global().Get("Error").New(err.Error())
	}
	return result
})

var jsLoad = js.FuncOf(func(this js.Value, args []js.Value) any {
	for _, p := range phonemizers {
		_ = p.Load()
	}
	return nil
})

var version string

func main() {
	js.Global().Set("phonemize", jsPhonemize)
	js.Global().Get("phonemize").Set("load", jsLoad)

	<-make(chan struct{})
}

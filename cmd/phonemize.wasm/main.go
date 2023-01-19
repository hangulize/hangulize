//go:build js

package main

import (
	"syscall/js"

	"github.com/hangulize/hangulize/phonemize/furigana"
	"github.com/hangulize/hangulize/phonemize/pinyin"
)

var version string

var jsFurigana = js.FuncOf(func(this js.Value, args []js.Value) any {
	word := args[0].String()
	result, err := furigana.P.Phonemize(word)
	if err != nil {
		return js.Global().Get("Error").New(err.Error())
	}
	return result
})

var jsPinyin = js.FuncOf(func(this js.Value, args []js.Value) any {
	word := args[0].String()
	result, err := pinyin.P.Phonemize(word)
	if err != nil {
		return js.Global().Get("Error").New(err.Error())
	}
	return result
})

func main() {
	jsFurigana.Set("load", js.FuncOf(func(this js.Value, args []js.Value) any {
		_ = furigana.P.Load()
		return nil
	}))

	jsPinyin.Set("load", js.FuncOf(func(this js.Value, args []js.Value) any {
		_ = pinyin.P.Load()
		return nil
	}))

	js.Global().Set("phonemizers", js.ValueOf(map[string]interface{}{
		"furigana": jsFurigana,
		"pinyin":   jsPinyin,
	}))

	<-make(chan struct{})
}

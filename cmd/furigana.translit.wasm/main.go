//go:build js

package main

import (
	"syscall/js"

	"github.com/hangulize/hangulize/translit/furigana"
)

var jsTransliterate = js.FuncOf(func(this js.Value, args []js.Value) any {
	word := args[1].String()

	result, err := furigana.T.Transliterate(word)
	if err != nil {
		return js.Global().Get("Error").New(err.Error())
	}
	return result
})

var version string

func main() {
	js.Global().Set("translit", jsTransliterate)
	js.Global().Get("translit").Set("version", version)

	// warm-up
	furigana.T.Transliterate("")

	<-make(chan struct{})
}

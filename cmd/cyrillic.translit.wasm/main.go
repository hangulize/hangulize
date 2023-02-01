//go:build js

package main

import (
	"syscall/js"

	"github.com/hangulize/hangulize"
	"github.com/hangulize/hangulize/translit/cyrillic"
)

var translits = map[string]hangulize.Translit{}

var jsTransliterate = js.FuncOf(func(this js.Value, args []js.Value) any {
	scheme := args[0].String()
	word := args[1].String()

	result, err := translits[scheme].Transliterate(word)
	if err != nil {
		return js.Global().Get("Error").New(err.Error())
	}
	return result
})

var version string

func main() {
	for _, t := range cyrillic.Ts {
		translits[t.Scheme()] = t
	}

	js.Global().Set("translit", jsTransliterate)
	js.Global().Get("translit").Set("version", version)

	<-make(chan struct{})
}

//go:build js

package main

import (
	"fmt"
	"syscall/js"

	"github.com/hangulize/hangulize"
)

// jsHangulize is a JavaScript function that transcribes a word into Hangul.
//
//	const result = await hangulize('ita', 'cappuccino')
var jsHangulize = js.FuncOf(func(this js.Value, args []js.Value) any {
	lang := args[0].String()
	word := args[1].String()

	return js.Global().Get("Promise").New(js.FuncOf(func(this js.Value, args []js.Value) any {
		resolve := args[0]
		reject := args[1]

		// Blocking code needs a new goroutine to avoid deadlock in a
		// JavaScript build.
		go func() {
			result, err := hangulize.Hangulize(lang, word)
			if err != nil {
				fmt.Println("error", err)
				reject.Invoke(js.Global().Get("Error").New(err.Error()))
			} else {
				fmt.Println(result)
				resolve.Invoke(result)
			}
		}()

		return nil
	}))
})

//go:build js

package main

import (
	"syscall/js"

	"github.com/hangulize/hangulize"
)

// jsHangulize wraps hangulize.Hangulize in JavaScript.
//
//	hangulize(lang: string, word: string, traceFn?: (trace) => void) => Promise<string>
var jsHangulize = js.FuncOf(func(this js.Value, args []js.Value) any {
	lang := args[0].String()
	word := args[1].String()

	traceFn := js.Undefined()
	if len(args) > 2 {
		traceFn = args[2]
	}

	return js.Global().Get("Promise").New(js.FuncOf(func(this js.Value, args []js.Value) any {
		resolve := args[0]
		reject := args[1]

		throw := func(err error) {
			reject.Invoke(js.Global().Get("Error").New(err.Error()))
		}

		// Blocking code needs a new goroutine to avoid deadlock in a
		// JavaScript build.
		go func() {
			spec, err := hangulize.LoadSpec(lang)
			if err != nil {
				throw(err)
				return
			}

			h := hangulize.New(spec)
			for _, translit := range hangulize.Translits() {
				h.UseTranslit(translit)
			}

			if !traceFn.IsUndefined() {
				h.Trace(func(t hangulize.Trace) {
					traceFn.Invoke(jsTrace(t))
				})
			}

			result, err := h.Hangulize(word)
			if err != nil {
				throw(err)
				return
			}

			resolve.Invoke(result)
		}()

		return nil
	}))
})

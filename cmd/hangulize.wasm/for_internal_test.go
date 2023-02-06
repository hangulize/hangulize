//go:build js

package main

import "syscall/js"

// promiseClass is the Promise class in JavaScript.
var promiseClass = js.Global().Get("Promise")

// await waits the given Promise and returns the result value.
func await(promise js.Value) js.Value {
	ch := make(chan js.Value)

	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		ch <- args[0]
		return nil
	}))

	return <-ch
}

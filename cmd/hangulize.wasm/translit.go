//go:build js

package main

import (
	"errors"
	"syscall/js"

	"github.com/hangulize/hangulize"
)

// promiseTranslit is a wrapper of an async JavaScript function to implement
// hangulize.Translit.
type promiseTranslit struct {
	scheme string
	fn     js.Value
}

func (t promiseTranslit) Scheme() string {
	return t.scheme
}

func (t promiseTranslit) Load() error {
	return nil
}

func (p promiseTranslit) Transliterate(word string) (string, error) {
	thenCh := make(chan string)
	defer close(thenCh)

	catchCh := make(chan error)
	defer close(catchCh)

	then := js.FuncOf(func(this js.Value, args []js.Value) any {
		result := args[0].String()
		thenCh <- result
		return nil
	})

	catch := js.FuncOf(func(this js.Value, args []js.Value) any {
		err := errors.New(args[0].Get("message").String())
		catchCh <- err
		return nil
	})

	p.fn.Invoke(word).Call("then", then).Call("catch", catch)

	select {
	case result := <-thenCh:
		return result, nil
	case err := <-catchCh:
		return "", err
	}
}

// jsUseTranslit wraps hangulize.UseTranslit in JavaScript.
//
//	useTranslit(scheme: string, fn: (word: string) => Promise<string>) => boolean
var jsUseTranslit = js.FuncOf(func(this js.Value, args []js.Value) any {
	scheme := args[0].String()
	fn := args[1]
	return hangulize.UseTranslit(promiseTranslit{scheme, fn})
})

// jsUnuseTranslit wraps hangulize.UnuseTranslit in JavaScript.
//
//	unuseTranslit(scheme: string) => boolean
var jsUnuseTranslit = js.FuncOf(func(this js.Value, args []js.Value) any {
	scheme := args[0].String()
	return hangulize.UnuseTranslit(scheme)
})

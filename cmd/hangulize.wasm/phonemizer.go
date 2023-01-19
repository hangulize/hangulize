//go:build js

package main

import (
	"errors"
	"syscall/js"

	"github.com/hangulize/hangulize"
)

// promisePhonemizer is a wrapper of an async JavaScript function to implement
// hangulize.Phonemizer.
type promisePhonemizer struct {
	id string
	fn js.Value
}

func (p promisePhonemizer) ID() string {
	return p.id
}

func (p promisePhonemizer) Load() error {
	return nil
}

func (p promisePhonemizer) Phonemize(word string) (string, error) {
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

// jsImportPhonemizer
var jsImportPhonemizer = js.FuncOf(func(this js.Value, args []js.Value) any {
	id := args[0].String()
	fn := args[1]
	return hangulize.ImportPhonemizer(promisePhonemizer{id, fn})
})

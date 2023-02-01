//go:build js

package main

import (
	"syscall/js"

	"github.com/hangulize/hangulize"
)

var version string

func main() {
	js.Global().Set("hangulize", jsHangulize)
	js.Global().Get("hangulize").Set("version", version)
	js.Global().Get("hangulize").Set("specs", jsSpecs(hangulize.ListLangs()))
	js.Global().Get("hangulize").Set("useTranslit", jsUseTranslit)

	<-make(chan struct{}, 0)
}

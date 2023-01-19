//go:build js

package main

import (
	"syscall/js"

	"github.com/hangulize/hangulize"
)

var version string

func main() {
	specs := []interface{}{}
	for _, lang := range hangulize.ListLangs() {
		spec, _ := hangulize.LoadSpec(lang)
		specs = append(specs, jsSpec(spec))
	}

	js.Global().Set("hangulize", jsHangulize)
	js.Global().Get("hangulize").Set("version", version)
	js.Global().Get("hangulize").Set("specs", specs)
	js.Global().Get("hangulize").Set("importPhonemizer", jsImportPhonemizer)

	<-make(chan struct{}, 0)
}

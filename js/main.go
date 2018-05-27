package main

import (
	"github.com/gopherjs/gopherjs/js"

	"github.com/sublee/hangulize2/hangulize"
)

func main() {
	js.Global.Set("hangulize", hangulize.Hangulize)
	js.Global.Set("__hangulize__", map[string]interface{}{
		"Hangulize": hangulize.Hangulize,
		"Version":   hangulize.Version,

		"ListLangs":     hangulize.ListLangs,
		"LoadSpec":      hangulize.LoadSpec,
		"NewHangulizer": hangulize.NewHangulizer,

		"ComposeHangul": hangulize.ComposeHangul,
	})
}

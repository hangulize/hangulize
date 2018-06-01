package main

import (
	"strings"

	"github.com/gopherjs/gopherjs/js"

	"github.com/sublee/hangulize2/hangulize"
)

func main() {
	js.Global.Set("hangulize", hangulize.Hangulize)
	js.Global.Set("__hangulize__", map[string]interface{}{
		"Hangulize": hangulize.Hangulize,
		"Version":   hangulize.Version,

		"ListLangs": hangulize.ListLangs,
		"LoadSpec":  hangulize.LoadSpec,

		"NewHangulizer": func(spec *hangulize.Spec) *js.Object {
			h := hangulize.NewHangulizer(spec)
			return js.MakeWrapper(h)
		},

		"ParseSpec": func(source string) *js.Object {
			r := strings.NewReader(source)
			spec, _ := hangulize.ParseSpec(r)
			return js.MakeWrapper(spec)
		},

		"ComposeHangul": hangulize.ComposeHangul,
	})
}

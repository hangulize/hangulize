package main

import (
	"strings"

	"github.com/gopherjs/gopherjs/js"

	"github.com/sublee/hangulize2/hangulize"
)

// Langs is the preloaded language informations.
var Langs = make(map[string]hangulize.Language)

func init() {
	for _, lang := range hangulize.ListLangs() {
		spec, _ := hangulize.LoadSpec(lang)
		Langs[lang] = spec.Lang
	}
}

func main() {
	exports := map[string]interface{}{
		"Hangulize": hangulize.Hangulize,
		"Version":   hangulize.Version,
		"Langs":     Langs,

		"LoadSpec": hangulize.LoadSpec,

		"ParseSpec": func(source string) *js.Object {
			r := strings.NewReader(source)
			spec, _ := hangulize.ParseSpec(r)
			return js.MakeWrapper(spec)
		},

		"NewHangulizer": func(spec *hangulize.Spec) *js.Object {
			h := hangulize.NewHangulizer(spec)
			return js.MakeWrapper(h)
		},
	}

	js.Global.Set("hangulize", hangulize.Hangulize)

	for attr, val := range exports {
		js.Global.Get("hangulize").Set(attr, val)
	}

	if js.Module != js.Undefined {
		js.Module.Set("exports", exports)
	}
}

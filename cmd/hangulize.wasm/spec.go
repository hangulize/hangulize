//go:build js

package main

import (
	"syscall/js"

	"github.com/hangulize/hangulize"
)

type object = map[string]interface{}
type array = []interface{}

// jsSpec converts a Spec as a JavaScript value.
func jsSpec(s *hangulize.Spec) js.Value {
	translit := array{}
	for _, m := range s.Lang.Translit {
		translit = append(translit, m)
	}

	lang := js.ValueOf(object{
		"id":       s.Lang.ID,
		"code2":    s.Lang.Codes[0],
		"code3":    s.Lang.Codes[1],
		"english":  s.Lang.English,
		"korean":   s.Lang.Korean,
		"script":   s.Lang.Script,
		"translit": translit,
	})

	authors := array{}
	for _, a := range s.Config.Authors {
		authors = append(authors, a)
	}
	config := js.ValueOf(object{
		"authors": authors,
		"stage":   s.Config.Stage,
	})

	test := array{}
	for _, exm := range s.Test {
		test = append(test, js.ValueOf(object{
			"word":   exm[0],
			"result": exm[1],
		}))
	}

	return js.ValueOf(object{
		"lang":   lang,
		"config": config,
		"test":   test,
		"source": s.Source,
	})
}

// jsSpecs creates an object of specs indexed by the language ID.
func jsSpecs(langs []string) js.Value {
	specs := make(map[string]interface{}, len(langs))
	for _, lang := range langs {
		spec, _ := hangulize.LoadSpec(lang)
		specs[lang] = jsSpec(spec)
	}
	return js.ValueOf(specs)
}

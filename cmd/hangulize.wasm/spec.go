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
	lang := js.ValueOf(object{
		"id":         s.Lang.ID,
		"code2":      s.Lang.Codes[0],
		"code3":      s.Lang.Codes[1],
		"english":    s.Lang.English,
		"korean":     s.Lang.Korean,
		"script":     s.Lang.Script,
		"phonemizer": s.Lang.Phonemizer,
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

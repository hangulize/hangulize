//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/hangulize/hangulize"
	"github.com/hangulize/hangulize/phonemize/furigana"
	"github.com/hangulize/hangulize/phonemize/pinyin"
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
			"word":        exm[0],
			"transcribed": exm[1],
		}))
	}

	return js.ValueOf(object{
		"lang":   lang,
		"config": config,
		"test":   test,
		"source": s.Source,
	})
}

// jsHangulize is a JavaScript function that transcribes a word into Hangul.
var jsHangulize = js.FuncOf(func(this js.Value, args []js.Value) any {
	lang := args[0].String()
	word := args[1].String()
	return hangulize.Hangulize(lang, word)
})

var version string

func main() {
	hangulize.UsePhonemizer(&furigana.P)
	hangulize.UsePhonemizer(&pinyin.P)

	specs := []interface{}{}
	for _, lang := range hangulize.ListLangs() {
		spec, _ := hangulize.LoadSpec(lang)
		specs = append(specs, jsSpec(spec))
	}

	js.Global().Set("hangulize", jsHangulize)
	js.Global().Get("hangulize").Set("version", version)
	js.Global().Get("hangulize").Set("specs", specs)

	<-make(chan bool)
}

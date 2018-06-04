package main

import (
	"strings"

	"github.com/gopherjs/gopherjs/js"

	"github.com/sublee/hangulize2/hangulize"
)

// packSpec packs a spec into {"spec": the-spec, "info": {"lang":, "config":,
// "test":, "source": ...}}.  It exposes some information
// to be used in JavaScript-side.
func packSpec(s *hangulize.Spec) *js.Object {
	// Append only lang, config, test, source.
	info := js.Global.Get("Object").New()

	lang := js.Global.Get("Object").New()
	lang.Set("id", s.Lang.ID)
	lang.Set("codes", s.Lang.Codes)
	lang.Set("english", s.Lang.English)
	lang.Set("korean", s.Lang.Korean)
	lang.Set("script", s.Lang.Script)

	config := js.Global.Get("Object").New()
	config.Set("authors", s.Config.Authors)
	config.Set("stage", s.Config.Stage)

	test := js.Global.Get("Object").New()
	for i, pair := range s.Test {
		o := js.Global.Get("Object").New()
		o.Set("word", pair.Left())
		o.Set("transcribed", pair.Right()[0])
		test.SetIndex(i, &o)
	}

	info.Set("lang", lang)
	info.Set("config", config)
	info.Set("test", test)
	info.Set("source", s.Source)

	// Result
	o := js.Global.Get("Object").New()
	o.Set("spec", js.MakeWrapper(s))
	o.Set("info", info)
	return o
}

var specs = make(map[string]*js.Object)

func init() {
	for _, lang := range hangulize.ListLangs() {
		spec, _ := hangulize.LoadSpec(lang)
		specs[lang] = packSpec(spec)
	}
}

func main() {
	exports := map[string]interface{}{
		"hangulize": hangulize.Hangulize,
		"version":   hangulize.Version,

		"specs": specs,

		"loadSpec": func(langID string) *js.Object {
			spec, ok := hangulize.LoadSpec(langID)
			if !ok {
				return nil
			}
			return packSpec(spec)
		},

		"parseSpec": func(source string) *js.Object {
			r := strings.NewReader(source)
			spec, err := hangulize.ParseSpec(r)
			if err != nil {
				return nil
			}
			return packSpec(spec)
		},

		"hangulizer": func(spec *js.Object) *js.Object {
			_spec := spec.Get("spec").Interface().(*hangulize.Spec)
			h := hangulize.NewHangulizer(_spec)
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

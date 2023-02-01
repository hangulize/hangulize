package main

import (
	"encoding/json"
	"fmt"

	"github.com/hangulize/hangulize"
)

type manifest struct {
	Version   string          `json:"version"`
	Specs     map[string]spec `json:"specs"`
	Translits []string        `json:"translits"`
}

type spec struct {
	Lang   lang      `json:"lang"`
	Config config    `json:"config"`
	Test   []example `json:"test"`
}

type lang struct {
	ID       string   `json:"id"`
	Code2    string   `json:"code2"`
	Code3    string   `json:"code3"`
	English  string   `json:"english"`
	Korean   string   `json:"korean"`
	Script   string   `json:"script"`
	Translit []string `json:"translit"`
}

type config struct {
	Authors []string `json:"authors"`
	Stage   string   `json:"stage"`
}

type example struct {
	Word   string `json:"word"`
	Result string `json:"result"`
}

// jsonSpec converts a Spec to be encoded as JSON.
func jsonSpec(s *hangulize.Spec) spec {
	lang := lang{
		s.Lang.ID,
		s.Lang.Codes[0],
		s.Lang.Codes[1],
		s.Lang.English,
		s.Lang.Korean,
		s.Lang.Script,
		s.Lang.Translit,
	}

	config := config{s.Config.Authors, s.Config.Stage}

	test := make([]example, 0, len(s.Test))
	for _, exm := range s.Test {
		test = append(test, example{exm[0], exm[1]})
	}

	return spec{lang, config, test}
}

var version string

func main() {
	langs := hangulize.ListLangs()
	specs := make(map[string]spec, len(langs))
	translitSet := make(map[string]bool)

	for _, lang := range langs {
		spec, _ := hangulize.LoadSpec(lang)
		specs[lang] = jsonSpec(spec)

		for _, m := range spec.Lang.Translit {
			translitSet[m] = true
		}
	}

	translits := make([]string, 0, len(translitSet))
	for id := range translitSet {
		translits = append(translits, id)
	}

	b, err := json.Marshal(manifest{version, specs, translits})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}

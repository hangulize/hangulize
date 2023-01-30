package main

import (
	"encoding/json"
	"fmt"

	"github.com/hangulize/hangulize"
)

type manifest struct {
	Version     string          `json:"version"`
	Specs       map[string]spec `json:"specs"`
	Phonemizers []string        `json:"phonemizers"`
}

type spec struct {
	Lang   lang      `json:"lang"`
	Config config    `json:"config"`
	Test   []example `json:"test"`
}

type lang struct {
	ID         string `json:"id"`
	Code2      string `json:"code2"`
	Code3      string `json:"code3"`
	English    string `json:"english"`
	Korean     string `json:"korean"`
	Script     string `json:"script"`
	Phonemizer string `json:"phonemizer"`
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
		s.Lang.Phonemizer,
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
	phonemizerSet := make(map[string]bool)

	for _, lang := range langs {
		spec, _ := hangulize.LoadSpec(lang)
		specs[lang] = jsonSpec(spec)

		if spec.Lang.Phonemizer != "" {
			phonemizerSet[spec.Lang.Phonemizer] = true
		}
	}

	phonemizers := make([]string, 0, len(phonemizerSet))
	for id := range phonemizerSet {
		phonemizers = append(phonemizers, id)
	}

	b, err := json.Marshal(manifest{version, specs, phonemizers})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}

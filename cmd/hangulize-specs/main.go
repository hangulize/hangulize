package main

import (
	"encoding/json"
	"fmt"

	"github.com/hangulize/hangulize"
)

type root struct {
	Version string `json:"version"`
	Specs   []spec `json:"specs"`
}

type spec struct {
	Lang   lang     `json:"lang"`
	Config config   `json:"config"`
	Test   []sample `json:"test"`
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

type sample struct {
	Word        string `json:"word"`
	Transcribed string `json:"transcribed"`
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

	test := make([]sample, 0, len(s.Test))
	for _, exm := range s.Test {
		test = append(test, sample{exm[0], exm[1]})
	}

	return spec{lang, config, test}
}

var version string

func main() {
	root := root{version, make([]spec, 0)}

	for _, lang := range hangulize.ListLangs() {
		spec, _ := hangulize.LoadSpec(lang)
		root.Specs = append(root.Specs, jsonSpec(spec))
	}

	b, err := json.Marshal(root)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}

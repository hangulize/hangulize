package hangulize

import (
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/sublee/hangulize2/hgl"
)

// Spec represents a transactiption specification for a language.
type Spec struct {
	Lang   Language
	Config Config

	Vars      map[string][]string
	Macros    map[string]string
	Rewrite   []hgl.Pair
	Hangulize []hgl.Pair

	Test []hgl.Pair
}

func (s *Spec) String() string {
	return fmt.Sprintf("<Spec lang=%s>", s.Lang.id)
}

// Language identifies a natural language.
type Language struct {
	id      string
	code    []string
	english string
	korean  string
	script  string
}

// Config keeps some configurations for a transactiption specification.
type Config struct {
	authors []string
	stage   string
	markers []rune
}

// ParseSpec parses a Spec from an HGL source.
func ParseSpec(r io.Reader) (*Spec, error) {
	h, err := hgl.Parse(r)

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse HGL source")
	}

	var sec hgl.Section
	var ok bool

	// lang
	sec, ok = h["lang"]
	if !ok {
		return nil, errors.New(`"lang" section required`)
	}
	lang := newLanguage(sec.(*hgl.DictSection))

	// config (optional)
	var config *Config

	if sec, ok := h["config"]; ok {
		config, err = newConfig(sec.(*hgl.DictSection))

		if err != nil {
			return nil, err
		}
	}

	// vars (optional)
	var vars map[string][]string
	if sec, ok := h["vars"]; ok {
		vars = sec.(*hgl.DictSection).Map()
	}

	// macros (optional)
	var macros map[string]string

	if sec, ok := h["macros"]; ok {
		macros, err = newMacros(sec.(*hgl.DictSection))

		if err != nil {
			return nil, err
		}
	}

	// rewrite
	sec, ok = h["rewrite"]
	if !ok {
		return nil, errors.New(`"rewrite" section required`)
	}
	rewrite := sec.(*hgl.ListSection).Array()

	// hangulize
	sec, ok = h["hangulize"]
	if !ok {
		return nil, errors.New(`"hangulize" section required`)
	}
	hangulize := sec.(*hgl.ListSection).Array()

	// test (optional)
	var test []hgl.Pair
	if sec, ok := h["test"]; ok {
		test = sec.(*hgl.ListSection).Array()
	}

	spec := Spec{
		*lang,
		*config,
		vars,
		macros,
		rewrite,
		hangulize,
		test,
	}
	return &spec, nil
}

func newLanguage(dict *hgl.DictSection) *Language {
	lang := Language{
		id:      dict.One("id"),
		code:    dict.All("code"),
		english: dict.One("english"),
		korean:  dict.One("korean"),
		script:  dict.One("script"),
	}
	return &lang
}

func newConfig(dict *hgl.DictSection) (*Config, error) {
	// A marker must be 1-character.
	stringMarkers := dict.All("markers")
	markers := make([]rune, len(stringMarkers))

	for i, stringMarker := range stringMarkers {
		if len(stringMarker) != 1 {
			err := fmt.Errorf("marker %#v must be 1-character", stringMarker)
			return nil, err
		}
		markers[i] = rune(stringMarker[0])
	}

	config := Config{
		authors: dict.All("authors"),
		stage:   dict.One("stage"),
		markers: markers,
	}
	return &config, nil
}

func newMacros(dict *hgl.DictSection) (map[string]string, error) {
	_map := dict.Map()
	macros := make(map[string]string, len(_map))

	for src, dst := range _map {
		if len(dst) != 1 {
			err := fmt.Errorf("macro %#v must has single target", src)
			return nil, err
		}
	}

	return macros, nil
}

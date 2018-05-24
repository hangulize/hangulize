package hangulize

import (
	"fmt"
	"io"

	"github.com/pkg/errors"

	"../hgl"
)

// Spec represents a transactiption specification for a language.
type Spec struct {
	Lang   Language
	Config Config

	Macros    map[string]string
	Vars      map[string][]string
	rewrite   *Rewriter
	hangulize *Rewriter

	Test []hgl.Pair
}

func (s *Spec) String() string {
	return fmt.Sprintf("<Spec lang=%s>", s.Lang.ID)
}

// ParseSpec parses a Spec from an HGL source.
func ParseSpec(r io.Reader) (*Spec, error) {
	var err error

	h, err := hgl.Parse(r)

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse HGL source")
	}

	// Every sections are optional.  An empty HGL source is also valid spec.

	// lang
	var lang Language

	if sec, ok := h["lang"]; ok {
		_lang, err := newLanguage(sec.(*hgl.DictSection))

		if err != nil {
			return nil, err
		}

		lang = *_lang
	}

	// config
	var config Config

	if sec, ok := h["config"]; ok {
		_config, err := newConfig(sec.(*hgl.DictSection))

		if err != nil {
			return nil, err
		}

		config = *_config
	}

	// macros
	var macros map[string]string

	if sec, ok := h["macros"]; ok {
		macros, err = newMacros(sec.(*hgl.DictSection))

		if err != nil {
			return nil, err
		}
	}

	// vars
	var vars map[string][]string
	if sec, ok := h["vars"]; ok {
		vars = sec.(*hgl.DictSection).Map()
	}

	// rewrite
	var rewrite *Rewriter
	var rewritePairs []hgl.Pair
	if sec, ok := h["rewrite"]; ok {
		rewritePairs = sec.(*hgl.ListSection).Array()
	}

	rewrite, err = NewRewriter(rewritePairs, macros, vars)
	if err != nil {
		return nil, err
	}

	// hangulize
	var hangulize *Rewriter
	var hangulizePairs []hgl.Pair
	if sec, ok := h["hangulize"]; ok {
		hangulizePairs = sec.(*hgl.ListSection).Array()
	}

	hangulize, err = NewRewriter(hangulizePairs, macros, vars)
	if err != nil {
		return nil, err
	}

	// test
	var test []hgl.Pair
	if sec, ok := h["test"]; ok {
		test = sec.(*hgl.ListSection).Array()
	}

	spec := Spec{
		lang,
		config,
		macros,
		vars,
		rewrite,
		hangulize,
		test,
	}
	return &spec, nil
}

// -----------------------------------------------------------------------------

// Language identifies a natural language.
type Language struct {
	ID      string    // Arbitrary, but identifiable language ID.
	Codes   [2]string // [0]: ISO 639-1 code, [1]: ISO 639-3 code
	English string    // The langauge name in English.
	Korean  string    // The langauge name in Korean.
	Script  string
}

func newLanguage(dict *hgl.DictSection) (*Language, error) {
	_codes := dict.All("codes")

	if len(_codes) != 2 {
		return nil, errors.New("codes must be 2; ISO 639-1 and 3")
	}

	var codes [2]string
	codes[0] = _codes[0]
	codes[1] = _codes[1]

	lang := Language{
		ID:      dict.One("id"),
		Codes:   codes,
		English: dict.One("english"),
		Korean:  dict.One("korean"),
		Script:  dict.One("script"),
	}
	return &lang, nil
}

// -----------------------------------------------------------------------------

// Config keeps some configurations for a transactiption specification.
type Config struct {
	Authors []string
	Stage   string
	Markers []rune
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
		Authors: dict.All("authors"),
		Stage:   dict.One("stage"),
		Markers: markers,
	}
	return &config, nil
}

// -----------------------------------------------------------------------------

func newMacros(dict *hgl.DictSection) (map[string]string, error) {
	_map := dict.Map()
	macros := make(map[string]string, len(_map))

	for src, dst := range _map {
		if len(dst) != 1 {
			err := fmt.Errorf("macro %#v must has single target", src)
			return nil, err
		}

		macros[src] = dst[0]
	}

	return macros, nil
}

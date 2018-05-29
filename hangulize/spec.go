package hangulize

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"

	"github.com/sublee/hangulize2/hgl"
)

// Spec represents a transactiption specification for a language.
type Spec struct {
	// Meta information sections.
	Lang   Language
	Config Config

	// Helper setting sections.
	Macros    map[string]string
	Vars      map[string][]string
	Normalize map[string][]string

	// Quantize/Transcribe
	Quantize   []*Rule
	Transcribe []*Rule

	// Test examples.
	Test []hgl.Pair

	// Prepared stuffs.
	normReplacer *strings.Replacer
	normLetters  []string
	norm         Normalizer
	letters      []string
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

	// -------------------------------------------------------------------------
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
		macros, err = sec.(*hgl.DictSection).Injective()

		if err != nil {
			return nil, err
		}
	}

	// vars
	var vars map[string][]string
	if sec, ok := h["vars"]; ok {
		vars = sec.(*hgl.DictSection).Map()
	}

	// normalize
	var normalize map[string][]string
	if sec, ok := h["normalize"]; ok {
		normalize = sec.(*hgl.DictSection).Map()
	}

	// quantize
	var quantizePairs []hgl.Pair
	if sec, ok := h["quantize"]; ok {
		quantizePairs = sec.(*hgl.ListSection).Array()
	}

	quantize, err := newRules(quantizePairs, macros, vars)
	if err != nil {
		return nil, err
	}

	// transcribe
	var transcribePairs []hgl.Pair
	if sec, ok := h["transcribe"]; ok {
		transcribePairs = sec.(*hgl.ListSection).Array()
	}

	transcribe, err := newRules(transcribePairs, macros, vars)
	if err != nil {
		return nil, err
	}

	// test
	var test []hgl.Pair
	if sec, ok := h["test"]; ok {
		test = sec.(*hgl.ListSection).Array()
	}

	// -------------------------------------------------------------------------

	// custom normalization
	var args []string
	for to, froms := range normalize {
		for _, from := range froms {
			args = append(args, from, to)
		}
	}
	normReplacer := strings.NewReplacer(args...)

	// letters in normalize
	var normLetters []string
	for to := range normalize {
		normLetters = append(normLetters, to)
	}

	// canonical normalizer
	norm, ok := GetNormalizer(lang.Script)
	_ = ok
	// if !ok {
	// 	return nil, fmt.Errorf("no normalizer for %#v", lang.Script)
	// }

	// unique/sorted letters in quantize/transcribe
	var letters []string

	rules := append(quantize, transcribe...)

	for _, rule := range rules {
		for _, let := range rule.From.letters {
			letters = append(letters, let)
		}
	}

	letters = set(letters)

	// -------------------------------------------------------------------------

	spec := Spec{
		lang,
		config,

		macros,
		vars,
		normalize,

		quantize,
		transcribe,

		test,

		normReplacer,
		normLetters,
		norm,
		letters,
	}
	return &spec, nil
}

// -----------------------------------------------------------------------------
// "lang" section

// Language identifies a natural language.
type Language struct {
	ID      string    // Arbitrary, but identifiable language ID.
	Codes   [2]string // [0]: ISO 639-1 code, [1]: ISO 639-3 code
	English string    // The langauge name in English.
	Korean  string    // The langauge name in Korean.
	Script  string
}

func (l *Language) String() string {
	return fmt.Sprintf("%s(%s)", l.ID, l.English)
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
// "config" section

// Config keeps some configurations for a transactiption specification.
type Config struct {
	Authors []string
	Stage   string
}

func newConfig(dict *hgl.DictSection) (*Config, error) {
	config := Config{
		Authors: dict.All("authors"),
		Stage:   dict.One("stage"),
	}
	return &config, nil
}

// -----------------------------------------------------------------------------
// "quantize"/"transcribe" section

func newRules(
	pairs []hgl.Pair,

	macros map[string]string,
	vars map[string][]string,

) ([]*Rule, error) {

	rules := make([]*Rule, len(pairs))

	for i, pair := range pairs {
		from, err := newPattern(pair.Left(), macros, vars)
		if err != nil {
			return nil, err
		}

		right := pair.Right()
		to := newRPattern(right[0], macros, vars)

		rules[i] = &Rule{from, to}
	}

	return rules, nil
}

package hangulize

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/pkg/errors"

	"github.com/hangulize/hgl"
	"github.com/hangulize/hre"
	"github.com/hangulize/stringset"
)

// Spec represents a transactiption specification for a language.
type Spec struct {
	// Meta information sections
	Lang   Language
	Config Config

	// Helper setting sections
	Macros    map[string]string
	Vars      map[string][]string
	Normalize map[string][]string

	// Rewrite/Transcribe
	Rewrite    []Rule
	Transcribe []Rule

	// Test examples
	Test [][2]string

	// Source code
	Source string

	// Prepared stuffs
	script script
	puncts stringset.StringSet

	// Custom normalization
	normReplacer *strings.Replacer
	normLetters  stringset.StringSet
}

func (s Spec) String() string {
	return s.Lang.ID
}

// GoString implements GoStringer for Spec.
func (s Spec) GoString() string {
	return fmt.Sprintf("hangulize.Spec{Lang.ID: %#v}", s.Lang.ID)
}

// ParseSpec parses a Spec from an HGL source.
func ParseSpec(r io.Reader) (*Spec, error) {
	var err error
	var sourceBuf bytes.Buffer

	// Use TeeReader to copy the source while parsing.
	tee := io.TeeReader(r, &sourceBuf)

	h, err := hgl.Parse(tee)

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse HGL source")
	}

	source := sourceBuf.String()

	// -------------------------------------------------------------------------
	// Every sections are optional. An empty HGL source is also valid spec.

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

	// rewrite
	var rewritePairs []hgl.Pair
	if sec, ok := h["rewrite"]; ok {
		rewritePairs = sec.(*hgl.ListSection).Pairs()
	}

	rewrite, err := newRules(rewritePairs, macros, vars)
	if err != nil {
		return nil, err
	}

	// transcribe
	var transcribePairs []hgl.Pair
	if sec, ok := h["transcribe"]; ok {
		transcribePairs = sec.(*hgl.ListSection).Pairs()
	}

	transcribe, err := newRules(transcribePairs, macros, vars)
	if err != nil {
		return nil, err
	}

	// test
	var test [][2]string
	if sec, ok := h["test"]; ok {
		for _, pair := range sec.(*hgl.ListSection).Pairs() {
			word := pair.Left()
			transcribed := pair.Right()[0]

			exm := [2]string{word, transcribed}

			test = append(test, exm)
		}
	}

	// -------------------------------------------------------------------------

	script, ok := getScript(lang.Script)
	if !ok {
		return nil, errors.Errorf("script not found: %s", lang.Script)
	}
	puncts := collectPuncts(rewrite, transcribe)

	// custom normalization
	var args []string
	for to, froms := range normalize {
		for _, from := range froms {
			args = append(args, from, to)
		}
	}
	normReplacer := strings.NewReplacer(args...)

	// letters in normalize
	normLetters := make(stringset.StringSet)
	for to := range normalize {
		normLetters[to] = true
	}

	// -------------------------------------------------------------------------

	spec := Spec{
		lang,
		config,

		macros,
		vars,
		normalize,

		rewrite,
		transcribe,

		test,

		source,

		script,
		puncts,

		normReplacer,
		normLetters,
	}
	return &spec, nil
}

// -----------------------------------------------------------------------------
// "lang" section

// Language identifies a natural language.
type Language struct {
	ID         string    // Arbitrary, but identifiable language ID.
	Codes      [2]string // [0]: ISO 639-1 code, [1]: ISO 639-3 code
	English    string    // The language name in English.
	Korean     string    // The language name in Korean.
	Script     string
	Phonemizer string
}

func (l Language) String() string {
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
		ID:         dict.One("id"),
		Codes:      codes,
		English:    dict.One("english"),
		Korean:     dict.One("korean"),
		Script:     dict.One("script"),
		Phonemizer: dict.One("phonemizer"),
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
// "rewrite"/"transcribe" section

func newRules(
	pairs []hgl.Pair,

	macros map[string]string,
	vars map[string][]string,

) ([]Rule, error) {

	rules := make([]Rule, len(pairs))

	for i, pair := range pairs {
		from, err := hre.NewPattern(pair.Left(), macros, vars)
		if err != nil {
			return nil, err
		}

		negAWidth, negBWidth := from.NegativeLookaroundWidths()
		if negAWidth == -1 || negBWidth == -1 {
			return nil, errors.Errorf(
				"%s contains unlimited negative lookaround", from)
		}

		right := pair.Right()
		to := hre.NewRPattern(right[0], macros, vars)

		rules[i] = Rule{i, from, to}
	}

	return rules, nil
}

// -----------------------------------------------------------------------------

// collectPuncts collects punctuation characters from rewrite/transcribe rules.
// It discards the punctuations that is used only for rewriting hints.
func collectPuncts(rewrite []Rule, transcribe []Rule) stringset.StringSet {
	var puncts []string
	rletters := make(map[string]bool)

	collectFrom := func(rule Rule) {
		for _, let := range rule.From.Letters() {
			ch, _ := utf8.DecodeRuneInString(let)

			// Collect only punctuation characters. (category P)
			if !unicode.IsPunct(ch) {
				continue
			}

			// Punctuations appearing in the above RPatterns should be
			// discarded. Bacause they are just hints for rewriting.
			if rletters[let] {
				continue
			}

			puncts = append(puncts, let)
		}
	}

	for _, rule := range rewrite {
		// Mark letters in RPatterns.
		for _, let := range rule.To.Letters() {
			rletters[let] = true
		}

		collectFrom(rule)
	}

	for _, rule := range transcribe {
		collectFrom(rule)
	}

	return stringset.NewStringSet(puncts...)
}

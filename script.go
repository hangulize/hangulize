package hangulize

import (
	"strings"
	"unicode"

	kagome "github.com/ikawaha/kagome/tokenizer"
	"golang.org/x/text/unicode/norm"
)

// script represents a writing system.
type script interface {
	Is(rune) bool
	Normalize(rune) rune
}

// scripts is the registry of Scripts by their name.
var scripts = map[string]script{
	// Latin is the default.
	"": &_Latin{},

	"cyrillic": &_Cyrillic{},
	"georgian": &_Georgian{},
	"greek":    &_Greek{},
	"kana":     &_Kana{},
	"latin":    &_Latin{},
}

// getScript chooses a script by the script name.
func getScript(name string) script {
	script, ok := scripts[name]
	if !ok {
		// Get the default.
		latin := scripts[""]
		return latin
	}
	return script
}

// -----------------------------------------------------------------------------

// _Latin represents the Latin or Roman script. Most langauges Hangulize
// supports use this script system. So it's the default script.
type _Latin struct{}

// Is checks whether the character is Latin or not.
func (_Latin) Is(ch rune) bool {
	return unicode.Is(unicode.Latin, ch)
}

// Normalize converts a Latin character into
// ISO basic Latin lower alphabet [a-z]:
//
//   Pokémon -> pokemon
//
func (_Latin) Normalize(ch rune) rune {
	props := norm.NFD.PropertiesString(string(ch))
	bin := props.Decomposition()
	if len(bin) != 0 {
		ch = rune(bin[0])
	}
	return unicode.ToLower(ch)
}

// -----------------------------------------------------------------------------

// _Cyrillic represents the Cyrillic script.
//
//   вулкан
//
type _Cyrillic struct{}

// Is checks whether the character is Cyrillic or not.
func (_Cyrillic) Is(ch rune) bool {
	return unicode.Is(unicode.Cyrillic, ch)
}

// Normalize converts character into lower case.
func (_Cyrillic) Normalize(ch rune) rune {
	return unicode.ToLower(ch)
}

// -----------------------------------------------------------------------------

// _Georgian represents the Georgian script.
//
//   ასომთავრული
//
type _Georgian struct{}

// Is checks whether the character is Georgian or not.
func (_Georgian) Is(ch rune) bool {
	return unicode.Is(unicode.Georgian, ch)
}

// Normalize does nothing. Georgian is unicase, which means, there's only one
// case for each letter.
func (_Georgian) Normalize(ch rune) rune {
	return ch
}

// -----------------------------------------------------------------------------

// _Greek represents the Greek script.
//
//   ελληνικά
//
type _Greek struct{}

// Is checks whether the character is Greek or not.
func (_Greek) Is(ch rune) bool {
	return unicode.Is(unicode.Greek, ch)
}

// Normalize converts character into lower case.
func (_Greek) Normalize(ch rune) rune {
	return unicode.ToLower(ch)
}

// -----------------------------------------------------------------------------

// TODO(sublee): Find out a Kanji to Kana dictionary to hangulize Japanese
// perfectly.

// _Kana represents the Kana script including Hiragana and Katakana.
//
//   ひらがな カタカナ
//
type _Kana struct {
	kagome *kagome.Tokenizer
}

// Is checks whether the character is either Hiragana or Katakana.
func (_Kana) Is(ch rune) bool {
	return unicode.Is(unicode.Hiragana, ch) || unicode.Is(unicode.Katakana, ch)
}

// Kagome caches a Kagome tokenizer because it is expensive.
func (k *_Kana) Kagome() *kagome.Tokenizer {
	if k.kagome == nil {
		t := kagome.New()
		k.kagome = &t
	}
	return k.kagome
}

// ReadKanji guesses Katakana from Kanji.
func (k *_Kana) ReadKanji(word string) string {
	const (
		kanjiMin = rune(0x4e00)
		kanjiMax = rune(0x9faf)
	)

	kanjiFound := false
	for _, ch := range word {
		if ch >= kanjiMin && ch <= kanjiMax {
			kanjiFound = true
			break
		}
	}

	// Don't initialize the Kagome tokenizer if there's no Kanji.
	if !kanjiFound {
		return word
	}

	tokens := k.Kagome().Tokenize(word)

	chunks := make([]string, 0)

	for _, tok := range tokens {
		switch tok.Class {
		case kagome.KNOWN:
			fs := tok.Features()
			pron := fs[7]
			chunks = append(chunks, pron)
		case kagome.UNKNOWN:
			chunks = append(chunks, tok.Surface)
		}
	}

	return strings.Join(chunks, " ")
}

// Normalize converts Hiragana to Katakana.
func (_Kana) Normalize(ch rune) rune {
	const (
		hiraganaMin = rune(0x3040)
		hiraganaMax = rune(0x309f)
	)

	if hiraganaMin <= ch && ch <= hiraganaMax {
		// hiragana to katakana
		return ch + 96
	}
	return ch
}

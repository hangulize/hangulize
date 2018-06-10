package hangulize

import (
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// Script represents a writing system.
type Script interface {
	Is(rune) bool
	Normalize(rune) rune
}

// scripts is the registry of Scripts by their name.
var scripts = map[string]Script{
	// Latin is the default.
	"": &Latin{},

	"cyrillic": &Cyrillic{},
	"georgian": &Georgian{},
	"greek":    &Greek{},
	"kana":     &Kana{},
	"latin":    &Latin{},
}

// GetScript chooses a script by the script name.
func GetScript(name string) Script {
	script, ok := scripts[name]
	if !ok {
		// Get the default.
		latin := scripts[""]
		return latin
	}
	return script
}

// -----------------------------------------------------------------------------

// Latin represents the Latin or Roman script. Most langauges Hangulize
// supports use this script system. So it's the default script.
type Latin struct{}

// Is checks whether the character is Latin or not.
func (Latin) Is(ch rune) bool {
	return unicode.Is(unicode.Latin, ch)
}

// Normalize converts a Latin character into
// ISO basic Latin lower alphabet [a-z]:
//
//   Pokémon -> pokemon
//
func (Latin) Normalize(ch rune) rune {
	props := norm.NFD.PropertiesString(string(ch))
	bin := props.Decomposition()
	if len(bin) != 0 {
		ch = rune(bin[0])
	}
	return unicode.ToLower(ch)
}

// -----------------------------------------------------------------------------

// Cyrillic represents the Cyrillic script.
//
//   вулкан
//
type Cyrillic struct{}

// Is checks whether the character is Cyrillic or not.
func (Cyrillic) Is(ch rune) bool {
	return unicode.Is(unicode.Cyrillic, ch)
}

// Normalize converts character into lower case.
func (Cyrillic) Normalize(ch rune) rune {
	return unicode.ToLower(ch)
}

// -----------------------------------------------------------------------------

// Georgian represents the Georgian script.
//
//   ასომთავრული
//
type Georgian struct{}

// Is checks whether the character is Georgian or not.
func (Georgian) Is(ch rune) bool {
	return unicode.Is(unicode.Georgian, ch)
}

// Normalize does nothing. Georgian is unicase, which means, there's only one
// case for each letter.
func (Georgian) Normalize(ch rune) rune {
	return ch
}

// -----------------------------------------------------------------------------

// Greek represents the Greek script.
//
//   ελληνικά
//
type Greek struct{}

// Is checks whether the character is Greek or not.
func (Greek) Is(ch rune) bool {
	return unicode.Is(unicode.Greek, ch)
}

// Normalize converts character into lower case.
func (Greek) Normalize(ch rune) rune {
	return unicode.ToLower(ch)
}

// -----------------------------------------------------------------------------

// TODO(sublee): Find out a Kanji to Kana dictionary to hangulize Japanese
// perfectly.

// Kana represents the Kana script including Hiragana and Katakana.
//
//   ひらがな カタカナ
//
type Kana struct{}

// Is checks whether the character is either Hiragana or Katakana.
func (Kana) Is(ch rune) bool {
	return unicode.Is(unicode.Hiragana, ch) || unicode.Is(unicode.Katakana, ch)
}

// Normalize converts Hiragana to Katakana.
func (Kana) Normalize(ch rune) rune {
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

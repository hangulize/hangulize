package scripts

import (
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// Latin represents the Latin or Roman script. Most languages Hangulize
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

// TransliteratePunct does nothing.
func (Latin) TransliteratePunct(punct rune) string {
	return string(punct)
}

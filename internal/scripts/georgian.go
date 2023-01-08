package scripts

import "unicode"

// Georgian represents the Georgian script.
//
//	ასომთავრული
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

// TransliteratePunct does nothing.
func (Georgian) TransliteratePunct(punct rune) string {
	return string(punct)
}

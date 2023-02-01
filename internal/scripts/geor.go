package scripts

import "unicode"

// Geor represents the Georgian script.
//
//	ასომთავრული
type Geor struct{}

// Is checks whether the character is Georgian or not.
func (Geor) Is(ch rune) bool {
	return unicode.Is(unicode.Georgian, ch)
}

// Normalize does nothing. Georgian is unicase, which means, there's only one
// case for each letter.
func (Geor) Normalize(ch rune) rune {
	return ch
}

// LocalizePunct does nothing.
func (Geor) LocalizePunct(punct rune) string {
	return string(punct)
}

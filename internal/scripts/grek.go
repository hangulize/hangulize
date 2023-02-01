package scripts

import "unicode"

// Grek represents the Greek script.
//
//	ελληνικά
type Grek struct{}

// Is checks whether the character is Greek or not.
func (Grek) Is(ch rune) bool {
	return unicode.Is(unicode.Greek, ch)
}

// Normalize converts character into lower case.
func (Grek) Normalize(ch rune) rune {
	return unicode.ToLower(ch)
}

// LocalizePunct does nothing.
func (Grek) LocalizePunct(punct rune) string {
	return string(punct)
}

package scripts

import "unicode"

// Greek represents the Greek script.
//
//	ελληνικά
type Greek struct{}

// Is checks whether the character is Greek or not.
func (Greek) Is(ch rune) bool {
	return unicode.Is(unicode.Greek, ch)
}

// Normalize converts character into lower case.
func (Greek) Normalize(ch rune) rune {
	return unicode.ToLower(ch)
}

// LocalizePunct does nothing.
func (Greek) LocalizePunct(punct rune) string {
	return string(punct)
}

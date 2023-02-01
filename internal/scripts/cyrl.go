package scripts

import "unicode"

// Cyrl represents the Cyrillic script.
//
//	вулкан
type Cyrl struct{}

// Is checks whether the character is Cyrillic or not.
func (Cyrl) Is(ch rune) bool {
	return unicode.Is(unicode.Cyrillic, ch) || unicode.IsMark(ch)
}

// Normalize converts character into lower case.
func (Cyrl) Normalize(ch rune) rune {
	return unicode.ToLower(ch)
}

// LocalizePunct does nothing.
func (Cyrl) LocalizePunct(punct rune) string {
	return string(punct)
}

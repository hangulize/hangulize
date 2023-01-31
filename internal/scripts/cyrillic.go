package scripts

import "unicode"

// Cyrillic represents the Cyrillic script.
//
//	вулкан
type Cyrillic struct{}

// Is checks whether the character is Cyrillic or not.
func (Cyrillic) Is(ch rune) bool {
	return unicode.Is(unicode.Cyrillic, ch) || unicode.IsMark(ch)
}

// Normalize converts character into lower case.
func (Cyrillic) Normalize(ch rune) rune {
	return unicode.ToLower(ch)
}

// LocalizePunct does nothing.
func (Cyrillic) LocalizePunct(punct rune) string {
	return string(punct)
}

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

// TransliteratePunct does nothing.
func (Cyrillic) TransliteratePunct(punct rune) string {
	return string(punct)
}

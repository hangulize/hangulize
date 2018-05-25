package hangulize

import (
	"strings"

	"golang.org/x/text/unicode/norm"
)

type Normalizer interface {
	Normalize(rune) rune
}

// Normalize normalizes various Roman letters into [a-z].  But it keeps
// the letters in except.
func Normalize(word string, norm Normalizer, except []string) string {
	var buf strings.Builder

	// Sort exception letters.
	exceptions := set(except)

	for _, ch := range word {
		if inSet(string(ch), exceptions) {
			buf.WriteRune(ch)
		} else {
			buf.WriteRune(norm.Normalize(ch))
		}
	}

	return buf.String()
}

// RomanNormalizer converts various Roman letters into [a-zA-Z].
type RomanNormalizer struct{}

func (RomanNormalizer) Normalize(ch rune) rune {
	props := norm.NFD.PropertiesString(string(ch))
	bin := props.Decomposition()
	if len(bin) == 0 {
		return ch
	}
	return rune(bin[0])
}

// KanaNormalizer converts Hiragana to Katakana.
type KanaNormalizer struct{}

func (KanaNormalizer) Normalize(ch rune) rune {
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

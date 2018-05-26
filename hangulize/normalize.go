package hangulize

import (
	"strings"

	"golang.org/x/text/unicode/norm"
)

// Normalize converts a word to normal form.  This behavior is called
// "normalization".  It takes a normalizer which normalizes a letter.  It
// doesn't normalize letters in array keep.
func Normalize(word string, norm Normalizer, keep []string) string {
	// Sort letters to keep.
	keepSet := set(keep)

	var buf strings.Builder

	for _, ch := range word {
		if inSet(string(ch), keepSet) {
			buf.WriteRune(ch)
		} else {
			buf.WriteRune(norm.Normalize(ch))
		}
	}

	return buf.String()
}

// Normalizer normalizes a letter.
type Normalizer interface {
	Normalize(rune) rune
}

// -----------------------------------------------------------------------------

// RomanNormalizer is a normalizer for Laion or Roman script.
type RomanNormalizer struct{}

// Normalize converts a Roman letter into ISO basic Latin alphabet [a-zA-Z].
func (RomanNormalizer) Normalize(ch rune) rune {
	props := norm.NFD.PropertiesString(string(ch))
	bin := props.Decomposition()
	if len(bin) == 0 {
		return ch
	}
	return rune(bin[0])
}

// -----------------------------------------------------------------------------

// TODO(sublee): Find out a Kanji to Kana dictionary to hangulize Japanese
// perfectly.

// KanaNormalizer is a normalizer for Kana script which is used in Japan.
type KanaNormalizer struct{}

// Normalize converts Hiragana to Katakana.
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

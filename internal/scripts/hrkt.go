package scripts

import "unicode"

// Hrkt represents the Japanese syllabaries including Hiragana and Katakana.
//
//	ひらがな カタカナ
type Hrkt struct{}

// Is checks whether the character is either Hiragana or Katakana.
func (Hrkt) Is(ch rune) bool {
	return (ch == 'ー' ||
		unicode.Is(unicode.Hiragana, ch) ||
		unicode.Is(unicode.Katakana, ch))
}

// Normalize converts Hiragana to Katakana.
func (Hrkt) Normalize(ch rune) rune {
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

// LocalizePunct converts a Japanese punctuation to fit in Korean.
func (Hrkt) LocalizePunct(punct rune) string {
	switch punct {
	case '。':
		return ". "
	case '、':
		return ", "
	case '：':
		return ": "
	case '！':
		return "! "
	case '？':
		return "? "
	case '〜':
		return "~"
	case '「':
		return " '"
	case '」':
		return "' "
	case '『':
		return " \""
	case '』':
		return "\" "
	}

	return string(punct)
}

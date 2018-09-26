package scripts

import "unicode"

// Kana represents the Kana script including Hiragana and Katakana.
//
//   ひらがな カタカナ
//
type Kana struct{}

// Is checks whether the character is either Hiragana or Katakana.
func (Kana) Is(ch rune) bool {
	return (ch == 'ー' ||
		unicode.Is(unicode.Hiragana, ch) ||
		unicode.Is(unicode.Katakana, ch))
}

// Normalize converts Hiragana to Katakana.
func (Kana) Normalize(ch rune) rune {
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

// TransliteratePunct converts a Japanese punctuation to fit in Korean.
func (Kana) TransliteratePunct(punct rune) string {
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

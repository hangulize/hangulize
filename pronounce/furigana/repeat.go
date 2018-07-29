package furigana

import "bytes"

// repeatKana resolves iteration marks.
func repeatKana(word string) string {
	var buf bytes.Buffer

	lastCh := rune(0)

	for _, ch := range word {
		if isHiragana(lastCh) {
			switch ch {
			case 'ゝ': // HIRAGANA ITERATION MARK
				ch = toSeion(lastCh)
			case 'ゞ': // HIRAGANA VOICED ITERATION MARK
				ch = toDakuon(lastCh)
			}
		} else if isKatakana(lastCh) {
			switch ch {
			case 'ヽ': // KATAKANA ITERATION MARK
				ch = toSeion(lastCh)
			case 'ヾ': // KATAKANA VOICED ITERATION MARK
				ch = toDakuon(lastCh)
			}
		}

		buf.WriteRune(ch)
		lastCh = ch
	}

	return buf.String()
}

func in(ch rune, min rune, max rune) bool {
	return min <= ch && ch <= max
}

func isHiragana(ch rune) bool {
	const (
		hiraganaMin = rune(0x3040)
		hiraganaMax = rune(0x309f)
	)
	return in(ch, hiraganaMin, hiraganaMax)
}

func isKatakana(ch rune) bool {
	const (
		katakanaMin = rune(0x30a0)
		katakanaMax = rune(0x30ff)
	)
	return in(ch, katakanaMin, katakanaMax)
}

// toSeion converts a Dakuon or Hankaduon to the corresponding Seion.
func toSeion(ch rune) rune {
	switch ch {

	case 'ゔ':
		fallthrough
	case 'ヴ':
		return ch - 78

	case 'ゞ':
		fallthrough
	case 'ヾ':
		return ch - 1

	}

	if ch%2 == 0 && (in(ch, 'か', 'ぢ') || in(ch, 'カ', 'ヂ')) {
		return ch - 1
	}
	if ch%2 == 1 && (in(ch, 'つ', 'ど') || in(ch, 'ツ', 'ド')) {
		return ch - 1
	}
	if ch%3 == 1 && (in(ch, 'は', 'ぽ') || in(ch, 'ハ', 'ポ')) {
		return ch - 1
	}
	if ch%3 == 2 && (in(ch, 'は', 'ぽ') || in(ch, 'ハ', 'ポ')) {
		return ch - 2
	}

	if in(ch, 'ヷ', 'ヺ') {
		return ch - 8
	}

	return ch
}

// toDakuon converts a Seion to the corresponding Dakuon.
func toDakuon(ch rune) rune {
	switch ch {

	case 'う':
		fallthrough
	case 'ウ':
		return ch + 78

	case 'ゝ':
		fallthrough
	case 'ヽ':
		return ch + 1

	}

	if ch%2 == 1 && (in(ch, 'か', 'ぢ') || in(ch, 'カ', 'ヂ')) {
		return ch + 1
	}
	if ch%2 == 0 && (in(ch, 'つ', 'ど') || in(ch, 'ツ', 'ド')) {
		return ch + 1
	}
	if ch%3 == 0 && (in(ch, 'は', 'ぽ') || in(ch, 'ハ', 'ポ')) {
		return ch + 1
	}

	if in(ch, 'ワ', 'ヲ') {
		return ch + 8
	}

	return ch
}

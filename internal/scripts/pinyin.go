package scripts

// Pinyin represents the Latin script for Chinese Pinyin.
//
//   漢語拼音
//
type Pinyin struct {
	Latin
}

// Normalize converts a Latin character for Pinyin into ISO basic Latin lower
// alphabet [a-z]. Especially, it converts "ü" to "v":
//
//   lüè -> lve
//
func (s Pinyin) Normalize(ch rune) rune {
	switch ch {
	case 'ü', 'Ü':
		return 'v'
	}
	return s.Latin.Normalize(ch)
}

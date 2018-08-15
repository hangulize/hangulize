package furigana

import (
	"bytes"
	"strings"
)

var longVowelRels = map[rune]string{
	'ァ': "ァアカガサザタダナハバパマャヤラヮワ",
	'ア': "ァアカガサザタダナハバパマャヤラヮワ",
	'ィ': "ィイェエキギケゲシジセゼチヂテデニネヒビピヘベペミメリレヰヱヸヹ",
	'イ': "ィイェエキギケゲシジセゼチヂテデニネヒビピヘベペミメリレヰヱヸヹ",
	'ゥ': "ゥウォオクグコゴスズソゾツヅトドヌノフブプホボポムモュユョヨルロヲヴヺ",
	'ウ': "ゥウォオクグコゴスズソゾツヅトドヌノフブプホボポムモュユョヨルロヲヴヺ",
	'ェ': "ェエケゲセゼテデネヘベペメレヱヹ",
	'エ': "ェエケゲセゼテデネヘベペメレヱヹ",
	'ォ': "ォオコゴソゾトドノホボポモョヨロヲヺ",
	'オ': "ォオコゴソゾトドノホボポモョヨロヲヺ",
}

// mergeLongVowels replaces Katakana long vowels with 'ー' after the given
// offset.
func mergeLongVowels(word string, offset int) string {
	var buf bytes.Buffer
	var prevCh rune

	i := 0
	for _, ch := range word {
		isLongVowel := false

		// Detect if the letter is a long vowel, only after the offset.
		if offset <= i {
			prior, ok := longVowelRels[ch]
			if ok {
				if strings.IndexRune(prior, prevCh) != -1 {
					isLongVowel = true
				}
			}
		}

		if isLongVowel {
			buf.WriteRune('ー')
		} else {
			buf.WriteRune(ch)
		}

		prevCh = ch
		i++
	}

	return buf.String()
}

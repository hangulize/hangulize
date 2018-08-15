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

func mergeLongVowels(word string, from int) string {
	var buf bytes.Buffer
	var prevCh rune

	i := 0
	for _, ch := range word {
		isLongVowel := false

		// Detect if the letter is a long vowel after the from.
		if from <= i {
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

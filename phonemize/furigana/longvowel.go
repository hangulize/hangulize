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

func markLongVowels(word string) string {
	var buf bytes.Buffer
	var prevCh rune

	for _, ch := range word {
		isLongVowel := false

		prior, ok := longVowelRels[ch]
		if ok {
			if strings.IndexRune(prior, prevCh) != -1 {
				isLongVowel = true
			}
		}

		if isLongVowel {
			buf.WriteRune('ー')
		} else {
			buf.WriteRune(ch)
		}

		prevCh = ch
	}

	return buf.String()
}

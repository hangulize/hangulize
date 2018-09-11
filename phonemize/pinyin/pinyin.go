/*
Package pinyin implements the hangulize.Phonemizer interface for Chinese
Hanzu. Hanzu has very broad characters so they need a dictionary to be
converted to a phonogram.
*/
package pinyin

import (
	"bytes"
	"strings"
	"unicode"

	goPinyin "github.com/mozillazg/go-pinyin"
	"golang.org/x/text/unicode/norm"
)

// P is the Pinyin phonemizer.
var P pinyinPhonemizer

// ----------------------------------------------------------------------------

type pinyinPhonemizer struct{}

func (pinyinPhonemizer) ID() string {
	return "pinyin"
}

func (p *pinyinPhonemizer) Phonemize(word string) string {
	var chunks []string
	var buf bytes.Buffer

	a := goPinyin.NewArgs()

	for _, ch := range word {
		// Normalize into CJK unified ideographs.
		if unicode.Is(unicode.Han, ch) {
			ch = []rune(norm.NFC.String(string(ch)))[0]
		}

		// Pick Pinyin.
		pyn := goPinyin.SinglePinyin(ch, a)

		if len(pyn) == 0 {
			buf.WriteRune(ch)
		} else {
			if buf.Len() != 0 {
				chunks = append(chunks, buf.String())
				buf.Reset()
			}
			chunks = append(chunks, pyn[0])
		}
	}
	if buf.Len() != 0 {
		chunks = append(chunks, buf.String())
	}

	return strings.Join(chunks, "\u200b")
}

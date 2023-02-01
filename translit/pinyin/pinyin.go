/*
Package pinyin implements the hangulize.Translit interface for Chinese Hanzu.
Hanzu has very broad characters so they need a dictionary to be converted to a
phonogram.
*/
package pinyin

import (
	"bytes"
	"strings"

	"github.com/hangulize/hangulize"
	goPinyin "github.com/mozillazg/go-pinyin"
	"golang.org/x/text/unicode/norm"
)

// T is a hangulize.Translit for Pinyin.
var T hangulize.Translit = &pinyin{}

// ----------------------------------------------------------------------------

type pinyin struct{}

func (pinyin) Scheme() string {
	return "pinyin"
}

func (p *pinyin) Transliterate(word string) (string, error) {
	// Normalize into CJK unified ideographs.
	word = norm.NFC.String(word)

	// Pick Pinyin.
	var chunks []string
	var buf bytes.Buffer

	a := goPinyin.NewArgs()

	for _, ch := range word {
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

	// U+200B: Zero Width Space
	return strings.Join(chunks, "\u200b"), nil
}

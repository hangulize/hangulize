/*
Package pinyin implements the hangulize.Pronouncer interface for Chinese
Hanzu. Hanzu has very broad characters so they need a dictionary to be
pronounced.
*/
package pinyin

import (
	"bytes"
	"strings"

	goPinyin "github.com/mozillazg/go-pinyin"
)

// P is the Pinyin pronouncer.
var P pinyinPronouncer

// ----------------------------------------------------------------------------

type pinyinPronouncer struct{}

func (pinyinPronouncer) ID() string {
	return "pinyin"
}

func (p *pinyinPronouncer) Pronounce(word string) string {
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

	return strings.Join(chunks, " ")
}

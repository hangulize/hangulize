/*
Package furigana implements the hangulize.Translit interface for Japanese
Kanji. Kanji has very broad characters so they need a dictionary to be
converted to Kana. This Translit uses IPADIC in Kagome to analyze Kanji.
*/
package furigana

import (
	"github.com/hangulize/hangulize"
	kagome "github.com/ikawaha/kagome.ipadic/tokenizer"
	"golang.org/x/text/unicode/norm"
)

// T is a hangulize.Translit for Furigana.
var T hangulize.Translit = &furigana{}

// ----------------------------------------------------------------------------

type furigana struct {
	kagome *kagome.Tokenizer
}

func (furigana) Scheme() string {
	return "furigana"
}

// ensureKagome caches a Kagome tokenizer because it is expensive.
func (p *furigana) ensureKagome() *kagome.Tokenizer {
	if p.kagome == nil {
		// It may take a while.
		k := kagome.New()
		p.kagome = &k
	}
	return p.kagome
}

func (p *furigana) Transliterate(word string) (string, error) {
	// Normalize into CJK unified ideographs.
	word = norm.NFC.String(word)

	// Resolve Kana repeatations.
	word = repeatKana(word)

	k := p.ensureKagome()
	tokens := k.Tokenize(word)

	tw := newTypewriter(tokens)
	word = tw.Typewrite()

	return word, nil
}

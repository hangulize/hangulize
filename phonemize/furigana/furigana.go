/*
Package furigana implements the hangulize.Phonemizer interface for Japanese
Kanji. Kanji has very broad characters so they need a dictionary to be
converted to Kana. This phonemizer uses IPADIC in Kagome to analyze Kanji.
*/
package furigana

import (
	kagome "github.com/ikawaha/kagome.ipadic/tokenizer"
)

// P is the Furigana phonemizer.
var P furiganaPhonemizer

// ----------------------------------------------------------------------------

type furiganaPhonemizer struct {
	kagome *kagome.Tokenizer
}

func (furiganaPhonemizer) ID() string {
	return "furigana"
}

// Kagome caches d Kagome tokenizer because it is expensive.
func (p *furiganaPhonemizer) Kagome() *kagome.Tokenizer {
	if p.kagome == nil {
		k := kagome.New()
		p.kagome = &k
	}
	return p.kagome
}

func (p *furiganaPhonemizer) Phonemize(word string) string {
	word = repeatKana(word)

	tokens := p.Kagome().Tokenize(word)
	tw := newTypewriter(tokens)
	word = tw.Typewrite()

	return word
}

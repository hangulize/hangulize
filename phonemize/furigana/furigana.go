/*
Package furigana implements the hangulize.Phonemizer interface for Japanese
Kanji. Kanji has very broad characters so they need a dictionary to be
converted to Kana. This phonemizer uses IPADIC in Kagome to analyze Kanji.
*/
package furigana

import (
	kagome "github.com/ikawaha/kagome.ipadic/tokenizer"
	"golang.org/x/text/unicode/norm"
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

func (p *furiganaPhonemizer) Load() error {
	if p.kagome == nil {
		k := kagome.New()
		p.kagome = &k
	}
	return nil
}

// Kagome caches a Kagome tokenizer because it is expensive.
func (p *furiganaPhonemizer) Kagome() *kagome.Tokenizer {
	return p.kagome
}

func (p *furiganaPhonemizer) Phonemize(word string) (string, error) {
	// Normalize into CJK unified ideographs.
	word = norm.NFC.String(word)

	// Resolve Kana repeatations.
	word = repeatKana(word)

	_ = p.Load()
	tokens := p.kagome.Tokenize(word)

	tw := newTypewriter(tokens)
	word = tw.Typewrite()

	return word, nil
}

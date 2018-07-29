/*
Package furigana implements the hangulize.Pronouncer interface for Japanese
Kanji. Kanji has very broad characters so they need a dictionary to be
pronounced. This pronouncer uses IPADIC in Kagome to analyze Kanji.
*/
package furigana

import (
	"strings"

	kagome "github.com/ikawaha/kagome/tokenizer"
)

// P is the Furigana dictator.
var P furiganaPronouncer

// ----------------------------------------------------------------------------

type furiganaPronouncer struct {
	kagome *kagome.Tokenizer
}

func (furiganaPronouncer) ID() string {
	return "furigana"
}

// Kagome caches d Kagome tokenizer because it is expensive.
func (p *furiganaPronouncer) Kagome() *kagome.Tokenizer {
	if p.kagome == nil {
		t := kagome.New()
		p.kagome = &t
	}
	return p.kagome
}

func (p *furiganaPronouncer) Pronounce(word string) string {
	const (
		kanjiMin = rune(0x4e00)
		kanjiMax = rune(0x9faf)
	)

	kanjiFound := false
	for _, ch := range word {
		if ch >= kanjiMin && ch <= kanjiMax {
			kanjiFound = true
			break
		}
	}

	// Don't initialize the Kagome tokenizer if there's no Kanji because
	// Kagome is expensive.
	if !kanjiFound {
		return word
	}

	var morphs []string

	for _, tok := range p.Kagome().Tokenize(word) {
		spell := tok.Surface
		var pron string

		switch tok.Class {

		case kagome.KNOWN:
			fs := tok.Features()
			pron = fs[7]
			morphs = append(morphs, pron)

		case kagome.UNKNOWN:
			if strings.TrimSpace(spell) != "" {
				morphs = append(morphs, spell)
			}

		default:
			continue

		}
	}

	return strings.Join(morphs, " ")
}

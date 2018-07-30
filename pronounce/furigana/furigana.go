/*
Package furigana implements the hangulize.Pronouncer interface for Japanese
Kanji. Kanji has very broad characters so they need a dictionary to be
pronounced. This pronouncer uses IPADIC in Kagome to analyze Kanji.
*/
package furigana

import (
	"bytes"
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
	if kanjiFound {
		word = p.analyze(word)
	}

	word = repeatKana(word)
	return word
}

type chunk struct {
	s   string
	sep bool
}

func (p *furiganaPronouncer) analyze(word string) string {
	var chunks []chunk

	prevWasSep := false

	for _, tok := range p.Kagome().Tokenize(word) {
		switch tok.Class {

		case kagome.KNOWN:
			fs := tok.Features()
			// 0: part-of-speech
			// 1: sub-class 1
			// 2: sub-class 2
			// 3: sub-class 3
			// 4: inflection
			// 5: conjugation
			// 6: root-form
			// 7: reading
			// 8: pronunciation

			class := fs[1]
			reading := fs[7]

			c := chunk{reading, !prevWasSep && (class == "固有名詞")}
			chunks = append(chunks, c)

		case kagome.UNKNOWN:
			surf := tok.Surface
			isSpace := strings.TrimSpace(surf) == ""

			var c *chunk
			if isSpace {
				// Whitespace is used as the sep.
				prevWasSep = true
				c = &chunk{surf, false}
			} else {
				c = &chunk{surf, !prevWasSep}
			}
			chunks = append(chunks, *c)

		default:
			continue

		}
	}

	var buf bytes.Buffer

	skipSep := true
	for _, c := range chunks {
		if c.sep {
			if skipSep {
				skipSep = false
			} else {
				// Separate this chunk by Nakaguro.
				buf.WriteRune('・')
			}
		}
		buf.WriteString(c.s)
	}

	return buf.String()
}

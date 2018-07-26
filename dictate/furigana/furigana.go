package furigana

import (
	"strings"

	kagome "github.com/ikawaha/kagome/tokenizer"
)

// D is the Furigana dictator.
var D furiganaDictator

// ----------------------------------------------------------------------------

type furiganaDictator struct {
	kagome *kagome.Tokenizer
}

func (furiganaDictator) ID() string {
	return "furigana"
}

// Kagome caches d Kagome tokenizer because it is expensive.
func (d *furiganaDictator) Kagome() *kagome.Tokenizer {
	if d.kagome == nil {
		t := kagome.New()
		d.kagome = &t
	}
	return d.kagome
}

func (d *furiganaDictator) Dictate(word string) [][2]string {
	const (
		furiganaMin = rune(0x4e00)
		furiganaMax = rune(0x9faf)
	)

	furiganaFound := false
	for _, ch := range word {
		if ch >= furiganaMin && ch <= furiganaMax {
			furiganaFound = true
			break
		}
	}

	// Don't initialize the Kagome tokenizer if there's no furigana.
	if !furiganaFound {
		return [][2]string{
			[2]string{word, word},
		}
	}

	tokens := d.Kagome().Tokenize(word)

	lexemes := make([][2]string, 0)

	for _, tok := range tokens {
		spell := tok.Surface
		var pron string

		switch tok.Class {

		case kagome.KNOWN:
			fs := tok.Features()
			pron = fs[7]

		case kagome.UNKNOWN:
			pron = spell

		default:
			continue

		}

		// Separate each lexeme by a space.
		pron += " "

		lexemes = append(lexemes, [2]string{spell, pron})
	}

	// Trim the last space.
	i := len(lexemes) - 1
	final := lexemes[i]
	final[1] = strings.TrimRight(final[1], " ")
	lexemes = append(lexemes[:i], final)

	return lexemes
}

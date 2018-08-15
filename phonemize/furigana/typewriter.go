package furigana

import (
	"bytes"
	"strings"

	kagome "github.com/ikawaha/kagome.ipadic/tokenizer"
)

type category int

const (
	illegal category = iota
	space
	filler
	punct
	morpheme
	properNoun
	personName
	unknownText
)

type typewriter struct {
	tokens  []kagome.Token
	cur     int
	lastCat category
}

func newTypewriter(tokens []kagome.Token) *typewriter {
	return &typewriter{tokens, -1, illegal}
}

func (t *typewriter) Typewrite() string {
	var buf bytes.Buffer

	for {
		str, cat, prevCat := t.read()

		if cat == illegal {
			break
		}

		// Prevent a redundant separator.
		switch prevCat {
		case space, filler, punct, illegal:
			buf.WriteString(str)
			continue
		}

		// Split between a first name and a last name.
		sep := ""
		if cat == personName && prevCat == personName {
			sep = " "
		}

		// Write the separator and string.
		buf.WriteString(sep)
		buf.WriteString(str)
	}

	return buf.String()
}

func (t *typewriter) read() (string, category, category) {
	var tok kagome.Token

	// Scan the next non-dummy token.
	for tok.Class == kagome.DUMMY {
		t.cur++

		if t.cur >= len(t.tokens) {
			return "", illegal, illegal
		}

		tok = t.tokens[t.cur]
	}

	str, cat := interpretToken(&tok)

	prevCat := t.lastCat
	t.lastCat = cat

	return str, cat, prevCat
}

func interpretToken(tok *kagome.Token) (string, category) {
	str := tok.Surface
	cat := unknownText

	if tok.Class == kagome.KNOWN {
		// 0: part-of-speech
		// 1: sub-class 1
		// 2: sub-class 2
		// 3: sub-class 3
		// 4: inflection
		// 5: conjugation
		// 6: root-form
		// 7: reading
		// 8: pronunciation
		fs := tok.Features()
		var (
			partOfSpeech  = fs[0]
			subClass1     = fs[1]
			subClass2     = fs[2]
			rootForm      = fs[6]
			pronunciation = fs[8]
		)

		str = pronunciation
		cat = morpheme

		switch partOfSpeech {

		case "フィラー":
			cat = filler

		case "記号":
			cat = punct

		case "助詞":
			// Keep the root form of particles.
			switch rootForm {
			case "は":
				str = "ハ"
			case "へ":
				str = "ヘ"
			}

		default:
			if subClass2 == "人名" {
				cat = personName
			} else if subClass1 == "固有名詞" {
				cat = properNoun
			}

		}
	} else {
		isSpace := strings.TrimSpace(str) == ""
		if isSpace {
			cat = space
		}
	}

	return str, cat
}

package furigana

import (
	"bytes"
	"strings"
	"unicode/utf8"

	kagome "github.com/ikawaha/kagome.ipadic/tokenizer"
)

type category int

// Categories in a typewriter.
const (
	illegal category = iota
	meta
	morpheme
	auxiliary
	properNoun
	personName
	unknown
)

// typewriter writes a whole pronunciation from the Kagome tokens.
type typewriter struct {
	tokens  []kagome.Token
	result  string
	cur     int
	lastCat category
}

// newTypewriter initializes a typewriter for the Kagome tokens.
func newTypewriter(tokens []kagome.Token) *typewriter {
	return &typewriter{tokens, "", -1, illegal}
}

// Typewrite returns a whole pronunciation from the Kagome tokens.
func (t *typewriter) Typewrite() string {
	// Re-use the cached result if already processed.
	if t.cur != -1 {
		return t.result
	}

	var buf bytes.Buffer

	for {
		sep, str := t.scanMorpheme()
		if str == "" {
			break
		}

		buf.WriteString(sep)
		buf.WriteString(str)
	}

	t.result = buf.String()
	return t.result
}

// scanMorpheme consumes the Kagome tokens one by one.
func (t *typewriter) scanMorpheme() (sep string, str string) {
	var buf bytes.Buffer

	// -------------------------------------------------------------------------
	// 1. The Core Morpheme

	tok := t.read()
	if tok == nil {
		return
	}

	str, cat := interpretToken(tok)

	// Merge long vowels in an unknown word. Because Kagome didn't detect the
	// pronunciation of this word.
	if cat == unknown {
		str = mergeLongVowels(str, 0)
	}

	if t.lastCat == meta || t.lastCat == illegal {
		// If here's a head of a word, any separator not required.
		sep = ""

	} else if cat == personName && t.lastCat == personName {
		// Split between a first name and a last name.
		sep = " "
	}

	// Remember this category.
	t.lastCat = cat

	// Keep the length of the core morpheme.
	coreLen := utf8.RuneCountInString(str)

	buf.WriteString(str)

	// -------------------------------------------------------------------------
	// 2. Following Auxiliary Morphemes

	// If the next tokens are auxiliary morphemes, merge them.
	for {
		tok := t.read()
		if tok == nil {
			break
		}

		str, cat := interpretToken(tok)
		if cat == auxiliary {
			buf.WriteString(str)
		} else {
			t.unread()
			break
		}
	}

	// Merge long vowels in auxiliary morphemes.
	str = buf.String()
	str = mergeLongVowels(str, coreLen)

	return sep, str
}

func (t *typewriter) read() *kagome.Token {
	var tok *kagome.Token

	// Scan the next non-dummy token.
	for tok == nil || tok.Class == kagome.DUMMY {
		t.cur++

		if t.cur >= len(t.tokens) {
			return nil
		}

		tok = &t.tokens[t.cur]
	}

	return tok
}

func (t *typewriter) unread() {
	t.cur--
}

// interpretToken picks a pronunciation and category from a Kagome token.
func interpretToken(tok *kagome.Token) (string, category) {
	str := tok.Surface
	cat := unknown

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

		case "フィラー", "記号":
			cat = meta

		case "助動詞":
			cat = auxiliary

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
			cat = meta
		}
	}

	return str, cat
}

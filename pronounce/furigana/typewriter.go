package furigana

import (
	"bytes"
	"strings"
	"unicode/utf8"

	kagome "github.com/ikawaha/kagome/tokenizer"
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

		// Choose a proper separator.
		sep := ""
		switch cat {

		// Insert a Nakaguro before a proper noun.
		case properNoun:
			sep = "・"

		// Insert a space between first/last names.
		case personName:
			if prevCat == personName {
				sep = " "
			}

		// Insert a Nakaguro before a morphem which starts with a vowel.
		case morpheme:
			ch, _ := utf8.DecodeRuneInString(str)
			startsWithVowel := ch%2 == 0 && in(ch, 'ア', 'オ')
			if startsWithVowel {
				sep = "・"
			}
		}

		// Write the separator and string.
		buf.WriteString(sep)
		buf.WriteString(str)
	}

	return buf.String()
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
			partOfSpeech = fs[0]
			subClass1    = fs[1]
			subClass2    = fs[2]
			reading      = fs[7]
		)

		str = reading

		if partOfSpeech == "フィラー" {
			cat = filler
		} else if partOfSpeech == "記号" {
			cat = punct
		} else if subClass2 == "人名" {
			cat = personName
		} else if subClass1 == "固有名詞" {
			cat = properNoun
		} else {
			cat = morpheme
		}
	} else {
		isSpace := strings.TrimSpace(str) == ""
		if isSpace {
			cat = space
		}
	}

	return str, cat
}

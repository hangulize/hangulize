package hangulize

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/suapapa/go_hangul"
)

const none = rune(0)

// ComposeHangul converts decomposed Jamo phonemes to composed Hangul
// syllables.
//
// Decomposed Jamo phonemes look like "ㅎㅏ-ㄴㄱㅡ-ㄹㄹㅏㅇㅣㅈㅡ". A Jaeum
// after a hyphen ("-ㄴ") means that it is a Jongseong (tail).
//
func ComposeHangul(word string) string {
	r := bufio.NewReader(strings.NewReader(word))
	var buf bytes.Buffer

	var lmt [3]rune // [lead, medial, tail]
	const (
		lead   = 0
		medial = 1
		tail   = 2
	)

	prevScore := -1
	score := -1

	isTail := false

	writeLetter := func() {
		if lmt[0] == none && lmt[1] == none && lmt[2] == none {
			return
		}

		// Fill missing Jamo.
		if lmt[0] == none {
			lmt[0] = 'ㅇ'
		}
		if lmt[1] == none {
			lmt[1] = 'ㅡ'
		}

		// Complete a letter.
		letter := hangul.Join(lmt[0], lmt[1], lmt[2])
		buf.WriteRune(letter)

		// Clear.
		lmt[0], lmt[1], lmt[2] = none, none, none
	}

	for {
		ch, _, err := r.ReadRune()

		if err != nil {
			break
		}

		// Hyphen is the prefix of a tail Jaeum.
		// Perhaps the next ch is a Jaeum.
		if ch == '-' {
			isTail = true
			continue
		}

		if !hangul.IsHangul(ch) {
			if prevScore != -1 {
				writeLetter()
			}

			buf.WriteRune(ch)
			prevScore = -1
			continue
		}

		isJaeum := hangul.IsJaeum(ch)
		isMoeum := hangul.IsMoeum(ch)

		if !isJaeum && !isMoeum {
			// Composed Hangul.
			writeLetter()

			lmt[0], lmt[1], lmt[2] = hangul.Split(ch)

			if lmt[2] == none {
				score = medial
			} else {
				score = tail
			}
		} else {
			// Decomposed Jamo.
			switch true {

			case isJaeum:
				if isTail {
					score = tail
				} else {
					score = lead
				}

			case isMoeum:
				score = medial

			}

			// Write a letter.
			if score <= prevScore {
				writeLetter()
			}

			if score != -1 {
				lmt[score] = ch
			}
		}

		prevScore = score
		isTail = false
	}

	// Write the final letter.
	if prevScore != -1 {
		writeLetter()
	}

	return buf.String()
}

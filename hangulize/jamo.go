package hangulize

import (
	"bufio"
	"strings"

	"github.com/suapapa/go_hangul"
)

const none = rune(0)

// AssembleJamo makes a complete Hangul word from a decomposed Jamo word.
//
// A decomposed Jamo word looks like "ㅎㅏ-ㄴㄱㅡ-ㄹㄹㅏㅇㅣㅈㅡ".  A Jaeum
// after a hyphen becomes a tail (Moeum).
//
func AssembleJamo(word string) string {
	r := bufio.NewReader(strings.NewReader(word))
	var buf strings.Builder

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

		// Classify Jamo.
		if hangul.IsJaeum(ch) {
			if isTail {
				score = tail
			} else {
				score = lead
			}
		} else if hangul.IsMoeum(ch) {
			score = medial
		}

		// Write a letter.
		if score <= prevScore {
			writeLetter()
			prevScore = -1
		}

		if score != -1 {
			lmt[score] = ch
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

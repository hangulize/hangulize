/*
Package jamo implements a Hangul composer. It converts decomposed Jamo phonemes
to composed Hangul syllables.

	fmt.Println(jamo.ComposeHangul("ㅈㅏㅁㅗ"))
	// Output: 자모
*/
package jamo

import (
	"bufio"
	"bytes"
	"strings"

	hangul "github.com/suapapa/go_hangul"
)

// ComposeHangul converts decomposed Jamo phonemes to composed Hangul
// syllables.
//
// Decomposed Jamo phonemes look like "ㅎㅏ-ㄴㄱㅡ-ㄹㄹㅏㅇㅣㅈㅡ". A Jaeum
// after a hyphen ("-ㄴ") means that it is a Jongseong (tail).
func ComposeHangul(word string) string {
	c := composer{
		r: bufio.NewReader(strings.NewReader(word)),
	}
	return c.Compose()
}

const (
	lead   = 0
	medial = 1
	tail   = 2
)

// composer is a state machine which converts decomposed Jamo phonemes to
// composed Hangul syllables.
type composer struct {
	r   *bufio.Reader
	buf bytes.Buffer // The output buffer.
	lmt [3]rune      // Buffered Jamos. [lead, medial. tail]
}

// read consumes 1 character. If the character is a tail Jamo, the second bool
// return value will be set as true.
func (c *composer) read() (rune, bool, error) {
	isTail := false

	for {
		ch, _, err := c.r.ReadRune()

		if err != nil {
			return 0, false, err
		}

		// Hyphen is the prefix of a tail Jaeum.
		// Perhaps the next ch is a Jaeum.
		if ch == '-' {
			isTail = true
			continue
		}

		return ch, isTail, nil
	}
}

// write writes a composed Hangul from the buffered Jamos into the output
// buffer.
func (c *composer) write() {
	if c.lmt == [3]rune{} {
		return
	}

	// Fill missing Jamo.
	if c.lmt[lead] == 0 {
		c.lmt[lead] = 'ㅇ'
	}
	if c.lmt[medial] == 0 {
		c.lmt[medial] = 'ㅡ'
	}

	// Complete a letter.
	letter := hangul.Join(c.lmt[lead], c.lmt[medial], c.lmt[tail])
	c.buf.WriteRune(letter)

	// Clear.
	c.lmt = [3]rune{}
}

// Compose converts decomposed Jamo phonemes to composed Hangul syllables.
func (c *composer) Compose() string {
	var (
		ch     rune
		isTail bool
		err    error
	)

	var isHangul, isMoeum, isComposed bool

	// Score values can be -1 for non-Hangul,
	// 0 for leads, 1 for medials, and 2 for tails.
	var score, prevScore int

	for {
		prevScore = score

		ch, isTail, err = c.read()

		if err != nil {
			break
		}

		isHangul, _, isMoeum, isComposed = analyzeHangul(ch)

		// Non-Hangul
		if !isHangul {
			if prevScore != -1 {
				c.write()
			}

			c.buf.WriteRune(ch)
			prevScore = -1

			continue
		}

		// Composed Hangul
		if isComposed {
			c.write()

			// Decompose it to merge with a tail later.
			c.lmt[lead], c.lmt[medial], c.lmt[tail] = hangul.Split(ch)

			if c.lmt[tail] == 0 {
				score = medial
			} else {
				score = tail
			}

			continue
		}

		// Decomposed Jamo
		if isMoeum {
			score = medial
		} else if isTail {
			score = tail
		} else {
			score = lead
		}

		// If cursor should be moved forward, flush the buffered letter.
		if score <= prevScore {
			c.write()
		}

		// Buffer the Jamo.
		if score != -1 {
			c.lmt[score] = ch
		}
	}

	// Write the final letter.
	if prevScore != -1 {
		c.write()
	}

	return c.buf.String()
}

// analyzeHangul analyzes a Hangul character to check if it is a Jaeum, a
// Moeum, or a composed Hangul.
func analyzeHangul(ch rune) (isHangul, isJaeum, isMoeum, isComposed bool) {
	isHangul = hangul.IsHangul(ch)
	if !isHangul {
		return
	}
	isJaeum = hangul.IsJaeum(ch)
	isMoeum = hangul.IsMoeum(ch)
	isComposed = !isJaeum && !isMoeum
	return
}

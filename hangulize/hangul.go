package hangulize

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/suapapa/go_hangul"
)

const null = rune(0)

func CompleteHangul(jamo string) string {
	r := bufio.NewReader(strings.NewReader(jamo))

	var buf strings.Builder
	var lmt [3]rune // [lead, medial, tail]

	isTail := false
	dirty := false

	flush := func() {
		if dirty {
			// Flush
			l, m, t := lmt[0], lmt[1], lmt[2]

			if l == null {
				l = 'ㅇ'
			}
			if m == null {
				m = 'ㅡ'
			}

			lmt[0], lmt[1], lmt[2] = null, null, null

			letter := hangul.Join(l, m, t)

			fmt.Println("l, m, t:", string(l), string(m), string(t))
			fmt.Println("letter:", string(letter))

			buf.WriteRune(letter)
		}
	}

	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			break
		}

		fmt.Println("ch:", string(ch))

		// Hyphen is the prefix of a tail Jaeum.
		// Perhaps the next ch is a Jaeum.
		if ch == '-' {
			isTail = true
			continue
		}

		if hangul.IsJaeum(ch) {
			if isTail {
				lmt[2] = ch
				flush()
				dirty = false
			} else {
				flush()

				lmt[0] = ch
				dirty = true
			}
		}
		if hangul.IsMoeum(ch) {
			if lmt[1] != null {
				flush()
			}

			lmt[1] = ch
			dirty = true
		}

		isTail = false
	}

	flush()

	return buf.String()
}

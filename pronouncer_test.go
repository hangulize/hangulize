package hangulize

import "github.com/hangulize/hangulize/pronounce/furigana"

func init() {
	// Use all pronouncers automatically for test.
	UsePronouncer(&furigana.P)
}

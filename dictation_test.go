package hangulize

import "github.com/hangulize/hangulize/dictate/furigana"

func init() {
	// Use all dictators in test.
	UseDictator(&furigana.D)
}

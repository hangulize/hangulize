/*
Package pronounce registers all standard pronouncers.
*/
package pronounce

import (
	"github.com/hangulize/hangulize"

	"github.com/hangulize/hangulize/pronounce/furigana"
)

func init() {
	hangulize.UsePronouncer(&furigana.P)
}

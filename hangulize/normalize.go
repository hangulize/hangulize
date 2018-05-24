package hangulize

import (
	"strings"

	"golang.org/x/text/unicode/norm"
)

func NormalizeRoman(word string) string {
	var buf strings.Builder

	var iter norm.Iter
	iter.InitString(norm.NFD, word)

	for !iter.Done() {
		bin := iter.Next()
		buf.WriteByte(bin[0])
	}

	return strings.ToLower(buf.String())
}

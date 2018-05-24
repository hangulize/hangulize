package hangulize

import (
	"strings"

	"golang.org/x/text/unicode/norm"
)

func NormalizeRoman(word string) string {
	var buf strings.Builder

	text := []rune(word)

	var iter norm.Iter
	iter.InitString(norm.NFD, word)

	i := 0
	for !iter.Done() {
		bin := iter.Next()

		// POC(sublee): Spanish
		switch string(text[i]) {
		case "Ñ":
			buf.WriteString("ñ")
		case "ñ":
			buf.WriteString("ñ")
		case "Ǘ":
			buf.WriteString("ü")
		case "ü":
			buf.WriteString("ü")
		case "Ü":
			buf.WriteString("ü")
		default:
			buf.WriteByte(bin[0])
		}

		i++
	}

	return strings.ToLower(buf.String())
}

package hangulize

import (
	"strings"

	"golang.org/x/text/unicode/norm"
)

// NormalizeRoman normalizes various Roman letters into [a-z].  But it keeps
// the letters in except.
func NormalizeRoman(word string, except []string) string {
	var buf strings.Builder

	// Sort exception letters.
	exceptions := set(except)

	// Normalize forms based on Unicode.
	var iter norm.Iter
	iter.InitString(norm.NFD, word)
	text := []rune(word)

	i := 0
	for !iter.Done() {
		bin := iter.Next()
		letter := string(text[i])

		isException := inSet(letter, exceptions)

		if isException {
			buf.WriteString(letter)
		} else {
			buf.WriteByte(bin[0])
		}

		i++
	}

	return buf.String()
}

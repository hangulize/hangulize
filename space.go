package hangulize

import (
	"strings"
	"unicode"
)

// hasSpace tests if the given word contains a space at least once.
func hasSpace(word string) bool {
	i := strings.IndexFunc(word, unicode.IsSpace)
	return i != -1
}

// hasSpaceOnly tests if the given word does not contain any letter except
// spaces.
func hasSpaceOnly(word string) bool {
	if word == "" {
		return false
	}
	for _, r := range word {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

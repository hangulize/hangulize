package hangulize

import (
	"fmt"
	"sort"
)

type stringSet map[string]bool

func (s *stringSet) String() string {
	return fmt.Sprint(s.Array())
}

// newStringSet creates a stringSet from the given strings.
// Duplicated string doesn't occur a failure.
func newStringSet(strs ...string) stringSet {
	set := make(stringSet, len(strs))
	for _, str := range strs {
		set[str] = true
	}
	return set
}

// Array returns a []string array containing strings in the set.
// Each string is unique and ordered in ascending order.
func (s *stringSet) Array() []string {
	strings := make([]string, len(*s))

	i := 0
	for str := range *s {
		strings[i] = str
		i++
	}

	sort.Strings(strings)
	return strings
}

// Has tests if the string is in the set.
func (s *stringSet) Has(str string) bool {
	return (*s)[str]
}

// HasRune tests if the rune is in the set.
func (s *stringSet) HasRune(ch rune) bool {
	return s.Has(string(ch))
}

/*
Package stringset implements a set of strings. StringSet has very simple
implementation. But it should be commonly unsed in both Hangulize and HRE. So
packed as a separate package.
*/
package stringset

import (
	"fmt"
	"sort"
)

// StringSet is a set of strings.
type StringSet map[string]bool

func (s *StringSet) String() string {
	return fmt.Sprint(s.Array())
}

// NewStringSet creates a stringSet from the given strings.
// Duplicated string doesn't occur a failure.
func NewStringSet(strs ...string) StringSet {
	set := make(StringSet, len(strs))
	for _, str := range strs {
		set[str] = true
	}
	return set
}

// Array returns a []string array containing strings in the set.
// Each string is unique and ordered in ascending order.
func (s *StringSet) Array() []string {
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
func (s *StringSet) Has(str string) bool {
	return (*s)[str]
}

// HasRune tests if the rune is in the set.
func (s *StringSet) HasRune(ch rune) bool {
	return s.Has(string(ch))
}

// Add inserts the string into the set.
func (s *StringSet) Add(str string) {
	(*s)[str] = true
}

// AddRune inserts the rune into the set.
func (s *StringSet) AddRune(ch rune) {
	s.Add(string(ch))
}

// Discard removes the string from the set.
func (s *StringSet) Discard(str string) bool {
	exists := (*s)[str]
	delete(*s, str)
	return exists
}

// DiscardRune removes the rune from the set.
func (s *StringSet) DiscardRune(ch rune) bool {
	return s.Discard(string(ch))
}

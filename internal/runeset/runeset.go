package runeset

import (
	"github.com/zyedidia/generic"
	"github.com/zyedidia/generic/hashset"
)

var (
	equals = generic.Equals[rune]
	hash   = generic.HashInt32
)

// Set is a rune set.
type Set struct {
	*hashset.Set[rune]
}

// New makes an empty set.
func New(capacity int) Set {
	return Set{hashset.New(uint64(capacity), equals, hash)}
}

// New makes a new set with the given values.
func Of(values ...rune) Set {
	return Set{hashset.Of(uint64(len(values)), equals, hash, values...)}
}

// Slice makes a slice containing the values in the set.
func (s Set) Slice() []rune {
	slice := make([]rune, 0, s.Size())
	s.Each(func(v rune) {
		slice = append(slice, v)
	})
	return slice
}

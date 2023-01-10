package runeset_test

import (
	"sort"
	"testing"

	"github.com/hangulize/hangulize/internal/runeset"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	s := runeset.New(1)
	assert.Equal(t, 0, s.Size())
}

func TestOf(t *testing.T) {
	s := runeset.Of(rune('a'), rune('b'), rune('c'))
	assert.Equal(t, 3, s.Size())
}

func TestSlice(t *testing.T) {
	s := runeset.Of(rune('a'), rune('b'), rune('c'))
	slice := s.Slice()
	sort.Slice(slice, func(i, j int) bool { return slice[i] < slice[j] })
	assert.Equal(t, []rune{'a', 'b', 'c'}, slice)
}

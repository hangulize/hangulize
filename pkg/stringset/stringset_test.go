package stringset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringSetUniqueness(t *testing.T) {
	s1 := NewStringSet("a", "a", "a")
	assert.Len(t, s1, 1)

	s2 := NewStringSet("a", "b", "c")
	assert.Len(t, s2, 3)
}

func TestStringSetArray(t *testing.T) {
	s := NewStringSet("z", "f", "a", "f")
	assert.Equal(t, []string{"a", "f", "z"}, s.Array())
}

func TestStringSetHas(t *testing.T) {
	s := NewStringSet("z", "a", "f")
	assert.True(t, s.Has("a"))
	assert.True(t, s.Has("z"))
	assert.False(t, s.Has("NOT_EXISTS"))
}

func TestStringSetAdd(t *testing.T) {
	s := NewStringSet()
	assert.Equal(t, []string{}, s.Array())

	s.Add("hello")
	assert.Equal(t, []string{"hello"}, s.Array())

	s.Add("world")
	assert.Equal(t, []string{"hello", "world"}, s.Array())
}

func TestStringSetDiscard(t *testing.T) {
	s := NewStringSet("a", "b", "c")
	assert.Equal(t, []string{"a", "b", "c"}, s.Array())

	assert.True(t, s.Discard("a"))
	assert.Equal(t, []string{"b", "c"}, s.Array())

	assert.False(t, s.Discard("NOT_EXISTS"))
	assert.Equal(t, []string{"b", "c"}, s.Array())
}

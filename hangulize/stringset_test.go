package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringSetUniqueness(t *testing.T) {
	s1 := newStringSet("a", "a", "a")
	assert.Len(t, s1, 1)

	s2 := newStringSet("a", "b", "c")
	assert.Len(t, s2, 3)
}

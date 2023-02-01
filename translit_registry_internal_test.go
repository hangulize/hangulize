package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type asisTranslit struct{}

func (asisTranslit) Scheme() string {
	return "asis"
}

func (asisTranslit) Transliterate(word string) (string, error) {
	return word, nil
}

func TestAdd(t *testing.T) {
	var ok bool

	registry := make(translitRegistry)

	// Successfully registered.
	ok = registry.Add(&asisTranslit{})
	assert.True(t, ok)

	// Already exists.
	ok = registry.Add(&asisTranslit{})
	assert.False(t, ok)

	// Successfully deregistered.
	ok = registry.Remove("asis")
	assert.True(t, ok)

	// Not exists.
	ok = registry.Remove("asis")
	assert.False(t, ok)
}

package hangulize_test

import (
	"testing"

	"github.com/hangulize/hangulize"
	"github.com/stretchr/testify/assert"
)

type asisPhonemizer struct{}

func (asisPhonemizer) ID() string {
	return "asis"
}

func (asisPhonemizer) Phonemize(word string) string {
	return word
}

func TestPhonemizerRegistry(t *testing.T) {
	var ok bool

	// Not exists.
	_, ok = hangulize.GetPhonemizer("asis")
	assert.False(t, ok)

	// Successfully registered.
	ok = hangulize.UsePhonemizer(&asisPhonemizer{})
	assert.True(t, ok)

	// Already exists.
	ok = hangulize.UsePhonemizer(&asisPhonemizer{})
	assert.False(t, ok)

	// Found.
	p, ok := hangulize.GetPhonemizer("asis")
	assert.True(t, ok)
	assert.IsType(t, &asisPhonemizer{}, p)

	// Successfully deregistered.
	ok = hangulize.UnusePhonemizer("asis")
	assert.True(t, ok)

	// Not exists.
	ok = hangulize.UnusePhonemizer("asis")
	assert.False(t, ok)
}

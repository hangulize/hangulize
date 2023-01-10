package hangulize_test

import (
	"testing"

	"github.com/hangulize/hangulize"
	"github.com/hangulize/hangulize/phonemize/furigana"
	"github.com/hangulize/hangulize/phonemize/pinyin"
	"github.com/stretchr/testify/assert"
)

// -----------------------------------------------------------------------------
// Use all phonemizers automatically for test.

func init() {
	hangulize.UsePhonemizer(&furigana.P)
	hangulize.UsePhonemizer(&pinyin.P)
}

// -----------------------------------------------------------------------------

type myPhonemizer struct{}

func (myPhonemizer) ID() string {
	return "my"
}

func (myPhonemizer) Phonemize(word string) string {
	return word
}

func TestPhonemizerRegistry(t *testing.T) {
	var ok bool

	// Not exists.
	_, ok = hangulize.GetPhonemizer("my")
	assert.False(t, ok)

	// Successfully registered.
	ok = hangulize.UsePhonemizer(&myPhonemizer{})
	assert.True(t, ok)

	// Already exists.
	ok = hangulize.UsePhonemizer(&myPhonemizer{})
	assert.False(t, ok)

	// Found.
	p, ok := hangulize.GetPhonemizer("my")
	assert.True(t, ok)
	assert.IsType(t, &myPhonemizer{}, p)

	// Successfully deregistered.
	ok = hangulize.UnusePhonemizer("my")
	assert.True(t, ok)

	// Not exists.
	ok = hangulize.UnusePhonemizer("my")
	assert.False(t, ok)
}

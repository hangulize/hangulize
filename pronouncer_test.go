package hangulize

import (
	"testing"

	"github.com/hangulize/hangulize/pronounce/furigana"
	"github.com/hangulize/hangulize/pronounce/pinyin"
	"github.com/stretchr/testify/assert"
)

// -----------------------------------------------------------------------------
// Use all pronouncers automatically for test.

func init() {
	UsePronouncer(&furigana.P)
	UsePronouncer(&pinyin.P)
}

// -----------------------------------------------------------------------------

type myPronouncer struct{}

func (myPronouncer) ID() string {
	return "my"
}

func (myPronouncer) Pronounce(word string) string {
	return word
}

func TestPronouncerRegistry(t *testing.T) {
	var ok bool

	// Not exists.
	_, ok = GetPronouncer("my")
	assert.False(t, ok)

	// Successfully registered.
	ok = UsePronouncer(&myPronouncer{})
	assert.True(t, ok)

	// Already exists.
	ok = UsePronouncer(&myPronouncer{})
	assert.False(t, ok)

	// Found.
	p, ok := GetPronouncer("my")
	assert.True(t, ok)
	assert.IsType(t, &myPronouncer{}, p)

	// Successfully deregistered.
	ok = UnusePronouncer("my")
	assert.True(t, ok)

	// Not exists.
	ok = UnusePronouncer("my")
	assert.False(t, ok)
}

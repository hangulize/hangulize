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

func (asisPhonemizer) Load() error {
	return nil
}

func (asisPhonemizer) Phonemize(word string) (string, error) {
	return word, nil
}

func TestPhonemizerRegistry(t *testing.T) {
	var ok bool

	// Not exists.
	_, ok = hangulize.DefaultPhonemizer("asis")
	assert.False(t, ok)

	// Successfully registered.
	ok = hangulize.ImportPhonemizer(&asisPhonemizer{})
	assert.True(t, ok)

	// Already exists.
	ok = hangulize.ImportPhonemizer(&asisPhonemizer{})
	assert.False(t, ok)

	// Found.
	p, ok := hangulize.DefaultPhonemizer("asis")
	assert.True(t, ok)
	assert.IsType(t, &asisPhonemizer{}, p)

	// Successfully deregistered.
	ok = hangulize.UnimportPhonemizer("asis")
	assert.True(t, ok)

	// Not exists.
	ok = hangulize.UnimportPhonemizer("asis")
	assert.False(t, ok)
}

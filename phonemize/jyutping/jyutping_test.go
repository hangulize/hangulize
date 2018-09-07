package jyutping

import (
	"testing"

	"github.com/hangulize/hangulize"
	"github.com/stretchr/testify/assert"
)

func TestPhonemizer(t *testing.T) {
	assert.Implements(t, (*hangulize.Phonemizer)(nil), &P)
}

func TestJyutping(t *testing.T) {
	assert.Equal(t, "hing\u200bsip", P.Phonemize("興燮"))
}

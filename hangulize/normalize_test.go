package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeRoman(t *testing.T) {
	roman := &RomanNormalizer{}

	assert.Equal(t, "hello", Normalize("Hello", roman, nil))
	assert.Equal(t, "cafe", Normalize("Café", roman, nil))
	assert.Equal(t, "melee", Normalize("Mêlée", roman, nil))
	assert.Equal(t, " cafe, hi! ", Normalize(" Café, Hi! ", roman, nil))
}

func TestNormalizeKana(t *testing.T) {
	kana := &KanaNormalizer{}

	assert.Equal(t, "ア", Normalize("あ", kana, nil))
	assert.Equal(t, "ァ", Normalize("ぁ", kana, nil))
}

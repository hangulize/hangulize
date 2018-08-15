package furigana

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeLongVowels(t *testing.T) {
	assert.Equal(t, "オー", mergeLongVowels("オウ", 0))
	assert.Equal(t, "オー", mergeLongVowels("オオ", 0))
	assert.Equal(t, "ケー", mergeLongVowels("ケェ", 0))
	assert.Equal(t, "ケー", mergeLongVowels("ケイ", 0))
}

func TestMergeLongVowelsOffset(t *testing.T) {
	assert.Equal(t, "ホーオー", mergeLongVowels("ホウオウ", 0))
	assert.Equal(t, "ホーオー", mergeLongVowels("ホウオウ", 1))
	assert.Equal(t, "ホウオー", mergeLongVowels("ホウオウ", 2))
	assert.Equal(t, "ホウオー", mergeLongVowels("ホウオウ", 3))
	assert.Equal(t, "ホウオウ", mergeLongVowels("ホウオウ", 4))
}

func TestMergeLongVowelsAdditionalSounds(t *testing.T) {
	// http://www.guidetojapanese.org/learn/grammar/katakana
	assert.Equal(t, "ヴァヴィヴェヴォ", mergeLongVowels("ヴァヴィヴェヴォ", 0))
	assert.Equal(t, "ウィウェウォ", mergeLongVowels("ウィウェウォ", 0))
	assert.Equal(t, "ファフィフェフォ", mergeLongVowels("ファフィフェフォ", 0))
	assert.Equal(t, "チェ", mergeLongVowels("チェ", 0))
	assert.Equal(t, "ディドゥ", mergeLongVowels("ディドゥ", 0))
	assert.Equal(t, "ティトゥ", mergeLongVowels("ティトゥ", 0))
	assert.Equal(t, "ジェシェ", mergeLongVowels("ジェシェ", 0))
}

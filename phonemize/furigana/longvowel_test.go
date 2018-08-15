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

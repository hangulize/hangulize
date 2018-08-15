package furigana

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeLongVowels(t *testing.T) {
	assert.Equal(t, "オー", mergeLongVowels("オウ", -1))
	assert.Equal(t, "オー", mergeLongVowels("オオ", -1))
	assert.Equal(t, "ケー", mergeLongVowels("ケェ", -1))
	assert.Equal(t, "ケー", mergeLongVowels("ケイ", -1))
}

func TestMergeLongVowelsOffset(t *testing.T) {
	assert.Equal(t, "オウオー", mergeLongVowels("オウオウ", 2))
}

package furigana

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkLongVowels(t *testing.T) {
	assert.Equal(t, "オー", markLongVowels("オウ"))
	assert.Equal(t, "オー", markLongVowels("オオ"))
	assert.Equal(t, "ケー", markLongVowels("ケェ"))
	assert.Equal(t, "ケー", markLongVowels("ケイ"))
}

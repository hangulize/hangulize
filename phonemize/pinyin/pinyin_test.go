package pinyin

import (
	"testing"

	"github.com/hangulize/hangulize"
	"github.com/stretchr/testify/assert"
)

func TestPhonemizer(t *testing.T) {
	assert.Implements(t, (*hangulize.Phonemizer)(nil), &P)
	assert.Equal(t, "pinyin", P.ID())
}

func TestPinyin(t *testing.T) {
	assert.Equal(t, "pin\u200byin", P.Phonemize("拼音"))
}

func TestCJKUnified(t *testing.T) {
	assert.Equal(t, "li", P.Phonemize("李"))
	assert.Equal(t, "le", P.Phonemize("樂"))
}

func TestNonHanzi(t *testing.T) {
	assert.Equal(t, "Abc", P.Phonemize("Abc"))
	assert.Equal(t, "아", P.Phonemize("아"))
}

func TestHanziAndNonHanzi(t *testing.T) {
	assert.Equal(t, "아\u200bpin\u200byin\u200bAbc", P.Phonemize("아拼音Abc"))
}

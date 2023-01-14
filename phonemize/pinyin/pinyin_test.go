package pinyin_test

import (
	"testing"

	"github.com/hangulize/hangulize"
	"github.com/hangulize/hangulize/phonemize/pinyin"
	"github.com/stretchr/testify/assert"
)

func TestPhonemizer(t *testing.T) {
	assert.Implements(t, (*hangulize.Phonemizer)(nil), &pinyin.P)
	assert.Equal(t, "pinyin", pinyin.P.ID())
}

func TestPinyin(t *testing.T) {
	assert.Equal(t, "pin\u200byin", pinyin.P.Phonemize("拼音"))
}

func TestCJKUnified(t *testing.T) {
	assert.Equal(t, "li", pinyin.P.Phonemize("李"))
	assert.Equal(t, "le", pinyin.P.Phonemize("樂"))
}

func TestNonHanzi(t *testing.T) {
	assert.Equal(t, "Abc", pinyin.P.Phonemize("Abc"))
	assert.Equal(t, "아", pinyin.P.Phonemize("아"))
}

func TestHanziAndNonHanzi(t *testing.T) {
	assert.Equal(t, "아\u200bpin\u200byin\u200bAbc", pinyin.P.Phonemize("아拼音Abc"))
}

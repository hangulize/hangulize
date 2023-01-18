package pinyin_test

import (
	"testing"

	"github.com/hangulize/hangulize"
	"github.com/hangulize/hangulize/phonemize/pinyin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPhonemizer(t *testing.T) {
	assert.Implements(t, (*hangulize.Phonemizer)(nil), &pinyin.P)
	assert.Equal(t, "pinyin", pinyin.P.ID())
}

func mustPhonemize(t *testing.T, word string) string {
	result, err := pinyin.P.Phonemize(word)
	require.NoError(t, err)
	return result
}

func TestPinyin(t *testing.T) {
	assert.Equal(t, "pin\u200byin", mustPhonemize(t, "拼音"))
}

func TestCJKUnified(t *testing.T) {
	assert.Equal(t, "li", mustPhonemize(t, "李"))
	assert.Equal(t, "le", mustPhonemize(t, "樂"))
}

func TestNonHanzi(t *testing.T) {
	assert.Equal(t, "Abc", mustPhonemize(t, "Abc"))
	assert.Equal(t, "아", mustPhonemize(t, "아"))
}

func TestHanziAndNonHanzi(t *testing.T) {
	assert.Equal(t, "아\u200bpin\u200byin\u200bAbc", mustPhonemize(t, "아拼音Abc"))
}

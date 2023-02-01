package pinyin_test

import (
	"testing"

	"github.com/hangulize/hangulize/translit/pinyin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustTransliterate(t *testing.T, word string) string {
	result, err := pinyin.T.Transliterate(word)
	require.NoError(t, err)
	return result
}

func TestPinyin(t *testing.T) {
	assert.Equal(t, "pin\u200byin", mustTransliterate(t, "拼音"))
}

func TestCJKUnified(t *testing.T) {
	assert.Equal(t, "li", mustTransliterate(t, "李"))
	assert.Equal(t, "le", mustTransliterate(t, "樂"))
}

func TestNonHanzi(t *testing.T) {
	assert.Equal(t, "Abc", mustTransliterate(t, "Abc"))
	assert.Equal(t, "아", mustTransliterate(t, "아"))
}

func TestHanziAndNonHanzi(t *testing.T) {
	assert.Equal(t, "아\u200bpin\u200byin\u200bAbc", mustTransliterate(t, "아拼音Abc"))
}

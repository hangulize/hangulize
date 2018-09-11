package pinyin

import (
	"testing"

	"github.com/hangulize/hangulize"
	"github.com/stretchr/testify/assert"
)

func TestPhonemizer(t *testing.T) {
	assert.Implements(t, (*hangulize.Phonemizer)(nil), &P)
}

func TestPinyin(t *testing.T) {
	assert.Equal(t, "pin\u200byin", P.Phonemize("拼音"))
}

func TestCJKUnified(t *testing.T) {
	assert.Equal(t, "li", P.Phonemize("李"))
	assert.Equal(t, "le", P.Phonemize("樂"))
}

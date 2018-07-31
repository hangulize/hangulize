package furigana

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSeion(t *testing.T) {
	assert.Equal(t, 'く', toSeion('ぐ'))
	assert.Equal(t, 'は', toSeion('ぱ'))
	assert.Equal(t, 'サ', toSeion('ザ'))
	assert.Equal(t, 'ウ', toSeion('ヴ'))
}

func TestToDakuon(t *testing.T) {
	assert.Equal(t, 'ぐ', toDakuon('く'))
	assert.Equal(t, 'ば', toDakuon('は'))
	assert.Equal(t, 'ザ', toDakuon('サ'))
	assert.Equal(t, 'ヴ', toDakuon('ウ'))
}

func TestRepeatHiragana(t *testing.T) {
	assert.Equal(t, "かか", repeatKana("かゝ"))
	assert.Equal(t, "かが", repeatKana("かゞ"))
	assert.Equal(t, "がか", repeatKana("がゝ"))
	assert.Equal(t, "がが", repeatKana("がゞ"))
}

func TestRepeatKatakana(t *testing.T) {
	assert.Equal(t, "カカ", repeatKana("カヽ"))
	assert.Equal(t, "カガ", repeatKana("カヾ"))
	assert.Equal(t, "ガカ", repeatKana("ガヽ"))
	assert.Equal(t, "ガガ", repeatKana("ガヾ"))
}

func TestRepeatKanaMixed(t *testing.T) {
	assert.Equal(t, "かヽ", repeatKana("かヽ"))
	assert.Equal(t, "カゝ", repeatKana("カゝ"))
}

func TestRepeatKanaMulti(t *testing.T) {
	assert.Equal(t, "かがか", repeatKana("かゞゝ"))
}

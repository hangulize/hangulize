package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatinNormalize(t *testing.T) {
	latin := &_Latin{}
	assert.Equal(t, "h", latin.Normalize("H"))
	assert.Equal(t, "e", latin.Normalize("é"))
}

func TestKanaNormalize(t *testing.T) {
	kana := &_Kana{}
	assert.Equal(t, "ア", kana.Normalize("あ"))
	assert.Equal(t, "ァ", kana.Normalize("ぁ"))
}

package furigana

import (
	"testing"

	kagome "github.com/ikawaha/kagome.ipadic/tokenizer"
	"github.com/stretchr/testify/assert"
)

func TestTypewriterAuxiliary(t *testing.T) {
	toks := kagome.New().Tokenize("食べよう")
	tw := newTypewriter(toks)

	assert.Equal(t, "タベヨー", tw.Typewrite())
}

func TestTypewriterUnknown(t *testing.T) {
	toks := kagome.New().Tokenize("ホウオウ")
	tw := newTypewriter(toks)

	assert.Equal(t, "ホーオー", tw.Typewrite())
}

func TestTypewriterParticles(t *testing.T) {
	toks := kagome.New().Tokenize("それはあちへ")
	tw := newTypewriter(toks)

	// ハ instead of ワ, ヘ instead of エ
	assert.Equal(t, "ソレハアチヘ", tw.Typewrite())
}

func TestTypewriterReuse(t *testing.T) {
	toks := kagome.New().Tokenize("東京")
	tw := newTypewriter(toks)

	assert.Equal(t, -1, tw.cur)
	assert.Equal(t, "トーキョー", tw.Typewrite())
	assert.NotEqual(t, -1, tw.cur)
	assert.Equal(t, "トーキョー", tw.Typewrite())
}

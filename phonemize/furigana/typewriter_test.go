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

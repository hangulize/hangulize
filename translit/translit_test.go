package translit_test

import (
	"testing"

	"github.com/hangulize/hangulize"
	"github.com/hangulize/hangulize/translit"
	"github.com/stretchr/testify/assert"
)

func TestInstallInstance(t *testing.T) {
	h := hangulize.New(&hangulize.Spec{})

	translits := h.Translits()
	assert.Empty(t, translits)

	ok := translit.Install(h)
	assert.True(t, ok)

	translits = h.Translits()
	assert.Contains(t, translits, "furigana")
	assert.Contains(t, translits, "pinyin")

	ok = translit.Install(h)
	assert.False(t, ok)
}

type fakeTranslit struct {
	scheme string
}

func (t fakeTranslit) Scheme() string {
	return t.scheme
}

func (t fakeTranslit) Load() error {
	return nil
}

func (t fakeTranslit) Transliterate(word string) (string, error) {
	return word, nil
}

func TestInstallRollback(t *testing.T) {
	h := hangulize.New(&hangulize.Spec{})

	fakePinyin := fakeTranslit{"pinyin"}
	h.UseTranslit(fakePinyin)

	translits := h.Translits()
	assert.NotContains(t, translits, "furigana")
	assert.Contains(t, translits, "pinyin")

	ok := translit.Install(h)
	assert.False(t, ok)

	translits = h.Translits()
	assert.NotContains(t, translits, "furigana")
	assert.Contains(t, translits, "pinyin")
	assert.Equal(t, fakePinyin, translits["pinyin"])
}

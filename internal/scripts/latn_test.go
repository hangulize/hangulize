package scripts_test

import (
	"testing"

	"github.com/hangulize/hangulize/internal/scripts"
	"github.com/stretchr/testify/assert"
)

func TestLatnIs(t *testing.T) {
	s := scripts.Latn{}
	assert.True(t, s.Is('A'))  // U+0041 Latin Capital Letter A
	assert.False(t, s.Is('Α')) // U+0391 Α Greek Capital Letter Alpha
	assert.False(t, s.Is('А')) // U+0410 Cyrillic Capital Letter A
	assert.False(t, s.Is('ა')) // U+10D0 Georgian Letter An
	assert.False(t, s.Is('ア')) // U+30A2 Katakana Letter A
	assert.False(t, s.Is('ㅏ')) // U+314F Hangul Letter A
}

func TestLatnNormalize(t *testing.T) {
	latin := scripts.Latn{}
	assert.Equal(t, 'h', latin.Normalize('H'))
	assert.Equal(t, 'e', latin.Normalize('é'))
}

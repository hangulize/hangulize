package scripts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeLatn(t *testing.T) {
	latin := Latn{}
	assert.Equal(t, 'h', latin.Normalize('H'))
	assert.Equal(t, 'e', latin.Normalize('é'))
}

func TestNormalizeHrkt(t *testing.T) {
	hrkt := Hrkt{}
	assert.Equal(t, 'ア', hrkt.Normalize('あ'))
	assert.Equal(t, 'ァ', hrkt.Normalize('ぁ'))
}

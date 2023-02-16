package subword_test

import (
	"testing"

	"github.com/hangulize/hangulize/internal/subword"
	"github.com/stretchr/testify/assert"
)

func TestBuilderEmpty(t *testing.T) {
	var swBuf subword.Builder

	assert.Equal(t, "", swBuf.String())
	assert.Len(t, swBuf.Subwords(), 0)
}

func TestBuilder1Subword(t *testing.T) {
	var swBuf subword.Builder

	swBuf.Write(subword.Subword{"hello", 1})

	assert.Equal(t, "hello", swBuf.String())
	assert.Len(t, swBuf.Subwords(), 1)
}

func TestBuilderMergeSameLevel(t *testing.T) {
	var swBuf subword.Builder

	swBuf.Write(subword.Subword{"hello", 1})
	swBuf.Write(subword.Subword{"world", 1})

	assert.Equal(t, "helloworld", swBuf.String())
	assert.Len(t, swBuf.Subwords(), 1)
}

func TestBuilderDifferentLevel(t *testing.T) {
	var swBuf subword.Builder

	swBuf.Write(subword.Subword{"hello", 1})
	swBuf.Write(subword.Subword{"world", 2})

	assert.Equal(t, "helloworld", swBuf.String())
	assert.Len(t, swBuf.Subwords(), 2)
}

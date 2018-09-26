package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubwordsBuilderEmpty(t *testing.T) {
	var swBuf Builder

	assert.Equal(t, "", swBuf.String())
	assert.Len(t, swBuf.Subwords(), 0)
}

func TestSubwordsBuilder1Subword(t *testing.T) {
	var swBuf Builder

	swBuf.Append(Subword{"hello", 1})

	assert.Equal(t, "hello", swBuf.String())
	assert.Len(t, swBuf.Subwords(), 1)
}

func TestSubwordsBuilderMergeSameLevel(t *testing.T) {
	var swBuf Builder

	swBuf.Append(Subword{"hello", 1})
	swBuf.Append(Subword{"world", 1})

	assert.Equal(t, "helloworld", swBuf.String())
	assert.Len(t, swBuf.Subwords(), 1)
}

func TestSubwordsBuilderDifferentLevel(t *testing.T) {
	var swBuf Builder

	swBuf.Append(Subword{"hello", 1})
	swBuf.Append(Subword{"world", 2})

	assert.Equal(t, "helloworld", swBuf.String())
	assert.Len(t, swBuf.Subwords(), 2)
}

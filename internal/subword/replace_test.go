package subword_test

import (
	"testing"

	"github.com/hangulize/hangulize/internal/subword"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplacer(t *testing.T) {
	replacer := subword.NewReplacer("Hello, world", 0, 1)
	replacer.Replace(0, 5, "Bye")
	sws := replacer.Subwords()

	require.Len(t, sws, 2)
	assert.Equal(t, subword.New("Bye", 1), sws[0])
	assert.Equal(t, subword.New(", world", 0), sws[1])
}

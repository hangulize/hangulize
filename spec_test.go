package hangulize_test

import (
	"bytes"
	"testing"

	"github.com/hangulize/hangulize"
	"github.com/stretchr/testify/assert"
)

func TestEmptySpec(t *testing.T) {
	spec := mustParseSpec(``)
	assert.Equal(t, "", spec.Lang.ID)
}

func TestSpecSource(t *testing.T) {
	spec := mustParseSpec(`
	rewrite:
		"hello" -> "world"
	`)
	assert.Equal(t, `
	rewrite:
		"hello" -> "world"
	`, spec.Source)
}

func TestUnlimitedNegativeLookaround(t *testing.T) {
	_, err := hangulize.ParseSpec(bytes.NewBufferString(`
		rewrite:
			"{~.*}@_@" -> "o<-<"
	`))
	assert.Error(t, err)
}

package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptySpec(t *testing.T) {
	spec := mustParseSpec(``)
	assert.Equal(t, "", spec.Lang.ID)
}

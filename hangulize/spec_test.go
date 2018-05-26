package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptySpec(t *testing.T) {
	spec := parseSpec(``)
	assert.Equal(t, "", spec.Lang.ID)
}

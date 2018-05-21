package hangulize

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func parseSpec(src string) *Spec {
	spec, err := ParseSpec(strings.NewReader(strings.TrimSpace(src)))
	if err != nil {
		panic(err)
	}
	return spec
}

func TestEmptySpec(t *testing.T) {
	spec := parseSpec(``)
	assert.Equal(t, "", spec.Lang.ID)
}

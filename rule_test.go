package hangulize

import (
	"testing"

	"github.com/hangulize/hre"
	"github.com/stretchr/testify/assert"
)

func TestRuleString(t *testing.T) {
	p, _ := hre.NewPattern("foo", nil, nil)
	rp := hre.NewRPattern("bar", nil, nil)
	r := Rule{p, rp}
	assert.Equal(t, `"foo" -> "bar"`, r.String())
}

package hgl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	src := strings.TrimSpace(`
	foo:
		hello = "world"
	`)
	hgl, err := Parse(strings.NewReader(src))
	if err != nil {
		t.Fail()
	}

	foo := hgl["foo"]
	assert.Equal(t, foo.Name(), "foo")
	assert.Equal(t, foo.Get("hello"), "world")
}

package hgl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	src := strings.TrimSpace(`
	foo:
		# 코멘트
		hello = "world"
	`)
	p := NewParser(strings.NewReader(src))

	hgl, err := p.Parse()
	if err != nil {
		t.Error(err)
	}

	foo := hgl["foo"].(Dict)
	assert.Equal(t, []string{"world"}, foo["hello"])
}

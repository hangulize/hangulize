package hgl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func _newParser(src string) *parser {
	return newParser(strings.NewReader(strings.TrimSpace(src)))
}

func TestParseSinglePairList(t *testing.T) {
	p := _newParser(`
	foo:
		# 코멘트
		hello -> world
	`)

	hgl, err := p.parse()
	if err != nil {
		t.Error(err)
	}

	foo := hgl["foo"].(*ListSection).Array()

	assert.Equal(t, "hello", foo[0].Left())
	assert.Equal(t, []string{"world"}, foo[0].Right())
}

func TestParseSinglePairDict(t *testing.T) {
	p := _newParser(`
	foo:
		# 코멘트
		hello = "world"
	`)

	hgl, err := p.parse()
	if err != nil {
		t.Error(err)
	}

	foo := hgl["foo"].(*DictSection).Map()

	assert.Equal(t, []string{"world"}, foo["hello"])
}

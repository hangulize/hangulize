package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hangulize/hangulize/internal/scripts"
)

func TestLocalize(t *testing.T) {
	s := Spec{}
	h := New(&s)
	p := procedure{h, nil}

	s.script = scripts.Hrkt{}

	assert.Equal(t, "foo.", p.localize("foo。"))
	assert.Equal(t, ", foo, bar", p.localize("、foo、bar"))
	assert.Equal(t, "foo, bar,", p.localize("foo、 bar、"))
	assert.Equal(t, "'foo' bar", p.localize("「foo」bar"))
	assert.Equal(t, "foo: bar!", p.localize("foo：bar！"))
	assert.Equal(t, "foo!?", p.localize("foo！？"))
}

func TestLocalizeZWSP(t *testing.T) {
	s := Spec{}
	h := New(&s)
	p := procedure{h, nil}

	assert.Equal(t, "foo", p.localize("f\u200Bo\u200Bo"))
}

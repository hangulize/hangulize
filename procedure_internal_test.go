package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hangulize/hangulize/internal/scripts"
)

func TestLocalize(t *testing.T) {
	s := &Spec{}
	s.script = scripts.Hrkt{}
	p := newProcedure(s, nil, nil)

	assert.Equal(t, "foo.", p.localize("foo。"))
	assert.Equal(t, ", foo, bar", p.localize("、foo、bar"))
	assert.Equal(t, "foo, bar,", p.localize("foo、 bar、"))
	assert.Equal(t, "'foo' bar", p.localize("「foo」bar"))
	assert.Equal(t, "foo: bar!", p.localize("foo：bar！"))
	assert.Equal(t, "foo!?", p.localize("foo！？"))
}

func TestLocalizeZWSP(t *testing.T) {
	p := newProcedure(&Spec{}, nil, nil)
	assert.Equal(t, "foo", p.localize("f\u200Bo\u200Bo"))
}

package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalizePuncts(t *testing.T) {
	s := Spec{}
	s.script = _Kana{}
	h := NewHangulizer(&s)
	p := pipeline{h, nil}

	assert.Equal(t, "foo.", p.localizePuncts("foo。"))
	assert.Equal(t, ", foo, bar, baz,", p.localizePuncts("、foo、bar、 baz、"))
	assert.Equal(t, "'foo' bar", p.localizePuncts("「foo」bar"))
	assert.Equal(t, "foo: bar!", p.localizePuncts("foo：bar！"))
	assert.Equal(t, "foo!?", p.localizePuncts("foo！？"))
}

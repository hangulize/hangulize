package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransliteratePuncts(t *testing.T) {
	s := Spec{}
	s.script = _Kana{}
	h := NewHangulizer(&s)
	p := pipeline{h, nil}

	assert.Equal(t, "foo.", p.transliteratePuncts("foo。"))
	assert.Equal(t, ", foo, bar", p.transliteratePuncts("、foo、bar"))
	assert.Equal(t, "foo, bar,", p.transliteratePuncts("foo、 bar、"))
	assert.Equal(t, "'foo' bar", p.transliteratePuncts("「foo」bar"))
	assert.Equal(t, "foo: bar!", p.transliteratePuncts("foo：bar！"))
	assert.Equal(t, "foo!?", p.transliteratePuncts("foo！？"))
}

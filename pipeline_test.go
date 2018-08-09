package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransliterate(t *testing.T) {
	s := Spec{}
	h := NewHangulizer(&s)
	p := pipeline{h, nil}

	s.script = _Kana{}

	assert.Equal(t, "foo.", p.transliterate("foo。"))
	assert.Equal(t, ", foo, bar", p.transliterate("、foo、bar"))
	assert.Equal(t, "foo, bar,", p.transliterate("foo、 bar、"))
	assert.Equal(t, "'foo' bar", p.transliterate("「foo」bar"))
	assert.Equal(t, "foo: bar!", p.transliterate("foo：bar！"))
	assert.Equal(t, "foo!?", p.transliterate("foo！？"))
}

func TestTransliterateZWSP(t *testing.T) {
	s := Spec{}
	h := NewHangulizer(&s)
	p := pipeline{h, nil}

	assert.Equal(t, "foo", p.transliterate("f\u200Bo\u200Bo"))
}
